package keyvalues

type KeyValues struct {
	name string
	data map[string]interface{} // Every value is either a string or a *KeyValues.
}

// Create a new KeyValues with a given name, defined by KeyValues.cpp:203.
func New(name string) *KeyValues {
	return &KeyValues{name: name, data: make(map[string]interface{})}
}

// Create a new KeyValues with a given name and first key, defined by KeyValues.cpp:214 :226 :238.
func NewSingle(name string, firstKey string, firstValue interface{}) *KeyValues {
	return &KeyValues{
		name: name,
		data: map[string]interface{}{
			firstKey: valueify(firstValue),
		},
	}
}

// Create a new KeyValues with a given name and first and second keys, defined by KeyValues.cpp:250 :263.
func NewDouble(name string, firstKey string, firstValue interface{}, secondKey string, secondValue interface{}) *KeyValues {
	return &KeyValues{
		name: name,
		data: map[string]interface{}{
			firstKey:  valueify(firstValue),
			secondKey: valueify(secondValue),
		},
	}
}
