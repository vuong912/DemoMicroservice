package controllers

import (
	"fmt"
	"log"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	var testcases = []struct {
		data     string
		expected bool
	}{
		{"", false},
		{"plainaddress", false},
		{"#@%^%#$@#$@#.com", false},
		{"@domain.com", false},
		{"Joe Smith <email@domain.com>", false},
		{"email.domain.com", false},
		{"email@domain@domain.com", false},
		{".email@domain.com", false},
		{"email.@domain.com", false},
		{"email..email@domain.com", false},
		{"あいうえお@domain.com", false},
		{"email@domain.com (Joe Smith)", false},
		{"email@domain", false},
		{"email@-domain.com", false},
		{"email@domain.web", false},
		{"email@111.222.333.44444", false},
		{"email@domain..com", false},
		{"email@domain.com", true},
		{"firstname.lastname@domain.com", true},
		{"email@subdomain.domain.com", true},
		{"firstname+lastname@domain.com", true},
		{"email@123.123.123.123", true},
		{"email@[123.123.123.123]", true},
		{"\"email\"@domain.com", true},
		{"1234567890@domain.com", true},
		{"email@domain-one.com", true},
		{"_______@domain.com", true},
		{"email@domain.name", true},
		{"email@domain.co.jp", true},
		{"firstname-lastname@domain.com", true},
	}
	var totalTestcase, countCorrectResult = len(testcases), 0
	for i, test := range testcases {
		if result := ValidateEmail(test.data); result != test.expected {
			t.Errorf("Incorrect %v \n\tTest: %v\n\tCorrect output: %v\n\tYour output: %v\n",
				i, test.data, test.expected, result)
		} else {
			countCorrectResult++
		}
	}
	fmt.Printf("=> %d/%d  %f%% \n", countCorrectResult, totalTestcase, float32(countCorrectResult*100)/float32(totalTestcase))
}

func TestValidatePhoneNumber(t *testing.T) {
	var testcases = []struct {
		data     string
		expected bool
	}{
		{"", false},
		{" ", false},
		{"032#325#6002", false},
		{"032b325a6002", false},
		{"0+3232560023", false},
		{"06546456565635", false},
		{"1234567890", false},
		{"+84234567890", false},
		{"032 325 6002", false},
		{"032.325.6002", false},
		{"0354545455", false},
		{"0132556498", false},
		{"0946552123", true},
		{"01256464523", true},
	}
	var totalTestcase, countCorrectResult = len(testcases), 0
	for i, test := range testcases {
		if result := ValidatePhoneNumber(test.data); result != test.expected {
			t.Errorf("Incorrect %v \n\tTest: %v\n\tCorrect output: %v\n\tYour output: %v\n",
				i, test.data, test.expected, result)
		} else {
			countCorrectResult++
		}
	}

	log.Printf("=> %d/%d  %f%% \n", countCorrectResult, totalTestcase, float32(countCorrectResult*100)/float32(totalTestcase))
}
