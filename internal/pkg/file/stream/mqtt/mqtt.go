package mqtt

import (
	"context"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.golang/paho"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/schema"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/schema/ecs"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/schema/hec"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/schema/raw"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/client"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/event"
)

const QoS = 2

type Options struct {
	Url      string
	Topic    string
	Username string
	Password string
	Schema   schema.Schema
}

type Mqtt struct {
	client *mqtt.Client
	opts   *Options
}

func New(opts *Options) *Mqtt {
	return &Mqtt{client.Mqtt(opts.Url, opts.Username, opts.Password), opts}
}

func (m Mqtt) String() string {
	return fmt.Sprintf("%s@%s/%s", m.opts.Schema, m.opts.Url, m.opts.Topic)
}

func (m Mqtt) Stream(evt *event.Event) error {
	var buf []byte
	var err error

	switch m.opts.Schema {
	case schema.Ecs:
		buf, err = ecs.Apply(evt)
	case schema.Hec:
		buf, err = hec.Apply(evt)
	case schema.Raw:
		buf, err = raw.Apply(evt)
	}

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	res, err := m.client.Publish(ctx, &mqtt.Publish{
		Topic:   m.opts.Topic,
		QoS:     QoS,
		Retain:  false,
		Payload: buf,
	})

	if res != nil && res.ReasonCode > 0 {
		log.Println(res.Properties.ReasonString)
	}

	return err
}
