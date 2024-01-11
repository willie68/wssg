package utils

import (
	"encoding/json"
	"fmt"
)

func ConvertArrIntToArrString(value any) []string {
	var props []string
	if value != nil {
		switch v := value.(type) {
		case []string:
			props = v
		case string:
			props = []string{v}
		case []interface{}:
			props = make([]string, 0)
			for _, s := range v {
				str := fmt.Sprintf("%v", s)
				props = append(props, str)
			}
		default:
			props = make([]string, 0)
		}
	}
	return props
}

func ConvertMapIntToMapString(value any) map[string]string {
	props := make(map[string]string)
	if value != nil {
		dt, err := json.Marshal(value)
		if err != nil {
			return props
		}
		err = json.Unmarshal(dt, &props)
		if err != nil {
			return props
		}
	}
	return props
}
