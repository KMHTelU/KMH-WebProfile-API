package utils

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Field JSON yang tidak boleh pernah keluar ke klien.
var sensitiveJSONFields = map[string]struct{}{
	"password_hash": {},
	"password":      {},
	"token_hash":    {},
}

// SanitizeForJSON mengubah nilai sql.Null* / uuid.NullUUID menjadi tipe JSON
// biasa (string, bool, null, dsb.) supaya respons API tidak berbentuk
// {"Bool":true,"Valid":true}.
//
// Field sensitif seperti password_hash ikut dihilangkan.
func SanitizeForJSON(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return sanitizeValue(reflect.ValueOf(value))
}

func sanitizeValue(value reflect.Value) interface{} {
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	if !value.IsValid() {
		return nil
	}

	if unwrapped, ok := unwrapNullable(value); ok {
		return unwrapped
	}

	// json.RawMessage (kolom JSONB) adalah []byte berisi teks JSON mentah.
	// Tanpa penanganan khusus ia jatuh ke cabang Slice generik dan berubah
	// menjadi deretan angka byte (mis. ["a"] → [91,34,97,34,93]).
	if value.CanInterface() {
		if raw, ok := value.Interface().(json.RawMessage); ok {
			return decodeRawJSON(raw)
		}
	}

	switch value.Kind() {
	case reflect.Struct:
		return sanitizeStruct(value)
	case reflect.Slice, reflect.Array:
		if value.Kind() == reflect.Slice && value.IsNil() {
			return []interface{}{}
		}
		out := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			out[i] = sanitizeValue(value.Index(i))
		}
		return out
	case reflect.Map:
		if value.IsNil() {
			return map[string]interface{}{}
		}
		out := make(map[string]interface{}, value.Len())
		for _, key := range value.MapKeys() {
			out[stringifyMapKey(key)] = sanitizeValue(value.MapIndex(key))
		}
		return out
	default:
		if value.CanInterface() {
			return value.Interface()
		}
		return nil
	}
}

// decodeRawJSON menerjemahkan JSONB mentah menjadi nilai JSON biasa agar
// tersisip apa adanya di respons (bukan array byte). Bila isinya bukan JSON
// valid, kembalikan sebagai string agar tidak menjatuhkan respons.
func decodeRawJSON(raw json.RawMessage) interface{} {
	if len(raw) == 0 {
		return nil
	}
	var out interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return string(raw)
	}
	return out
}

func unwrapNullable(value reflect.Value) (interface{}, bool) {
	if !value.CanInterface() {
		return nil, false
	}

	switch typed := value.Interface().(type) {
	case sql.NullString:
		if !typed.Valid {
			return nil, true
		}
		return typed.String, true
	case sql.NullBool:
		if !typed.Valid {
			return nil, true
		}
		return typed.Bool, true
	case sql.NullTime:
		if !typed.Valid {
			return nil, true
		}
		return typed.Time.UTC().Format(time.RFC3339Nano), true
	case sql.NullInt16:
		if !typed.Valid {
			return nil, true
		}
		return typed.Int16, true
	case sql.NullInt32:
		if !typed.Valid {
			return nil, true
		}
		return typed.Int32, true
	case sql.NullInt64:
		if !typed.Valid {
			return nil, true
		}
		return typed.Int64, true
	case sql.NullFloat64:
		if !typed.Valid {
			return nil, true
		}
		return typed.Float64, true
	case sql.NullByte:
		if !typed.Valid {
			return nil, true
		}
		return typed.Byte, true
	case uuid.NullUUID:
		if !typed.Valid {
			return nil, true
		}
		return typed.UUID.String(), true
	case uuid.UUID:
		return typed.String(), true
	case time.Time:
		if typed.IsZero() {
			return nil, true
		}
		return typed.UTC().Format(time.RFC3339Nano), true
	}

	return nil, false
}

func sanitizeStruct(value reflect.Value) interface{} {
	typ := value.Type()

	// time.Time sudah ditangani di unwrapNullable, tapi jaga-jaga bila
	// masuk lewat path lain.
	if typ == reflect.TypeOf(time.Time{}) {
		t := value.Interface().(time.Time)
		if t.IsZero() {
			return nil
		}
		return t.UTC().Format(time.RFC3339Nano)
	}

	out := make(map[string]interface{}, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, skip, omitEmpty := jsonFieldName(field)
		if skip {
			continue
		}
		if _, sensitive := sensitiveJSONFields[name]; sensitive {
			continue
		}

		fieldValue := value.Field(i)
		if omitEmpty && isEmptyValue(fieldValue) {
			continue
		}

		out[name] = sanitizeValue(fieldValue)
	}
	return out
}

func jsonFieldName(field reflect.StructField) (name string, skip bool, omitEmpty bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", true, false
	}
	if tag == "" {
		return field.Name, false, false
	}

	name = tag
	if comma := indexByte(tag, ','); comma >= 0 {
		name = tag[:comma]
		omitEmpty = containsToken(tag[comma+1:], "omitempty")
	}
	if name == "" {
		name = field.Name
	}
	return name, false, omitEmpty
}

func isEmptyValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return value.IsNil()
	}
	return false
}

func stringifyMapKey(key reflect.Value) string {
	if key.Kind() == reflect.String {
		return key.String()
	}
	if key.CanInterface() {
		if asString, ok := key.Interface().(string); ok {
			return asString
		}
		return key.Type().String()
	}
	return "key"
}

func indexByte(value string, target byte) int {
	for i := 0; i < len(value); i++ {
		if value[i] == target {
			return i
		}
	}
	return -1
}

func containsToken(value, token string) bool {
	start := 0
	for start <= len(value) {
		end := start
		for end < len(value) && value[end] != ',' {
			end++
		}
		if value[start:end] == token {
			return true
		}
		if end == len(value) {
			return false
		}
		start = end + 1
	}
	return false
}
