package models

import (
	"encoding/json"
	"bytes"
)

// An alias of Buffer that json encoder will marshal to a string and unmarshal from a string.
type OutputHolder bytes.Buffer

func (holder *OutputHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(holder.String())
}

func (holder *OutputHolder) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	holder.Reset()
	_, err := holder.WriteString(s)
	return err
}

func (holder *OutputHolder) WriteString(s string) (int, error) {
	return (*bytes.Buffer)(holder).WriteString(s)
}

func (holder *OutputHolder) Reset() {
	(*bytes.Buffer)(holder).Reset()
}


func (holder *OutputHolder) String() string {
	return (*bytes.Buffer)(holder).String()
}