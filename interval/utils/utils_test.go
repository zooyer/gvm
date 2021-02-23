package utils

import (
	"testing"
)

func TestT(t *testing.T) {
	t.Log(GetAbsEnv("CBDDD"))
	t.Log(GetAbsEnv("JAVA_HOME"))
	t.Log(SetAbsEnv("ABCC", "DDDD"))
	return
	out, err := Command("cmd", "echo", `"abc"`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}
