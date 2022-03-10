// Code generated by metactl. DO NOT EDIT.
// source: metactl

package meta

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/jageros/hawox/example/mygame/protos/meta/sess"

	pb "github.com/jageros/hawox/example/mygame/protos/pb"
)

var metaData = make(map[pb.MsgID]IMeta)

var NoMetaErr = errors.New("NoMetaErr")

type IMeta interface {
	GetMsgID() pb.MsgID
	EncodeArg(interface{}) ([]byte, error)
	DecodeArg([]byte) (interface{}, error)
	EncodeReply(interface{}) ([]byte, error)
	DecodeReply([]byte) (interface{}, error)
	Handle(session sess.ISession, arg interface{}) (interface{}, error)
}

func registerMeta(meta IMeta) {
	metaData[meta.GetMsgID()] = meta
}

func GetMeta(msgId pb.MsgID) (IMeta, error) {
	if m, ok := metaData[msgId]; ok {
		return m, nil
	} else {
		return nil, NoMetaErr
	}
}

func Call(session sess.ISession, msgid pb.MsgID, data []byte) ([]byte, error) {
	im, err := GetMeta(msgid)
	if err != nil {
		return nil, err
	}
	arg, err := im.DecodeArg(data)
	if err != nil {
		return nil, err
	}
	
	var resp interface{}
	err = catchPanic(func() error {
		resp, err = im.Handle(session, arg)
		return err
	})
	if err != nil {
		return nil, err
	}

	return im.EncodeReply(resp)
}

func catchPanic(f func() error) (err error) {
	defer func() {
		err1 := recover()
		if err1 != nil {
			fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			err = errors.New(fmt.Sprintf("%!s(MISSING) call err: %!v(MISSING)", fn, err1))
		}
	}()

	err = f()

	return
}
