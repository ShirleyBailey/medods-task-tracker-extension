package task

import (
	"fmt"
	"time"
)

// RecurrenceType defines how a task repeats.
type RecurrenceType string

const (
	// RecurrenceDaily repeats every N days (Interval field required, >= 1).
	RecurrenceDaily RecurrenceType = "daily"
	// RecurrenceMonthly repeats on a fixed day of each month (DayOfMonth field required, 1-30).
	RecurrenceMonthly RecurrenceType = "monthly"
	// RecurrenceSpecificDates creates tasks only on the listed dates.
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	// RecurrenceEvenOdd creates tasks on even or odd calendar days of the month.
	RecurrenceEvenOdd RecurrenceType = "even_odd"
)

// RecurrenceRule describes the periodicity settings for a task template.
type RecurrenceRule struct {
	// Type is the recurrence strategy.
	Type RecurrenceType `json:"type"`
	// Interval is used by RecurrenceDaily: repeat every Interval days (>= 1).
	Interval int `json:"interval,omitempty"`
	// DayOfMonth is used by RecurrenceMonthly: the calendar day (1-30).
	DayOfMonth int `json:"day_of_month,omitempty"`
	// Dates is used by RecurrenceSpecificDates: explicit list of dates (YYYY-MM-DD).
	Dates []time.Time `json:"dates,omitempty"`
	// EvenDays is used by RecurrenceEvenOdd: true = even days, false = odd days.
	EvenDays bool `json:"even_days,omitempty"`
}

// Validate checks that the rule is internally consistent.
func (r RecurrenceRule) Validate() error {
	switch r.Type {
	case RecurrenceDaily:
		if r.Interval < 1 {
			return fmt.Errorf("daily recurrence requires interval >= 1")
		}
	case RecurrenceMonthly:
		if r.DayOfMonth < 1 || r.DayOfMonth > 30 {
			return fmt.Errorf("monthly recurrence requires day_of_month between 1 and 30")
		}
	case RecurrenceSpecificDates:
		if len(r.Dates) == 0 {
			return fmt.Errorf("specific_dates recurrence requires at least one date")
		}
	case RecurrenceEvenOdd:
		// no extra fields required
	default:
		return fmt.Errorf("unknown recurrence type: %q", r.Type)
	}
	return nil
}

// GenerateDates returns all dates within [start, end] that match the rule.
// Dates are returned in ascending order, truncated to midnight UTC.
func (r RecurrenceRule) GenerateDates(start, end time.Time) ([]time.Time, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	start = truncateToDay(start)
	end = truncateToDay(end)

	if end.Before(start) {
		return nil, fmt.Errorf("end_date must be on or after start_date")
	}

	switch r.Type {
	case RecurrenceDaily:
		return generateDaily(start, end, r.Interval), nil
	case RecurrenceMonthly:
		return generateMonthly(start, end, r.DayOfMonth), nil
	case RecurrenceSpecificDates:
		return generateSpecificDates(start, end, r.Dates), nil
	case RecurrenceEvenOdd:
		return generateEvenOdd(start, end, r.EvenDays), nil
	}
	return nil, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func generateDaily(start, end time.Time, interval int) []time.Time {
	var dates []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, interval) {
		dates = append(dates, d)
	}
	return dates
}

func generateMonthly(start, end time.Time, dayOfMonth int) []time.Time {
	var dates []time.Time
	// iterate month by month
	cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !cur.After(endMonth) {
		// clamp day to last day of month so e.g. day=30 works in February
		lastDay := daysInMonth(cur.Year(), cur.Month())
		day := dayOfMonth
		if day > lastDay {
			day = lastDay
		}
		candidate := time.Date(cur.Year(), cur.Month(), day, 0, 0, 0, 0, time.UTC)
		if !candidate.Before(start) && !candidate.After(end) {
			dates = append(dates, candidate)
		}
		cur = cur.AddDate(0, 1, 0)
	}
	return dates
}

func generateSpecificDates(start, end time.Time, dates []time.Time) []time.Time {
	var result []time.Time
	for _, d := range dates {
		d = truncateToDay(d)
		if !d.Before(start) && !d.After(end) {
			result = append(result, d)
		}
	}
	return result
}

func generateEvenOdd(start, end time.Time, evenDays bool) []time.Time {
	var dates []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		isEven := d.Day()%2 == 0
		if isEven == evenDays {
			dates = append(dates, d)
		}
	}
	return dates
}

func daysInMonth(year int, month time.Month) int {
	// day 0 of next month = last day of current month
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
