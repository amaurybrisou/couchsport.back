package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ParseBody(i interface{}, body io.Reader) (interface{}, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &i)

	if err != nil {
		return nil, err
	}

	return i, nil
}
