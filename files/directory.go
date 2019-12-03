package files

import (
	"os"
	"runtime"
	"strings"
)

//Directory desc
//@Struct Directory desc: Virtual Directory
//@Member (string) Base Directory
type Directory struct {
	rootPath string
	wildCard string
}

//Initial desc
//@Method Initial desc: initialization Directory
func (slf *Directory) Initial() {

	if runtime.GOOS == "windows" {
		slf.wildCard = "\\"
	} else {
		slf.wildCard = "/"
	}

	currPath, _ := os.Getwd()
	slf.WithRoot(currPath)
}

//WithRoot desc
//@Method WithRoot desc: Setting Root path
//@Param  (string) path
func (slf *Directory) WithRoot(path string) {
	slf.rootPath = path
	if strings.HasSuffix(slf.rootPath, slf.wildCard) {
		slf.rootPath = slf.rootPath[:len(slf.rootPath)-1]
	}
}

//GetFullPathName desc
//@Method GetFullPathName desc: Return Full path and file name
//@Return (string) Full path and file name
func (slf *Directory) GetFullPathName(filePath string) string {
	if strings.HasPrefix(filePath, slf.rootPath) {
		return filePath
	}

	return slf.rootPath + slf.wildCard + filePath
}
