package temp

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func lookExtensions(path, dir string) (string, error) {
	if filepath.Base(path) == path {
		path = filepath.Join(".", path)
	}
	if dir == "" {
		return exec.LookPath(path)
	}
	if filepath.VolumeName(path) != "" {
		return exec.LookPath(path)
	}
	if len(path) > 1 && os.IsPathSeparator(path[0]) {
		return exec.LookPath(path)
	}
	dirandpath := filepath.Join(dir, path)
	// We assume that LookPath will only add file extension.
	lp, err := exec.LookPath(dirandpath)
	if err != nil {
		return "", err
	}
	ext := strings.TrimPrefix(lp, dirandpath)
	return path + ext, nil
}

//func bak() {
//	if filepath.Base(name) != name {
//		path, err := exec.LookPath(name)
//		if err != nil {
//			return "", err
//		}
//		name = path
//	}
//	if runtime.GOOS == "windows" {
//		if name, err = lookExtensions(name, ""); err != nil {
//			return
//		}
//	}
//
//	ir, iw, err := os.Pipe()
//	if err != nil {
//		return
//	}
//	defer ir.Close()
//	defer iw.Close()
//
//	or, ow, err := os.Pipe()
//	if err != nil {
//		return
//	}
//	defer or.Close()
//	defer ow.Close()
//
//	er, ew, err := os.Pipe()
//	if err != nil {
//		return
//	}
//	defer er.Close()
//	defer ew.Close()
//
//	// TODO
//	if _, err = iw.WriteString(strings.Join(append(args, ""), " ")); err != nil {
//		return
//	}
//	//if _, err = iw.WriteString("exit\r\n"); err != nil {
//	//	return
//	//}
//
//	var attr = os.ProcAttr{
//		Env:   []string{},
//		Files: []*os.File{ir, ow, ew},
//	}
//
//	process, err := os.StartProcess(name, []string{"/k"}, &attr)
//	if err != nil {
//		return
//	}
//	defer process.Release()
//
//	go func() {
//		defer process.Kill()
//
//		time.Sleep(time.Second * 5)
//		var buf = make([]byte, 4096*4)
//		n, e := or.Read(buf)
//		if e != nil {
//			err = e
//			return
//		}
//		n, e = or.Read(buf)
//		if e != nil {
//			err = e
//			return
//		}
//
//		out = string(buf[:n])
//		fmt.Println("out:", out)
//	}()
//
//	state, err := process.Wait()
//	if err != nil {
//		return
//	}
//
//	if !state.Success() {
//		return "", errors.New(state.String())
//	}
//
//	return
//}


//func Getenv(key string) string {
//	out, _ := cmd("bash", "-c", "echo $"+key)
//	return out
//}


//func SHELL() string {
//	return os.Getenv("SHELL")
//}

//func Setenv(key, value string) {
//	var exists bool
//	var shell = SHELL()
//
//	var filename string
//	switch shell {
//	case "/bin/zsh":
//		filename = "~/.zshrc"
//	case "/bin/bash":
//		filename = "~/.bashrc"
//	default:
//		filename = "~/.bashrc"
//	}
//
//	data, err := ioutil.ReadFile(filename)
//	if err == nil {
//		for _, line := range strings.Split(string(data), "\n") {
//			if exists = strings.HasPrefix(line, fmt.Sprintf("export %s=", key)); exists {
//				break
//			}
//		}
//	}
//
//	if !exists {
//		file, err := os.OpenFile("/etc/profile", os.O_APPEND, os.ModeAppend)
//		if err != nil {
//			panic(err)
//		}
//		defer file.Close()
//
//		if _, err = file.WriteString(fmt.Sprintf("\nexport %s=%s\n", key, value)); err != nil {
//			panic(err)
//		}
//	}
//
//}
