package util

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// StrPointer returns a pointer of the string
func StrPointer(s string) *string {
	return &s
}

// IntPointer returns a pointer of the int
func IntPointer(i int32) *int32 {
	return &i
}

// Int64Pointer returns a pointer of the int64
func Int64Pointer(i int64) *int64 {
	return &i
}

// BoolPointer returns a point of the bool
func BoolPointer(b bool) *bool {
	return &b
}

// PointerToBool returns the bool from a pointer
func PointerToBool(flag *bool) bool {
	if flag == nil {
		return false
	}
	return *flag
}

// PointerToString returns the string from a pointer
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PointerToInt32 returns the int32 from a pointer
func PointerToInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// IntstrPointer returns a point of the intOrString
func IntstrPointer(i int) *intstr.IntOrString {
	is := intstr.FromInt(i)
	return &is
}

// MergeStringMaps merges two maps of strings
func MergeStringMaps(l map[string]string, l2 map[string]string) map[string]string {
	merged := make(map[string]string)
	if l == nil {
		l = make(map[string]string)
	}
	for lKey, lValue := range l {
		merged[lKey] = lValue
	}
	for lKey, lValue := range l2 {
		merged[lKey] = lValue
	}
	return merged
}

// MergeMultipleStringMaps merges multiple maps of strings
func MergeMultipleStringMaps(stringMaps ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, stringMap := range stringMaps {
		merged = MergeStringMaps(merged, stringMap)
	}
	return merged
}

// EmptyTypedStrSlice .
func EmptyTypedStrSlice(s ...string) []interface{} {
	ret := make([]interface{}, len(s))
	for i := 0; i < len(s); i++ {
		ret[i] = s[i]
	}
	return ret
}

// EmptyTypedFloatSlice .
func EmptyTypedFloatSlice(f ...float64) []interface{} {
	ret := make([]interface{}, len(f))
	for i := 0; i < len(f); i++ {
		ret[i] = f[i]
	}
	return ret
}

// ContainsString .
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString .
func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// GetPodSecurityContextFromSecurityContext .
func GetPodSecurityContextFromSecurityContext(sc *corev1.SecurityContext) *corev1.PodSecurityContext {
	if sc == nil || *sc == (corev1.SecurityContext{}) {
		return &corev1.PodSecurityContext{}
	}
	return &corev1.PodSecurityContext{
		RunAsGroup:   sc.RunAsGroup,
		RunAsNonRoot: sc.RunAsNonRoot,
		RunAsUser:    sc.RunAsUser,
		FSGroup:      sc.RunAsGroup,
	}
}
