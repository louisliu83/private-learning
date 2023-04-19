package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func MD5(d []byte) string {
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func MD5OfFile(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()

	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
