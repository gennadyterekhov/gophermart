package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"
)

type TestHTTPServer struct {
	Server *httptest.Server
}

// deprecated
var TestServer *httptest.Server

func NewTestHTTPServer(routerInterface chi.Router) *TestHTTPServer {
	return &TestHTTPServer{
		Server: httptest.NewServer(
			routerInterface,
		),
	}
}

// deprecated
func InitTestServer(routerInterface chi.Router) {
	TestServer = httptest.NewServer(
		routerInterface,
	)
}

func (ts *TestHTTPServer) SendGet(
	path string,
	token string,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, ts.Server.URL+path, strings.NewReader(""))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	response, err := ts.Server.Client().Do(req)
	if err != nil {
		panic(err)
	}
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}
	return response.StatusCode, bodyAsBytes
}

func (ts *TestHTTPServer) SendPostWithoutToken(
	path string,
	requestBody *bytes.Buffer,
) int {
	code, _ := ts.SendPostAndReturnBody(path, "application/json", "", requestBody)

	return code
}

func (ts *TestHTTPServer) SendPost(
	path string,
	contentType string,
	token string,
	requestBody *bytes.Buffer,
) int {
	code, _ := ts.SendPostAndReturnBody(path, contentType, token, requestBody)

	return code
}

func (ts *TestHTTPServer) SendPostAndReturnBody(
	path string,
	contentType string,
	token string,
	requestBody *bytes.Buffer,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodPost, ts.Server.URL+path, requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)

	response, err := ts.Server.Client().Do(req)
	if err != nil {
		panic(err)
	}
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}
	return response.StatusCode, bodyAsBytes
}

func getBodyAsBytes(reader io.Reader) ([]byte, error) {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return readBytes, nil
}
