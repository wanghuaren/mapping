package uts

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

func Struct2Map(_struct interface{}) map[string]interface{} {
	_type := reflect.TypeOf(_struct)
	_value := reflect.ValueOf(_struct)
	_itemMap := StructReflect2Map(_type, _value)
	return _itemMap
}

func StructReflect2Map(_type reflect.Type, _value reflect.Value) map[string]interface{} {
	if _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
		_value = _value.Elem()
	}
	_itemMap := map[string]interface{}{}
	for n := 0; n < _type.NumField(); n++ {
		_fValue := _value.Field(n)
		_fType := _type.Field(n)
		if _fType.Name == "unknownFields" || _fType.Name == "sizeCache" {
			continue
		}
		switch _fType.Type.Kind() {
		case reflect.String:
			_itemMap[_fType.Name] = _fValue.String()
		case reflect.Int:
			_itemMap[_fType.Name] = _fValue.Int()
		case reflect.Int8:
			_itemMap[_fType.Name] = int8(_fValue.Int())
		case reflect.Int16:
			_itemMap[_fType.Name] = int16(_fValue.Int())
		case reflect.Int32:
			_itemMap[_fType.Name] = int32(_fValue.Int())
		case reflect.Int64:
			_itemMap[_fType.Name] = int64(_fValue.Int())
		case reflect.Float32, reflect.Float64:
			_itemMap[_fType.Name] = _fValue.Float()
		case reflect.Bool:
			_itemMap[_fType.Name] = _fValue.Bool()
		case reflect.Struct:
			_itemMap[_fType.Name] = Struct2Map(_fValue.Interface())
		case reflect.Slice:
			_itemMap[_fType.Name] = forMatSlice(_fValue)
		case reflect.Ptr:
			if _fValue.IsNil() {
				_itemMap[_fType.Name] = "nil"
			} else {
				_v := _fValue.Elem()
				_itemMap[_fType.Name] = StructReflect2Map(_v.Type(), _v)
			}
		}
	}
	return _itemMap
}

func forMatSlice(_sliceValue reflect.Value) []interface{} {
	if _sliceValue.IsNil() {
		return nil
	} else {
		_ret := []interface{}{}

		for i := 0; i < _sliceValue.Len(); i++ {
			_mValue := _sliceValue.Index(i)
			switch _mValue.Kind() {
			case reflect.Struct:
				_ret = append(_ret, Struct2Map(_mValue.Interface()))
			case reflect.Slice:
				_ret = append(_ret, forMatSlice(_mValue))
			case reflect.String:
				_ret = append(_ret, _mValue.String())
			case reflect.Int:
				_ret = append(_ret, _mValue.Int())
			case reflect.Int8:
				_ret = append(_ret, _mValue.Int())
			case reflect.Int16:
				_ret = append(_ret, _mValue.Int())
			case reflect.Int32:
				_ret = append(_ret, _mValue.Int())
			case reflect.Int64:
				_ret = append(_ret, _mValue.Int())
			case reflect.Float32, reflect.Float64:
				_ret = append(_ret, _mValue.Float())
			case reflect.Bool:
				_ret = append(_ret, _mValue.Bool())
			case reflect.Ptr:
				if _mValue.IsNil() {
					_ret = append(_ret, nil)
				} else {
					_ret = append(_ret, forMatInterface(_mValue))
				}
			}
		}
		return _ret
	}
}

func forMatInterface(_interfaceValue reflect.Value) interface{} {
	_mValue := _interfaceValue.Elem()
	switch _mValue.Kind() {
	case reflect.Struct:
		return Struct2Map(_mValue.Interface())
	case reflect.Slice:
		return forMatSlice(_mValue)
	case reflect.String:
		return _mValue.String()
	case reflect.Int:
		return _mValue.Int()
	case reflect.Int8:
		return _mValue.Int()
	case reflect.Int16:
		return _mValue.Int()
	case reflect.Int32:
		return _mValue.Int()
	case reflect.Int64:
		return _mValue.Int()
	case reflect.Float32, reflect.Float64:
		return _mValue.Float()
	case reflect.Bool:
		return _mValue.Bool()
	}
	return nil
}

