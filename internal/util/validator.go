package util

import (
	"regexp"
	"time"

	"github.com/nozzlium/halosuster/internal/constant"
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

func ValidateUserEmployeeID(
	employeeId uint64,
) error {
	if employeeId < uint64(
		6150000000000,
	) {
		return constant.ErrBadInput
	}

	employeeId /= 1000
	currYear := time.Now().Year()

	if monthDigits := employeeId % 100; monthDigits < 1 ||
		monthDigits > 12 {
		return constant.ErrBadInput
	}
	employeeId /= 100

	if yearDigits := employeeId % 10000; yearDigits < 2000 ||
		yearDigits > uint64(currYear) {
		return constant.ErrBadInput
	}
	employeeId /= 10000

	if genderDigit := employeeId % 10; genderDigit < 1 ||
		genderDigit > 2 {
		return constant.ErrBadInput
	}
	employeeId /= 10

	if employeeId != 615 {
		return constant.ErrNotFound
	}

	return nil
}

func ValidateNurseEmployeeID(
	employeeId uint64,
) error {
	if employeeId < uint64(
		3030000000000,
	) {
		return constant.ErrBadInput
	}

	employeeId /= 1000
	currYear := time.Now().Year()

	if monthDigits := employeeId % 100; monthDigits < 1 ||
		monthDigits > 12 {
		return constant.ErrBadInput
	}
	employeeId /= 100

	if yearDigits := employeeId % 10000; yearDigits < 2000 ||
		yearDigits > uint64(currYear) {
		return constant.ErrBadInput
	}
	employeeId /= 10000

	if genderDigit := employeeId % 10; genderDigit < 1 ||
		genderDigit > 2 {
		return constant.ErrBadInput
	}
	employeeId /= 10

	if employeeId != 303 {
		return constant.ErrNotFound
	}

	return nil
}
