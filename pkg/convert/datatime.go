package convert

import "time"

// GetCurrentDateTime получить текущую дату и время в формате "2006-01-02 15:04:05".
func GetCurrentDateTime() string {
	return time.Now().Format(time.DateTime)
}
