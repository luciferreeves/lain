package format

import (
	"lain/types"
	"time"
)

func FormatEmailDate(date time.Time, dateFormat types.DateFormat, timeFormat types.TimeFormat, prettyDates bool, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	date = date.In(loc)
	now := time.Now().In(loc)

	if prettyDates {
		return formatPrettyDate(date, now, timeFormat)
	}

	return formatFullDate(date, dateFormat, timeFormat)
}

func formatPrettyDate(date, now time.Time, timeFormat types.TimeFormat) string {
	diff := now.Sub(date)

	// Today - show time only
	if date.Year() == now.Year() && date.YearDay() == now.YearDay() {
		return formatTime(date, timeFormat)
	}

	// Yesterday
	yesterday := now.AddDate(0, 0, -1)
	if date.Year() == yesterday.Year() && date.YearDay() == yesterday.YearDay() {
		return "Yesterday"
	}

	// This week - show day name and time
	if diff.Hours() < 168 { // 7 days
		dayName := date.Format("Mon")
		timeStr := formatTime(date, timeFormat)
		return dayName + " " + timeStr
	}

	// This year - show month and day
	if date.Year() == now.Year() {
		return date.Format("Jan 2")
	}

	// Older - show full date
	return date.Format("Jan 2, 2006")
}

func formatFullDate(date time.Time, dateFormat types.DateFormat, timeFormat types.TimeFormat) string {
	dateStr := formatDate(date, dateFormat)
	timeStr := formatTime(date, timeFormat)
	return dateStr + " " + timeStr
}

func formatDate(date time.Time, dateFormat types.DateFormat) string {
	switch dateFormat {
	case types.YearMonthDayDashed:
		return date.Format("2006-01-02")
	case types.YearMonthDaySlashed:
		return date.Format("2006/01/02")
	case types.YearMonthDayDotted:
		return date.Format("2006.01.02")
	case types.DayMonthYearDashed:
		return date.Format("02-01-2006")
	case types.DayMonthYearSlashed:
		return date.Format("02/01/2006")
	case types.DayMonthYearDotted:
		return date.Format("02.01.2006")
	case types.DayMonthYearDottedShort:
		return date.Format("2.1.06")
	default:
		return date.Format("2006-01-02")
	}
}

func formatTime(date time.Time, timeFormat types.TimeFormat) string {
	switch timeFormat {
	case types.ShortHoursAndMinutes24Hours:
		return date.Format("15:4")
	case types.FullHoursAndMinutes24Hours:
		return date.Format("15:04")
	case types.ShortHoursAndMinutes12Hours:
		return date.Format("3:4 PM")
	case types.FullHoursAndMinutes12Hours:
		return date.Format("03:04 PM")
	default:
		return date.Format("15:04")
	}
}
