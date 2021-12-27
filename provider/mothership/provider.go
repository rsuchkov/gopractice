package mothership

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rsuchkov/gopractice/model"
)

type Provider struct {
	client http.Client
	server string
}

// New returns a new Provider instance.
func New(server string) *Provider {
	provider := &Provider{
		client: http.Client{},
		server: server,
	}
	return provider
}

func (p *Provider) SendMetric(metric model.Metric) (statsuCode int, retErr error) {
	data := url.Values{}
	endpoint := fmt.Sprintf("%supdate/%s/%s/%f", p.server, metric.MetricType, metric.Name, metric.Value)
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return 0, err
	}
	request.Header.Add("application-type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := p.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(body))
	return response.StatusCode, nil
}
