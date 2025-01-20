package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Add these structs to your model package
type OpenTokenPublicKeyResponse struct {
    Data OpenTokenPublicKey `json:"data"`
}

type OpenTokenPublicKey struct {
    ID        int    `json:"id"`
    PublicKey string `json:"publicKey"`
}

// Add this new function
func GetOpenTokenPublicKey(ctx context.Context, endpoint Endpoint, tokenID int) (*OpenTokenPublicKey, error) {
    ctx, span := modelTracer.Start(ctx, "http.getPublicKey")
    defer span.End()

    client := &http.Client{
        Timeout:   time.Second * 10,
        Transport: otelhttp.NewTransport(http.DefaultTransport),
    }

    // Create URL with query parameter
    url := fmt.Sprintf("%s/api/v1/opentoken/publickey?tid=%d", endpoint.APIEndpoint, tokenID)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        logrus.Errorln(err)
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", fmt.Sprintf("shelltimeCLI@%s", commitID))
    req.Header.Set("Authorization", "CLI "+endpoint.Token)

    logrus.Traceln("http: ", req.URL.String())

    resp, err := client.Do(req)
    if err != nil {
        logrus.Errorln(err)
        return nil, err
    }
    defer resp.Body.Close()

    logrus.Traceln("http: ", resp.Status)

    if resp.StatusCode != http.StatusOK {
        buf, err := io.ReadAll(resp.Body)
        if err != nil {
            logrus.Errorln(err)
            return nil, err
        }
        var errResp errorResponse
        if err := json.Unmarshal(buf, &errResp); err != nil {
            logrus.Errorln("Failed to parse error response:", err)
            return nil, err
        }
        logrus.Errorln("Error response:", errResp.ErrorMessage)
        return nil, errors.New(errResp.ErrorMessage)
    }

    var response OpenTokenPublicKeyResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        logrus.Errorln("Failed to decode response:", err)
        return nil, err
    }

    return &response.Data, nil
}
