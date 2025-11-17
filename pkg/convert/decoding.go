package convert

import (
	"bytes"
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

func FixEncoding(s string) (string, error) {
	// Декодируем из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	return decoder.String(s)
}

// Win1251ToUTF8Map - ручная карта преобразования Windows-1251 → UTF-8
func Win1251ToUTF81(s string) string {
	win1251 := []byte(s)
	utf8 := make([]rune, len(win1251))

	for i, b := range win1251 {
		switch {
		case b >= 0x00 && b <= 0x7F:
			utf8[i] = rune(b) // ASCII совпадает
		case b >= 0xC0 && b <= 0xFF:
			// Кириллица в Windows-1251
			utf8[i] = rune(b) + 0x350
		default:
			utf8[i] = rune(b) // Оставляем как есть
		}
	}
	return string(utf8)
}

// DecodeCP866 декодирование русского текста.
func DecodeCP866(data []byte) (string, error) {
	decoder := charmap.CodePage866.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), decoder)
	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// CP1251ToUTF8 конвертирует строку из Windows-1251 в UTF-8
func CP1251ToUTF8(s string) (string, error) {
	reader := transform.NewReader(
		bytes.NewReader([]byte(s)),
		charmap.Windows1251.NewDecoder(),
	)
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
