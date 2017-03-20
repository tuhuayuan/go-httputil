package httputil

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// URLValuesUnpacker TODO 自定义解析接口
type URLValuesUnpacker interface {
}

// UnpackRequest 把请求参数映射到制定结构体
func UnpackRequest(r *http.Request, ptrStruct interface{}) *ValidateErrors {
	err := r.ParseForm()
	if err != nil {
		return &ValidateErrors{
			ValidateError{
				Code:    -1,
				Message: "parse request form error " + err.Error(),
			},
		}
	}
	return UnpackURLValues(r.Form, ptrStruct)
}

// UnpackURLValues 把map[string][]string映射到制定结构体
func UnpackURLValues(values url.Values, ptrStruct interface{}) *ValidateErrors {
	errs := &ValidateErrors{}

	structType := reflect.TypeOf(ptrStruct)
	if structType.Kind() != reflect.Ptr {
		errs.Add([]string{}, -1, "must unpack to a point")
		return errs
	}
	strcutValue := reflect.ValueOf(ptrStruct)

	mapURLValues(strcutValue, values, errs)
	if errs.Len() == 0 {
		return nil
	}
	return errs
}

// mapURLValues 映射url.Values到reflect.Value
// 支持Slice、Struct、Ptr
func mapURLValues(structValue reflect.Value, form url.Values, errs *ValidateErrors) {
	// 获取指针指向
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	formStructType := structValue.Type()

	for i := 0; i < formStructType.NumField(); i++ {
		fieldType := formStructType.Field(i)
		fieldValue := structValue.Field(i)

		if fieldType.Type.Kind() == reflect.Ptr && fieldType.Anonymous {
			// 递归匿名指针字段
			fieldValue.Set(reflect.New(fieldType.Type.Elem()))
			mapURLValues(fieldValue.Elem(), form, errs)
			// 设置空值
			if reflect.DeepEqual(fieldValue.Elem().Interface(), reflect.Zero(fieldValue.Elem().Type()).Interface()) {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			}
		} else if fieldType.Type.Kind() == reflect.Struct {
			// 递归结构体
			mapURLValues(fieldValue, form, errs)
		} else {
			tagValue, ok := fieldType.Tag.Lookup("form")
			if !ok {
				continue
			}
			omit := false
			mapName := fieldType.Name
			tags := strings.SplitN(tagValue, ",", 2)
			// TODO: tags support ,alias
			if len(tags) >= 1 {
				mapName = tags[0]
			}
			if len(tags) >= 2 {
				if tags[1] == "omit" || tags[1] == "optional" {
					omit = true
				}
			}

			urlValues, ok := form[mapName]
			if !ok {
				if !omit {
					errs.Add([]string{mapName}, -1, "field missing")
				}
				continue
			}
			valuesLen := len(urlValues)
			if fieldType.Type.Kind() == reflect.Slice {
				sliceType := fieldType.Type.Elem()
				slice := reflect.MakeSlice(fieldType.Type, valuesLen, valuesLen)
				for i := 0; i < valuesLen; i++ {
					err := setProperType(sliceType, slice.Index(i), urlValues[i])
					if err != nil {
						errs.Add([]string{mapName}, -1, err.Error())
					}
				}
				fieldValue.Set(slice)
			} else {
				err := setProperType(fieldType.Type, fieldValue, urlValues[0])
				if err != nil {
					errs.Add([]string{mapName}, -1, err.Error())
				}
			}
		}
	}
}

// setProperType 支持设置以下原生类型以及它们的指针
// Int Int8 Int16 Int32 Int64
// Uint Uint8 Uint16 Uint32 Uint64
// bool
// float32 float64
// string
func setProperType(fieldType reflect.Type, fieldValue reflect.Value, val string) (err error) {
	if !fieldValue.CanSet() {
		return errors.New("field can't set")
	}
	// 这里只解引用一次，不去递归处理，所以不支持指针的指针。。。。。。
	if fieldType.Kind() == reflect.Ptr {
		underType := fieldType.Elem()
		underValue := reflect.New(underType)
		fieldValue.Set(underValue)
		// fieldValue 必定是CanSet
		fieldValue = underValue.Elem()
		fieldType = underType
	}
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intVal int64
		if intVal, err = strconv.ParseInt(val, 10, 64); err == nil {
			fieldValue.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var uintVal uint64
		if uintVal, err = strconv.ParseUint(val, 10, 64); err == nil {
			fieldValue.SetUint(uintVal)
		}
	case reflect.Bool:
		var boolVal bool
		if boolVal, err = strconv.ParseBool(val); err == nil {
			fieldValue.SetBool(boolVal)
		}
	case reflect.Float32:
		var floatVal float64
		if floatVal, err = strconv.ParseFloat(val, 32); err == nil {
			fieldValue.SetFloat(floatVal)
		}
	case reflect.Float64:
		var floatVal float64
		floatVal, err = strconv.ParseFloat(val, 64)
		if err != nil {
			fieldValue.SetFloat(floatVal)
		}
	case reflect.String:
		fieldValue.SetString(val)
	default:
		err = errors.New("not supported kind " + fieldType.Kind().String())
	}
	return err
}
