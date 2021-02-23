package temp

//const script = `#!/bin/sh
//export PATH="%s:%s:$PATH"
//export GOHOME="/usr/local/go"
//`
//
//var defaultEnviron = [][2]string{
//	{"GOHOME", `"/usr/local/go"`},
//	{"PATH", fmt.Sprintf(`"%s:$PATH"`, paths.AbsThisDir())},
//}

//func genShell(goroot string) (string, error) {
//	dir := paths.AbsThisDir()
//
//	return fmt.Sprintf(
//		script,
//		dir,
//		goroot+"/bin",
//	), nil
//}

//func GenShell(goroot string) (string, error) {
//	shell, err := genShell(goroot)
//	if err != nil {
//		return "", err
//	}
//
//	if goroot != "" {
//		shell += fmt.Sprintf("export GOROOT=\"%s\"\n", goroot)
//	}
//
//	return shell, nil
//}

//func GenShellFile(filename, goroot string) (err error) {
//	shell, err := GenShell(goroot)
//	if err != nil {
//		return
//	}
//
//	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
//	if err != nil {
//		return
//	}
//	defer file.Close()
//
//	if _, err = file.WriteString(shell); err != nil {
//		return
//	}
//
//	return
//}
