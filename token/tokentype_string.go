// Code generated by "stringer -type=TokenType"; DO NOT EDIT

package token

import "fmt"

const _TokenType_name = "ILLEGALEOFNEWLINEIDENTINTSTRINGSYMBOLASSIGNPLUSMINUSBANGASTERISKSLASHLTGTEQNOT_EQCOMMASEMICOLONDOTCOLONLPARENRPARENLBRACERBRACEDEFENDIFTHENELSETRUEFALSERETURN"

var _TokenType_index = [...]uint8{0, 7, 10, 17, 22, 25, 31, 37, 43, 47, 52, 56, 64, 69, 71, 73, 75, 81, 86, 95, 98, 103, 109, 115, 121, 127, 130, 133, 135, 139, 143, 147, 152, 158}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
