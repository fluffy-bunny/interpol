package interpol

import (
	"bytes"
)

// Templater is the parser, which was set up for specific data
// type. It has its own internal state and buffer.
//
// Each goroutine should initiate its own Templater - it supports
// only synchronous operations.
//
// Use New() to create the Templater.
type Templater struct {
	buf *bytes.Buffer

	gfs getterFuncSpawner
}

// New creates new Templater configured for passed data type.
//
// Pass (optionally) empty struct that will be passed for each Exec(),
// i.e. use New(map[string]string{}) if you'll use the Exec() with
// map[string]string as data source container.
func New(d interface{}) (*Templater, error) {
	gfs, err := getterSpawnerSelector(d)
	if err != nil {
		return nil, err
	}

	return &Templater{buf: bytes.NewBuffer([]byte{}), gfs: gfs}, nil
}

// Exec interpolates the values from data into the string.
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

	// Loop over the string
	for i := 0; i < length; i++ {
		// Nothing to do on pos 0
		if i == 0 {
			continue
		}

		// Non-special character; loop over
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

	// Dump leftovers if there's an
	if lastWrite != length {
		t.buf.WriteString(s[lastWrite:length])
	}
	return t.buf.String(), nil
}
