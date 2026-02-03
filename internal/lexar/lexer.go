package lexar

import "umineko_quote/internal/lexar/ast"

type Lexer struct {
	input    string
	pos      int
	line     int
	col      int
	inDialog bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, line: 1, col: 1}
}

func (l *Lexer) NextToken() ast.Token {
	if l.inDialog {
		return l.nextDialogToken()
	}
	return l.nextLineToken()
}

func (l *Lexer) Tokenize() []ast.Token {
	var tokens []ast.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == ast.TokenEOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekN(n int) string {
	end := l.pos + n
	if end > len(l.input) {
		end = len(l.input)
	}
	return l.input[l.pos:end]
}

func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.peek() == ' ' || l.peek() == '\t' {
		l.advance()
	}
}

func (l *Lexer) makeToken(typ ast.TokenType, value string) ast.Token {
	return ast.Token{Type: typ, Value: value, Line: l.line, Column: l.col - len(value)}
}

func (l *Lexer) nextLineToken() ast.Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return l.makeToken(ast.TokenEOF, "")
	}

	ch := l.peek()
	startLine, startCol := l.line, l.col

	if ch == '\n' {
		l.advance()
		return ast.Token{Type: ast.TokenNewline, Value: "\n", Line: startLine, Column: startCol}
	}

	if ch == '\r' {
		l.advance()
		if l.peek() == '\n' {
			l.advance()
		}
		return ast.Token{Type: ast.TokenNewline, Value: "\n", Line: startLine, Column: startCol}
	}

	if ch == ';' {
		return l.scanComment()
	}

	if ch == '*' {
		return l.scanLabel()
	}

	if ch == ',' {
		l.advance()
		return ast.Token{Type: ast.TokenComma, Value: ",", Line: startLine, Column: startCol}
	}

	if ch == '-' || (ch >= '0' && ch <= '9') {
		return l.scanNumber()
	}

	if ch == '"' {
		return l.scanQuotedString()
	}

	if ch == '`' {
		l.inDialog = true
		l.advance()
		return ast.Token{Type: ast.TokenBacktick, Value: "`", Line: startLine, Column: startCol}
	}

	if ch == '[' {
		return l.scanInlineCommand()
	}

	if ch == '#' {
		return l.scanHashColour()
	}

	if isIdentStart(ch) {
		tok := l.scanIdentifier()
		if tok.Value == "d" || tok.Value == "d2" {
			l.inDialog = true
		}
		return tok
	}

	l.advance()
	return l.nextLineToken()
}

func (l *Lexer) nextDialogToken() ast.Token {
	for (l.peek() == ' ' || l.peek() == '\t') && l.peekN(2)[1:2] != "" {
		next := l.peekN(2)
		if len(next) >= 2 && (next[1] == '`' || next[1] == '[' || next[1] == '\n' || next[1] == '\r') {
			l.advance()
		} else if len(next) == 1 {
			l.advance()
		} else {
			break
		}
	}

	if l.pos >= len(l.input) {
		l.inDialog = false
		return l.makeToken(ast.TokenEOF, "")
	}

	ch := l.peek()
	startLine, startCol := l.line, l.col

	if ch == '`' {
		l.advance()
		next := l.peek()
		if next == '\n' || next == '\r' || next == 0 {
			l.inDialog = false
		}
		return ast.Token{Type: ast.TokenBacktick, Value: "`", Line: startLine, Column: startCol}
	}

	if ch == '\n' || ch == '\r' {
		l.inDialog = false
		return l.nextLineToken()
	}

	if ch == '[' {
		return l.scanInlineCommand()
	}

	if ch == '{' {
		return l.scanFormatTag()
	}

	return l.scanDialogText()
}

func (l *Lexer) scanComment() ast.Token {
	startLine, startCol := l.line, l.col
	l.advance()
	start := l.pos
	for l.peek() != '\n' && l.peek() != '\r' && l.peek() != 0 {
		l.advance()
	}
	return ast.Token{Type: ast.TokenComment, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func (l *Lexer) scanLabel() ast.Token {
	startLine, startCol := l.line, l.col
	l.advance()
	start := l.pos
	for isIdentChar(l.peek()) {
		l.advance()
	}
	return ast.Token{Type: ast.TokenLabel, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func (l *Lexer) scanNumber() ast.Token {
	startLine, startCol := l.line, l.col
	start := l.pos
	if l.peek() == '-' {
		l.advance()
	}
	for l.peek() >= '0' && l.peek() <= '9' {
		l.advance()
	}
	return ast.Token{Type: ast.TokenNumber, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func (l *Lexer) scanQuotedString() ast.Token {
	startLine, startCol := l.line, l.col
	l.advance()
	start := l.pos
	for l.peek() != '"' && l.peek() != '\n' && l.peek() != 0 {
		l.advance()
	}
	value := l.input[start:l.pos]
	if l.peek() == '"' {
		l.advance()
	}
	return ast.Token{Type: ast.TokenString, Value: value, Line: startLine, Column: startCol}
}

func (l *Lexer) scanHashColour() ast.Token {
	startLine, startCol := l.line, l.col
	start := l.pos
	l.advance()
	for isHexDigit(l.peek()) {
		l.advance()
	}
	return ast.Token{Type: ast.TokenString, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func (l *Lexer) scanIdentifier() ast.Token {
	startLine, startCol := l.line, l.col
	start := l.pos
	for isIdentChar(l.peek()) {
		l.advance()
	}
	return ast.Token{Type: ast.TokenCommand, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func (l *Lexer) scanInlineCommand() ast.Token {
	startLine, startCol := l.line, l.col
	l.advance()
	start := l.pos
	depth := 1
	for depth > 0 && l.peek() != 0 && l.peek() != '\n' {
		ch := l.advance()
		if ch == '[' {
			depth++
		} else if ch == ']' {
			depth--
		}
	}
	value := l.input[start : l.pos-1]
	return ast.Token{Type: ast.TokenInlineCommand, Value: value, Line: startLine, Column: startCol}
}

func (l *Lexer) scanFormatTag() ast.Token {
	startLine, startCol := l.line, l.col
	l.advance()
	start := l.pos
	depth := 1
	for depth > 0 && l.peek() != 0 && l.peek() != '\n' {
		ch := l.advance()
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
		}
	}
	value := l.input[start : l.pos-1]
	return ast.Token{Type: ast.TokenFormatTag, Value: value, Line: startLine, Column: startCol}
}

func (l *Lexer) scanDialogText() ast.Token {
	startLine, startCol := l.line, l.col
	start := l.pos
	for {
		ch := l.peek()
		if ch == 0 || ch == '\n' || ch == '\r' || ch == '`' || ch == '[' || ch == '{' {
			break
		}
		l.advance()
	}
	return ast.Token{Type: ast.TokenText, Value: l.input[start:l.pos], Line: startLine, Column: startCol}
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentChar(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '_'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
