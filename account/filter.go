package account

import "encoding/json"

type Filter map[string]interface{}

func BuildFilter(jsonData []byte) (filter Filter, err error) {
	err = json.Unmarshal([]byte(jsonData), &filter)
	return filter, err
}

func BuildEmptyFilter() Filter {
	return make(Filter)
}