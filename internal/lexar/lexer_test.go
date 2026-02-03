package lexar

import (
	"testing"

	"umineko_quote/internal/lexar/ast"
)

func tokenize(input string) []ast.Token {
	return NewLexer(input).Tokenize()
}

func TestLexer_SimpleCommand(t *testing.T) {
	input := "new_episode 1"
	tokens := tokenize(input)

	expected := []struct {
		typ   ast.TokenType
		value string
	}{
		{ast.TokenCommand, "new_episode"},
		{ast.TokenNumber, "1"},
		{ast.TokenEOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("token count: got %d, want %d", len(tokens), len(expected))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.typ {
			t.Errorf("token[%d] type: got %v, want %v", i, tokens[i].Type, exp.typ)
		}
		if tokens[i].Value != exp.value {
			t.Errorf("token[%d] value: got %q, want %q", i, tokens[i].Value, exp.value)
		}
	}
}

func TestLexer_PresetDefine(t *testing.T) {
	input := `preset_define 1,1,-1,#FF0000,0,0,0,0,0`
	tokens := tokenize(input)

	if tokens[0].Type != ast.TokenCommand || tokens[0].Value != "preset_define" {
		t.Errorf("expected preset_define command, got %v %q", tokens[0].Type, tokens[0].Value)
	}
	if tokens[1].Type != ast.TokenNumber || tokens[1].Value != "1" {
		t.Errorf("expected number 1, got %v %q", tokens[1].Type, tokens[1].Value)
	}
	if tokens[2].Type != ast.TokenComma {
		t.Errorf("expected comma, got %v", tokens[2].Type)
	}
}

func TestLexer_Comment(t *testing.T) {
	input := "; this is a comment"
	tokens := tokenize(input)

	if len(tokens) < 2 {
		t.Fatalf("expected at least 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != ast.TokenComment {
		t.Errorf("expected comment, got %v", tokens[0].Type)
	}
	if tokens[0].Value != " this is a comment" {
		t.Errorf("comment value: got %q, want %q", tokens[0].Value, " this is a comment")
	}
}

func TestLexer_Label(t *testing.T) {
	input := "*episode1_start"
	tokens := tokenize(input)

	if tokens[0].Type != ast.TokenLabel {
		t.Errorf("expected label, got %v", tokens[0].Type)
	}
	if tokens[0].Value != "episode1_start" {
		t.Errorf("label value: got %q, want %q", tokens[0].Value, "episode1_start")
	}
}

func TestLexer_DialogueLine(t *testing.T) {
	input := "d `Hello, world!`"
	tokens := tokenize(input)

	expected := []struct {
		typ   ast.TokenType
		value string
	}{
		{ast.TokenCommand, "d"},
		{ast.TokenBacktick, "`"},
		{ast.TokenText, "Hello, world!"},
		{ast.TokenBacktick, "`"},
		{ast.TokenEOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("token count: got %d, want %d\ntokens: %+v", len(tokens), len(expected), tokens)
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.typ {
			t.Errorf("token[%d] type: got %v, want %v", i, tokens[i].Type, exp.typ)
		}
		if tokens[i].Value != exp.value {
			t.Errorf("token[%d] value: got %q, want %q", i, tokens[i].Value, exp.value)
		}
	}
}

func TestLexer_DialogueWithInlineCommand(t *testing.T) {
	input := `d [lv 0*"10"*"10100001"]` + "`Test`"
	tokens := tokenize(input)

	var foundInline bool
	var inlineValue string
	for _, tok := range tokens {
		if tok.Type == ast.TokenInlineCommand {
			foundInline = true
			inlineValue = tok.Value
			break
		}
	}

	if !foundInline {
		t.Fatal("expected to find inline command token")
	}
	if inlineValue != `lv 0*"10"*"10100001"` {
		t.Errorf("inline command value: got %q, want %q", inlineValue, `lv 0*"10"*"10100001"`)
	}
}

func TestLexer_DialogueWithFormatTag(t *testing.T) {
	input := "d `This is {c:FF0000:red text} here.`"
	tokens := tokenize(input)

	var foundTag bool
	var tagValue string
	for _, tok := range tokens {
		if tok.Type == ast.TokenFormatTag {
			foundTag = true
			tagValue = tok.Value
			break
		}
	}

	if !foundTag {
		t.Fatal("expected to find format tag token")
	}
	if tagValue != "c:FF0000:red text" {
		t.Errorf("format tag value: got %q, want %q", tagValue, "c:FF0000:red text")
	}
}

func TestLexer_NestedFormatTags(t *testing.T) {
	input := "d `{p:1:{c:FF0000:nested}}`"
	tokens := tokenize(input)

	var foundTag bool
	var tagValue string
	for _, tok := range tokens {
		if tok.Type == ast.TokenFormatTag {
			foundTag = true
			tagValue = tok.Value
			break
		}
	}

	if !foundTag {
		t.Fatal("expected to find format tag token")
	}
	if tagValue != "p:1:{c:FF0000:nested}" {
		t.Errorf("format tag value: got %q, want %q", tagValue, "p:1:{c:FF0000:nested}")
	}
}

func TestLexer_SpecialCharTag(t *testing.T) {
	input := "d `Line one{n}Line two`"
	tokens := tokenize(input)

	var tags []string
	for _, tok := range tokens {
		if tok.Type == ast.TokenFormatTag {
			tags = append(tags, tok.Value)
		}
	}

	if len(tags) != 1 || tags[0] != "n" {
		t.Errorf("expected format tag 'n', got %v", tags)
	}
}

func TestLexer_RealDialogueLine(t *testing.T) {
	input := `d2 [lv 0*"10"*"10100001"]` + "`" + `{p:1:I'll make you understand!}` + "`[@]"

	tokens := tokenize(input)

	expectedTypes := []ast.TokenType{
		ast.TokenCommand,       // d2
		ast.TokenInlineCommand, // lv 0*"10"*"10100001"
		ast.TokenBacktick,      // `
		ast.TokenFormatTag,     // p:1:I'll make you understand!
		ast.TokenBacktick,      // `
		ast.TokenInlineCommand, // @
		ast.TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("token count: got %d, want %d\ntokens: %+v", len(tokens), len(expectedTypes), tokens)
	}

	for i, expType := range expectedTypes {
		if tokens[i].Type != expType {
			t.Errorf("token[%d] type: got %v, want %v (value: %q)", i, tokens[i].Type, expType, tokens[i].Value)
		}
	}

	if tokens[0].Value != "d2" {
		t.Errorf("command value: got %q, want 'd2'", tokens[0].Value)
	}
	if tokens[1].Value != `lv 0*"10"*"10100001"` {
		t.Errorf("inline value: got %q", tokens[1].Value)
	}
	if tokens[3].Value != "p:1:I'll make you understand!" {
		t.Errorf("format tag value: got %q", tokens[3].Value)
	}
	if tokens[5].Value != "@" {
		t.Errorf("click wait value: got %q", tokens[5].Value)
	}
}
