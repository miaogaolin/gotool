package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func SortStructSlice(arr interface{}, asc bool, keys ...string) error {
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	tempArr := make([]reflect.Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		tempArr[i] = v.Index(i)
	}
	for i := 0; i < v.Len(); i++ {
		for j := i + 1; j < v.Len(); j++ {
			count, err := Compare(tempArr[i], tempArr[j], keys...)
			if err != nil {
				return err
			}

			if asc {
				if count > 0 {
					tempArr[i], tempArr[j] = tempArr[j], tempArr[i]
				}
			} else {
				if count < 0 {
					tempArr[i], tempArr[j] = tempArr[j], tempArr[i]
				}
			}
		}
	}
	newv := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		newv.Index(i).Set(tempArr[i])
	}
	reflect.Copy(v, newv)
	return nil
}

func Compare(v1 reflect.Value, v2 reflect.Value, keys ...string) (int, error) {
	if v1.Kind() == reflect.Ptr || v1.Kind() == reflect.Interface {
		v1 = v1.Elem()
	}
	if v2.Kind() == reflect.Ptr || v2.Kind() == reflect.Interface {
		v2 = v2.Elem()
	}
	if v1.Kind() != v2.Kind() {
		return 0, errors.New("比较对象类型不同，无法比较")
	}

	if len(keys) > 0 {
		for _, key := range keys {
			var count int
			var err error
			if v1.Kind().String() == "map" && v2.Kind().String() == "map" {
				keyv := reflect.ValueOf(key)
				count, err = CompareBase(v1.MapIndex(keyv), v2.MapIndex(keyv))
			} else {
				count, err = CompareBase(v1.FieldByName(key), v2.FieldByName(key))
			}

			if err != nil {
				return 0, err
			}
			if count != 0 {
				return count, nil
			}
		}
	} else {
		count, err := CompareBase(v1, v2)
		if err != nil {
			return 0, err
		}
		if count != 0 {
			return count, nil
		}
	}
	return 0, nil
}

func CompareBase(value1 reflect.Value, value2 reflect.Value) (int, error) {
	if value1.Kind() == reflect.Ptr || value1.Kind() == reflect.Interface {
		value1 = value1.Elem()
	}
	if value2.Kind() == reflect.Ptr || value2.Kind() == reflect.Interface {
		value2 = value2.Elem()
	}
	if fmt.Sprint(value1) == fmt.Sprint(value2) {
		return 0, nil
	}
	switch value1.Type().Kind() {
	case reflect.String:
		v1v := value1.String()
		v2v := value2.String()
		index := len(v1v)
		if len(v1v) > len(v2v) {
			index = len(v2v)
		}
		for i := 0; i < index; i++ {
			if v1v[i] == v2v[i] {
				continue
			} else {
				if v1v[i] > v2v[i] {
					return 1, nil
				} else {
					return -1, nil
				}
			}
		}
		return len(v1v) - len(v2v), nil
	case reflect.Bool:
		if value1.Bool() {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		if value1.Int() > value2.Int() {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Float32, reflect.Float64:
		if value1.Float() > value2.Float() {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		if value1.Uint() > value2.Uint() {
			return 1, nil
		} else {
			return -1, nil
		}
	default:
		return 0, errors.New("无法比较")
	}
}
