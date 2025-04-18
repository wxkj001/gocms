package middleware

import (
	"bytes"
	"net/http"
	"strings"

	cryptobin "github.com/deatil/go-cryptobin/cryptobin/rsa"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type dataEncrypt struct {
	config *viper.Viper
}

func NewDataEncrypt(config *viper.Viper) *dataEncrypt {
	return &dataEncrypt{config: config}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {

	return w.body.Write(b)
}
func (c *dataEncrypt) Encrypt(ctx *gin.Context) {
	encrypt := c.config.GetBool("web.encrypt.enable")
	// 如果不需要加密，直接放行
	if !encrypt {
		ctx.Next()
		return
	}
	// 保存原始的 ResponseWriter
	writer := ctx.Writer
	// 创建自定义的 ResponseWriter
	responseWriter := &bodyWriter{
		ResponseWriter: writer,
		body:           &bytes.Buffer{},
	}
	// 替换原始的 ResponseWriter
	ctx.Writer = responseWriter

	// 继续处理请求
	ctx.Next()

	// 获取响应内容
	responseBody := responseWriter.body.Bytes()

	// 如果有响应内容，进行加密处理
	if len(responseBody) > 0 {
		rsa := cryptobin.NewRSA()
		privateKey := c.config.GetString("web.encrypt.publicKey")

		// 加密响应内容
		newBody := rsa.FromBytes(responseBody).
			FromPublicKey([]byte(strings.Join([]string{
				"-----BEGIN PUBLIC KEY-----",
				privateKey,
				"-----END PUBLIC KEY-----",
			}, "\n"))).
			EncryptECB().
			ToBase64String()

		// 重置响应头和状态码
		for k := range writer.Header() {
			writer.Header().Del(k)
		}
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)

		// 写入加密后的内容
		writer.Write([]byte(newBody))

	}
	return
}
