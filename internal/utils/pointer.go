package utils

// it converts boolean value to
// boolean pointer value.
func BoolToPoint(val bool) *bool {

	return &val
}

func StringToPoint(val string) *string {
	return &val
}

func IntToPoint(val int) *int {

	return &val
}
