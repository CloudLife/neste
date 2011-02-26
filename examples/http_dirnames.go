package main

import (
	"http"
	"log"
	"os"
	"bitbucket.org/fzzbt/neste"
	"path"
)

const (
	idir        = "/home/fuzzybyte"
	templateDir = "templates"
)

var tm *neste.Manager = neste.New(templateDir, nil)

type fileInfoRow struct {
	Name string // File name
	Size int64  // File size in bytes
}

type dBase struct {
	Title   string
	Content string
}

type dIndex struct {
	FileRows []fileInfoRow
}

// Returns up to count file names residing in a given directory.
// Returned file names will not include hidden files, if hidden is false.
func getDirnames(dir string, count int, hidden bool) (dirnames []string, err os.Error) {
	dirFile, err := os.Open(dir, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}

	defer dirFile.Close()

	rawDirnames, err := dirFile.Readdirnames(count)
	if err != nil {
		return nil, err
	}

	if hidden {
		// Return file names as they are, including the hidden files.
		return rawDirnames, nil
	}

	for _, v := range rawDirnames {
		// Ignore hidden files that start with a period (.)
		if v[0] != '.' {
			dirnames = append(dirnames, v)
		}
	}

	return dirnames, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	dirnames, err := getDirnames(idir, 1024, false)
	if err != nil {
		panic(err)
	}

	fileRows := make([]fileInfoRow, len(dirnames))
	for i, v := range dirnames {
		fileRows[i].Name = v

		fi, _ := os.Lstat(path.Join(idir, v))
		fileRows[i].Size = fi.Size
	}

	dBase := &dBase{
		Title: "Index of " + idir}

	dIndex := &dIndex{
		FileRows: fileRows}

	renderIndex(w, dBase, dIndex)
}

func renderIndex(w http.ResponseWriter, dBase *dBase, dIndex *dIndex) {
	dBase.Content, _ = tm.GetFile("index.html").Render(dIndex)
	tm.GetFile("base.html").Execute(w, dBase)
}

func initTemplates() {
	_, err := tm.AddFilenl("base.html")
	if err != nil {
		panic(err)
	}

	_, err = tm.AddFile("index.html")
	if err != nil {
		panic(err)
	}
}

func main() {
	initTemplates()
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}

}

