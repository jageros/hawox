/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    utils
 * @Date:    2021/5/28 3:29 下午
 * @package: http
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/logx"
	"net/http"
)

func DecodeUrlVal(c *gin.Context, key string) (string, bool) {
	v, ok := c.GetQuery(key)
	if !ok {
		ErrInterrupt(c, errcode.InvalidParam)
	}

	return v, ok
}

func BindQueryArgs(c *gin.Context) (map[string]interface{}, bool) {
	arg := map[string]interface{}{}
	err := c.BindQuery(&arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return nil, false
	}
	return arg, true
}

func BindJsonArgs(c *gin.Context) (map[string]interface{}, bool) {
	arg := map[string]interface{}{}
	err := c.BindJSON(&arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return nil, false
	}
	return arg, true
}

func PkgMsgWrite(c *gin.Context, data interface{}) {
	code := errcode.Success
	dataMap := gin.H{"code": code.Code(), "msg": code.ErrMsg()}
	if data != nil {
		dataMap["data"] = data
	}
	c.JSON(http.StatusOK, dataMap)
}

func ErrInterrupt(c *gin.Context, err errcode.IErr) {
	logx.Errorf("ErrorInterrupt %s", err.Error())
	c.JSON(http.StatusOK, gin.H{"code": err.Code(), "msg": err.ErrMsg()})
	c.Abort()
}