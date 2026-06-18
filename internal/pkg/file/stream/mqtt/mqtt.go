package mqtt

import (
	"context"
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.golang/paho"

	"go.foxforensics.eu/fox/v4/internal/pkg/file/schema"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/schema/ecs"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/schema/hec"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/schema/raw"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/client"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
)

type Options struct {
	Url      string
	Topic    string
	Username string
	Password string
	QoS      byte
	Schema   schema.Schema
}

type Mqtt struct {
	client *mqtt.Client
	opts   *Options
}

func Create(opts *Options) (*Mqtt, error) {
	v, err := client.Mqtt(opts.Url, opts.Username, opts.Password)
	return &Mqtt{v, opts}, err
}

func (m Mqtt) String() string {
	return fmt.Sprintf("%s@%s/%s", m.opts.Schema, m.opts.Url, m.opts.Topic)
}

func (m Mqtt) Stream(ctx context.Context, evt *event.Event) error {
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

	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	res, err := m.client.Publish(ctx, &mqtt.Publish{
		Topic:   m.opts.Topic,
		QoS:     m.opts.QoS,
		Retain:  false,
		Payload: buf,
	})

	if res != nil && res.ReasonCode > 0 {
		return errors.New(res.Properties.ReasonString)
	}

	return err
}

func (m Mqtt) Close() error {
	return m.client.Disconnect(&mqtt.Disconnect{
		ReasonCode: 0,
	})
}
