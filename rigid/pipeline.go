package rigid

// Pipeline stores source input plus a composed transformer.
// Then(...) is a free function (not a method) to keep full type changes.
type Pipeline[In, Out any] struct {
	source []In
	trans  Transformer[In, Out]
	err    error
}

// From starts a typed pipeline from a source slice.
func From[T any](source []T) Pipeline[T, T] {
	return Pipeline[T, T]{
		source: source,
		trans:  Identity[T](),
	}
}

// Then appends a next stage to an existing pipeline.
func Then[A, B, C any](p Pipeline[A, B], next Transformer[B, C]) Pipeline[A, C] {
	if p.err != nil {
		return Pipeline[A, C]{err: p.err}
	}
	if next == nil {
		return Pipeline[A, C]{err: ErrNilTransformer}
	}
	return Pipeline[A, C]{
		source: p.source,
		trans:  Compose(p.trans, next),
	}
}

// Execute materializes the pipeline.
func (p Pipeline[In, Out]) Execute() ([]Out, error) {
	if p.err != nil {
		return nil, p.err
	}
	if p.trans == nil {
		return nil, ErrNilTransformer
	}
	return Run(p.source, p.trans)
}

// ToSlice is an alias for Execute and mirrors the default package terminology.
func (p Pipeline[In, Out]) ToSlice() ([]Out, error) {
	return p.Execute()
}
