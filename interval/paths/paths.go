package paths

import (
	"github.com/zooyer/gvm/interval/debug"
	"os"
	"path/filepath"
	"runtime"
)

func Home(path ...string) string {
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(append([]string{home}, path...)...)
	}
	debug.Println("paths: user home dir:", err.Error())
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(append([]string{"C:\\Users\\Administrator"}, path...)...)
	}
	return filepath.Join(append([]string{"/root"}, path...)...)

}

func AbsThisFile() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

func AbsThisDir() string {
	return filepath.Dir(AbsThisFile())
}

func GvmRunCom() string {
	return Home(".gvmrc")
}
