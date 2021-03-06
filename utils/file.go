package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"mime/multipart"
	"os"
)

// HashFileByAlgo 根据算法获取文件hash
func HashFileByAlgo(fh *multipart.FileHeader, algo string) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "nil", err
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	hashcode := getHash(algo)

	if _, err := io.Copy(hashcode, file); err != nil {
		log.Println(err)
		return "", err
	}

	return hex.EncodeToString(hashcode.Sum(nil)), nil
}

func getHash(algo string) hash.Hash {
	switch algo {
	case "md5":
		return md5.New()
	case "sha256":
		return sha256.New()
	case "sha512":
		return sha512.New()
	case "sha1":
		return sha1.New()
	default:
		return sha256.New()
	}
}

// PathExists
// check file or dir exists
// 检查文件或者文件夹是否存在
func PathExists(dst string) (bool, error) {
	_, err := os.Stat(dst)
	if err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
