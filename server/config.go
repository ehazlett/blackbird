package server

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	api "github.com/ehazlett/blackbird/api/v1"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	configTemplate = `# blackbird proxy config
*:{{ $.HTTPPort }} {
    status 200 /healthz
}
{{ range $server := .Servers }}
{{ if $server.TLS }}https{{ else }}http{{ end }}://{{ $server.Host }}:{{ if $server.TLS }}{{ $.HTTPSPort }}{{ else }}{{ $.HTTPPort }}{{ end }} {
    proxy {{ $server.Path }} { {{ if ne $server.Preset "" }}
	{{ $server.Preset }}{{ end }}
	policy {{policyname $server.Policy }}
	try_duration 1s
	{{ range $upstream := $server.Upstreams }}upstream {{ $upstream }}
	{{ end }}
    }
    {{ with $t := $server.Timeouts }} {{ if $t }}timeouts {{duration $t }} {{ end }} {{ end }}
}{{ end }}
`
)

type proxyConfig struct {
	HTTPPort  int
	HTTPSPort int
	Servers   []*api.Server
}

func policyName(v api.Policy) string {
	return strings.ToLower(v.String())
}

func duration(v *types.Duration) time.Duration {
	d, err := types.DurationFromProto(v)
	if err != nil {
		logrus.Error(errors.Wrap(err, "error converting proto duration"))
	}
	return d
}

func (s *Server) generateConfig() ([]byte, error) {
	t := template.New("blackbird").Funcs(template.FuncMap{
		"policyname": policyName,
		"duration":   duration,
	})
	tmpl, err := t.Parse(configTemplate)
	if err != nil {
		return nil, err
	}

	servers, err := s.datastore.Servers()
	if err != nil {
		return nil, err
	}
	srvs := []*api.Server{}
	for _, srv := range servers {
		srvs = append(srvs, srv)
	}

	config := &proxyConfig{
		HTTPPort:  s.config.HTTPPort,
		HTTPSPort: s.config.HTTPSPort,
		Servers:   srvs,
	}

	var c bytes.Buffer
	if err := tmpl.Execute(&c, config); err != nil {
		return nil, err
	}

	return c.Bytes(), nil
}
