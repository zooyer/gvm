package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func Download(url, filename string) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	bar := pb.Default.Start64(res.ContentLength)
	defer bar.Finish()

	writer := bar.NewProxyWriter(file)
	if _, err = io.Copy(writer, res.Body); err != nil {
		return
	}

	return
}

func Untargz(filename, dir string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	bar := pb.Default.Start64(stat.Size())
	defer bar.Finish()

	gz, err := gzip.NewReader(bar.NewProxyReader(file))
	if err != nil {
		return
	}
	defer gz.Close()

	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if header.Typeflag != tar.TypeDir {
			filename := filepath.Join(dir, header.Name)
			if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}

			if _, err = io.Copy(file, reader); err != nil {
				file.Close()
				return err
			}

			file.Close()
		}
	}

	return nil
}

func Unzip(filename, dir string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	bar := pb.Default.Start64(stat.Size())
	defer bar.Finish()

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return
	}

	var total int64
	for _, f := range reader.File {
		total += f.FileInfo().Size()
	}
	bar.SetTotal(total)

	for _, f := range reader.File {
		file, err := f.Open()
		if err != nil {
			return err
		}

		if !f.FileInfo().IsDir() {
			filename := filepath.Join(dir, f.Name)
			if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				return err
			}

			writer, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			if _, err = io.Copy(bar.NewProxyWriter(writer), file); err != nil {
				writer.Close()
				file.Close()
				return err
			}

			writer.Close()
		}

		file.Close()
	}

	return
}

func Command(name string, args ...string) (out string, err error) {
	cmd := exec.Command(name, args...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Run(name string, args ...string) {
	out, err := Command(name, args...)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}

func Goenv(key string) string {
	if out, err := Command("go", "env", key); err == nil {
		return strings.TrimRight(out, " \f\t\r\n")
	}
	return os.Getenv(key)
}

func GetAbsEnv(key string) (val string, err error) {
	return getAbsEnv(key)
}

func SetAbsEnv(key, val string) (err error) {
	return setAbsEnv(key, val)
}

func AddPath(path string, home string) (err error) {
	p, err := GetAbsEnv("PATH")
	if err != nil {
		return
	}
	var ps []string
	var sep string
	switch runtime.GOOS {
	case "windows":
		sep = ";"
		ps = strings.Split(p, ";")
	case "linux":
		fallthrough
	default:
		sep = ":"
		ps = strings.Split(p, ":")
	}

	var exists bool
	for i, p := range ps {
		if strings.HasPrefix(p, home) {
			if runtime.GOOS == "windows" {
				if strings.HasSuffix(p, "\\bin") {
					exists = true
					ps[i] = path
				}
			} else {
				if strings.HasSuffix(p, "/bin") {
					exists = true
					ps[i] = path
				}
			}
		}
		if p == path {
			return nil
		}
	}

	if !exists {
		ps = append([]string{path}, ps...)
	}

	return SetAbsEnv("PATH", strings.Join(ps, sep))
}
