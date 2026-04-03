package output

import (
	"fmt"
	"io"
)

type TextFormatter struct{}

func (f *TextFormatter) Format(w io.Writer, data any) error {
	_, err := fmt.Fprintln(w, data)
	return err
}
