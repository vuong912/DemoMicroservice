package controllers

import (
	"errors"
	"regexp"

	"github.com/DemoMicroservice/EmployeeService/models"
)

func ValidateEmail(data string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(data)
}
func ValidatePhoneNumber(data string) bool {
	//re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	//re := regexp.MustCompile(`^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$`)
	re := regexp.MustCompile(`(09|01[2|6|8|9])+([0-9]{8})\b`)
	return re.MatchString(data)
}

func ValidateEmployee(employee *models.Employee) error {
	if employee == nil {
		return errors.New("Employee require")
	}
	if employee.Name == "" {
		return errors.New("Name require")
	}
	if len(employee.Name) > 100 {
		return errors.New("Length of name must less than or equal to 100")
	}
	if !ValidatePhoneNumber(employee.PhoneNumber) {
		return errors.New("Invalid phone number")
	}
	if !ValidateEmail(employee.Email) {
		return errors.New("Invalid phone number")
	}
	if employee.IdBranch == "" {
		return errors.New("Branch require")
	}
	if employee.Rangewage <= 0 {
		return errors.New("Rangewage must greater than 0")
	}
	return nil
}
