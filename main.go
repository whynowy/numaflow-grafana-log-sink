package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// grafanaLogSink is a sinker implementation that sends data to grafana
type grafanaLogSink struct {
	url        string
	httpClient *http.Client
}

func newGrafanaSink(url string) *grafanaLogSink {
	return &grafanaLogSink{
		url: url,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 1,
		},
	}
}

func (gs *grafanaLogSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()

	var buffer bytes.Buffer
	buffer.WriteString(`{"streams": [{"stream": {}, "values": [`)

	fmt.Println(buffer.String())
	for d := range datumStreamCh {

		log := fmt.Sprintf(`["%s", "%s"]`)

		buffer.WriteString(string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))

	}
	buffer.WriteString(`]}]}`)

	// logs := []byte(fmt.Sprintf("{"streams": [{"stream": {"Language": "Go", "source": "Code" }, "values": [["%s", "This is my log line"]]}]}", strconv.FormatInt(time.Now().UnixNano(), 10)))

	req, err := http.NewRequest("POST", gs.url, bytes.NewBuffer(logs))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("299424", "your-api-key")
	client.Do(req)

	return result
}

func main() {
	err := sinksdk.NewServer(grafanaLogSink("host", "path")).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
