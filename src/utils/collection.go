package utils

func Contains(sArray []string,s string) bool {
	for _, ss := range sArray {
		if ss == s {return true}
	}
	return false
}
