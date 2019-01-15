package utils

import "encoding/json"

func MapToJson(mapItem map[string]interface{}) string {
	d, _ := json.Marshal(mapItem)
	return string(d)
}
