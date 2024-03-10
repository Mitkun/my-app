package utils

func StringFromPointer(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func GetStrPt(s string) *string {
	return &s
}
