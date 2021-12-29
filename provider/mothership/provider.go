package mothership

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/rsuchkov/gopractice/model"
)

type Provider struct {
	client    http.Client
	serverURI string
}

// New returns a new Provider instance.
func New(serverURI string) (*Provider, error) {
	_, err := url.ParseRequestURI(serverURI)
	if err != nil {
		return nil, fmt.Errorf("incorrect server uri: %s", serverURI)
	}
	provider := &Provider{
		client:    http.Client{},
		serverURI: serverURI,
	}

	return provider, nil
}

func (p *Provider) SendMetric(metric model.Metric) (retErr error) {
	data := url.Values{}

	endpoint := fmt.Sprintf("update/%s/%s/%f", metric.MetricType, metric.Name, metric.Value)
	endpoint = path.Join(p.serverURI, endpoint)

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add("application-type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode received (%d), but (%d) expected", response.StatusCode, http.StatusOK)
	}

	return nil
}
