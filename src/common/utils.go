package common

import (
	"cluster-membership-go/src/model"
	"fmt"
)

//ArrayToString convert an array to string, using 'sep' argument as separator
func ArrayToString(array []string, sep string) string {
	if len(array) == 0 {
		return ""
	}
	var ans = array[0]
	for i := 1; i < len(array); i++ {
		ans += sep + array[i]
	}
	return ans
}

//CheckNonEmpty array
func CheckNonEmpty(array []string, name string) error {
	for i := 0; i < len(array); i++ {
		if len(array[i]) == 0 {
			return &ArgumentParsingError{ErrorS: fmt.Sprintf("The %s argument has an empty string", name)}
		}
	}
	return nil
}

//EqualsArrays compare arrays of any type
func EqualsArrays(a []model.Node, b []model.Node) bool {

	var al = len(a)
	var bl = len(b)

	if al != bl {
		return false
	}

	for i := 0; i < al; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true

}
