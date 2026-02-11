package rigid

import (
	"errors"
	"testing"
)

func FuzzRunIdentityRoundTrip(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{255, 0, 1})

	f.Fuzz(func(t *testing.T, input []byte) {
		got, err := Run(input, Identity[byte]())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slicesEqual(got, input) {
			t.Fatalf("mismatch: got=%v want=%v", got, input)
		}
	})
}

func FuzzRecoverHandlesErrors(f *testing.F) {
	f.Add([]byte{}, uint8(0))
	f.Add([]byte{1, 2, 3, 4}, uint8(2))

	f.Fuzz(func(t *testing.T, input []byte, failAt uint8) {
		if len(input) == 0 {
			return
		}
		boom := errors.New("boom")
		idx := int(failAt) % len(input)
		p := Then(From(input), Where(func(v byte) (bool, error) {
			if int(v)%len(input) == idx {
				return false, boom
			}
			return true, nil
		})).Recover(func(error) Pipeline[byte, byte] {
			return From([]byte{9, 10})
		})

		got, err := p.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) == 0 {
			t.Fatalf("expected recovered output")
		}
	})
}
