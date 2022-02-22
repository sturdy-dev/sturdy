package cloud

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type transport struct {
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

const dataLimit = 256 * 1024 // 256kb

func marshal(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to json: %w", err)
	}

	if len(jsonData) < dataLimit {
		return json.Marshal(transport{
			ContentType: "application/json",
			Data:        jsonData,
		})
	}

	b := &bytes.Buffer{}
	b64 := base64.NewEncoder(base64.StdEncoding, b)
	deflate, _ := zlib.NewWriterLevel(b64, zlib.BestCompression)
	_, _ = deflate.Write(jsonData)
	_ = deflate.Close()
	_ = b64.Close()

	return json.Marshal(transport{
		ContentType: "application/octet-stream",
		Data:        b.Bytes(),
	})
}

func unmarshal(data []byte, dist interface{}) error {
	var msg transport
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	dataReader := bytes.NewReader(msg.Data)
	switch msg.ContentType {
	case "application/json":
		return json.NewDecoder(dataReader).Decode(&dist)
	case "application/octet-stream":
		b64 := base64.NewDecoder(base64.StdEncoding, dataReader)
		deflate, _ := zlib.NewReader(b64)
		defer deflate.Close()
		return json.NewDecoder(deflate).Decode(&dist)
	default:
		return fmt.Errorf("unsupported content type: %s", msg.ContentType)
	}
}
