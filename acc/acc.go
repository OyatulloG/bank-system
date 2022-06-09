package acc

import (
	"dbConnection/helper"
	"fmt"
	"time"
)

func FirstOrLastName(title string) (bool, string) {
	//takes input for FirstName or LastName,
	//performs validation for correct Name standards:
	//	(1) Name should contain only alphabetical letters
	//	(2) Name should start with a capital letter
	//returns (TRUE, FirstName||LastName) if input is valid
	//returns (FALSE, "" empty string) if input is not valid

	name := helper.InputText(title)

	if helper.IsCleanAlpha(name) && helper.IsFirstCharUpper(name) {
		return true, name
	}
	return false, ""
}

func DateOfBirth() (bool, time.Time) {
	//takes user input for year, month, day and validates them,
	//returns (TRUE, YYYY-MM-DD in time.Time format) if date is valid
	//returns (FALSE, time.Time{} zero value = 0001-01-01) if date is not valid

	fmt.Println("\tDate of Birth")

	//year
	year := helper.InputText("Year")
	yearInt := helper.IsValidYear(year)
	if yearInt == -1 {
		fmt.Println("Input for year is not valid")
		return false, time.Time{}
	}

	//month
	month := helper.InputText("Month")
	monthInt := helper.IsValidMonth(month)
	if monthInt == -1 {
		fmt.Println("Input for month is not valid")
		return false, time.Time{}
	}

	//day
	day := helper.InputText("Day")
	dayInt := helper.IsValidDay(day, monthInt, yearInt)
	if dayInt == -1 {
		fmt.Println("Input for day is not valid")
		return false, time.Time{}
	}

	month = helper.DateAdjuster(monthInt)
	day = helper.DateAdjuster(dayInt)
	date := year + "-" + month + "-" + day
	return true, helper.ConvStrToTime(date)
}

func Gender() (bool, string) {
	//displays choice for gender: male and female,
	//takes user input,
	//returns (TRUE, 'male' or 'female') if input matches,
	//returns (FALSE, "" empty string) if input does not match
	fmt.Println("\t\t Gender (choose option 1 or 2)")
	fmt.Println("\t (1) Male")
	fmt.Println("\t (2) Female")

	gender := helper.InputText("\t Gender")
	switch gender {
	case "1":
		return true, "male"
	case "2":
		return true, "female"
	default:
		return false, ""
	}
}

func Username() (bool, string) {
	//takes input for username,
	//returns (TRUE, username) if input is valid
	//returns (FALSE, "" empty string) if input is not valid

	username := helper.InputText("Username")
	if !helper.IsValidUsername(username) {
		return false, ""
	}
	return true, username
}

func Password() (bool, string) {
	//takes input for password,
	//returns (TRUE, password) if input is valid
	//returns (FALSE, "" empty string) if input is not valid

	password := helper.InputText("Password")
	if !helper.IsValidPassword(password) {
		return false, ""
	}
	return true, password
}
