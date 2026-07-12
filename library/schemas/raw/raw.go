package raw

import (
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
)

func Apply(evt *event.Event) ([]byte, error) {
	return []byte(evt.AsCEF()), nil
}
