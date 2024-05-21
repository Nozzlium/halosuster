package util

import (
	"regexp"

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
	employeeId string,
) error {
	regex := "^[615]{3}[1-2]{1}(200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])[0-9]{3,5}$"
	idStringRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !idStringRegex.MatchString(
		employeeId,
	) {
		return constant.ErrBadInput
	}

	regex = "^[615]{3}[0-9]{10,12}$"
	roleStringRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !roleStringRegex.MatchString(
		employeeId,
	) {
		return constant.ErrNotFound
	}

	return nil
}

func ValidateIdentityNumber(
	identityNumber string,
) error {
	regex := "^[0-9]{16}$"
	identityNumberRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !identityNumberRegex.MatchString(
		identityNumber,
	) {
		return constant.ErrBadInput
	}

	return nil
}

func ValidatePhoneNumber(
	phoneNumber string,
) error {
	regex := "^[+]{1}[62]{2}[0-9]{7,12}$"
	phoneNumberRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !phoneNumberRegex.MatchString(
		phoneNumber,
	) {
		return constant.ErrBadInput
	}

	return nil
}

func ValidateIsANurse(
	employeeId string,
) error {
	regex := "^[303]{3}[0-9]{10,12}$"
	roleStringRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !roleStringRegex.MatchString(
		employeeId,
	) {
		return constant.ErrNotFound
	}

	return nil
}

func ValidateGeneralEmployeeID(
	employeeId string,
) error {
	regex := "^[0-9]{3}[1-2]{1}(200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])[0-9]{3,5}$"
	idStringRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !idStringRegex.MatchString(
		employeeId,
	) {
		return constant.ErrBadInput
	}

	return nil
}
