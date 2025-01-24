package go3s

import (
	"testing"
)

type I[D any] interface {
	Get() D
}

type A struct {
	B string
}

func TestX(t *testing.T) {
	println(1)
}
