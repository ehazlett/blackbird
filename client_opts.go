package blackbird

import (
	"context"
	"fmt"
	"time"

	api "github.com/ehazlett/blackbird/api/v1"
	"github.com/gogo/protobuf/types"
)

type AddOpts func(ctx context.Context, srv *api.Server) error

func WithPath(path string) AddOpts {
	return func(ctx context.Context, srv *api.Server) error {
		if path == "" {
			return fmt.Errorf("path cannot be empty")
		}
		srv.Path = path
		return nil
	}
}

func WithTLS(ctx context.Context, srv *api.Server) error {
	srv.TLS = true
	return nil
}

func WithPolicy(p api.Policy) AddOpts {
	return func(ctx context.Context, srv *api.Server) error {
		srv.Policy = p
		return nil
	}
}

func WithUpstreams(upstreams ...string) AddOpts {
	return func(ctx context.Context, srv *api.Server) error {
		srv.Upstreams = upstreams
		return nil
	}
}

func WithTimeouts(d time.Duration) AddOpts {
	return func(ctx context.Context, srv *api.Server) error {
		srv.Timeouts = types.DurationProto(d)
		return nil
	}
}

func WithPreset(preset string) AddOpts {
	return func(ctx context.Context, srv *api.Server) error {
		srv.Preset = preset
		return nil
	}
}
