package workerpool

import (
	"testing"
)

func TestDoPool(t *testing.T) {
	dir := "./data"
	poolSize := 2
	DoPools(poolSize, dir)
}
