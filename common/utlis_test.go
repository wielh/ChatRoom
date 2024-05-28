package common

import (
	"fmt"
	"testing"
)

func TestToCommonID(t *testing.T) {
	type inputStruct struct {
		username string
		userType int
	}
	type testCase struct {
		input          inputStruct
		expectedOutput string
	}

	getName := func(t *testCase) string {
		return fmt.Sprintf("input:%+v", t.input)
	}

	testCases := []testCase{
		{input: inputStruct{"abc", 0}, expectedOutput: "googleID:abc"},
		{input: inputStruct{"abc", 1}, expectedOutput: ""},
		{input: inputStruct{"abc", -1}, expectedOutput: ""},
	}

	for _, tt := range testCases {
		t.Run(getName(&tt), func(t *testing.T) {
			if output := ToCommonID(tt.input.username, UserType(tt.input.userType)); output != tt.expectedOutput {
				t.Errorf("input = %+v, output = %v, expectedOutput %v", tt.input, output, tt.expectedOutput)
			}
		})
	}
}

func TestToSpecificID(t *testing.T) {
	type outputStruct struct {
		username string
		userType int
	}
	type testCase struct {
		input          string
		expectedOutput outputStruct
	}

	getName := func(t *testCase) string {
		return fmt.Sprintf("input:%+v", t.input)
	}

	testCases := []testCase{
		{input: "googleID:abc", expectedOutput: outputStruct{username: "googleID", userType: 0}},
		{input: "abc:abc", expectedOutput: outputStruct{username: "", userType: -1}},
	}

	for _, tt := range testCases {
		t.Run(getName(&tt), func(t *testing.T) {
			username, userType := ToSpecificID(tt.input)
			output := outputStruct{username: username, userType: int(userType)}
			if output != tt.expectedOutput {
				t.Errorf("input = %+v, output = %v, expectedOutput %v", tt.input, output, tt.expectedOutput)
			}
		})
	}
}

func TestTokenVerify(t *testing.T) {
	type inputStruct struct {
		username string
		userID   string
		userType int
	}

	type outputStruct struct {
		userID   string
		errorNil bool
	}

	type testCase struct {
		input          inputStruct
		expectedOutput outputStruct
	}

	getName := func(t *testCase) string {
		return fmt.Sprintf("input:%+v", t.input)
	}

	testCases := []testCase{
		{input: inputStruct{userID: "abc", userType: 0, username: "abc"}, expectedOutput: outputStruct{userID: "googleID:abc", errorNil: true}},
		{input: inputStruct{userID: "abc", userType: 1, username: "abc"}, expectedOutput: outputStruct{userID: "", errorNil: false}},
		{input: inputStruct{userID: "xyz", userType: 0, username: "xyz"}, expectedOutput: outputStruct{userID: "googleID:xyz", errorNil: true}},
		{input: inputStruct{userID: "xyz", userType: 1, username: "xyz"}, expectedOutput: outputStruct{userID: "", errorNil: false}},
	}

	for _, tt := range testCases {
		t.Run(getName(&tt), func(t *testing.T) {
			token := CreateToken(tt.input.userID, int32(tt.input.userType), tt.input.username)
			claims, err := ParseJWT(token)
			userID := ""
			if claims["userID"] != nil {
				userID = fmt.Sprintf("%v", claims["userID"])
			}

			output := outputStruct{userID: userID, errorNil: (err == nil)}
			if output != tt.expectedOutput {
				t.Errorf("input = %+v, output = %v, expectedOutput %v", tt.input, output, tt.expectedOutput)
			}
		})
	}
}
