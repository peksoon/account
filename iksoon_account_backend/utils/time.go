package utils

import (
	"fmt"
	"time"
)

// ParseDateTimeKST KST 시간대로 날짜 파싱
func ParseDateTimeKST(dateString string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, fmt.Errorf("KST 시간대 로드 오류: %v", err)
	}

	// 다양한 날짜 형식 지원
	formats := []string{
		"2006-01-02T15:04:05",       // ISO 형식
		"2006-01-02 15:04:05",       // 일반적인 형식
		"2006-01-02",                // 날짜만
		"2006-01-02T15:04:05Z",      // UTC 형식
		"2006-01-02T15:04:05+09:00", // 타임존 포함
	}

	for _, format := range formats {
		if parsedDate, err := time.ParseInLocation(format, dateString, location); err == nil {
			return parsedDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("지원되지 않는 날짜 형식: %s", dateString)
}

// FormatDateTimeKST KST 시간대로 날짜 포맷
func FormatDateTimeKST(t time.Time) string {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// 오류 시 UTC 사용
		return t.UTC().Format("2006-01-02 15:04:05")
	}

	return t.In(location).Format("2006-01-02 15:04:05")
}

// GetCurrentKST 현재 KST 시간 반환
func GetCurrentKST() time.Time {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// 오류 시 UTC 사용
		return time.Now().UTC()
	}

	return time.Now().In(location)
}

// FormatDateKST KST 시간대로 날짜만 포맷 (시간 제외)
func FormatDateKST(t time.Time) string {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// 오류 시 UTC 사용
		return t.UTC().Format("2006-01-02")
	}

	return t.In(location).Format("2006-01-02")
}

// StartOfDayKST 해당 날짜의 시작 시간 (00:00:00) KST 반환
func StartOfDayKST(t time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		location = time.UTC
	}

	year, month, day := t.In(location).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}

// EndOfDayKST 해당 날짜의 끝 시간 (23:59:59) KST 반환
func EndOfDayKST(t time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		location = time.UTC
	}

	year, month, day := t.In(location).Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, location)
}

// StartOfMonthKST 해당 월의 시작일 KST 반환
func StartOfMonthKST(t time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		location = time.UTC
	}

	year, month, _ := t.In(location).Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, location)
}

// EndOfMonthKST 해당 월의 마지막일 KST 반환
func EndOfMonthKST(t time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		location = time.UTC
	}

	year, month, _ := t.In(location).Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, location)
	return firstOfNextMonth.AddDate(0, 0, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
}
