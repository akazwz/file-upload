package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
func PathExists(dst string) bool {
	_, err := os.Stat(dst)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func MergeChunkFile(dir string) (int64, error) {
	start := time.Now().UnixMicro()
	// 按照文件名index排序读取文件夹内的文件
	files, _ := ioutil.ReadDir(dir)
	sort.Slice(files, func(i, j int) bool {
		// 获取文件 index
		filename := files[i].Name()
		index := strings.Split(filename, "-")[0]

		indexInt, _ := strconv.Atoi(index)
		nextInt, _ := strconv.Atoi(strings.Split(files[j].Name(), "-")[0])
		return indexInt < nextInt
	})

	// 创建完整文件
	completeFile, err := os.Create(fmt.Sprintf("%s/complete", dir))
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		/* 无用文件, 跳过 */
		if file.Name() == ".BD_Store" {
			continue
		}

		// 读取 chunk file
		bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return 0, err
		}

		// 完整文件写入数据
		_, err = completeFile.Write(bytes)
		if err != nil {
			return 0, err
		}
	}
	end := time.Now().UnixMicro()
	timeSend := end - start
	return timeSend, nil
}
