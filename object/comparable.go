package object

var comparableModule = newModule("Comparable", comparableMethodSet, nil)

func init() {
	classes.Set("Comparable", comparableModule)
}

var comparableMethodSet = map[string]RubyMethod{}
