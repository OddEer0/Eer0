package elogger

import "encoding/json"

func (v Value) MarshalJSON() ([]byte, error) {
	r, ok := v.Val.(LogValuer)
	if ok {
		return json.Marshal(r.LogValue())
	}

	return json.Marshal(v.Val)
}
