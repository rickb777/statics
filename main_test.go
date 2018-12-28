package main

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNonExistentStaticDir(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/garbagedir"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	ignore := ""
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	g.Expect(func() { main() }).To(Panic(), "stat static/test-files/garbagedir: no such file or directory")
}

func TestBadPackage(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := ""
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	init := false
	flagInit = &init

	g.Expect(main).To(Panic(), "**invalid Package Name")
}

func TestBadOutputDir(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/test-inner"
	flagStaticDir = &i

	o := ""
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	init := false
	flagInit = &init

	g.Expect(main).To(Panic(), "**invalid Output Directory")
}

func TestBadStaticDir(t *testing.T) {
	g := NewGomegaWithT(t)

	i := ""
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	init := false
	flagInit = &init

	g.Expect(main).To(Panic(), "**invalid Static File Directoy '.'")
}

func TestGenerateInitFile(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	init := true
	flagInit = &init

	main()

	b, err := ioutil.ReadFile("static/test-files/test.go")
	g.Expect(err, nil)

	expected := `//go:generate statics -i=static/test-files/teststart -o=static/test-files/test.go -pkg=test -group=Assets

package test

import "github.com/rickb777/statics/static"

// newStaticAssets initializes a new *static.Files instance for use
func newStaticAssets(config *static.Config) (*static.Files, error) {

	return static.New(config, &static.DirFile{})
}
`

	g.Expect(string(b), expected)
}

func TestIgnore(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	ignore := ".*.txt"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	g.Expect(main).NotTo(Panic())
}

func TestBadIgnore(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	ignore := "([12.gitignore"
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	g.Expect(main, "**Error Compiling Regex:error parsing regexp: missing closing ]: `[12.gitignore`")
}

func TestGenerateFilePrefix(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	ignore := ""
	flagIgnore = &ignore

	prefix := "static/"
	flagPrefix = &prefix

	init := false
	flagInit = &init

	g.Expect(main).NotTo(Panic())
}

func TestGenerateFile(t *testing.T) {
	g := NewGomegaWithT(t)

	i := "static/test-files/teststart"
	flagStaticDir = &i

	o := "static/test-files/test.go"
	flagOuputFile = &o

	p := "test"
	flagPkg = &p

	gr := "Assets"
	flagGroup = &gr

	ignore := ""
	flagIgnore = &ignore

	prefix := ""
	flagPrefix = &prefix

	init := false
	flagInit = &init

	g.Expect(main).NotTo(Panic())
}
