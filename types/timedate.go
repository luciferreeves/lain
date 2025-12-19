package types

type TimeFormat string

const (
	ShortHoursAndMinutes24Hours TimeFormat = "7:30"
	FullHoursAndMinutes24Hours  TimeFormat = "07:30"
	ShortHoursAndMinutes12Hours TimeFormat = "7:30 PM"
	FullHoursAndMinutes12Hours  TimeFormat = "07:30 PM"
)

type DateFormat string

const (
	YearMonthDayDashed      DateFormat = "2025-12-20"
	YearMonthDaySlashed     DateFormat = "2025/12/20"
	YearMonthDayDotted      DateFormat = "2025.12.20"
	DayMonthYearDashed      DateFormat = "20-12-2025"
	DayMonthYearSlashed     DateFormat = "20/12/2025"
	DayMonthYearDotted      DateFormat = "20.12.2025"
	DayMonthYearDottedShort DateFormat = "7.7.25"
)
