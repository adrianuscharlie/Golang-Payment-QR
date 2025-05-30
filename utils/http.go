package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestMeta struct {
	ClientSecret string
	ExtraParam1  string
	ExtraParam2  string
	ExtraParam3  string
	Token        string
}

func SendRequest(method string, url string, body interface{}, meta RequestMeta) ([]byte, string, error) {
	timeStamp := time.Now().Add(5 * time.Minute).Format(time.RFC3339)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	signature, err := SignatureHeader(meta.ClientSecret, timeStamp)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Signature request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-TIMESTAMP", timeStamp)
	httpReq.Header.Set("X-PARTNER-ID", meta.ExtraParam1)
	httpReq.Header.Set("X-EXTERNAL-ID", meta.ExtraParam2)
	httpReq.Header.Set("CHANNEL-ID", meta.ExtraParam3)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", meta.Token))
	httpReq.Header.Set("X-SIGNATURE", signature)

	// Build header string for logging
	var headerStr string
	for key, values := range httpReq.Header {
		for _, value := range values {
			headerStr += key + ":" + value + "\n"
		}
	}

	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, headerStr, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()
	bodyResp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, headerStr, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyResp, headerStr, nil
}
