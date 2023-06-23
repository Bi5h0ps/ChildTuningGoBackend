package provider

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() HttpClient {
	return HttpClient{&http.Client{}}
}

func (c *HttpClient) Post(url string, requestBody *bytes.Buffer, requestWriter *multipart.Writer, ctx *gin.Context) ([]byte, error) {
	// Create a new HTTP request with the POST method
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return []byte{}, err
	}
	// Set the Content-Type header with the boundary from the multipart writer
	req.Header.Set("Content-Type", requestWriter.FormDataContentType())
	// Set the request body with the multipart payload
	req.Body = io.NopCloser(requestBody)
	// Send the request and retrieve the response
	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return respBody, nil
}
