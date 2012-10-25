package keyvalues

import "bytes"

func (kv *KeyValues) String() string {
	var buf bytes.Buffer

	kv.WriteTo(&buf)

	return string(buf.Bytes())
}
