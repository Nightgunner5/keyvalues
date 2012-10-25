package keyvalues

import "fmt"

// Convert all types to strings except *KeyValues.
func valueify(in interface{}) interface{} {
	if _, ok := in.(*KeyValues); ok {
		return in
	}

	if b, ok := in.(bool); ok {
		if b {
			return "1"
		}
		return "0"
	}

	return fmt.Sprint(in)
}
