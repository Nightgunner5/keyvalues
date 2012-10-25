package keyvalues

import (
	"fmt"
	"io"
	"strings"
)

func (kv *KeyValues) WriteTo(w io.Writer) (n int64, err error) {
	var c int
	c, err = fmt.Fprintf(w, "%q ", kv.name)
	n += int64(c)
	if err != nil {
		return
	}

	c, err = kv.writeIndented(w, 0)
	n += int64(c)
	return
}

func (kv *KeyValues) writeIndented(w io.Writer, tabs int) (n int, err error) {
	var c int

	indent := strings.Repeat("\t", tabs)

	c, err = fmt.Fprint(w, "{\n")
	n += c
	if err != nil {
		return
	}

	for key, value := range kv.data {
		c, err = fmt.Fprintf(w, "%s\t%q ", indent, key)
		n += c
		if err != nil {
			return
		}

		if s, ok := value.(string); ok {
			c, err = fmt.Fprintf(w, "%q\n", s)
		} else if k, ok := value.(*KeyValues); ok {
			c, err = k.writeIndented(w, tabs+1)
		} else {
			panic(fmt.Sprintf("Invariant failed! %T in data.", value))
		}
		n += c
		if err != nil {
			return
		}
	}

	c, err = fmt.Fprintf(w, "%s}\n", indent)
	n += c
	return
}
