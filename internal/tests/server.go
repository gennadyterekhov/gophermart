package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

var TestServer *httptest.Server

func InitTestServer(routerInterface chi.Router) {
	TestServer = httptest.NewServer(
		routerInterface,
	)
}

func SendGet(
	t *testing.T,
	ts *httptest.Server,
	path string,
	token string,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, ts.URL+path, strings.NewReader(""))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	response, err := ts.Client().Do(req)
	require.NoError(t, err)
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	return response.StatusCode, bodyAsBytes
}

func SendPostWithoutToken(
	t *testing.T,
	ts *httptest.Server,
	path string,
	requestBody *bytes.Buffer,
) int {
	code, _ := SendPostAndReturnBody(t, ts, path, "application/json", requestBody)

	return code
}

func SendPost(
	t *testing.T,
	ts *httptest.Server,
	path string,
	contentType string,
	requestBody *bytes.Buffer,
) int {
	code, _ := SendPostAndReturnBody(t, ts, path, contentType, requestBody)

	return code
}

func SendPostAndReturnBody(
	t *testing.T,
	ts *httptest.Server,
	path string,
	contentType string,
	requestBody *bytes.Buffer,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, requestBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	response, err := ts.Client().Do(req)
	require.NoError(t, err)
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	return response.StatusCode, bodyAsBytes
}

func getBodyAsBytes(reader io.Reader) ([]byte, error) {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return readBytes, nil
}