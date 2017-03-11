package object

type Object struct {
}

func (o *Object) PatternMatch(other RubyObject) RubyObject {
	return NIL
}
