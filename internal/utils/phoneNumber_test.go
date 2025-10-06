package utils

import "testing"

func TestPhoneNumberValidation(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "with prefix (+)",
			Input:    "+1 (415) 123-4567",
			Expected: "+14151234567",
		},
		{
			Name:     "without prefix",
			Input:    "0917-123-4567",
			Expected: "09171234567",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ValidateAndFormatPhoneNumber(tc.Input)
			if err != nil {
				t.Errorf("expetced not error but got %v", err.Error())
				return
			}

			if result != tc.Expected {
				t.Errorf("expected to be equal: %v but got: %v", tc.Expected, result)
				return
			}
		})
	}

	t.Run("fail because contain character", func(t *testing.T) {
		expected := ""
		result, err := ValidateAndFormatPhoneNumber("abc123")
		if err == nil {
			t.Error("expecting error but success")
			return
		}

		if result != expected {
			t.Errorf("expcted an empty string but got: %v", result)
			return
		}
	})
}
