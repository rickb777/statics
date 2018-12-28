package static

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var testDirFile *DirFile

func getGOPATH() string {
	gopath := os.Getenv("GOPATH")

	//if len(gopath) == 0 {
	//	panic("$GOPATH environment is not setup correctly, ending transmission!")
	//}

	return gopath
}

func TestMain(m *testing.M) {

	// setup
	testDirFile = &DirFile{
		Path:    "/static/test-files/teststart",
		Name:    "teststart",
		Size:    170,
		Mode:    os.FileMode(2147484141),
		ModTime: 1446650128,
		IsDir:   true,
		Compressed: `
`,
		Files: []*DirFile{{
			Path:    "/static/test-files/teststart/symlinkeddir",
			Name:    "symlinkeddir",
			Size:    136,
			Mode:    os.FileMode(2147484141),
			ModTime: 1446648191,
			IsDir:   true,
			Compressed: `
`,
			Files: []*DirFile{{
				Path:    "/static/test-files/teststart/symlinkeddir/realdir",
				Name:    "realdir",
				Size:    136,
				Mode:    os.FileMode(2147484141),
				ModTime: 1446648584,
				IsDir:   true,
				Compressed: `
`,
				Files: []*DirFile{{
					Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir",
					Name:    "doublesymlinkeddir",
					Size:    136,
					Mode:    os.FileMode(2147484141),
					ModTime: 1447163695,
					IsDir:   true,
					Compressed: `
`,
					Files: []*DirFile{{
						Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/doublesymlinkedfile.txt",
						Name:    "doublesymlinkedfile.txt",
						Size:    5,
						Mode:    os.FileMode(420),
						ModTime: 1446648265,
						IsDir:   false,
						Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
						Files: []*DirFile{},
					},
						{
							Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/triplesymlinkeddir",
							Name:    "triplesymlinkeddir",
							Size:    102,
							Mode:    os.FileMode(2147484141),
							ModTime: 1447163709,
							IsDir:   true,
							Compressed: `
`,
							Files: []*DirFile{{
								Path:    "/static/test-files/teststart/symlinkeddir/realdir/doublesymlinkeddir/triplesymlinkeddir/triplefile.txt",
								Name:    "triplefile.txt",
								Size:    5,
								Mode:    os.FileMode(420),
								ModTime: 1447163511,
								IsDir:   false,
								Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
								Files: []*DirFile{},
							},
							},
						},
					},
				},
					{
						Path:    "/static/test-files/teststart/symlinkeddir/realdir/realdirfile.txt",
						Name:    "realdirfile.txt",
						Size:    5,
						Mode:    os.FileMode(420),
						ModTime: 1446648207,
						IsDir:   false,
						Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
						Files: []*DirFile{},
					},
				},
			},
				{
					Path:    "/static/test-files/teststart/symlinkeddir/symlinkeddirfile.txt",
					Name:    "symlinkeddirfile.txt",
					Size:    5,
					Mode:    os.FileMode(420),
					ModTime: 1446647769,
					IsDir:   false,
					Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
					Files: []*DirFile{},
				},
			},
		},
			{
				Path:    "/static/test-files/teststart/plainfile.txt",
				Name:    "plainfile.txt",
				Size:    10,
				Mode:    os.FileMode(420),
				ModTime: 1446650128,
				IsDir:   false,
				Compressed: `
H4sIAAAJbogA/ypIzMnMS0ksSeQCBAAA//9+mKzVCgAAAA==
`,
				Files: []*DirFile{},
			},
			{
				Path:    "/static/test-files/teststart/symlinkedfile.txt",
				Name:    "symlinkedfile.txt",
				Size:    5,
				Mode:    os.FileMode(420),
				ModTime: 1446647746,
				IsDir:   false,
				Compressed: `
H4sIAAAJbogA/0pJLEnkAgQAAP//gsXB5gUAAAA=
`,
				Files: []*DirFile{},
			},
		},
	}

	os.Exit(m.Run())

	// teardown
}

