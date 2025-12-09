package goext

import "time"

func MinTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func MaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func MonthsBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	years := b.Year() - a.Year()
	months := int(b.Month()) - int(a.Month())
	total := years*12 + months

	return total
}
