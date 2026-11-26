package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGroupedParametersAndNestedControls(t *testing.T) {
	cases := map[string]bool{
		"func ok(a,b,c int) {}": true,
		"func bad(a,b,c,d int) {}": false,
		"func bad() { if true { for { if true { for {} } } } }": false,
		"func bad() { _ = func(a,b,c,d int) {} }": false,
	}
	for code, want := range cases {
		path := filepath.Join(t.TempDir(), "sample.go")
		if err := os.WriteFile(path, []byte("package sample\n"+code), 0600); err != nil {
			t.Fatal(err)
		}
		if got := check(path); got != want {
			t.Errorf("check(%q) = %v", code, got)
		}
	}
}
