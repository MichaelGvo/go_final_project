package nextdate

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("20060102", s)
	return t, err
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("null in repeat")
	}

	startTime, err := parseTime(date)
	if err != nil {
		return "", fmt.Errorf("incorrect date: %w", err)
	}

	switch repeat[0] {
	case 'y':
		return plusYear(now, startTime)
	case 'd':
		return plusDay(now, startTime, repeat)
	case 'w':
		return plusWeek(now, repeat)
	case 'm':
		return plusMonth(now, startTime, repeat)
	default:
		return "", fmt.Errorf("invalid repeat format: %s", repeat)
	}
}

func plusYear(now time.Time, startTime time.Time) (string, error) {
	oneYearLater := startTime.AddDate(1, 0, 0)
	for oneYearLater.Before(now) {
		oneYearLater = oneYearLater.AddDate(1, 0, 0)
	}
	oneYearLaterStr := oneYearLater.Format("20060102")
	return oneYearLaterStr, nil

}

func plusDay(now time.Time, startTime time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Printf("failed to parse number: %v", err)
	}
	if days > 400 {
		fmt.Printf("Current number is too much: %v", days)
	}
	someDaysLater := startTime.AddDate(0, 0, days)
	for someDaysLater.Before(now) {
		someDaysLater = someDaysLater.AddDate(0, 0, days)
	}
	someDaysLaterStr := someDaysLater.Format("2006-01-02")
	return someDaysLaterStr, nil
}
func plusWeek(now time.Time, repeat string) (string, error) {
	partsOfRepeat := strings.Split(repeat, " ")
	allDays := partsOfRepeat[1]
	eachDaySeparated := strings.Split(allDays, ",")
	uniqueElements := make(map[int]bool)

	var correctNumberDays []string
	for _, part := range eachDaySeparated {
		num, err := strconv.Atoi(part)
		if err == nil && num <= 7 && !uniqueElements[num] {
			uniqueElements[num] = true
			correctNumberDays = append(correctNumberDays, part)
			if len(correctNumberDays) == 7 {
				break
			}
		}
	}

	var weekdays []time.Weekday

	for _, str := range correctNumberDays {
		num, err := strconv.Atoi(str)
		if err == nil {
			weekday := time.Weekday(num)
			weekdays = append(weekdays, weekday)
		}
	}
	sort.Slice(weekdays, func(i, j int) bool {
		return weekdays[i] < weekdays[j]
	})

	someDayOfWeek := now.AddDate(0, 0, int(now.Weekday()-weekdays[0]))
	someDayOfWeekStr := someDayOfWeek.Format("2006-01-02")
	return someDayOfWeekStr, nil
}

func plusMonth(now time.Time, startTime time.Time, repeat string) (string, error) {
	partsOfRepeat := strings.Split(repeat, " ")
	allDays := partsOfRepeat[1]
	allMonths := partsOfRepeat[2]

	monthsMap := make(map[int]bool)
	if len(allMonths) != 0 {
		for _, m := range strings.Split(allMonths, ",") {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("invalid format of month: %s", allMonths)
			}
			monthsMap[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthsMap[i] = true
		}
	}

	partsOfAllDays := strings.Split(allDays, ",")
	daysMap := make(map[int]bool)
	if len(allDays) != 0 {
		for _, d := range partsOfAllDays {
			day, err := strconv.Atoi(d)
			if err != nil || day < -2 || day > 31 {
				return "", fmt.Errorf("invalid format of month: %s", allDays)
			}
			daysMap[day] = true

		}
	}

	for neededDate := startTime; ; neededDate = neededDate.AddDate(0, 0, 1) {
		day := neededDate.Day()
		month := int(neededDate.Month())

		if daysMap[day] || daysMap[day-getDaysInMonth(neededDate.Year()-1, neededDate.Month())] {

			if monthsMap[month] {
				if neededDate.After(now) {
					return neededDate.Format("20060102"), nil
				}
			}
		}
	}
}

func getDaysInMonth(year int, month time.Month) int {

	date := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	return date.Day()
}
