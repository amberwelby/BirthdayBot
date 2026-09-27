package src

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetToday() (string, string) {
	date := strings.Split(time.Now().Format(time.DateOnly), "-")
	month := date[1]
	day := date[2]

	return month, day
}

func ValidateYear(year string) (string, error) {
	if year != "NA" {
		int_year, err := strconv.Atoi(year)
		if err != nil {
			return "", fmt.Errorf("Year string conversion: %s", err)
		}
		if int_year < 1900 || int_year > time.Now().Year() {
			return "", fmt.Errorf("Year out of range 1900-%v", time.Now().Year())
		}
	}

	return year, nil
}

func ValidateMonth(month string) (string, int, error) {
	int_month := 0

	if month != "" {
		int_month, err := strconv.Atoi(month)
		if err != nil {
			return "", 0, fmt.Errorf("Month string conversion: %s", err)
		}
		if int_month < 1 || int_month > 12 {
			return "", 0, fmt.Errorf("Month out of range 01-12")
		}
	}

	return month, int_month, nil
}

func ValidateDate(input string) (string, string, string, error) {
	date := strings.Split(input, "/")
	year := "NA"
	month := ""
	day := ""
	if len(date) == 3 {
		year = date[0]
		month = date[1]
		day = date[2]
	} else if len(date) == 2 {
		month = date[0]
		day = date[1]
	}

	year, err := ValidateYear(year)
	if err != nil {
		return "", "", "", err		
	}

	month, int_month, err := ValidateMonth(month)
	if err != nil {
		return "", "", "", err		
	}

	if day != "" {
		int_day, err := strconv.Atoi(day)
		if err != nil {
			return "", "", "", fmt.Errorf("Day string conversion: %s", err)
		}

		month_end := 31 // Jan, Mar, May, Jul, Aug, Oct, Dec
		if int_month == 2 {
			month_end = 29
		} // Feb
		if int_month == 4 || int_month == 6 || int_month == 9 || int_month == 11 {
			month_end = 30
		} // Apr, Jun, Sept, Nov

		if int_day < 1 || int_day > month_end {
			return "", "", "", fmt.Errorf("Date out of range 1-%v", month_end)
		}
	}

	return month, day, year, nil
}
