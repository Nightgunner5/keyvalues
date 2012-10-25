package keyvalues

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

// #include, #base, and [$conditionals] are unsupported.
func FromReader(r io.Reader) (kv *KeyValues, err error) {
	in := bufio.NewReader(r)

	var tok string
	var special bool
	tok, special, err = token(in)
	if err != nil {
		return
	}

	if special {
		err = fmt.Errorf("Expected normal token, but got special token %q", tok)
		return
	}
	name := tok

	tok, special, err = token(in)
	if err != nil {
		return
	}

	if !special || tok != "{" {
		err = fmt.Errorf("Expected special token \"{\", but got \"\"", tok)
		return
	}

	return key(name, in)
}

func key(name string, in *bufio.Reader) (kv *KeyValues, err error) {
	kv = New(name)

	for {
		var tok string
		var special bool
		tok, special, err = token(in)
		if err != nil {
			return
		}
		if special {
			switch tok {
			case "{":
				err = fmt.Errorf("Missing key name")
				return
			case "}":
				return
			default:
				err = fmt.Errorf("Conditionals are unsupported, but found %q", tok)
				return
			}
		}

		name = tok

		tok, special, err = token(in)
		if err != nil {
			return
		}

		if special {
			switch tok {
			case "{":
				var val *KeyValues
				val, err = key(name, in)
				if err != nil {
					return
				}

				kv.data[name] = val

			case "}":
				err = fmt.Errorf("Missing value for key %q", name)
				return

			default:
				err = fmt.Errorf("Invalid position for conditional %q", tok)
				return
			}
		} else {
			kv.data[name] = tok
		}
	}
	panic("unreachable")
}

// From KeyValues.cpp:373
func token(in *bufio.Reader) (tok string, special bool, err error) {
	var r rune = ' '

	for unicode.IsSpace(r) || r == '/' {
		if r == '/' {
			var next rune
			if next, _, err = in.ReadRune(); next == '/' {
				in.ReadString('\n')
			} else {
				in.UnreadRune()
				break
			}
		}

		r, _, err = in.ReadRune()
		if err != nil {
			return
		}
	}

	// read quoted strings specially
	if r == '"' {
		tok, err = in.ReadString('"')
		tok = tok[:len(tok)-1]
		return
	}

	if r == '{' || r == '}' {
		// it's a control char, just add this one char and stop reading
		tok = string([]rune{r})
		special = true
		return
	}

	if r == '[' {
		special = true
	}

	// read in the token until we hit a whitespace or a control character
	runes := []rune{r}
	for r, _, err = in.ReadRune(); err == nil; r, _, err = in.ReadRune() {
		// break if any control character appears in non quoted tokens
		if r == '"' || r == '{' || r == '}' {
			in.UnreadRune()
			break
		}

		if unicode.IsSpace(r) {
			break
		}

		runes = append(runes, r)

		if r == ']' && special {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}
	tok = string(runes)
	return
}
