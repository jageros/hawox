/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    publish
 * @Date:    2021/7/2 3:51 下午
 * @package: nsq
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package nsq

import (
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/protos/meta"
	"github.com/jageros/hawox/protos/pbf"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"sync"
	"time"
)

type Producer struct {
	ctx contextx.Context
	opt *Config
	pd  *nsq.Producer
	cfg *nsq.Config
	clk *sync.Mutex
}

func (p *Producer) getNodeAddr() (string, error) {
	addrs := fmt.Sprintf(p.opt.Addrs, ";")
	idx := rand.Intn(len(addrs))
	url := fmt.Sprintf("http://%s/nodes", addrs[idx])
	resp, err := httpc.Request(httpc.GET, url, httpc.FORM, nil, nil)
	if err != nil {
		return "", err
	}
	pds := resp["producers"].([]interface{})
	pdn := len(pds)
	if pdn <= 0 {
		return "", errcode.New(101, "无可用NSQ节点")
	}
	idx = rand.Intn(len(pds))
	pd := pds[idx].(map[string]interface{})
	addr := fmt.Sprintf("%v:%v", pd["broadcast_address"], pd["tcp_port"])
	logx.Debugf("NsqAddr=%s", addr)
	return addr, nil
}

func (p *Producer) connectToNsqd() error {
	p.clk.Lock()
	defer p.clk.Unlock()

	if p.pd != nil {
		err := p.pd.Ping()
		if err == nil {
			return nil
		}
		p.pd.Stop()
	}

	addr, err := p.getNodeAddr()
	if err != nil {
		return err
	}
	pd, err := nsq.NewProducer(addr, p.cfg)
	if err != nil {
		return err
	}

	p.pd = pd
	return nil
}

func (p *Producer) PushProtoMsg(msgId int32, arg interface{}, target *pbf.Target) error {
	start := time.Now()
	im, err := meta.GetMeta(msgId)
	if err != nil {
		return err
	}
	data, err := im.EncodeArg(arg)
	if err != nil {
		return err
	}
	resp := &pbf.Response{
		MsgID:   msgId,
		Code:    errcode.Success.Code(),
		Payload: data,
	}
	data2, err := resp.Marshal()
	if err != nil {
		return err
	}
	msg := &pbf.QueueMsg{
		Data:    data2,
		Targets: target,
	}
	err = p.Push(msg)
	take := time.Now().Sub(start)
	if take > p.opt.WarnTime {
		logx.Warnf("Nsq Push Msg take: %s", take.String())
	}
	return err
}

func (p *Producer) Push(msg *pbf.QueueMsg) error {
	data, err := msg.Marshal()
	if err != nil {
		return err
	}
	err = p.pd.Publish(p.opt.Topic, data)
	if err != nil {
		err = p.connectToNsqd()
		if err != nil {
			return err
		}
		err = p.pd.Publish(p.opt.Topic, data)
	}
	return err
}

func NewProducer(g contextx.Context, opfs ...func(cfg *Config)) (*Producer, error) {
	p := &Producer{
		ctx: g,
		opt: defaultConfig(),
		clk: &sync.Mutex{},
	}

	for _, opf := range opfs {
		opf(p.opt)
	}

	p.cfg = nsq.NewConfig()
	err := p.connectToNsqd()
	p.run()
	return p, err
}

func (p *Producer) run() {
	p.ctx.Go(func(ctx contextx.Context) error {
		<-ctx.Done()
		p.pd.Stop()
		return ctx.Err()
	})
}