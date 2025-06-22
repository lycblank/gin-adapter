package ginadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/authboss/v3"
)

var _ authboss.WrappingResponseWriter = &AuthbossResponseWriter{}

// 包裹authboss中间件
func WarpAuthboss(f func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h, recoverWriter := makeAuthBossHandler(c)
		gin.WrapH(f(h))(c)
		recoverWriter()
	}
}

func makeAuthBossHandler(c *gin.Context) (http.Handler, func()) {
	writer := c.Writer
	f := func(w http.ResponseWriter, r *http.Request) {
		newResponseWriter := &AuthbossResponseWriter{
			ResponseWriter: c.Writer,
			AuthbossWriter: w,
		}
		c.Writer = newResponseWriter
		c.Next()
	}
	recoverWriter := func() {
		c.Writer = writer
	}
	return http.HandlerFunc(f), recoverWriter
}

type AuthbossResponseWriter struct {
	gin.ResponseWriter
	AuthbossWriter http.ResponseWriter
}

func (arw *AuthbossResponseWriter) Unwrap() http.ResponseWriter {
	return arw.AuthbossWriter
}
