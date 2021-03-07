package golang

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zooyer/gvm/interval/conf"
	"github.com/zooyer/gvm/interval/debug"
	"github.com/zooyer/gvm/interval/utils"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strings"
)

type Version struct {
	Name    string
	Version string
	URL     string
	Kind    string
	OS      string
	Arch    string
	Size    string
	Sha256  string
}

type Filter func(version Version) bool

var os = map[string]string{
	"macOS":   "darwin",
	"Linux":   "linux",
	"Windows": "windows",
	"FreeBSD": "freebsd",
}

var arch = map[string]string{
	"ARMv6":   "arm",
	"ARMv8":   "arm64",
	"ppc64le": "ppc64le",
	"s390x":   "s390x",
	"x86":     "386",
	"x86-64":  "amd64",
}

var suffix = map[string]string{
	"darwin":  "tar.gz",
	"linux":   "tar.gz",
	"windows": "zip",
}

var home = map[string]string{
	"darwin":  "/usr/local/go",
	"linux":   "/usr/local/go",
	"windows": "C:\\Program Files\\go",
}

var decoder = map[string]func(filename, dir string) error{
	"darwin":  utils.Untargz,
	"linux":   utils.Untargz,
	"windows": utils.Unzip,
}

var defaultFilter = func(version Version) bool {
	if version.Kind != "Archive" {
		return false
	}
	if os[version.OS] != runtime.GOOS {
		return false
	}
	if arch[version.Arch] != runtime.GOARCH {
		return false
	}
	return true
}

var defaultVersions = []string{
	"go1.3",
	"go1.4",
	"go1.5",
	"go1.6",
	"go1.7",
	"go1.8",
	"go1.9",
	"go1.10",
	"go1.11",
	"go1.12",
	"go1.13",
	"go1.14",
	"go1.15",
	"go1.16",
	"go1.17",
	"go1.18",
	"go1.19",
	"go1.20",
}

var client = http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, conf.Timeout.Duration()/3)
		},
	},
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       conf.Timeout.Duration() * 2 / 3,
}

func getHTML() (html []byte, err error) {
	res, err := client.Get("https://golang.org/dl")
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	return ioutil.ReadAll(res.Body)
}

func parse(html []byte) (versions []Version, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		debug.Println("golang: parse error:", err.Error())
		return
	}

	doc.Find("div.expanded>table>tbody>tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() != 6 {
			return
		}

		node := td.Find("a")
		uri, _ := node.Attr("href")
		name := node.Text()
		field := strings.Split(name, ".")

		// 过滤掉1.9以前的版本
		if len(field[1]) == 1 && field[1][0] < '9' {
			//return
		}

		version := strings.Join(field[:3], ".")
		if c := field[2][0]; c < '0' || c > '9' {
			if index := strings.Index(field[1], "-"); index >= 0 {
				field[1] = field[1][:index]
			}
			version = strings.Join(field[:2], ".")
		}

		kind := td.Eq(1).Text()
		os := td.Eq(2).Text()
		arch := td.Eq(3).Text()
		size := td.Eq(4).Text()
		sha256 := td.Eq(5).Text()

		versions = append(versions, Version{
			Name:    name,
			Version: version,
			URL:     uri,
			Kind:    kind,
			OS:      os,
			Arch:    arch,
			Size:    size,
			Sha256:  sha256,
		})
	})

	return
}

func GoVersions(filter ...Filter) (versions []Version) {
	defer func() {
		var vs = make([]Version, 0, len(versions))
		for _, v := range versions {
			var ok = true
			for _, filter := range filter {
				ok = filter(v)
			}
			if ok {
				vs = append(vs, v)
			}
		}
		versions = vs
	}()

	html, err := getHTML()
	if err != nil {
		debug.Println("golang: get html error:", err.Error())
		return
	}

	versions, err = parse(html)
	if err != nil {
		debug.Println("golang: parse error:", err.Error())
		return
	}

	return versions
}

func GoVersionsList() []string {
	var m = make(map[string]bool)
	var versions []string
	for _, v := range GoVersions(defaultFilter) {
		if !m[v.Version] {
			m[v.Version] = true
			versions = append(versions, v.Version)
		}
	}

	if len(versions) == 0 {
		versions = defaultVersions
	}

	sort.Slice(versions, func(i, j int) bool {
		a := strings.Split(versions[i], ".")
		b := strings.Split(versions[j], ".")
		var length = len(a)
		if len(a) > len(b) {
			length = len(b)
		}
		for i := 0; i < length; i++ {
			if len(a[i]) != len(b[i]) {
				return len(a[i]) < len(b[i])
			}
			if a[i] != b[i] {
				return a[i] < b[i]
			}
		}
		return len(a) < len(b)
	})

	return versions
}

func Suffix() string {
	return suffix[runtime.GOOS]
}

func Filename(version string) string {
	return fmt.Sprintf("%s.%s-%s.%s", version, runtime.GOOS, runtime.GOARCH, suffix[runtime.GOOS])
}

func Decode(filename, dir string) error {
	if decode, exists := decoder[runtime.GOOS]; exists {
		return decode(filename, dir)
	}
	return errors.New("not support os: " + runtime.GOOS)
}

func DefaultGoHome() string {
	return home[runtime.GOOS]
}
