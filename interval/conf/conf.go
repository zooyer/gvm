package conf

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type unmarshaler interface {
	UnmarshalEnv(data []byte) error
}

type Duration time.Duration

var Debug = false

var Timeout = 20 * Duration(time.Second)

var addr = map[string]interface{}{
	"debug":   &Debug,
	"timeout": &Timeout,
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d *Duration) UnmarshalEnv(data []byte) (err error) {
	var index int
	for index = range data {
		if data[index] < '0' || data[index] > '9' {
			break
		}
	}

	var i time.Duration
	if err = json.Unmarshal(data[:index], &i); err != nil {
		return
	}

	var dd = (*time.Duration)(d)

	var unit = string(data[index:])
	switch unit {
	case "ms":
		*dd = i * time.Millisecond
	case "s":
		*dd = i * time.Second
	case "min":
		*dd = i * time.Minute
	case "h":
		*dd = i * time.Hour
	}

	return
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

func BindEnv(key string, v interface{}) (err error) {
	if val := os.Getenv(key); val != "" {
		switch value := v.(type) {
		case unmarshaler:
			return value.UnmarshalEnv([]byte(val))
		case *bool, *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64:
			return json.Unmarshal([]byte(val), value)
		case *string:
			*value = val
		default:
			return json.Unmarshal([]byte(val), value)
		}
	}

	return
}

func init() {
	for key, addr := range addr {
		BindEnv("GVM_"+strings.ToUpper(key), addr)
	}
}
