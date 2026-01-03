package client

import (
	"fmt"
	"net/http"
	"time"

	app "github.com/cuhsat/fox/v4/internal"
)

// UserAgent fox
var UserAgent = fmt.Sprintf("fox %s", app.Version)

// Default client
var Default = &http.Client{
	Timeout: time.Second * 30,
}
