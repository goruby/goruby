package evaluator

import (
	"bytes"
	"fmt"
	"go/token"
)

type runtime interface {
	addFile(*token.File)
	pushToStack(name string, position int)
	popFromStack()
	String() string
}

func newRuntime() runtime {
	return &_runtime{}
}

type _runtime struct {
	currentFile    *token.File
	currentContext *executionContext
	stack          []*executionContext
}

func (r *_runtime) addFile(file *token.File) {
	r.currentFile = file
	if len(r.stack) == 0 {
		ctx := &executionContext{file: file}
		ctx.addStackFrame("<main>", 1)
		r.currentContext = ctx
		r.stack = append(r.stack, ctx)
	}
}

func (r *_runtime) pushToStack(name string, position int) {
	ctx := &executionContext{file: r.currentFile}
	ctx.addStackFrame(name, position)
	r.currentContext = ctx
	r.stack = append(r.stack, ctx)
}

func (r *_runtime) popFromStack() {
	r.stack = r.stack[:len(r.stack)-1]
	r.currentContext = r.stack[len(r.stack)-1]
}

func (r *_runtime) String() string {
	var out bytes.Buffer
	for _, ctx := range r.stack {
		fmt.Fprintln(&out, ctx.String())
	}
	return out.String()
}

type executionContext struct {
	file  *token.File
	frame frame
}

func (ec *executionContext) addStackFrame(name string, position int) error {
	epos := ec.file.Position(token.Pos(position))
	if epos.Filename == "" || !epos.IsValid() {
		return fmt.Errorf("Error finding context within file %q", ec.file.Name())
	}
	ec.frame = frame{context: name, position: epos}
	return nil
}

func (ec *executionContext) String() string {
	var out bytes.Buffer
	fmt.Fprintln(&out, ec.frame.String())
	return out.String()
}

type frame struct {
	context  string
	position token.Position
}

func (ef frame) String() string {
	return fmt.Sprintf("%s:in `%s'", ef.position, ef.context)
}
