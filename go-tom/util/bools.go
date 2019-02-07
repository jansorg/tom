package util

func FalseP() *bool {
	f := false
	return &f
}

func TrueP() *bool {
	f := true
	return &f
}
