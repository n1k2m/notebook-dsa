package event

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	Year    int
	Month   int
	Day     int
	Hour    int
	Minute  int
	Weekday int
	Desc    string
	Place   string
}

var weekdayNames = [7]string{
	"Воскресенье", "Понедельник", "Вторник",
	"Среда", "Четверг", "Пятница", "Суббота",
}

func WeekdayName(w int) string {
	return weekdayNames[w%7]
}

// CalcWeekday вычисляет день недели по алгоритму Зеллера
// 0=Вс, 1=Пн, ..., 6=Сб
func CalcWeekday(y, m, d int) int {
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return int(t.Weekday())
}

func DaysInMonth(y, m int) int {
	// Берём первый день следующего месяца и откатываемся на 1 день
	t := time.Date(y, time.Month(m+1), 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}

func IsValidDate(y, m, d int) bool {
	if y < 1900 || y > 2100 {
		return false
	}
	if m < 1 || m > 12 {
		return false
	}
	if d < 1 || d > DaysInMonth(y, m) {
		return false
	}
	return true
}

func IsValidTime(h, min int) bool {
	return h >= 0 && h <= 23 && min >= 0 && min <= 59
}

func Cmp(a, b *Event) int {
	if a.Year != b.Year {
		return a.Year - b.Year
	}
	if a.Month != b.Month {
		return a.Month - b.Month
	}
	if a.Day != b.Day {
		return a.Day - b.Day
	}
	if a.Hour != b.Hour {
		return a.Hour - b.Hour
	}
	return a.Minute - b.Minute
}

// HashKey возвращает числовой ключ для хеш-таблицы
func HashKey(e *Event) uint64 {
	return uint64(e.Year)*10_000_000_000 +
		uint64(e.Month)*100_000_000 +
		uint64(e.Day)*1_000_000 +
		uint64(e.Hour)*10_000 +
		uint64(e.Minute)
}

// Equal проверяет совпадение по дате и времени
func Equal(a, b *Event) bool {
	return a.Year == b.Year && a.Month == b.Month &&
		a.Day == b.Day && a.Hour == b.Hour && a.Minute == b.Minute
}

// NewEvent создаёт событие с вычисленным днём недели
func NewEvent(y, m, d, h, mi int, desc, place string) Event {
	return Event{
		Year:    y,
		Month:   m,
		Day:     d,
		Hour:    h,
		Minute:  mi,
		Weekday: CalcWeekday(y, m, d),
		Desc:    desc,
		Place:   place,
	}
}

func Print(e *Event) {
	fmt.Printf("%04d-%02d-%02d %02d:%02d  [%s]\n",
		e.Year, e.Month, e.Day, e.Hour, e.Minute,
		WeekdayName(e.Weekday))
	fmt.Printf("Место: %s\n", e.Place)
	fmt.Printf("Описание: %s\n", e.Desc)
	fmt.Println("─────────────────────────────────────")
}

// PlaceContains проверяет, содержит ли место подстроку
func PlaceContains(e *Event, substr string) bool {
	return strings.Contains(
		strings.ToLower(e.Place),
		strings.ToLower(substr),
	)
}