func TestStaticNew(t *testing.T) {
	g := NewGomegaWithT(t)

	config := &Config{
		UseStaticFiles: true,
		AbsPkgPath:     "$PWD/foo/bar",
	}

	staticFiles, err := New(config, testDirFile)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(staticFiles).NotTo(BeNil())

	go func(sf *Files) {

		http.Handle("/static/", http.StripPrefix("/", http.FileServer(sf.FS())))
		http.ListenAndServe("127.0.0.1:13006", nil)

	}(staticFiles)

	time.Sleep(5000)

	f, err := staticFiles.GetHTTPFile("/static/test-files/teststart/plainfile.txt")
	g.Expect(err).NotTo(HaveOccurred())

	fis, err := f.Readdir(-1)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(Equal("not a directory"))

	fi, err := f.Stat()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fi.Name()).To(Equal("plainfile.txt"))
	g.Expect(fi.Size()).To(Equal(int64(10)))
	g.Expect(fi.IsDir()).To(Equal(false))
	g.Expect(fi.Mode()).To(Equal(os.FileMode(420)))
	g.Expect(fi.ModTime()).To(Equal(time.Unix(1446650128, 0)))
	g.Expect(fi.Sys()).NotTo(BeNil())

	err = f.Close()
	g.Expect(err).NotTo(HaveOccurred())

	f, err = staticFiles.GetHTTPFile("/static/test-files/teststart")
	g.Expect(err).NotTo(HaveOccurred())

	fi, err = f.Stat()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fi.Name()).To(Equal("teststart"))
	g.Expect(fi.Size()).To(Equal(int64(170)))
	g.Expect(fi.IsDir()).To(Equal(true))
	g.Expect(fi.Mode()).To(Equal(os.FileMode(2147484141)))
	g.Expect(fi.ModTime()).To(Equal(time.Unix(1446650128, 0)))
	g.Expect(fi.Sys()).NotTo(BeNil())

	var j int

	for err != io.EOF {
		fis, err = f.Readdir(2)
		g.Expect(fis).To(HaveLen(2 - j))
		j++
	}

	err = f.Close()
	g.Expect(err).NotTo(HaveOccurred())

	b, err := staticFiles.ReadFile("/static/test-files/teststart/plainfile.txt")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(b)).To(Equal("palindata\n"))

	b, err = staticFiles.ReadFile("nonexistantfile")
	g.Expect(err).To(HaveOccurred())

	bs, err := staticFiles.ReadFiles("/static/test-files/teststart", false)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(bs)).To(Equal(2))
	g.Expect(string(bs["/static/test-files/teststart/plainfile.txt"])).To(Equal("palindata\n"))
	g.Expect(string(bs["/static/test-files/teststart/symlinkedfile.txt"])).To(Equal("data\n"))

	bs, err = staticFiles.ReadFiles("/static/test-files/teststart", true)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(bs)).To(Equal(6))

	bs, err = staticFiles.ReadFiles("nonexistantdir", false)
	g.Expect(err).To(HaveOccurred())

	fis, err = staticFiles.ReadDir("/static/test-files/teststart")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(fis)).To(Equal(3))
	g.Expect(fis[0].Name()).To(Equal("plainfile.txt"))
	g.Expect(fis[1].Name()).To(Equal("symlinkeddir"))
	g.Expect(fis[2].Name()).To(Equal("symlinkedfile.txt"))

	fis, err = staticFiles.ReadDir("nonexistantdir")
	g.Expect(err).To(HaveOccurred())

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:13006/static/test-files/teststart/plainfile.txt", nil)
	g.Expect(err).NotTo(HaveOccurred())

	resp, err := client.Do(req)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

	bytes, err := ioutil.ReadAll(resp.Body)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(bytes)).To(Equal("palindata\n"))

	defer resp.Body.Close()
}

