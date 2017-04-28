package object

import "testing"

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		object RubyObject
		value  bool
	}{
		{
			NIL,
			false,
		},
		{
			TRUE,
			true,
		},
		{
			FALSE,
			false,
		},
		{
			&String{Value: "foo"},
			true,
		},
		{
			&Integer{Value: 0},
			true,
		},
		{
			&String{Value: ""},
			true,
		},
	}

	for _, tt := range tests {
		result := IsTruthy(tt.object)

		if result != tt.value {
			t.Logf("Expected result to be %t, got %t", tt.value, result)
			t.Fail()
		}
	}
}
