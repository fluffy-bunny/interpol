package interpol

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInterpol(t *testing.T) {
	Convey("With templater", t, func() {
		tm := Templater{}
		tm.buf = bytes.NewBuffer([]byte{})

		Convey("Interpolation, reverse func", func() {
			// This func gets the variable name, reverses it and returns back.
			tm.gfs = func(_ interface{}) (getterFunc, error) {
				return func(data string) (result []byte, err error) {
					for i := len(data) - 1; i >= 0; i-- {
						result = append(result, data[i])
					}
					return result, nil
				}, nil
			}

			Convey("Simple interpolation", func() {
				str := `foo {{bar}}baz}`

				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, `foo rabbaz}`)
			})

			Convey("Multiple interpolation, bad parens", func() {
				str := `{{foo}} {{bar} }}{{bazz}}`

				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, `oof  }rabzzab`)
			})

		})
		Convey("Edge cases", func() {
			tm.gfs = func(_ interface{}) (getterFunc, error) {
				return func(data string) ([]byte, error) {
					return []byte(data), nil
				}, nil
			}

			Convey("Empty string", func() {
				got, err := tm.Exec(``, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, ``)
			})
			Convey("String without variables", func() {
				str := `some string  `

				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, str)
			})
			Convey("String with odd parens", func() {
				str := `some string { { }}}}}} `

				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, str)
			})
			Convey("Variable @ start", func() {
				str := `{{qwe}} foo bar`
				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, `qwe foo bar`)
			})
			Convey("Variable @ end", func() {
				str := `foo bar {{qwe}}`
				got, err := tm.Exec(str, map[string]string{})
				So(err, ShouldBeNil)
				So(got, ShouldEqual, `foo bar qwe`)
			})
		}) // Edge cases final

		Convey("Interpolation, map[string]string", func() {
			tm, err := New(map[string]string{})
			So(err, ShouldBeNil)

			data := map[string]string{
				"var1": "look!",
				"var2": "string",
			}

			str := `{{var1}} some {{var2}}'s coming!`

			got, err := tm.Exec(str, data)
			So(err, ShouldBeNil)
			So(got, ShouldEqual, `look! some string's coming!`)

		})
	})
}
