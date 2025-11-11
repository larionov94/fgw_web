package convert

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func Win1251ToUTF8(str string) (string, error) {
	tr := transform.NewReader(strings.NewReader(str), charmap.Windows1251.NewDecoder())
	buf, err := io.ReadAll(tr)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
