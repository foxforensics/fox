package mqtt

import (
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"go.foxforensics.dev/fox/v4/internal/pkg/types/client"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/event"
)

type Mqtt struct {
	url   string
	topic string

	client mqtt.Client
}

func New(url, topic string) Mqtt {
	return Mqtt{url, topic, client.Mqtt(url)}
}

func (m Mqtt) String() string {
	return fmt.Sprintf("%s/%s", m.url, m.topic)
}

func (m Mqtt) Stream(e *event.Event) error {
	var res mqtt.Token

	if !m.client.IsConnected() {
		res = m.client.Connect()
	}

	if !res.WaitTimeout(client.Timeout) {
		return errors.New("mqtt connect timeout")
	}

	if res.Error() != nil {
		return res.Error()
	}

	res = m.client.Publish(m.topic, 2, false, e.ToJSONL())

	if !res.WaitTimeout(client.Timeout) {
		return errors.New("mqtt publish timeout")
	}

	return res.Error()
}
