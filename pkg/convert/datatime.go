package convert

import (
	"time"
)

// GetCurrentDateTime получить текущую дату и время в формате "2006-01-02 15:04:05".
func GetCurrentDateTime() string {
	return time.Now().Format(time.DateTime)
}

// FormatDateTime - функция форматирования даты в формате ДД.ММ.ГГГГ ЧЧ:ММ
func FormatDateTime(dateTime string) string {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return dateTime
	}

	return t.Format("02.01.2006 15:04:05")
}

func FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "не указано"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("02.01.2006 15:04:05")
}
