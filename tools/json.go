package tools

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJsonStringE(v any) (string, error) {
	arr, err := Marshal(v)
	return string(arr), err
}

func ToJsonString(v any) string {
	s, _ := ToJsonStringE(v)
	return s
}

func ToJsonStringWithMaxLen(v any, max int) string {
	arr, _ := Marshal(v)
	if len(arr) > max {
		arr = arr[:max]
	}
	return string(arr)
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(input []byte, v any) error {
	return json.Unmarshal(input, v)
}
