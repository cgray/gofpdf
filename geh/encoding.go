package geh

import (
	"encoding/gob"
	"encoding/json"

	"github.com/d1ngd0/gofpdf/bp"
)

func EncodeMany(v ...interface{}) ([]byte, error) {
	var err error
	w := bp.GetBuffer()
	encoder := gob.NewEncoder(w)

	for x := 0; x < len(v); x++ {
		if err == nil {
			err = encoder.Encode(v[x])
		}
	}

	return w.Bytes(), err
}

// GobDecode decodes the specified byte buffer into the receiving template.
func DecodeMany(buf []byte, v ...interface{}) error {
	r := bp.GetFilledBuffer(buf)
	defer bp.PutBuffer(r)

	decoder := gob.NewDecoder(r)
	for x := 0; x < len(v); x++ {
		if err := decoder.Decode(v[x]); err != nil {
			return err
		}
	}

	return nil
}

func EncodeManyJSON(v ...interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// GobDecode decodes the specified byte buffer into the receiving template.
func DecodeManyJSON(buf []byte, v ...interface{}) error {
	return json.Unmarshal(buf, &v)
}
