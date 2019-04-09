package ginhelper

import (
	"github.com/gentwolf-shen/gohelper/convert"
	"github.com/gentwolf-shen/gohelper/dict"
	"github.com/gin-gonic/gin"
)

func AllowCrossDomainAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowCrossDomain(c, c.Request.Header.Get("Origin"))
	}
}

func AllowCrossDomain(domains []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("Origin")
		bl := false
		for _, domain := range domains {
			if domain == host {
				bl = true
				break
			}
		}

		if bl {
			allowCrossDomain(c, c.Request.Header.Get("Origin"))
		}
	}
}

func allowCrossDomain(c *gin.Context, host string) {
	c.Header("Access-Control-Allow-Origin", host)
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Max-Age", "3600")
		c.AbortWithStatus(200)
	}
}

func ShowNoContent(c *gin.Context) {
	c.Status(204)
}

func ShowNetError(c *gin.Context) {
	ShowError(c, 5000000)
}

func ShowNoAuth(c *gin.Context) {
	ShowError(c, 4010000)
}

func ShowParamError(c *gin.Context) {
	ShowError(c, 4000001)
}

func ShowNotFound(c *gin.Context) {
	ShowError(c, 4040000)
}

func ShowError(c *gin.Context, errorCode int) {
	if errorCode == 0 {
		ShowSuccess(c, nil)
	} else {
		msg := ErrorMessage{}
		msg.Code = errorCode
		msg.Message = dict.Get(convert.ToStr(errorCode))

		ShowMsg(c, errorCode/10000, msg)
	}
}

func ShowErrorMsg(c *gin.Context, errorCode int, errMsg interface{}) {
	msg := ErrorMessage{}
	msg.Code = errorCode
	msg.Message = errMsg

	ShowMsg(c, errorCode/10000, msg)
}

func ShowSuccess(c *gin.Context, msg interface{}) {
	if msg == nil {
		msg := ErrorMessage{}
		msg.Code = 0
		msg.Message = "success"
		ShowMsg(c, 200, msg)
	} else {
		ShowMsg(c, 200, msg)
	}
}

func ShowMsg(c *gin.Context, httpCode int, msg interface{}) {
	c.JSON(httpCode, msg)
}

type ErrorMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
