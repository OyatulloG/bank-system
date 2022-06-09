package helper

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//helper
func InputText(title string) string {
	//HELPER FUNCTION:
	//takes user input,
	//cleans leading and trailing white spaces,
	//returns formatted input back

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s: ", title)
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)
	return input
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

//converters
func ConvStrToTime(dateStr string) time.Time {
	//HELPER:
	//CONVERTER:
	//converts date in String format to time.Time format,
	//returns date in time.Time format if conversion is successful,
	//panic err if conversion is unsuccessful

	date, err := time.Parse("2006-01-02", dateStr)
	CheckError(err)
	return date
}

func ConvStrToFloat64(money string) float64 {
	//HELPER:
	//CONVERTER:
	//converts money str to float64 for representing currency,
	//returns Money in float64 format if conversion is successful,
	//returns -1 if conversion is unsuccessful

	moneyFloat, err := strconv.ParseFloat(money, 64)
	CheckError(err)
	return moneyFloat
}

//validators
func IsWhiteSpaceInStr(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//checks for presence of white space inside 'str',
	//returns TRUE if white space is found,
	//returns FALSE if white space is not found

	return strings.Contains(str, " ")
}

func IsLengthInLimit(str string, limit int) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//calculates the length of str,
	//length of str should be more than or equal to 'limit',
	//return TRUE if length >= limit
	//return FALSE if length < limit
	return len(str) >= limit
}

func IsLetterInStr(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//check for presense of alphabetic character in 'str',
	//return TRUE if there is at least 1 an alpha letter
	//return FALSE if there is not an alpha letter at all
	for _, r := range str {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func IsNumInStr(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//check for presence of number in 'str',
	//return TRUE if there is at least 1 number
	//return FALSE if there is not any number at all
	for _, r := range str {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func IsSpecialCharInStr(str string) bool {
	//HELPER FUNCTION:
	//VALIDATION:
	//check for presence of special char in str,
	//return TRUE if there is at least 1 special char,
	//return FALSE if there is not any special char at all
	specialChars := []rune{'!', '@', '#', '$', '%', '^', '&'}
	for _, char := range str {
		for _, c := range specialChars {
			if char == c {
				return true
			}
		}
	}
	return false
}

func IsCleanAlpha(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns TRUE if str contains only alphabetical chars,
	//returns FALSE if str contains a number, a whitespace or a symbol
	return !IsNumInStr(str) && !IsWhiteSpaceInStr(str) &&
		!IsSpecialCharInStr(str)
}

func IsCleanNum(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns TRUE if str contains only numerical chars,
	//returns FALSE if str contains an alpha letter, a whitespace
	//								or a symbol
	return !IsLetterInStr(str) && !IsWhiteSpaceInStr(str) &&
		!IsSpecialCharInStr(str) && len(str) != 0
}

func IsFirstCharUpper(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns TRUE if the first char is uppercase,
	//returns FALSE if the first char is lowercase
	return unicode.IsUpper(rune(str[0]))
}

func IsValidUsername(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns TRUE if str does not contain whitespace and len(str) >= 4
	//returns FALSE if str contains whitespace and len(str) < 4

	return !IsWhiteSpaceInStr(str) && IsLengthInLimit(str, 4)
}

func IsValidPassword(str string) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//(1) password should contain at least 6 chars,
	//(2) 						  at least 1 symbol,
	//(3) 						  at least 1 digit
	//returns TRUE if all requirements are satisfied
	//returns FALSE if all or one of requirements is not satisfied

	return IsLengthInLimit(str, 6) && IsSpecialCharInStr(str) &&
		IsNumInStr(str)
}

//validators of date
func DateAdjuster(date int) string {
	//HELPER:
	//VALIDATOR:
	//adjusts format of dates in range (1,9),
	//changes format from "d" to "dd"
	if date >= 1 && date <= 9 {
		return "0" + strconv.Itoa(date)
	}
	return strconv.Itoa(date)
}

func IsLeapYear(year int) bool {
	//HELPER FUNCTION:
	//VALIDATOR:
	//Year is leap if:
	// (1) year is divided by 400
	// (2) year is divided by 4 and not divided by 100
	//returns TRUE if year is leap
	//return FALSE if year is not leap
	if year%400 == 0 {
		return true
	}

	if year%4 == 0 && year%100 != 0 {
		return true
	}

	return false
}

func IsValidYear(year string) int {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns yearINT if 1900<year<=NOW.YEAR()
	//returns -1 if year is not valid

	if !IsCleanNum(year) {
		return -1
	}

	//convert string to int
	yearInt, err := strconv.Atoi(year)
	CheckError(err)

	if yearInt < 1900 || yearInt > time.Now().Year() {
		return -1
	}

	return yearInt
}

func IsValidMonth(month string) int {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns monthINT if 1 <= month <= 12
	//returns -1 if month is not valid

	if !IsCleanNum(month) {
		return -1
	}

	//convert string to int
	monthInt, err := strconv.Atoi(month)
	CheckError(err)

	if monthInt < 1 || monthInt > 12 {
		return -1
	}

	return monthInt
}

func IsValidDay(day string, month, year int) int {
	//HELPER FUNCTION:
	//VALIDATOR:
	//returns dayINT if 1 <= day <= 28-31 depending on month
	//returns -1 if day is invalid

	february := 28
	if IsLeapYear(year) {
		february = 29
	}

	if !IsCleanNum(day) {
		return -1
	}

	//convert string to int
	dayInt, err := strconv.Atoi(day)
	CheckError(err)

	//January, March, May, July, August, October, December
	if (month == 1 || month == 3 || month == 5 || month == 7 ||
		month == 8 || month == 10 || month == 12) &&
		(dayInt < 1 || dayInt > 31) {
		return -1
	}
	//April, June, September, November
	if (month == 4 || month == 6 || month == 9 || month == 11) &&
		(dayInt < 1 || dayInt > 30) {
		return -1
	}
	//February
	if month == 2 && (dayInt < 1 || dayInt > february) {
		return -1
	}

	return dayInt
}
