package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rsuchkov/gopractice/service/serverstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestHandler(t *testing.T) {
	type args struct {
		url        string
		statusCode int
		method     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Wrong method test",
			args: args{
				url:        "/update/gauge/dummy/1",
				statusCode: http.StatusMethodNotAllowed,
				method:     http.MethodGet,
			},
		},
		{
			name: "Wrong url",
			args: args{
				url:        "/notfound",
				statusCode: http.StatusNotFound,
				method:     http.MethodPost,
			},
		},
		{
			name: "Wrong metric type",
			args: args{
				url:        "/update/unknown/dummy/1",
				statusCode: http.StatusNotImplemented,
				method:     http.MethodPost,
			},
		},
		{
			name: "Wrong value",
			args: args{
				url:        "/update/gauge/dummy/42fake",
				statusCode: http.StatusBadRequest,
				method:     http.MethodPost,
			},
		},
		{
			name: "Float value",
			args: args{
				url:        "/update/gauge/dummy/42.42",
				statusCode: http.StatusOK,
				method:     http.MethodPost,
			},
		},
		{
			name: "Integer value",
			args: args{
				url:        "/update/gauge/dummy/42",
				statusCode: http.StatusOK,
				method:     http.MethodPost,
			},
		},
		{
			name: "Counter",
			args: args{
				url:        "/update/counter/dummy/42",
				statusCode: http.StatusOK,
				method:     http.MethodPost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := memory.New()
			if err != nil {
				log.Fatal(err)
				return
			}
			svc, err := serverstats.New(serverstats.WithStatsStorage(st))
			if err != nil {
				log.Fatal(err)
				return
			}

			r := NewRouter(svc)
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, _ := testRequest(t, ts, tt.args.method, tt.args.url)
			defer resp.Body.Close()
			assert.Equal(t, tt.args.statusCode, resp.StatusCode)

		})
	}
}
