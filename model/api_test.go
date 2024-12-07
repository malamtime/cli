package model

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vmihailenco/msgpack/v5"
)

type apiTestSuite struct {
    suite.Suite
}

func (s *apiTestSuite) TestDoSendData() {
    // Test successful request
    s.T().Run("successful request", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Verify headers
            assert.Equal(t, "application/msgpack", r.Header.Get("Content-Type"))
            assert.Contains(t, r.Header.Get("User-Agent"), "shelltimeCLI@")
            assert.Equal(t, "CLI testToken", r.Header.Get("Authorization"))

            // Decode request body
            var payload PostTrackArgs
            err := msgpack.NewDecoder(r.Body).Decode(&payload)
            assert.NoError(t, err)

            // Verify payload
            assert.Equal(t, int64(1000), payload.CursorID)
            assert.Len(t, payload.Data, 1)
            assert.Equal(t, "test_shell", payload.Data[0].Shell)

            w.WriteHeader(http.StatusNoContent)
        }))
        defer server.Close()

        endpoint := Endpoint{
            Token:       "testToken",
            APIEndpoint: server.URL,
        }

        trackingData := []TrackingData{
            {
                Shell:     "test_shell",
                SessionID: 123,
                Command:   "test_command",
                StartTime: time.Now().Unix(),
                EndTime:   time.Now().Unix(),
            },
        }

        err := doSendData(context.Background(), endpoint, time.Unix(0, 1000), trackingData)
        assert.NoError(t, err)
    })

    // Test error response
    s.T().Run("error response", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"code": 400, "error": "test error"}`))
        }))
        defer server.Close()

        endpoint := Endpoint{
            Token:       "testToken",
            APIEndpoint: server.URL,
        }

        trackingData := []TrackingData{
            {
                Shell:     "test_shell",
                SessionID: 123,
                Command:   "test_command",
                StartTime: time.Now().Unix(),
                EndTime:   time.Now().Unix(),
            },
        }

        err := doSendData(context.Background(), endpoint, time.Unix(0, 1000), trackingData)
        assert.Error(t, err)
        assert.Equal(t, "test error", err.Error())
    })
}

func (s *apiTestSuite) TestSendLocalDataToServer() {
    s.T().Run("no token configured", func(t *testing.T) {
        config := ShellTimeConfig{
            Token: "",
        }

        err := SendLocalDataToServer(context.Background(), config, time.Now(), nil)
        assert.NoError(t, err)
    })

    s.T().Run("multiple endpoints", func(t *testing.T) {
        requestCount := 0
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestCount++
            w.WriteHeader(http.StatusNoContent)
        }))
        defer server.Close()

        config := ShellTimeConfig{
            Token:       "mainToken",
            APIEndpoint: server.URL,
            Endpoints: []Endpoint{
                {
                    Token:       "token1",
                    APIEndpoint: server.URL,
                },
                {
                    Token:       "token2",
                    APIEndpoint: server.URL,
                },
            },
        }

        trackingData := []TrackingData{
            {
                Shell:     "test_shell",
                SessionID: 123,
                Command:   "test_command",
                StartTime: time.Now().Unix(),
                EndTime:   time.Now().Unix(),
            },
        }

        err := SendLocalDataToServer(context.Background(), config, time.Now(), trackingData)
        assert.NoError(t, err)
        assert.Equal(t, 3, requestCount) // Main endpoint + 2 additional endpoints
    })

    s.T().Run("partial failure", func(t *testing.T) {
        successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusNoContent)
        }))
        defer successServer.Close()

        failureServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"code": 400, "error": "test error"}`))
        }))
        defer failureServer.Close()

        config := ShellTimeConfig{
            Token:       "mainToken",
            APIEndpoint: successServer.URL,
            Endpoints: []Endpoint{
                {
                    Token:       "token1",
                    APIEndpoint: failureServer.URL,
                },
            },
        }

        trackingData := []TrackingData{
            {
                Shell:     "test_shell",
                SessionID: 123,
                Command:   "test_command",
                StartTime: time.Now().Unix(),
                EndTime:   time.Now().Unix(),
            },
        }

        err := SendLocalDataToServer(context.Background(), config, time.Now(), trackingData)
        assert.Error(t, err)
        assert.Equal(t, "test error", err.Error())
    })
}

func TestAPITestSuite(t *testing.T) {
    suite.Run(t, new(apiTestSuite))
}
