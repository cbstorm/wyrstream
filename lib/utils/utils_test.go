package utils_test

import (
	"testing"

	"github.com/cbstorm/wyrstream/lib/utils"
)

func TestStringRand(t *testing.T) {
	str_rand := utils.StringRand(30)
	l := len(str_rand)
	t.Log(l)
}
