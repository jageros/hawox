/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    pull
 * @Date:    2021/7/5 6:49 下午
 * @package: kafka
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"strings"
	"time"
)

type Consumer struct {
	ctx     contextx.Context
	cg      sarama.ConsumerGroup
	cfg     *Config
	handler func(data []byte)
}

func (c *Consumer) OnMessageHandler(f func(data []byte)) {
	c.handler = f
}

func NewConsumer(ctx contextx.Context, opfs ...func(cfg *Config)) (*Consumer, error) {
	csr := &Consumer{
		ctx: ctx,
		cfg: defaultConfig(),
		handler: func(data []byte) {
			logx.Warn().Msg("Kafka Consumer receive msg, but handler not set")
		},
	}

	for _, opf := range opfs {
		opf(csr.cfg)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Version = sarama.V0_11_0_2
	addrs := strings.Split(csr.cfg.Addrs, ";")
	cli, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	cg, err := sarama.NewConsumerGroupFromClient(csr.cfg.GroupId, cli)

	if err != nil {
		return nil, err
	}

	csr.cg = cg

	csr.run()

	return csr, nil
}

func (c *Consumer) run() {
	c.ctx.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				if c.cg != nil {
					err := c.cg.Close()
					if err != nil {
						return err
					}
				}
				return ctx.Err()
			default:
				err := c.cg.Consume(ctx, []string{c.cfg.Topic}, c)
				if err != nil {
					logx.Err(err).Msg("kafka consume")
					return err
				}
			}
		}
	})

}

func (c *Consumer) Setup(assignment sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(assignment sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(assignment sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		start := time.Now()
		if msg == nil {
			logx.Info().Msg("kafka ConsumeClaim recv msg=nil")
			continue
		}

		c.handler(msg.Value)

		assignment.MarkMessage(msg, "") // 确认消息
		take := time.Now().Sub(start)
		if take >= c.cfg.WarnTime {
			logx.Warn().Str("take", take.String()).Msg("kafka consume msg")
		}
	}
	return nil
}
