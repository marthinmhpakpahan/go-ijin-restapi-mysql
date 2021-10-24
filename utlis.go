package main

import (
    "crypto/md5"
    "encoding/hex"
    "time"
    "strconv"
)

func generateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getCurrentTimestamp() int64 {
	now := time.Now()
	sec := now.Unix()  
	return sec
}

func int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func stringToInt64(s string) (int64, error) {
	no, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return no, err
}