package ast

import (
	"reflect"
	"sync"
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
	done := make(chan struct{})
	defer close(done)

	w1 := nodeWalk(a)
	w2 := nodeWalk(b)

	equalCh := walkNodes(done, w1, w2)
	for equal := range equalCh {
		if !equal {
			done <- struct{}{}
			return false
		}
	}

	return true
}

func walkNodes(done <-chan struct{}, nc ...<-chan Node) <-chan bool {
	out := make(chan bool)
	go func() {
		defer close(out)
		if len(nc) < 2 {
			out <- false
			return
		}
		for n := range mergeNodeWalks(done, nc...) {
			out <- compare(n...)
		}
	}()

	return out
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

func mergeNodeWalks(done <-chan struct{}, nc ...<-chan Node) <-chan []Node {
	out := make(chan []Node)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
				nodes, ok := consolidateNodeWalks(nc...)
				if !ok {
					return
				}
				out <- nodes
			}
		}
	}()
	return out
}

func consolidateNodeWalks(nc ...<-chan Node) ([]Node, bool) {
	var wg sync.WaitGroup
	out := make(chan Node)
	closed := make(chan struct{}, len(nc))

	output := func(c <-chan Node) {
		defer wg.Done()
		n, ok := <-c
		if !ok {
			closed <- struct{}{}
			out <- nil
		}
		out <- n
	}

	wg.Add(len(nc))
	for _, c := range nc {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
		close(closed)
	}()

	var nodes []Node
	for n := range out {
		nodes = append(nodes, n)
	}
	var closedChans []struct{}
	for c := range closed {
		closedChans = append(closedChans, c)
	}

	return nodes, len(closedChans) != len(nc)
}

func nodeWalk(root Node) <-chan Node {
	out := make(chan Node)
	var visitor Visitor
	visitor = VisitorFunc(func(n Node) Visitor {
		out <- n
		return visitor
	})
	go func() {
		defer close(out)
		Walk(visitor, root)
	}()
	return out
}
