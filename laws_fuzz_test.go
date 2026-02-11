package iterlinq

import (
	"errors"
	"testing"
)

type fuzzRecord struct {
	ID   byte
	Flag bool
}

func FuzzToSliceRoundTripBytes(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 4})
	f.Add([]byte{255, 0, 128, 64})

	f.Fuzz(func(t *testing.T, input []byte) {
		got, err := FromSlice(input).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slicesEqual(got, input) {
			t.Fatalf("round trip mismatch: got=%v want=%v", got, input)
		}
	})
}

func FuzzWhereMatchesManualFilter(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 4, 5, 6})
	f.Add([]byte{9, 9, 9, 8, 8, 7})

	f.Fuzz(func(t *testing.T, input []byte) {
		want := make([]byte, 0, len(input))
		for _, v := range input {
			if v%2 == 0 {
				want = append(want, v)
			}
		}

		got, err := FromSlice(input).Where(func(v byte) bool { return v%2 == 0 }).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slicesEqual(got, want) {
			t.Fatalf("where mismatch: got=%v want=%v", got, want)
		}
	})
}

func FuzzTakeSkipReconstructsInput(f *testing.F) {
	f.Add([]byte{}, uint8(0))
	f.Add([]byte{1, 2, 3, 4, 5}, uint8(2))
	f.Add([]byte{42}, uint8(8))

	f.Fuzz(func(t *testing.T, input []byte, n uint8) {
		k := int(n)

		head, err := FromSlice(input).Take(k).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error from take: %v", err)
		}
		tail, err := FromSlice(input).Skip(k).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error from skip: %v", err)
		}

		reconstructed := append(append([]byte{}, head...), tail...)
		if !slicesEqual(reconstructed, input) {
			t.Fatalf("reconstruction mismatch: got=%v want=%v", reconstructed, input)
		}
	})
}

func FuzzDistinctByFirstOccurrenceWins(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 1, 3, 2, 4})
	f.Add([]byte{9, 9, 9, 9})

	f.Fuzz(func(t *testing.T, input []byte) {
		records := make([]fuzzRecord, 0, len(input))
		for _, b := range input {
			records = append(records, fuzzRecord{ID: b % 16, Flag: b%2 == 0})
		}

		got, err := FromSlice(records).DistinctBy(func(r fuzzRecord) any { return r.ID }).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		seen := map[byte]struct{}{}
		want := make([]fuzzRecord, 0, len(records))
		for _, r := range records {
			if _, ok := seen[r.ID]; ok {
				continue
			}
			seen[r.ID] = struct{}{}
			want = append(want, r)
		}

		if !slicesEqual(got, want) {
			t.Fatalf("distinct mismatch: got=%v want=%v", got, want)
		}
	})
}

func FuzzRecoverHandlesSourceError(f *testing.F) {
	f.Add([]byte{}, uint8(0))
	f.Add([]byte{1, 2, 3, 4}, uint8(2))
	f.Add([]byte{7, 8}, uint8(10))

	f.Fuzz(func(t *testing.T, input []byte, failAt uint8) {
		if len(input) == 0 {
			_, err := FromSlice(input).Recover(func(error) Sequence[byte] {
				return FromSlice([]byte{1, 2})
			}).ToSlice()
			if err != nil {
				t.Fatalf("unexpected error on empty input: %v", err)
			}
			return
		}

		failIdx := int(failAt) % len(input)
		fallback := []byte{200, 201}
		boom := errors.New("boom")

		got, err := FromFunc(func(yield func(byte) bool) error {
			for i, v := range input {
				if i == failIdx {
					return boom
				}
				if !yield(v) {
					return nil
				}
			}
			return nil
		}).Recover(func(err error) Sequence[byte] {
			if err != boom {
				t.Fatalf("unexpected error passed to recover: %v", err)
			}
			return FromSlice(fallback)
		}).ToSlice()
		if err != nil {
			t.Fatalf("unexpected error after recover: %v", err)
		}

		want := append(append([]byte{}, input[:failIdx]...), fallback...)
		if !slicesEqual(got, want) {
			t.Fatalf("recover mismatch: got=%v want=%v", got, want)
		}
	})
}
