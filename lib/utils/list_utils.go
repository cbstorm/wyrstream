package utils

import (
	"fmt"
)

func Map[K any, V any](list []K, f func(e K, idx int) V) []V {
	n := make([]V, len(list))
	for i, e := range list {
		n[i] = f(e, i)
	}
	return n
}

func Flat[K any](list [][]K) []K {
	out := make([]K, 0)
	for i := range list {
		out = append(out, list[i]...)
	}
	return out
}

func ForEach[K any](list []K, f func(e K, idx int)) {
	for i, e := range list {
		f(e, i)
	}
}

func Filter[K any](list []K, f func(e K, idx int) bool) *[]K {
	n := make([]K, 0)
	for i, e := range list {
		if f(e, i) {
			n = append(n, e)
		}
	}
	return &n
}

func Find[K any](list []K, f func(e K) bool) (K, bool) {
	var n K
	for _, e := range list {
		if f(e) {
			return e, true
		}
	}

	return n, false
}

func FindIndex[K any](list []K, f func(e K) bool) int {
	for i, e := range list {
		if f(e) {
			return i
		}
	}
	return -1
}

func Includes[K int | string | bool](list []K, compare K) bool {
	for _, e := range list {
		if e == compare {
			return true
		}
	}
	return false
}

func IncludesEnum[K any](list []K, compare K) bool {
	for _, e := range list {
		if fmt.Sprintf("%v", e) == fmt.Sprintf("%v", compare) {
			return true
		}
	}
	return false
}

func Join(list []string, separator string) string {
	var n string
	for i, e := range list {
		if i == 0 {
			n = e
		} else {
			n = fmt.Sprintf("%s%s%s", n, separator, e)
		}
	}
	return n
}

func Reduce[K any, V any](list []K, f func(a V, b K) V, inititalValue V) V {
	for _, e := range list {
		inititalValue = f(inititalValue, e)
	}
	return inititalValue
}

func KeyBy[K any](list []K, f func(a K) string) map[string]K {
	result := make(map[string]K)
	for _, e := range list {
		key := f(e)
		result[fmt.Sprintf("%v", key)] = e
	}
	return result
}

func Shift[K any](list *[]K) K {
	item := (*list)[0]
	*list = (*list)[1:]
	return item
}

func UnShift[K any](list *[]K, elem K) int {
	r := make([]K, len(*list)+1)
	for i := 1; i < len(r); i++ {
		r[i] = (*list)[i-1]
	}
	r[0] = elem
	*list = r
	return len(*list)
}

func Pop[K any](list *[]K) K {
	val := (*list)[len(*list)-1]
	(*list) = (*list)[:len(*list)-1]
	return val
}

func Reverse[K any](list []K) {
	start := 0
	end := len(list) - 1
	for start < end {
		if start == end {
			continue
		}
		temp := list[start]
		list[start] = list[end]
		list[end] = temp
		start++
		end--
	}
}

func Every[K any](list *[]K, f func(a K) bool) bool {
	for _, v := range *list {
		if !f(v) {
			return false
		}
	}
	return true
}

func GroupBy[K any](list []K, f func(a K) string) map[string][]K {
	out := make(map[string][]K)
	for _, v := range list {
		key := f(v)
		if out[key] == nil {
			out[key] = []K{}
		}
		out[key] = append(out[key], v)
	}
	return out
}
