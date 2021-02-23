package golang

import "testing"

func TestGoVersions(t *testing.T) {
	t.Log(GoVersionsList())
}
