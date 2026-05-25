package raw

import "go.foxforensics.dev/fox/v4/internal/pkg/types/event"

func Apply(evt *event.Event) ([]byte, error) {
	return []byte(evt.ToCEF()), nil
}