func TestLocalNew(t *testing.T) {
	g := NewGomegaWithT(t)

	config := &Config{
		UseStaticFiles: false,
		AbsPkgPath:     "$PWD/..",
	}

	staticFiles, err := New(config, testDirFile)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(staticFiles).NotTo(BeNil())

	go func(sf *Files) {

		http.Handle("/static/test-files/", http.StripPrefix("/", http.FileServer(sf.FS())))
		http.ListenAndServe("127.0.0.1:13007", nil)

	}(staticFiles)

	time.Sleep(5000)

	f, err := staticFiles.GetHTTPFile("/static/test-files/teststart/plainfile.txt")
	g.Expect(err).NotTo(HaveOccurred())

	fis, err := f.Readdir(-1)
	g.Expect(err).To(HaveOccurred())
	// g.Expect(err.Error()).To(Equal("readdirent: invalid argument") // this is "not a directory" in linux but readdirent: invalid argument in osx

	fi, err := f.Stat()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fi.Name()).To(Equal("plainfile.txt"))
	// g.Expect(fi.Size()).To(Equal(int64(10))  // commented out as size can differ on different file systems when cloned
	g.Expect(fi.IsDir()).To(Equal(false))
	// g.Expect(fi.Mode()).To(Equal(os.FileMode(420)) // commented out as permissions can be different based on when & where you cloned
	// g.Expect(fi.ModTime()).To(Equal(time.Unix(1446650128, 0)) // commented out as file mod times will be different based on when you cloned
	g.Expect(fi.Sys()).NotTo(BeNil())

	err = f.Close()
	g.Expect(err).NotTo(HaveOccurred())

	f, err = staticFiles.GetHTTPFile("/static/test-files/teststart")
	g.Expect(err).NotTo(HaveOccurred())

	fi, err = f.Stat()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fi.Name()).To(Equal("teststart"))
	// g.Expect(fi.Size()).To(Equal(int64(170)) // commented out as size can differ on different file systems when cloned
	g.Expect(fi.IsDir()).To(Equal(true))
	// g.Expect(fi.Mode()).To(Equal(os.FileMode(2147484141))  // commented out as permissions can be different based on when & where you cloned
	// g.Expect(fi.ModTime()).To(Equal(time.Unix(1446650128, 0))  // commented out as file mod times will be different based on when you cloned
	g.Expect(fi.Sys()).NotTo(BeNil())

	var j int

	for err != io.EOF {
		fis, err = f.Readdir(2)
		g.Expect(fis).To(HaveLen(2 - j))
		j++
	}

	err = f.Close()
	g.Expect(err).NotTo(HaveOccurred())

	b, err := staticFiles.ReadFile("/static/test-files/teststart/plainfile.txt")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(b)).To(Equal("palindata\n"))

	b, err = staticFiles.ReadFile("nonexistantfile")
	g.Expect(err).To(HaveOccurred())

	bs, err := staticFiles.ReadFiles("/static/test-files/teststart", false)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(bs)).To(Equal(2))
	g.Expect(string(bs["/static/test-files/teststart/plainfile.txt"])).To(Equal("palindata\n"))
	g.Expect(string(bs["/static/test-files/teststart/symlinkedfile.txt"])).To(Equal("data\n"))

	bs, err = staticFiles.ReadFiles("/static/test-files/teststart", true)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(bs)).To(Equal(6))

	bs, err = staticFiles.ReadFiles("nonexistantdir", false)
	g.Expect(err).To(HaveOccurred())

	fis, err = staticFiles.ReadDir("/static/test-files/teststart")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(fis)).To(Equal(3))
	g.Expect(fis[0].Name()).To(Equal("plainfile.txt"))
	g.Expect(fis[1].Name()).To(Equal("symlinkeddir"))
	g.Expect(fis[2].Name()).To(Equal("symlinkedfile.txt"))

	fis, err = staticFiles.ReadDir("nonexistantdir")
	g.Expect(err).To(HaveOccurred())

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:13007/static/test-files/teststart/plainfile.txt", nil)
	g.Expect(err).NotTo(HaveOccurred())

	resp, err := client.Do(req)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

	bytes, err := ioutil.ReadAll(resp.Body)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(bytes)).To(Equal("palindata\n"))

	defer resp.Body.Close()
}

func TestBadLocalAbsPath(t *testing.T) {
	g := NewGomegaWithT(t)

	config := &Config{
		UseStaticFiles: false,
		AbsPkgPath:     "../github.com/rickb777/statics",
	}

	staticFiles, err := New(config, testDirFile)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(Equal("AbsPkgPath is required when not using static files otherwise the static package has no idea where to grab local files from when your package is used from within another package."))
	g.Expect(staticFiles).To(BeNil())
}
