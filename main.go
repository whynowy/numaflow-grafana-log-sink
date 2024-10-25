package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// grafanaSink is a sinker implementation that sends data to grafana
type grafanaSink struct {
	client *goapi.GrafanaHTTPAPI
}

func newGrafanaSink(host, basePath string) *grafanaSink {
	cfg := &goapi.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: host,
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: basePath,
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{"http"},
		// APIKey is an optional API key or service account token.
		APIKey: os.Getenv("API_ACCESS_TOKEN"),
		// BasicAuth is optional basic auth credentials.
		BasicAuth: url.UserPassword("admin", "admin"),
		// OrgID provides an optional organization ID.
		// OrgID is only supported with BasicAuth since API keys are already org-scoped.
		OrgID: 1,
		// TLSConfig provides an optional configuration for a TLS client
		TLSConfig: &tls.Config{},
		// NumRetries contains the optional number of attempted retries
		NumRetries: 3,
		// RetryTimeout sets an optional time to wait before retrying a request
		RetryTimeout: 0,
		// RetryStatusCodes contains the optional list of status codes to retry
		// Use "x" as a wildcard for a single digit (default: [429, 5xx])
		RetryStatusCodes: []string{"420", "5xx"},
		// HTTPHeaders contains an optional map of HTTP headers to add to each request
		HTTPHeaders: map[string]string{},
	}
	return &grafanaSink{
		client: goapi.NewHTTPClientWithConfig(strfmt.Default, cfg),
	}
}

func (l *grafanaSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		l.client.
			fmt.Println("User Defined Sink:", string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
		// if we are not able to write to sink and if we have a fallback sink configured
		// we can use sinksdk.ResponseFallback(id)) to write the message to fallback sink
	}
	return result
}

func main() {
	err := sinksdk.NewServer(newGrafanaSink("host", "path")).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
