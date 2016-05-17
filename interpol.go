package interpol

import (
	"bytes"
)

type Templater struct {
	buf bytes.Buffer

	gfs getterFuncSpawner
}

// func used to retrieve the value from an object (struct or map)
type getterFunc func(string) ([]byte, error)

type getterFuncSpawner func(interface{}) (getterFunc, error)

func (t *Templater) Exec(s string, data interface{}) (string, error) {
	t.buf.Reset()

	// Create a getterFunc for the data using the getterFuncSpawner
	gf, err := t.gfs(data)
	if err != nil {
		return ``, err
	}

	// blockStart is the current block's starting position
	// (either simple text or variable's name)
	var blockStart int
	// position at which last write was executed
	var lastWrite int

	// true if reading the variable's name
	var isVarBlock bool

	length := len(s)

	for i := 0; i < length; i++ {
		if i == 0 {
			continue
		}

		if s[i] != '{' && s[i] != '}' {
			continue
		}

		if s[i] == '{' && s[i-1] == '{' {
			// Variable name started, write everything before
			t.buf.WriteString(s[blockStart : i-1])
			isVarBlock = true
			blockStart = i + 1
			lastWrite = i + 1
			continue
		}

		if isVarBlock && s[i] == '}' && s[i-1] == '}' {
			// Variable name finished, look up and write away
			value, err := gf(string(s[blockStart : i-1]))
			if err != nil {
				return ``, err
			}

			lastWrite = i + 1
			t.buf.Write(value)
			isVarBlock = false
			blockStart = i + 1
		}
	}
	if lastWrite != length {
		t.buf.WriteString(s[lastWrite:length])
	}
	return t.buf.String(), nil
}
