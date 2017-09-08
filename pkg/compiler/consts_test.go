// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package compiler

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/syzkaller/pkg/ast"
)

func TestExtractConsts(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", "consts.txt"))
	if err != nil {
		t.Fatalf("failed to read input file: %v", err)
	}
	desc := ast.Parse(data, "test", nil)
	if desc == nil {
		t.Fatalf("failed to parse input")
	}
	info := ExtractConsts(desc, func(pos ast.Pos, msg string) {
		t.Fatalf("%v: %v", pos, msg)
	})
	wantConsts := []string{"CONST1", "CONST10", "CONST11", "CONST12", "CONST13",
		"CONST14", "CONST15", "CONST16",
		"CONST2", "CONST3", "CONST4", "CONST5",
		"CONST6", "CONST7", "CONST8", "CONST9", "__NR_bar", "__NR_foo"}
	if !reflect.DeepEqual(info.Consts, wantConsts) {
		t.Fatalf("got consts:\n%q\nwant:\n%q", info.Consts, wantConsts)
	}
	wantIncludes := []string{"foo/bar.h", "bar/foo.h"}
	if !reflect.DeepEqual(info.Includes, wantIncludes) {
		t.Fatalf("got includes:\n%q\nwant:\n%q", info.Includes, wantIncludes)
	}
	wantIncdirs := []string{"/foo", "/bar"}
	if !reflect.DeepEqual(info.Incdirs, wantIncdirs) {
		t.Fatalf("got incdirs:\n%q\nwant:\n%q", info.Incdirs, wantIncdirs)
	}
	wantDefines := map[string]string{
		"CONST1": "1",
		"CONST2": "FOOBAR + 1",
	}
	if !reflect.DeepEqual(info.Defines, wantDefines) {
		t.Fatalf("got defines:\n%q\nwant:\n%q", info.Defines, wantDefines)
	}
}

func TestConstErrors(t *testing.T) {
	name := "consts_errors.txt"
	em := ast.NewErrorMatcher(t, filepath.Join("testdata", name))
	desc := ast.Parse(em.Data, name, em.ErrorHandler)
	if desc == nil {
		em.DumpErrors(t)
		t.Fatalf("parsing failed")
	}
	ExtractConsts(desc, em.ErrorHandler)
	em.Check(t)
}