func Map2Struct(_struct interface{}, mapValue interface{}, skipErrPrint ...bool) {
	_value := reflect.ValueOf(_struct).Elem()
	_type := reflect.TypeOf(_struct).Elem()

	switch mapValue := mapValue.(type) {
	case map[string]string:
		for i := 0; i < _type.NumField(); i++ {
			_fType := _type.Field(i)
			_fieldName := _fType.Name
			var _fieldValue string
			_hasField := false
			_fieldValue, _hasField = mapValue[_fieldName]
			if !_hasField {
				_fieldValue, _hasField = mapValue[strings.ToLower(_fieldName)]
			}
			if _hasField {
				switch _fType.Type.Kind() {
				case reflect.String:
					_value.Field(i).SetString(_fieldValue)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					var _mValue int64 = 0
					_v, err := strconv.ParseInt(_fieldValue, 10, 64)
					if len(skipErrPrint) > 0 {
						if !ChkErrNormal(err) {
							_mValue = _v
						}
					} else if !ChkErr(err) {
						_mValue = _v
					}
					_value.Field(i).SetInt(_mValue)
				case reflect.Float32, reflect.Float64:
					var _mValue float64 = 0
					_v, err := strconv.ParseFloat(_fieldValue, 64)
					if len(skipErrPrint) > 0 {
						if !ChkErrNormal(err) {
							_mValue = _v
						}
					} else if !ChkErr(err) {
						_mValue = _v
					}
					_value.Field(i).SetFloat(_mValue)
				case reflect.Bool:
					var _mValue bool = true
					if _fieldValue == "" || _fieldValue == "0" {
						_mValue = false
					}
					_value.Field(i).SetBool(_mValue)
				default:
					Log("Struct 字段 " + _fieldName + " 类型不对")
				}
			} else {
				// Log("Struct 字段 " + _fieldName + " 没有对应的值")
			}
		}
	case map[string]interface{}:
		for i := 0; i < _type.NumField(); i++ {
			_fType := _type.Field(i)
			_fieldName := _fType.Name
			var _fieldValue interface{}
			_hasField := false
			_fieldValue, _hasField = mapValue[_fieldName]
			if !_hasField {
				_fieldValue, _hasField = mapValue[strings.ToLower(_fieldName)]
			}
			if _hasField {
				switch _fType.Type.Kind() {
				case reflect.String:
					var _mValue = ""
					switch reflect.TypeOf(_fieldValue).Kind() {
					case reflect.String:
						_mValue = _fieldValue.(string)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						_mValue = strconv.FormatInt(_fieldValue.(int64), 10)
					case reflect.Float32, reflect.Float64:
						_mValue = strconv.FormatFloat(_fieldValue.(float64), 'g', 32, 64)
					case reflect.Bool:
						_mValue = strconv.FormatBool(_fieldValue.(bool))
					default:
						Log("Struct 字段 " + _fieldName + " 类型不对")
					}
					_value.Field(i).SetString(_mValue)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					var _mValue int64 = 0
					switch reflect.TypeOf(_fieldValue).Kind() {
					case reflect.String:
						_v, err := strconv.ParseInt(_fieldValue.(string), 10, 64)
						if len(skipErrPrint) > 0 {
							if !ChkErrNormal(err) {
								_mValue = _v
							}
						} else if !ChkErr(err) {
							_mValue = _v
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						_mValue = _fieldValue.(int64)
					case reflect.Float32, reflect.Float64:
						_mValue = int64(math.Floor(_fieldValue.(float64)))
					case reflect.Bool:
						if _fieldValue.(bool) {
							_mValue = 1
						}
					default:
						Log("Struct 字段 " + _fieldName + " 类型不对")
					}
					_value.Field(i).SetInt(_mValue)
				case reflect.Float32, reflect.Float64:
					var _mValue float64 = 0
					switch reflect.TypeOf(_fieldValue).Kind() {
					case reflect.String:
						_v, err := strconv.ParseFloat(_fieldValue.(string), 64)
						if len(skipErrPrint) > 0 {
							if !ChkErrNormal(err) {
								_mValue = _v
							}
						} else if !ChkErr(err) {
							_mValue = _v
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						_mValue = float64(_fieldValue.(int64))
					case reflect.Float32, reflect.Float64:
						_mValue = _fieldValue.(float64)
					case reflect.Bool:
						if _fieldValue.(bool) {
							_mValue = 1
						}
					default:
						Log("Struct 字段 " + _fieldName + " 类型不对")
					}
					_value.Field(i).SetFloat(_mValue)
				case reflect.Bool:
					var _mValue bool = false
					switch reflect.TypeOf(_fieldValue).Kind() {
					case reflect.String:
						if len(_fieldValue.(string)) > 0 {
							_mValue = true
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if _fieldValue.(int64) != 0 {
							_mValue = true
						}
					case reflect.Float32, reflect.Float64:
						if _fieldValue.(float64) != 0 {
							_mValue = true
						}
					case reflect.Bool:
						_mValue = _fieldValue.(bool)
					default:
						Log("Struct 字段 " + _fieldName + " 类型不对")
					}
					_value.Field(i).SetBool(_mValue)
				}
			}
		}
	}

}

// func EnGobDB(_data interface{}) []byte
// func DeGobDB(_data interface{}, _b []byte)

func InterfaceCompare(a interface{}, b interface{}) bool {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	switch aValue.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if bValue.Type().Kind() != reflect.Int && bValue.Type().Kind() != reflect.Int8 && bValue.Type().Kind() != reflect.Int16 && bValue.Type().Kind() != reflect.Int32 && bValue.Type().Kind() != reflect.Int64 {
			return false
		} else if aValue.Int() != bValue.Int() {
			return false
		}
	case reflect.String:
		if bValue.Type().Kind() != reflect.String {
			return false
		} else if aValue.String() != bValue.String() {
			return false
		}
	case reflect.Bool:
		if bValue.Type().Kind() != reflect.Bool {
			return false
		} else if aValue.Bool() != bValue.Bool() {
			return false
		}
	default:
		Log("类型 " + aValue.Type().Kind().String() + " 未增加对比")
		return false
	}
	return true
}

func InterfaceCompareString(a interface{}, b string) bool {
	aValue := reflect.ValueOf(a)
	switch aValue.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b == "" {
			return false
		}
		value, err := strconv.ParseInt(b, 10, 64)
		if ChkErrNormal(err) {
			return false
		} else if aValue.Int() != value {
			return false
		}
	case reflect.String:
		if aValue.String() != b {
			return false
		}
	case reflect.Bool:
		if b == "" {
			return false
		}
		value, err := strconv.ParseBool(b)
		if ChkErrNormal(err) {
			return false
		} else if aValue.Bool() != value {
			return false
		}
	default:
		Log("类型 " + aValue.Type().Kind().String() + " 未增加对比")
		return false
	}
	return true
}
