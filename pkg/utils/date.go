package utils

import "time"

func ToYearISO(date *time.Time) *string {
	if date.IsZero() {
		return nil
	}
	year := date.Format(FormatYearISO)
	return &year
}

func ToDateISO(date *time.Time) *string {
	if date.IsZero() {
		return nil
	}
	res := date.Format(FormatDateISO)
	return &res
}

func ToTimeHHMM(dt *time.Time) *string {
	if dt.IsZero() {
		return nil
	}
	res := dt.Format(FormatTimeHHMM)
	return &res
}

func ToTimeHHMMSS(dt *time.Time) *string {
	if dt.IsZero() {
		return nil
	}
	res := dt.Format(FormatTimeHHMMSS)
	return &res
}

func ToDateTimeISO(dt *time.Time) *string {
	if dt.IsZero() {
		return nil
	}
	res := dt.Format(FormatDateTimeISO)
	return &res
}

func ToDateTimeSQL(dt *time.Time) *string {
	if dt.IsZero() {
		return nil
	}
	res := dt.Format(FormatDateTimeSQL)
	return &res
}

func ToDateCompact(dt *time.Time) *string {
	if dt.IsZero() {
		return nil
	}
	res := dt.Format(FormatDateCompact)
	return &res
}

func ParseDateTimeFromSQlOrISO(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	dt, err := time.Parse(FormatDateTimeISO, *s)
	if err != nil {
		dt, err = time.Parse(FormatDateTimeSQL, *s)
	}
	return &dt, err
}
