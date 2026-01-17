package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_JsonMap_JsonInto_JsonAs(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"n":123,"s":"x"}`))
	}))
	t.Cleanup(srv.Close)

	// Usando http.Get para obter um Response válido que implementa io.ReadCloser
	httpResp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("http.Get() err=%v", err)
	}
	resp := &Response{Response: httpResp}

	m, err := resp.JsonMap()
	if err != nil {
		t.Fatalf("JsonMap() err=%v", err)
	}
	if m["ok"] != true {
		t.Fatalf("expected ok=true; got=%v", m["ok"])
	}

	type payload struct {
		OK bool   `json:"ok"`
		N  int    `json:"n"`
		S  string `json:"s"`
	}

	httpResp2, _ := http.Get(srv.URL)
	resp2 := &Response{Response: httpResp2}
	var p payload
	if err := resp2.JsonInto(&p); err != nil {
		t.Fatalf("JsonInto() err=%v", err)
	}
	if p.OK != true || p.N != 123 || p.S != "x" {
		t.Fatalf("unexpected payload: %+v", p)
	}

	httpResp3, _ := http.Get(srv.URL)
	resp3 := &Response{Response: httpResp3}
	p2, err := JsonAs[payload](resp3)
	if err != nil {
		t.Fatalf("JsonAs() err=%v", err)
	}
	if p2.OK != true || p2.N != 123 || p2.S != "x" {
		t.Fatalf("unexpected payload: %+v", p2)
	}
}

func TestResponse_JsonSlice(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
	}))
	t.Cleanup(srv.Close)

	httpResp, _ := http.Get(srv.URL)
	resp := &Response{Response: httpResp}

	arr, err := resp.JsonSlice()
	if err != nil {
		t.Fatalf("JsonSlice() err=%v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("expected len=2; got=%d", len(arr))
	}
}

func TestResponse_Ok(t *testing.T) {
	resp := &Response{
		Response: &http.Response{
			StatusCode: 200,
		},
	}
	if !resp.Ok() {
		t.Error("expected Ok() true for 200")
	}

	resp.StatusCode = 404
	if resp.Ok() {
		t.Error("expected Ok() false for 404")
	}
}

func TestResponse_Json_Error(t *testing.T) {
	// Simulate error in Body Read
	resp := &Response{
		Response: &http.Response{
			Body: io.NopCloser(errorReader{}),
		},
	}
	_, err := resp.Json()
	if err == nil {
		t.Error("expected error when reading body")
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
