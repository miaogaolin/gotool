package rest

import (
	"encoding/json"
	"net/http"
	"github.com/miaogaolin/gotool/errorx"
	"github.com/miaogaolin/gotool/logx"

	"github.com/gin-gonic/gin"
)

type Err struct {
	Error string `json:"error"`
}

func ErrJson(err error) []byte {
	if err == nil {
		return nil
	}
	logx.Label("ws request").Error(err)
	marshal, _ := json.Marshal(Err{Error: err.Error()})
	return marshal
}

func Success(c *gin.Context, httpCode int, resp interface{}) {
	switch d := resp.(type) {
	case string:
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, d)
	case []byte:
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, string(d))
	default:
		c.JSON(httpCode, d)
	}

	c.Abort()
}

func Error(c *gin.Context, httpCode int, req interface{}, err interface{}) {
	logx.WithField("request", req).Error(err)
	switch e := err.(type) {
	case error:
		c.JSON(httpCode, Err{Error: e.Error()})
	case errorx.Error:
		c.JSON(httpCode, Err{Error: e.ApiMessage})
	}

	c.Abort()
}
