package transformer

import "fmt"

// Format represents a transformer output format.
type Format int

const (
	FormatPlainText Format = iota
	FormatHTML
)

// Factory creates and caches transformer instances.
type Factory struct {
	presets      *PresetContext
	transformers map[Format]Transformer
}

// NewFactory creates a new transformer factory with the given preset context.
func NewFactory(presets *PresetContext) *Factory {
	f := &Factory{
		presets:      presets,
		transformers: make(map[Format]Transformer),
	}

	// Register default transformers
	f.transformers[FormatPlainText] = NewPlainTextTransformer()
	f.transformers[FormatHTML] = NewHtmlTransformer(presets)

	return f
}

// Get returns a transformer by format.
func (f *Factory) Get(format Format) (Transformer, error) {
	if t, ok := f.transformers[format]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("unknown transformer format: %d", format)
}

// MustGet returns a transformer by format, panics if not found.
func (f *Factory) MustGet(format Format) Transformer {
	t, err := f.Get(format)
	if err != nil {
		panic(err)
	}
	return t
}

// Register adds a custom transformer to the factory.
func (f *Factory) Register(format Format, t Transformer) {
	f.transformers[format] = t
}

// Presets returns the preset context used by this factory.
func (f *Factory) Presets() *PresetContext {
	return f.presets
}
