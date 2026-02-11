package rigid

// Err returns the monadic error currently attached to the pipeline, if any.
func (p Pipeline[In, Out]) Err() error {
	return p.err
}

// HasError reports whether the pipeline is in an error state.
func (p Pipeline[In, Out]) HasError() bool {
	return p.err != nil
}

// MapError maps a pipeline error (attached or execution-time).
func (p Pipeline[In, Out]) MapError(mapper func(error) error) Pipeline[In, Out] {
	if mapper == nil {
		p.err = ErrNilErrorMapper
		p.trans = nil
		return p
	}
	if p.err != nil {
		p.err = mapper(p.err)
		return p
	}
	if p.trans == nil {
		p.err = ErrNilTransformer
		return p
	}
	original := p.trans
	p.trans = TransformerFunc[In, Out](func(in []In) ([]Out, error) {
		out, err := original.Transform(in)
		if err != nil {
			return nil, mapper(err)
		}
		return out, nil
	})
	return p
}

// Recover handles a pipeline error (attached or execution-time) by switching to a fallback pipeline.
func (p Pipeline[In, Out]) Recover(handler func(error) Pipeline[In, Out]) Pipeline[In, Out] {
	if handler == nil {
		p.err = ErrNilErrorHandler
		p.trans = nil
		return p
	}
	if p.err != nil {
		return handler(p.err)
	}
	if p.trans == nil {
		return Pipeline[In, Out]{err: ErrNilTransformer}
	}
	original := p.trans
	p.trans = TransformerFunc[In, Out](func(in []In) ([]Out, error) {
		out, err := original.Transform(in)
		if err == nil {
			return out, nil
		}
		return handler(err).Execute()
	})
	return p
}

// OrElse replaces an errored pipeline with the provided fallback pipeline.
func (p Pipeline[In, Out]) OrElse(fallback Pipeline[In, Out]) Pipeline[In, Out] {
	return p.Recover(func(error) Pipeline[In, Out] { return fallback })
}
