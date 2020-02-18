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
		port  string
	}{
		{image: "redis", err: nil, port: "6379/tcp"},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ports, stop, err := Run(test.image)
			defer stop()

			if !errors.Is(err, test.err) {
				t.Fatalf("unexpected err; expected: %s; actual: %s", test.err, err)
				return
			}
			if _, ok := ports[test.port]; !ok {
				t.Fatalf("unexpected err; expected: %s; actual: %s", test.err, err)
				return
			}
		})
	}
}
