package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"strings"
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

// FromBuffer accepts a b64 image string (having content-type specified) and returns a
// base64 encoded string.
func B64ImageToFile(b64 string) (image.Image, error) {
	i := strings.Index(b64, ",")
	if i < 0 {
		return nil, fmt.Errorf("no comma in image")
	}

	start := strings.Index(b64, ":")
	end := strings.Index(b64, ";")

	mime := b64[start+1 : end]
	// pass reader to NewDecoder
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64[i+1:]))

	switch mime {
	case "image/png":
		return png.Decode(dec)
	case "image/jpg":
	case "image/jpeg":
		return jpeg.Decode(dec)
	case "image/gif":
		return gif.Decode(dec)
	default:
		return nil, fmt.Errorf("image format not allowed: %s", mime)
	}
	return nil, nil
}
