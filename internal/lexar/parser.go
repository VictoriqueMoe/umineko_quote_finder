package lexar

import (
	"strconv"
	"strings"
	"unicode"

	"umineko_quote/internal/lexar/ast"
)

type Parser struct {
	tokens []ast.Token
	pos    int
}

func Parse(input string) *ast.Script {
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()
	parser := &Parser{tokens: tokens, pos: 0}
	return parser.parse()
}

func (p *Parser) parse() *ast.Script {
	var lines []ast.Line
	for !p.isAtEnd() {
		line := p.parseLine()
		if line != nil {
			lines = append(lines, line)
		}
	}
	return &ast.Script{Lines: lines}
}

func (p *Parser) peek() ast.Token {
	if p.pos >= len(p.tokens) {
		return ast.Token{Type: ast.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() ast.Token {
	tok := p.peek()
	if !p.isAtEnd() {
		p.pos++
	}
	return tok
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.TokenEOF
}

func (p *Parser) skipNewlines() {
	for p.peek().Type == ast.TokenNewline {
		p.advance()
	}
}

func (p *Parser) parseLine() ast.Line {
	p.skipNewlines()
	if p.isAtEnd() {
		return nil
	}

	tok := p.peek()

	switch tok.Type {
	case ast.TokenComment:
		p.advance()
		return &ast.CommentLine{Text: tok.Value, Pos: tok}

	case ast.TokenLabel:
		p.advance()
		return &ast.LabelLine{Name: tok.Value, Pos: tok}

	case ast.TokenCommand:
		return p.parseCommand(tok)

	default:
		p.advance()
		return nil
	}
}

func (p *Parser) parseCommand(tok ast.Token) ast.Line {
	p.advance()

	switch tok.Value {
	case "d", "d2":
		return p.parseDialogue(tok)
	case "preset_define":
		return p.parsePresetDefine(tok)
	case "new_episode":
		return p.parseEpisodeMarker(tok, "episode")
	case "new_tea":
		return p.parseEpisodeMarker(tok, "tea")
	case "new_ura":
		return p.parseEpisodeMarker(tok, "ura")
	default:
		return p.parseGenericCommand(tok)
	}
}

func (p *Parser) parseDialogue(cmdTok ast.Token) *ast.DialogueLine {
	var elements []ast.DialogueElement

	for {
		tok := p.peek()
		if tok.Type == ast.TokenEOF || tok.Type == ast.TokenNewline {
			break
		}

		switch tok.Type {
		case ast.TokenBacktick:
			p.advance()
		case ast.TokenText:
			p.advance()
			elements = append(elements, &ast.PlainText{Text: tok.Value, Pos: tok})
		case ast.TokenFormatTag:
			p.advance()
			elements = append(elements, p.parseFormatTagContent(tok))
		case ast.TokenInlineCommand:
			p.advance()
			elements = append(elements, p.parseInlineCommandContent(tok))
		default:
			p.advance()
		}
	}

	return &ast.DialogueLine{Command: cmdTok.Value, Content: elements, Pos: cmdTok}
}

func (p *Parser) parseFormatTagContent(tok ast.Token) ast.DialogueElement {
	tagName, param, content := p.parseFormatTag(tok.Value)

	if param == "" && content == "" {
		switch tagName {
		case "n", "0", "qt", "ob", "eb", "os", "es", "-", "t", "parallel":
			return &ast.SpecialChar{Name: tagName, Pos: tok}
		}
	}

	var nested []ast.DialogueElement
	if content != "" {
		nested = p.parseNestedContent(content)
	}

	return &ast.FormatTag{Name: tagName, Param: param, Content: nested, Pos: tok}
}

func (p *Parser) parseNestedContent(content string) []ast.DialogueElement {
	var elements []ast.DialogueElement
	lexer := NewLexer("d `" + content + "`")

	for {
		tok := lexer.NextToken()
		if tok.Type == ast.TokenEOF {
			break
		}

		switch tok.Type {
		case ast.TokenText:
			if tok.Value != "" {
				elements = append(elements, &ast.PlainText{Text: tok.Value, Pos: tok})
			}
		case ast.TokenFormatTag:
			elements = append(elements, p.parseFormatTagContent(tok))
		case ast.TokenInlineCommand:
			elements = append(elements, p.parseInlineCommandContent(tok))
		case ast.TokenCommand, ast.TokenBacktick:
		}
	}

	return elements
}

func (p *Parser) parseInlineCommandContent(tok ast.Token) ast.DialogueElement {
	cmd, args := p.parseInlineCommand(tok.Value)

	switch cmd {
	case "@", "\\", "|":
		return &ast.ClickWait{Type: cmd, Pos: tok}
	case "!w":
		dur, _ := strconv.Atoi(args)
		return &ast.TimedWait{Skippable: false, Duration: dur, Pos: tok}
	case "!d":
		dur, _ := strconv.Atoi(args)
		return &ast.TimedWait{Skippable: true, Duration: dur, Pos: tok}
	case "lv":
		return p.parseVoiceCommand(args, tok)
	default:
		return &ast.InlineCommand{Command: cmd, Args: args, Pos: tok}
	}
}

func (p *Parser) parseVoiceCommand(args string, tok ast.Token) *ast.VoiceCommand {
	parts := strings.Split(args, "*")
	vc := &ast.VoiceCommand{Pos: tok}

	if len(parts) >= 1 {
		vc.Channel, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	}
	if len(parts) >= 2 {
		vc.CharacterID = strings.Trim(parts[1], `"`)
	}
	if len(parts) >= 3 {
		vc.AudioID = strings.Trim(parts[2], `"`)
	}

	return vc
}

func (p *Parser) parsePresetDefine(tok ast.Token) *ast.PresetDefineLine {
	preset := &ast.PresetDefineLine{Pos: tok}

	var nums []int
	var colour string

	for {
		next := p.peek()
		if next.Type == ast.TokenEOF || next.Type == ast.TokenNewline {
			break
		}

		if next.Type == ast.TokenNumber {
			p.advance()
			n, _ := strconv.Atoi(next.Value)
			nums = append(nums, n)
		} else if next.Type == ast.TokenString && strings.HasPrefix(next.Value, "#") {
			p.advance()
			if colour == "" {
				colour = next.Value
			}
		} else if next.Type == ast.TokenComma {
			p.advance()
		} else {
			p.advance()
		}
	}

	if len(nums) >= 1 {
		preset.ID = nums[0]
	}
	if len(nums) >= 2 {
		preset.FontID = nums[1]
	}
	if len(nums) >= 3 {
		preset.Size = nums[2]
	}
	if len(nums) > 3 {
		preset.Flags = nums[3:]
	}
	preset.Colour = colour

	return preset
}

func (p *Parser) parseEpisodeMarker(tok ast.Token, markerType string) *ast.EpisodeMarkerLine {
	marker := &ast.EpisodeMarkerLine{Type: markerType, Pos: tok}

	if p.peek().Type == ast.TokenNumber {
		numTok := p.advance()
		marker.Episode, _ = strconv.Atoi(numTok.Value)
	}

	return marker
}

func (p *Parser) parseGenericCommand(tok ast.Token) *ast.CommandLine {
	cmd := &ast.CommandLine{Command: tok.Value, Pos: tok}

	for {
		next := p.peek()
		if next.Type == ast.TokenEOF || next.Type == ast.TokenNewline {
			break
		}
		if next.Type == ast.TokenComma {
			p.advance()
			continue
		}
		cmd.Args = append(cmd.Args, p.advance())
	}

	return cmd
}

func (*Parser) parseFormatTag(value string) (tagName, param, content string) {
	if !strings.Contains(value, ":") {
		return value, "", ""
	}

	parts := strings.SplitN(value, ":", 3)
	tagName = parts[0]

	if len(parts) == 2 {
		content = parts[1]
	} else if len(parts) == 3 {
		// If the candidate param contains '{', it's actually nested content,
		// not a parameter. Rejoin as content-only (tag:content).
		if strings.Contains(parts[1], "{") {
			content = parts[1] + ":" + parts[2]
		} else {
			param = parts[1]
			content = parts[2]
		}
	}

	return tagName, param, content
}

func (*Parser) parseInlineCommand(value string) (cmd, args string) {
	value = strings.TrimSpace(value)

	if len(value) == 1 {
		return value, ""
	}

	if strings.HasPrefix(value, "!w") || strings.HasPrefix(value, "!d") {
		return value[:2], value[2:]
	}

	idx := strings.IndexFunc(value, unicode.IsSpace)
	if idx == -1 {
		return value, ""
	}
	return value[:idx], strings.TrimSpace(value[idx+1:])
}
