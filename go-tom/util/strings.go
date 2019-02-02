package util

func MapStrings(data []string) map[string]bool {
	mapping := make(map[string]bool)

	for _, d := range data {
		mapping[d] = true
	}

	return mapping
}

func StringP(s string) *string {
	return &s
}
