package handler

import "time"

func ConvertTime(t time.Time) string {
	return ConvertTime2String(ConvertUtc2Jst(t))
}

func ConvertDate(t time.Time) string {
	return ConvertDate2String(t)
}

func ConvertTime2String(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func ConvertDate2String(t time.Time) string {
	return t.Format("2006-01-02")
}

func ConvertUtc2Jst(t time.Time) time.Time {
	return t.Add(-9 * time.Hour)
}

