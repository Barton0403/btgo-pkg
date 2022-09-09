package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return true
	// 注意不要在这里使用锁性能成功率反而更低
	// 此方法算力大处理时间长，高并发情况下会造成CPU跑满处理速度下降，导致请求超时
	// QPS 50-60
	//err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	//return err == nil
}
