package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/net/client"
	"go.foxforensics.eu/fox/v4/internal/net/schema"
	"go.foxforensics.eu/fox/v4/internal/net/schema/ecs"
	"go.foxforensics.eu/fox/v4/internal/net/schema/hec"
	"go.foxforensics.eu/fox/v4/internal/net/schema/raw"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
)

type Options struct {
	Url    string
	Token  string
	Schema schema.Schema
}

type Http struct {
	client *http.Client
	opts   *Options
}

func Create(opts *Options) (*Http, error) {
	u, err := url.Parse(opts.Url)

	if err != nil {
		return nil, err
	}

	if u.Scheme == "http" {
		slog.Warn("data will be streamed unencrypted!")
	}

	return &Http{client.Http(), opts}, nil
}

func (h Http) String() string {
	return fmt.Sprintf("%s@%s", h.opts.Schema, h.opts.Url)
}

func (h Http) Stream(ctx context.Context, evt *event.Event) error {
	var buf []byte
	var err error

	switch h.opts.Schema {
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

	u, err := url.Parse(h.opts.Url)

	if err != nil {
		return err
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return errors.New("unsupported scheme")
	}

	if u.Scheme == "http" {
		slog.Warn("data will be streamed unencrypted!")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.opts.Url, bytes.NewReader(buf))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", client.Name())

	switch h.opts.Schema {
	case schema.Ecs, schema.Hec:
		req.Header.Set("Content-Type", "application/json")
	default:
		req.Header.Set("Content-Type", "text/plain")
	}

	// add authorization for Splunk
	if h.opts.Schema == schema.Hec && len(h.opts.Token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", h.opts.Token))
	}

	res, err := h.client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	// drain body
	_, _ = io.Copy(io.Discard, res.Body)

	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return nil
}

func (h Http) Close() error {
	return nil // stateless
}
