package ctx

import "encoding/json"

type Tags map[string]json.RawMessage

func (t Tags) Set(key string, val interface{}) {
	switch val := val.(type) {
	case json.RawMessage:
		t[key] = val
	default:
		j, err := json.Marshal(val)
		if err != nil {
			j, _ = json.Marshal(err.Error()) // we put the error if marshal fails
		}
		t[key] = j
	}
}
