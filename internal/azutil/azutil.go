package azutil

import "strings"

func DerefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func DerefInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func DerefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func ExtractRG(id string) string {
	parts := strings.Split(id, "/")
	for i, p := range parts {
		if strings.EqualFold(p, "resourceGroups") && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func ExtractName(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
