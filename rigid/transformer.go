package rigid

import "errors"

// Transformer converts a collection of In into a collection of Out.
// Errors stop the pipeline immediately.
type Transformer[In, Out any] interface {
	Transform([]In) ([]Out, error)
}

// TransformerFunc adapts a function to a Transformer.
type TransformerFunc[In, Out any] func([]In) ([]Out, error)

func (f TransformerFunc[In, Out]) Transform(in []In) ([]Out, error) {
	return f(in)
}

var (
	// ErrNilTransformer indicates a nil transformer was supplied to Compose/Then.
	ErrNilTransformer = errors.New("nil transformer")
	// ErrNilPredicate indicates a nil predicate was supplied.
	ErrNilPredicate = errors.New("nil predicate")
	// ErrNilSelector indicates a nil selector/transform function was supplied.
	ErrNilSelector = errors.New("nil selector")
	// ErrNoMatch is returned when no element satisfies the predicate.
	ErrNoMatch = errors.New("no element satisfies the predicate")
	// ErrNilErrorMapper indicates a nil mapper was supplied to MapError.
	ErrNilErrorMapper = errors.New("nil error mapper")
	// ErrNilErrorHandler indicates a nil handler was supplied to Recover.
	ErrNilErrorHandler = errors.New("nil error handler")
	// ErrNilComparator indicates a nil comparator was supplied for ordered pipelines.
	ErrNilComparator = errors.New("nil comparator")
	// ErrNilSortFunc indicates a nil sorter was supplied via UsingSorter or UsingSortFunc.
	ErrNilSortFunc = errors.New("nil sort function")
)

func failed[In, Out any](err error) Transformer[In, Out] {
	return TransformerFunc[In, Out](func([]In) ([]Out, error) { return nil, err })
}

// Run executes a transformer against an input collection.
func Run[In, Out any](in []In, t Transformer[In, Out]) ([]Out, error) {
	if t == nil {
		return nil, ErrNilTransformer
	}
	return t.Transform(in)
}

// Compose combines two transformers into one.
func Compose[A, B, C any](ab Transformer[A, B], bc Transformer[B, C]) Transformer[A, C] {
	if ab == nil || bc == nil {
		return failed[A, C](ErrNilTransformer)
	}
	return TransformerFunc[A, C](func(in []A) ([]C, error) {
		mid, err := ab.Transform(in)
		if err != nil {
			return nil, err
		}
		return bc.Transform(mid)
	})
}

// Identity returns a transformer that leaves a collection unchanged.
func Identity[T any]() Transformer[T, T] {
	return TransformerFunc[T, T](func(in []T) ([]T, error) {
		// Return a copy so downstream mutation never aliases caller memory.
		out := make([]T, len(in))
		copy(out, in)
		return out, nil
	})
}
