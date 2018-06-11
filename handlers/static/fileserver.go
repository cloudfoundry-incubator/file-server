package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type fileServer struct {
	root http.FileSystem
}

func NewFileServer(dir string) http.Handler {
	return &fileServer{
		root: http.Dir(dir),
	}
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if containsDotDot(upath) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	tgzPath := path.Clean(upath)

	file, fileStats := f.validateFile(tgzPath, w)
	if file == nil {
		return
	}
	defer file.Close()

	cksFile, _ := f.validateFile(fmt.Sprintf("%s.sha1", tgzPath), w)
	if cksFile == nil {
		return
	}
	defer cksFile.Close()

	cksBytes, err := ioutil.ReadAll(cksFile)
	if err != nil {
		http.Error(w, "Cannot read sha1 file", http.StatusInternalServerError)
		return
	}
	sha1sum := strings.TrimSpace(string(cksBytes))

	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, sha1sum))

	http.ServeContent(w, r, fileStats.Name(), fileStats.ModTime(), file)

}

// validateFile checks that a file can be found and is not a directory. It
// responds with an HTTP error and nil file
func (f *fileServer) validateFile(p string, w http.ResponseWriter) (ret http.File, stat os.FileInfo) {
	file, err := f.root.Open(p)
	if err != nil {
		http.Error(w, fmt.Sprintf("File not found: %s", filepath.Base(p)), http.StatusNotFound)
		return nil, nil
	}
	defer func() {
		if ret == nil {
			file.Close()
		}
	}()

	d, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot stat file: %s", filepath.Base(p)), http.StatusInternalServerError)
		return nil, nil
	}

	if d.IsDir() {
		http.Error(w, "Unauthorized to list the directory", http.StatusUnauthorized)
		return nil, nil
	}

	return file, d
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
