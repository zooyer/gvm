package main

import (
	"fmt"
	"github.com/zooyer/gvm/interval/debug"
	"github.com/zooyer/gvm/interval/golang"
	"github.com/zooyer/gvm/interval/paths"
	"github.com/zooyer/gvm/interval/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var helps = `Usage: gvm [command] [args]

Description:
	gvm is the go version manager

Commands:
	set       - set go version
	use       - use go version
	info      - show the go info
	list      - list all go versions
	help      - show the help manual
	install   - install go versions
	uninstall - uninstall go versions`

var usage = func(command string) string {
	return helps
}

var usages = map[string]func() string{
	"set": func() string {
		return fmt.Sprintf("show: %s set go1.9.2", os.Args[0])
	},
	"use": func() string {
		return fmt.Sprintf("show: %s use go1.9.2", os.Args[0])
	},
	"info": func() string {
		return fmt.Sprintf("show: %s info", os.Args[0])
	},
	"install": func() string {
		return fmt.Sprintf("show: %s install go1.9.2", os.Args[0])
	},
	"uninstall": func() string {
		return fmt.Sprintf("show: %s uninstall go1.9.2", os.Args[0])
	},
	"list": func() string {
		var buf strings.Builder
		buf.WriteString("> \033[1;32mcurrented\033[0m\n")
		buf.WriteString("+ \033[1;36minstalled\033[0m\n")
		buf.WriteString("- \033[1;37muninstalled\033[0m")
		return buf.String()
	},
	"help": func() string {
		if len(os.Args) > 1 {
			return usage(os.Args[2])
		}
		return helps
	},
}

var command string

var config struct {
	GoHome string `yaml:"GOHOME" json:"GOHOME"`
	GoRoot string `yaml:"GOROOT" json:"GOROOT"`
	GoPath string `yaml:"GOPATH" json:"GOPATH"`
}

func init() {
	usage = func(command string) string {
		if fn, exists := usages[command]; exists {
			return fn()
		}
		return helps
	}
}

func show(command string) {
	if fn, exists := usages[command]; exists {
		fmt.Println(fn())
	} else {
		fmt.Println(helps)
	}
}

func initGvmRunCom() (err error) {
	var filename = paths.GvmRunCom()

	if _, err = os.Lstat(filename); err != nil {
		if !os.IsNotExist(err) {
			return
		}

		if err = utils.SetAbsEnv("PATH", fmt.Sprintf(`%s:$PATH`, paths.AbsThisDir())); err != nil {
			return
		}
		if err = utils.SetAbsEnv("GOHOME", golang.DefaultGoHome()); err != nil {
			return
		}
	}

	return nil
}

func initShellRunCom() (err error) {
	for _, filename := range []string{".bashrc", ".zshrc"} {
		filename = paths.Home(filename)

		if _, err = os.Lstat(filename); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return
		}

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var env = fmt.Sprintf("source %s", paths.GvmRunCom())
		var exists bool
		for _, line := range strings.Split(string(data), "\n") {
			if line = strings.TrimRight(line, "\r"); strings.HasPrefix(line, env) {
				exists = true
				break
			}
		}

		if !exists {
			file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return err
			}

			if _, err = file.WriteString(env); err != nil {
				return err
			}

			file.Close()
		}
	}

	return nil
}

func initLinuxMacEnviron() (err error) {
	if err = initGvmRunCom(); err != nil {
		return
	}
	return initShellRunCom()
}

func initWindowsEnviron() (err error) {
	if dir := paths.AbsThisDir(); !strings.HasPrefix(dir, os.TempDir()) {
		val, err := utils.GetAbsEnv("PATH")
		if err == nil && !strings.Contains(val, dir) {
			val = dir + ";" + val
			if err = utils.SetAbsEnv("PATH", val); err != nil {
				return err
			}
		}
	}

	var val string
	if val, err = utils.GetAbsEnv("GOHOME"); err != nil {
		return
	}

	if val == "" {
		if err = utils.SetAbsEnv("GOHOME", golang.DefaultGoHome()); err != nil {
			return
		}
	}

	return
}

func initEnviron() (err error) {
	if runtime.GOOS == "windows" {
		return initWindowsEnviron()
	}
	return initLinuxMacEnviron()
}

func showEnv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GOHOME: %s\n", config.GoHome))
	sb.WriteString(fmt.Sprintf("GOROOT: %s\n", config.GoRoot))
	sb.WriteString(fmt.Sprintf("GOPATH: %s", config.GoPath))
	return sb.String()
}

func initConfig() {
	if config.GoHome, _ = utils.GetAbsEnv("GOHOME"); config.GoHome == "" {
		if config.GoHome = os.Getenv("GOHOME"); config.GoHome == "" {
			config.GoHome = golang.DefaultGoHome()
		}
	}
	if config.GoPath = utils.Goenv("GOPATH"); config.GoPath == "" {
		if config.GoPath = os.Getenv("GOPATH"); config.GoPath == "" {
			config.GoPath = paths.Home("go")
		}
	}
	if config.GoRoot, _ = utils.GetAbsEnv("GOROOT"); config.GoRoot == "" {
		if config.GoRoot = utils.Goenv("GOROOT"); config.GoRoot == "" {
			if config.GoRoot = os.Getenv("GOROOT"); config.GoRoot == "" {
				config.GoRoot = runtime.GOROOT()
			}
		}
	}
}

func init() {
	var err error

	// init config
	initConfig()

	// init environ
	if err = initEnviron(); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println(helps)
		os.Exit(0)
	}

	command = os.Args[1]
	debug.Println(showEnv())
}

func exists(version string) bool {
	binary := filepath.Join(config.GoHome, version, "bin", "go")
	out, err := utils.Command(binary, "version")
	if err != nil {
		return false
	}
	return strings.HasPrefix(out, fmt.Sprintf("go version %s %s/%s", version, runtime.GOOS, runtime.GOARCH))
}

func args(index int) string {
	if len(os.Args) > index+2 {
		return os.Args[index+2]
	}
	show(os.Args[1])
	os.Exit(1)
	return ""
}

func set() {
	var version = args(0)

	if !exists(version) {
		fmt.Println(version, "not found, will be install")
		install()
		if !exists(version) {
			fmt.Println(version, "install failed")
			return
		}
	}

	filename := filepath.Join(config.GoHome, version)
	if err := utils.SetAbsEnv("GOROOT", filename); err != nil {
		panic(err)
	}

	if err := utils.AddPath(filepath.Join(filename, "bin"), config.GoHome); err != nil {
		panic(err)
	}

	fmt.Println("GOHOME:", config.GoHome)
	fmt.Println("GOROOT:", filename)
	fmt.Println("GOPATH:", config.GoPath)
}

func use() {
	var version = args(0)

	if !exists(version) {
		fmt.Println(version, "not found, will be install")
		install()
		if !exists(version) {
			fmt.Println(version, "install failed")
			return
		}
	}

	filename := filepath.Join(config.GoHome, version)
	// TODO 设置环境变量

	source := fmt.Sprintf("export GOROOT=\"%s\"\n", filename)
	if err := ioutil.WriteFile(filepath.Join(config.GoHome, "source.sh"), []byte(source), 0755); err != nil {
		panic(err)
	}

	fmt.Println("GOHOME:", config.GoHome)
	fmt.Println("GOROOT:", filename)
	fmt.Println("GOPATH:", config.GoPath)
}

func info() {
	if out, err := utils.Command("go", "version"); err == nil && out != "" {
		fmt.Println("GOOS:", runtime.GOOS)
		fmt.Println("GOARCH:", runtime.GOARCH)
		fmt.Println("GOVERSION:", runtime.Version())
	}

	fmt.Println("GOHOME:", config.GoHome)
	fmt.Println("GOROOT:", config.GoRoot)
	fmt.Println("GOPATH:", config.GoPath)
}

func list() {
	var version string
	if root, err := utils.GetAbsEnv("GOROOT"); err == nil {
		_, version = filepath.Split(strings.TrimSpace(root))
	}

	var buf strings.Builder

	for _, ver := range golang.GoVersionsList() {
		var line = ver
		if exists(ver) {
			if filepath.Clean(config.GoRoot) == filepath.Join(config.GoHome, ver) {
				line = fmt.Sprintf("> \033[1;32m%s\033[0m", ver)
			} else {
				line = fmt.Sprintf("+ \033[1;36m%s\033[0m", ver)
			}
		} else {
			line = fmt.Sprintf("- \033[1;37m%s\033[0m", ver)
		}
		if ver == version {
			line += "(system)"
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	fmt.Print(buf.String())
}

func install() {
	if len(os.Args) < 3 {
		show(command)
		os.Exit(1)
	}

	for _, version := range os.Args[2:] {
		if exists(version) {
			fmt.Println(version, "already installed")
			return
		}

		var dir = filepath.Join(config.GoHome, version)
		var filename = dir + "." + golang.Suffix()
		var url = fmt.Sprintf("https://dl.google.com/go/%s", golang.Filename(version))

		fmt.Println(version, "installing: ")
		if err := utils.Download(url, filename); err != nil {
			panic(err)
		}

		fmt.Println(version, "unpacking: ")
		if err := golang.Decode(filename, config.GoHome); err != nil {
			panic(err)
		}

		if err := os.Rename(filepath.Join(config.GoHome, "go"), dir); err != nil {
			panic(err)
		}

		fmt.Println(version, "installed")
	}
}

func uninstall() {
	if len(os.Args) < 3 {
		show(command)
		os.Exit(1)
	}

	for _, version := range os.Args[2:] {
		if exists(version) {
			if err := os.RemoveAll(filepath.Join(config.GoHome, version)); err != nil {
				panic(err)
			}
			_ = os.RemoveAll(filepath.Join(config.GoHome, version+".tar.gz"))
		}

		fmt.Println(version, "uninstalled")
	}
}

func help() {
	if fn, exists := usages[command]; exists {
		fmt.Println(fn())
	} else {
		fmt.Println(helps)
	}
	os.Exit(0)
}

func main() {
	switch command {
	case "set":
		set()
	case "use":
		use()
	case "info":
		info()
	case "list":
		list()
	case "help":
		help()
	case "install":
		install()
	case "uninstall":
		uninstall()
	default:
		fmt.Println(helps)
		os.Exit(1)
	}
}
