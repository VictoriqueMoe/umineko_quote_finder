package lexar

import (
	"fmt"

	"umineko_quote/internal/lexar/ast"
)

type (
	Severity int

	ValidationError struct {
		Severity Severity
		Line     int
		Column   int
		Message  string
	}

	validator struct {
		errors []ValidationError
	}
)

const (
	SeverityWarning Severity = iota
	SeverityError
)

var knownFormatTags = map[string]bool{
	"p": true, "preset": true,
	"c": true, "color": true, "colour": true,
	"i": true, "italic": true,
	"ruby": true, "h": true,
	"y": true, "n": true,
	"0": true, "qt": true,
	"ob": true, "eb": true,
	"os": true, "es": true,
	"-": true, "t": true,
	"parallel": true,
	"f":        true, "a": true, "e": true,
	"m": true, "b": true, "o": true, "g": true,
	"nobr": true, "nobreak": true,
	"loghint": true,
	"\u2010":  true,
}

func Validate(script *ast.Script) []ValidationError {
	v := &validator{}
	v.walk(script)
	return v.errors
}

func (v *validator) addWarning(pos ast.Token, format string, args ...any) {
	v.errors = append(v.errors, ValidationError{
		Severity: SeverityWarning,
		Line:     pos.Line,
		Column:   pos.Column,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (v *validator) addError(pos ast.Token, format string, args ...any) {
	v.errors = append(v.errors, ValidationError{
		Severity: SeverityError,
		Line:     pos.Line,
		Column:   pos.Column,
		Message:  fmt.Sprintf(format, args...),
	})
}

func (v *validator) walk(script *ast.Script) {
	for _, line := range script.Lines {
		switch l := line.(type) {
		case *ast.DialogueLine:
			v.validateDialogue(l)
		case *ast.EpisodeMarkerLine:
			if l.Episode == 0 {
				v.addError(l.Pos, "episode marker missing episode number")
			}
		}
	}
}

func (v *validator) validateDialogue(d *ast.DialogueLine) {
	v.validateElements(d.Content)
}

func (v *validator) validateElements(elements []ast.DialogueElement) {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.FormatTag:
			v.validateFormatTag(el)
		case *ast.VoiceCommand:
			v.validateVoiceCommand(el)
		}
	}
}

func (v *validator) validateFormatTag(tag *ast.FormatTag) {
	if tag.Name == "" {
		v.addWarning(tag.Pos, "empty format tag name (expected tag name before ':')")
	} else if !knownFormatTags[tag.Name] {
		v.addWarning(tag.Pos, "unknown format tag %q", tag.Name)
	}
	v.validateElements(tag.Content)
}

func (v *validator) validateVoiceCommand(vc *ast.VoiceCommand) {
	if vc.CharacterID == "" {
		v.addError(vc.Pos, "voice command missing character ID")
	}
	if vc.AudioID == "" {
		v.addError(vc.Pos, "voice command missing audio ID")
	}
}

func (e ValidationError) String() string {
	level := "warning"
	if e.Severity == SeverityError {
		level = "error"
	}
	return fmt.Sprintf("%s line %d:%d: %s", level, e.Line, e.Column, e.Message)
}
