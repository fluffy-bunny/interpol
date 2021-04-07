# interpol
Named interpolation library for Go.

this variant of the libray looks for C# style keys.

i.e. ${key} vs {{key}}
```
func TestInterpolate(t *testing.T) {
	var replaceMap = map[string]string{
		"A": "John",
		"B": "Jane",
	}
	var templateString = `${A}`

	tm, err := interpol.New(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, tm)

	got, err := tm.Exec(templateString, replaceMap)
	assert.NoError(t, err)
	assert.Equal(t, "John", got)

	templateString = `${A} ${B}`
	got, err = tm.Exec(templateString, replaceMap)
	assert.NoError(t, err)
	assert.Equal(t, "John Jane", got)

	templateString = `${A} ${B} {}`
	got, err = tm.Exec(templateString, replaceMap)
	assert.NoError(t, err)
	assert.Equal(t, "John Jane {}", got)
}
```
