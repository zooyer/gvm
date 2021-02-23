package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var v = viper.New()

var Debug = false

var Timeout = 20 * time.Second

var addr = map[string]interface{}{
	"debug":   &Debug,
	"timeout": &Timeout,
}

func init() {
	var err error
	v.SetEnvPrefix("GVM")

	for key, addr := range addr {
		if err = v.BindEnv(key); err != nil {
			fmt.Println("gvm bind env key", key, "error:", err.Error())
		}
		if err = v.UnmarshalKey(key, addr); err != nil {
			fmt.Println("gvm conf unmarshal key", key, "error:", err.Error())
		}
	}
}
