/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    errcode
 * @Date:    2021/5/28 3:35 下午
 * @package: errcode
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package errcode

import (
	"errors"
	"fmt"
)

var (
	UnknownErrCode      = err{0, "未知错误"}
	InternalErr         = err{-1, "服务器内部错误"}
	Success             = err{200, "successful"} // 成功
	VerifyErr           = err{401, "验证失败"}
	MetaCoderNotFound   = err{402, "meta解码器未注册"}
	ProtoMsgIdNoHandles = err{403, "该协议未注册"}
	ServiceNotFound     = err{404, "未找到该服务"}
	InvalidParam        = err{412, "无效参数"}
	Overload            = err{503, "请求超载"}
)

// IErr 自定义错误接口
type IErr interface {
	Error() string
	ErrMsg() string
	Code() int32
	WithErr(err_ error) IErr
	WithMsg(msg string) IErr
}

type err struct {
	code   int32
	errMsg string
}

func (e err) Error() string {
	return fmt.Sprintf("%d#%s", e.code, e.ErrMsg())
}

func (e err) Code() int32 {
	return e.code
}

func (e err) ErrMsg() string {
	return e.errMsg
}

func (e err) WithErr(err_ error) IErr {
	errMsg := fmt.Sprintf("%s;%s", e.errMsg, err_.Error())
	return New(e.code, errMsg)
}

func (e err) WithMsg(msg string) IErr {
	errMsg := fmt.Sprintf("%s;%s", e.errMsg, msg)
	return New(e.code, errMsg)
}

func (e err) Equal(err IErr) bool {
	return e.Code() == err.Code() && e.ErrMsg() == e.ErrMsg()
}

// =========

// New 创建一个错误码，业务逻辑上的错误，错误码使用1000-1999
func New(code int32, errMsg string) IErr {
	return err{
		code:   code,
		errMsg: errMsg,
	}
}

func WithErrcode(code int32, err_ error) IErr {
	err2 := err{
		code: code,
	}
	if er, ok := err_.(IErr); ok {
		err2.errMsg = fmt.Sprintf("%d_%s", er.Code(), er.ErrMsg())
	} else if err_ != nil {
		err2.errMsg = err_.Error()
	}
	return err2
}

func Errors(errs ...error) error {
	var errMsg string
	for _, err := range errs {
		if err != nil {
			if errMsg == "" {
				errMsg = err.Error()
			} else {
				errMsg = errMsg + "|" + err.Error()
			}
		}
	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}
