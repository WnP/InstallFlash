package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery" // parse xml made easy
)

var (
	// use to clean installed files if an error append
	installed []string
	// main configuration
	config = []byte(`[
  {
    "name": "glibc",
    "repo": "http://mirrors.kernel.org/archlinux/core/os/x86_64/",
    "rules": [
      {
        "src": "^usr/lib/ld-.+.so$",
        "dest-dir": "/usr/local/lib/",
	"file-mode": "0755"
      }
    ]
  },
  {
    "name": "flashplugin",
    "repo": "http://mirrors.kernel.org/archlinux/extra/os/x86_64/",
    "rules": [
      {
        "src": "usr/lib/mozilla/plugins/libflashplayer.so",
        "dest-dir": "/usr/lib/mozilla/plugins/",
	"file-mode": "0755"
      }
    ]
  }
]`)
)

func main() {
	var pkgs []pkg
	err := json.Unmarshal(config, &pkgs)
	check(err, fmt.Errorf("Can't load configuration: "))

	for _, p := range pkgs {
		p.install()
	}

	fmt.Printf("\nIf you're using grsecurity kernel do:\n\n    $ paxctl -c -m  /usr/lib/firefox-<version>/plugin-container\n\n")
}

type rule struct {
	Src      string      `json:"src"`      // path to file in archive to be installed
	DestDir  string      `json:"dest-dir"` // destination directory
	FileMode os.FileMode `json:"file-mode` // file's mode and permission bits
}

type pkg struct {
	Name  string `json:"name"`  // package name without version
	Url   string `json:"repo"`  // repository url, will be replaced by pakage latest download url
	Rules []rule `json:"rules"` // install rules
}

// Main routine to manage package installation
func (p *pkg) install() {

	log.Printf("installing %s required component(s)... Please wait\n", p.Name)

	err := p.getUrl() // get latest url
	check(err)

	file := p.download() // download package archive

	xz := xzReader(file) // extract archive

	p.installFiles(xz) // install from archive

	log.Printf("Installed %s required component(s) from %s", p.Name, p.Url) // informe user
}

// install files from tar archive with the provided rules
func (p *pkg) installFiles(r io.ReadCloser) {

	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		check(err, fmt.Errorf("tar extraction failed: "))

		for _, r := range p.Rules {
			if hdr.Name == r.Src || matchRe(hdr.Name, r.Src) {

				// get content
				var content []byte
				_, err = tr.Read(content)
				check(err)

				// get file name
				splits := strings.Split(hdr.Name, "/")
				fname := splits[len(splits)-1]
				fname = r.DestDir + fname

				// install file
				err = ioutil.WriteFile(fname, content, r.FileMode)
				check(err)

				// record filename for removing it in case of crash
				installed = append(installed, fname)
			}
		}
	}
}

// Get package's latest version url from repository
func (p *pkg) getUrl() error {
	var re = regexp.MustCompile(fmt.Sprintf("^%s-.+-x86_64\\.pkg\\.tar\\.xz$", p.Name))

	// get repository root tree
	doc, err := goquery.NewDocument(p.Url)
	check(err, fmt.Errorf("Can't get repository url: "))

	// rexep matching against `a` tags
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {

		href, exists := s.Attr("href")

		if exists { // don't want to perform rexep if not exists

			// don't match if more than one found
			if res := re.FindAllString(href, -1); len(res) == 1 {
				p.Url = p.Url + res[0]
				return true // break each loop
			}
		}
		return true
	})
	return nil
}

// Dowload package archive from repository
func (p *pkg) download() io.ReadCloser {
	rpipe, wpipe := io.Pipe() // don't use temporay file, pipe direct to other process

	go func() {
		// download from url
		resp, err := http.Get(p.Url)
		defer resp.Body.Close()

		check(err, fmt.Errorf("Can't download from url: %s", p.Url))

		// io.Copy reads 32kb (maximum) from input and writes them to output,
		// then repeats. don't worry about memory.
		// cf. http://golang.org/src/io/io.go?s=12247:12307#L340
		n, err := io.Copy(wpipe, resp.Body)
		check(err, fmt.Errorf("Fail to copy from url: %s", p.Url))

		wpipe.CloseWithError(err)

		log.Printf("archive %s: %d bytes downloaded", p.Name, n)
	}()
	return rpipe
}

// extract xz archive using xz executable
func xzReader(r io.Reader) io.ReadCloser {
	var stderr = bytes.NewBuffer([]byte{})
	rpipe, wpipe := io.Pipe() // don't use temporay file, pipe direct to other process

	cmd := exec.Command("xz", "-dc")
	cmd.Stdin = r
	cmd.Stdout = wpipe
	cmd.Stderr = stderr

	go func() {
		err := cmd.Run()
		check(err, fmt.Errorf("extraction failed: %s", stderr.String()))

		wpipe.CloseWithError(err)
	}()

	return rpipe
}

// check error, only the first one is checked other are here for comunication purpose
// use the second err do define your error message and the first one for the process error
func check(err ...interface{}) {
	if err[0] != nil {
		clean()
		if len(err) == 2 {
			// first message is the custom one
			log.Fatal(err[1], err[0])
		} else {
			// print all error(s) in provided order
			log.Fatal(err...)
		}
	}
}

func matchRe(fileName string, reString string) bool {
	var re = regexp.MustCompile(reString)
	if res := re.FindAllString(fileName, -1); len(res) > 0 {
		return true
	}
	return false
}

// TODO: Clean installed files
func clean() {
	var err error
	for _, v := range installed {
		err = os.Remove(v)
		if err != nil {
			log.Printf("Can't remove installed file: %s: %s", v, err)
		}
	}
}
