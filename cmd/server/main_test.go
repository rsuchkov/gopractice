package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rsuchkov/gopractice/service/serverstats"
	"github.com/rsuchkov/gopractice/storage/memory"
)

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
				url:        "/update/gauge/dummy/1/",
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
				url:        "/update/unknown/dummy/1/",
				statusCode: http.StatusNotImplemented,
				method:     http.MethodPost,
			},
		},
		{
			name: "Wrong value",
			args: args{
				url:        "/update/gauge/dummy/42fake/",
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
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			st, err := memory.New()
			if err != nil {
				fmt.Println(err)
				return
			}
			svc, err := serverstats.New(serverstats.WithStatsStorage(st))
			if err != nil {
				fmt.Println(err)
				return
			}
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Handler(svc, w, r)
			})
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()

		})
	}
}
