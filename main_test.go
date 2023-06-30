package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func runTestServer() *httptest.Server {
	return httptest.NewServer(startServer())
}

func Test_healthCheck(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("it should return 200 when health is ok", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/health", ts.URL))

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})
}

func Test_getAllAirports(t *testing.T) {
	reader, _ := os.ReadFile("./test/resources/airports.csv")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/csv")
		if r.RequestURI == "/airports.csv" {
			w.Write(reader)
		}
	}))

	_ = os.Setenv("AIRPORT_CSV_URL", server.URL+"/airports.csv")

	ts := runTestServer()
	defer ts.Close()

	t.Run("it should return 200 when airports are found", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/airports/all", ts.URL))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		file, _ := os.ReadFile("./test/resources/airports-expected.json")
		assert.Equal(t, string(compact(file[:])), string(compact(bodyBytes[:])))

	})
}

func compact(b []byte) []byte {
	var dst = &bytes.Buffer{}
	_ = json.Compact(dst, b)
	return dst.Bytes()
}
