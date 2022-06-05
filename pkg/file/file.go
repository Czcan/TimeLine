package file

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func GetSize(f multipart.File) (int, error) {
	data, err := ioutil.ReadAll(f)
	return len(data), err
}

func CheckSize(f multipart.File, max int) bool {
	size, err := GetSize(f)
	if err != nil {
		return false
	}
	return size <= max
}

func CheckFileExt(name string, exts ...string) bool {
	ext := path.Ext(name)
	for _, val := range exts {
		if val == ext {
			return true
		}
	}
	return false
}

func IsNotExistMkDir(path string) error {
	if ok := Exists(path); !ok {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func SaveUploadFile(file *multipart.FileHeader, path string, filename string) (string, error) {
	if err := IsNotExistMkDir(path); err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	uploadPath := filepath.Join(path, filename+".jpg")
	dst, err := os.OpenFile(uploadPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return uploadPath, nil
}
