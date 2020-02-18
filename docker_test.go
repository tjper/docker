package docker

import (
	"errors"
	"strconv"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		image string
		err   error
	}{
		{image: "redis", err: nil},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, stop, err := Run(test.image)
			if !errors.Is(err, test.err) {
				t.Fatalf("unexpected err; expected: %s; actual: %s", test.err, err)
			}
			err = stop()
			if !errors.Is(err, test.err) {
				t.Fatalf("unexpected err; expected: %s; actual: %s", test.err, err)
			}
		})
	}
}
