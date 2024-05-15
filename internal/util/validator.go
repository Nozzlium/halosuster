package util

import (
	"regexp"
	"time"
)

func ValidateURL(url string) bool {
	regex := `^[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	urlRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return false
	}

	return urlRegex.MatchString(url)
}

// func ValidateUserEmployeeID(
// 	employeeId string,
// ) bool {
// 	regex := `^[615]{3}[1-2]{1}(200[0-9]|201[0-9]|202[1-4])(0[1-9]|1[0-2])[0-9]{3}$`
// 	idRegex, err := regexp.Compile(
// 		regex,
// 	)
// 	if err != nil {
// 		return false
// 	}
//
// 	return idRegex.MatchString(
// 		employeeId,
// 	)
// }

func ValidateNurseEmployeeID(
	employeeId string,
) bool {
	regex := `^[303]{3}[1-2]{1}(200[0-9]|201[0-9]|202[1-4])(0[1-9]|1[0-2])[0-9]{3}$`
	idRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return false
	}

	return idRegex.MatchString(
		employeeId,
	)
}

func ValidateUserEmployeeID(
	employeeId uint64,
) bool {
	if employeeId < uint64(
		6150000000000,
	) {
		return false
	}

	employeeId /= 1000
	currYear := time.Now().Year()

	if monthDigits := employeeId % 100; monthDigits < 1 ||
		monthDigits > 12 {
		return false
	}
	employeeId /= 100

	if yearDigits := employeeId % 10000; yearDigits < 2000 ||
		yearDigits > uint64(currYear) {
		return false
	}
	employeeId /= 10000

	if genderDigit := employeeId % 10; genderDigit < 1 ||
		genderDigit > 2 {
		return false
	}
	employeeId /= 10

	return employeeId == 615
}
