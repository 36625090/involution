package utils

import (
	"bytes"
	"encoding/json"
)

func JSONDump(s interface{})string{
	data, _ := json.Marshal(s)
	return string(data)
}

func JSONPrettyDump(s interface{}) string {
	data, _ := json.Marshal(s)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return ""
	}
	return prettyJSON.String()
}
