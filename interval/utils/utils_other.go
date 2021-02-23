//+build !windows

package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/zooyer/gvm/interval/paths"
)

// 设置永久环境变量
func GvmEnviron() (env [][2]string, err error) {
	data, err := ioutil.ReadFile(paths.GvmRunCom())
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line = strings.Trim(line, " \t\f\r"); line == "" || line[0] == '#' {
			continue
		}
		if !strings.HasPrefix(line, "export ") {
			continue
		}
		line = line[7:]
		if index := strings.Index(line, "="); index > 0 {
			key := strings.TrimSpace(line[:index])
			if len(key) > 1 && key[0] == '"' && key[len(key)-1] == '"' {
				if v, err := strconv.Unquote(key); err == nil {
					key = v
				}
			}
			val := strings.TrimSpace(line[index+1:])
			if len(val) > 1 && val[0] == '"' && val[len(val)-1] == '"' {
				if v, err := strconv.Unquote(val); err == nil {
					val = v
				}
			}
			env = append(env, [2]string{key, val})
		}
	}
	return
}

func GetGvmEnv(key string) (val string, err error) {
	env, err := GvmEnviron()
	if err != nil {
		return
	}

	for _, env := range env {
		if env[0] == key {
			return env[1], nil
		}
	}

	return
}

func SetGvmEnv(key, val string) (err error) {
	env, err := GvmEnviron()
	if err != nil && !os.IsNotExist(err) {
		return
	}

	var exists bool
	for i := range env {
		if env[i][0] == key {
			env[i][1] = val
			exists = true
		}
	}

	if !exists {
		env = append(env, [2]string{key, val})
	}

	var buf bytes.Buffer
	for _, env := range env {
		buf.WriteString(fmt.Sprintf("export %s=\"%s\"\n", env[0], env[1]))
	}

	if err = ioutil.WriteFile(paths.GvmRunCom(), buf.Bytes(), 0644); err != nil {
		return
	}

	return
}

func GetGvmEnvByShell(key string) (val string, err error) {
	out, err := Command("bash", "-c", "source ~/.gvmrc && echo $"+key)
	if err != nil {
		return
	}
	return strings.TrimRight(out, " \t\f\r\n"), nil
}

func getAbsEnv(key string) (val string, err error) {
	if val, err = GetGvmEnv(key); err == nil {
		return
	}
	return GetGvmEnvByShell(key)
}

func setAbsEnv(key, val string) (err error) {
	return SetGvmEnv(key, val)
}
