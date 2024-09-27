package convert

import "encoding/json"

func interfaceToBytes(i interface{}) []byte {
	var val []byte

	switch d := i.(type) {
	case []byte:
		val = d
	case string:
		val = []byte(d)
	default:
		b, err := json.Marshal(d)
		if err == nil {
			val = b
		}
	}

	return val
}
