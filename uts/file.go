package uts

import (
	"io"
	"os"
)

func SaveFile(dstPath string, content string) {
	var err error
	f, err := os.Create(dstPath)
	if !ChkErr(err) {
		f.WriteString(content)
	}
	defer f.Close()
}

func CopyFile(dst string, src string, fileFunc ...func(f *os.File) string) {
	s, err := os.Stat(src)
	if ChkErr(err) {
		return
	}

	if s.IsDir() {
		if _, err := os.Stat(dst); err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dst, os.ModePerm)
			ChkErrNormal(err)
		}

		fileInfoList, err := os.ReadDir(src)
		if err != nil {

		} else {
			for i := range fileInfoList {
				file := fileInfoList[i]
				file_path := src + "/" + file.Name()
				if file.IsDir() {
					CopyFile(dst+"/"+file.Name(), src+"/"+file.Name())
				} else {
					copy(dst+"/"+file.Name(), file_path, fileFunc...)
				}
			}
		}
	} else {
		if _, err := os.Stat(dst); err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dst, os.ModePerm)
			ChkErr(err)
		}
		copy(dst+"/"+s.Name(), src)
	}
}

func copy(dstFile string, srcFile string, fileFunc ...func(f *os.File) string) {
	dstF, err := os.Create(dstFile)
	if ChkErr(err) {
		return
	}
	defer dstF.Close()

	srcF, err := os.Open(srcFile)
	if ChkErr(err) {
		return
	}
	defer srcF.Close()

	if len(fileFunc) > 0 {
		ff := fileFunc[0]
		contenStr := ff(srcF)
		dstF.WriteString(contenStr)
	} else {
		_, err = io.Copy(dstF, srcF)
		ChkErr(err)
	}
}
