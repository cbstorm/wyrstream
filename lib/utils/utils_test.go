package utils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cbstorm/wyrstream/lib/utils"
)

func TestStringRand(t *testing.T) {
	str_rand := utils.StringRand(30)
	l := len(str_rand)
	t.Log(l)
}

func TestPop(t *testing.T) {
	l := make([]int, 3)
	l[0] = 1
	l[1] = 2
	l[2] = 3
	v := utils.Pop(&l)
	if v != 3 || len(l) != 2 {
		t.FailNow()
	}
}

func TestReverse(t *testing.T) {
	s := "abc-xyz-qwe-srt"
	sp := strings.Split(s, "")
	utils.Reverse(sp)
	sr := strings.Join(sp, "")
	fmt.Printf("%v", sr)
	if sr != "trs-ewq-zyx-cba" {
		t.FailNow()
	}
}

func TestInShift(t *testing.T) {
	l := make([]int, 3)
	ll := utils.UnShift(&l, 1)
	if ll != 4 || l[0] != 1 {
		t.FailNow()
	}
}
