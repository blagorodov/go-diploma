package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"diploma/internal/logger"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

var hashKey = []byte("secrethashkey")

var ErrWrongTokenFormat = errors.New("wrong token format")
var ErrWrongToken = errors.New("wrong token")

func GetID(c *gin.Context) (string, error) {
	var token string
	token, err := c.Cookie("token")
	if err != nil {
		token = c.GetHeader("Authorization")
	}
	return DecodeID(token)
}

func SetID(c *gin.Context, userID string) {
	if userID != "" {
		token := EncodeID(userID)
		c.SetCookie("token", token, 0, "", "", false, true)
		c.Writer.Header().Set("Authorization", token)
	}
}

func DecodeID(token string) (string, error) {
	parts := strings.Split(token, ":")

	if len(parts) != 2 {
		return "", ErrWrongTokenFormat
	}
	id := parts[0]
	key, _ := hex.DecodeString(parts[1])

	h := hmac.New(sha256.New, hashKey)
	h.Write([]byte(id))
	dst := h.Sum(nil)
	fmt.Println(dst)
	fmt.Println(key)
	if !hmac.Equal(dst, key) {
		return "", ErrWrongToken
	}
	return id, nil
}

func EncodeID(UserID string) string {
	h := hmac.New(sha256.New, hashKey)
	h.Write([]byte(UserID))
	dst := h.Sum(nil)

	return fmt.Sprintf("%s:%x", UserID, dst)
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func init() {
	key, err := generateRandom(16)
	if err != nil {
		logger.Log("auth::init")
		logger.Log(err.Error())
		panic(err)
	}
	hashKey = key
}
