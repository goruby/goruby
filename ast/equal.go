package ast

import (
	"reflect"
)

// Equal returns true if the nodes are identical in content, false otherwise
func Equal(a, b Node) bool {
	if a == b {
		return true
	}
	if reflect.DeepEqual(a, b) {
		return true
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	if _, ok := a.(literal); ok {
		return a.String() == b.String()
	}
	return compareTrees(a, b)
}

func compareTrees(a, b Node) bool {
	res := make(chan bool)

	go func() {
		defer close(res)
		ch1 := WalkEmit(a)
		ch2 := WalkEmit(b)
		for {
			x, ok1 := <-ch1
			y, ok2 := <-ch2
			if !ok1 && !ok2 {
				break
			}
			if !compare(x, y) {
				res <- false
				return
			}
		}
		res <- true
	}()

	return <-res
}

func compare(nodes ...Node) bool {
	if len(nodes) < 2 {
		return false
	}
	var prev, next Node
	for _, n := range nodes {
		if prev == nil {
			prev = n
			continue
		}
		if next != nil {
			prev = next
		}
		next = n
		if reflect.TypeOf(prev) != reflect.TypeOf(next) {
			return false
		}
		if prev.String() != next.String() {
			return false
		}
	}
	return true
}
