package set

import "reflect"

// InSlice 是否存在切片中
func InSlice(mix interface{}, slice interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		forElement := reflect.ValueOf(slice)
		for i := 0; i < forElement.Len(); i++ {
			if reflect.DeepEqual(mix, forElement.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
