package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Response struct {
	*http.Response
	Logger   *zap.Logger
	Once     sync.Once
	BodyData []byte
	BodyErr  error
}

// Ok verifica se o status da resposta é 2xx.
func (r *Response) Ok() bool {
	if r == nil || r.Response == nil {
		return false
	}
	return r.Response.StatusCode >= 200 && r.Response.StatusCode < 300
}

// Json lê todo o body da resposta via io.ReadAll e retorna os bytes.
// O resultado é cacheado; chamadas subsequentes retornam o mesmo conteúdo.
func (r *Response) Json() ([]byte, error) {
	if r == nil || r.Response == nil || r.Body == nil {
		return nil, nil
	}
	r.Once.Do(func() {
		defer r.Body.Close()
		r.BodyData, r.BodyErr = io.ReadAll(r.Body)
	})

	if r.Logger != nil {
		r.Logger.Info(
			"goxios response: json",
			zap.String("response_body", string(r.BodyData)),
			zap.Int("status", r.Response.StatusCode),
			zap.String("url", r.Response.Request.URL.String()),
			zap.String("method", r.Response.Request.Method),
		)
	}

	return r.BodyData, r.BodyErr
}

// JsonAny faz unmarshal do body JSON em um tipo dinâmico.
func (r *Response) JsonAny() (any, error) {
	if r == nil || r.Response == nil || r.Body == nil {
		return nil, nil
	}
	b, err := r.Json()
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// JsonSlice faz unmarshal do body JSON em []any.
func (r *Response) JsonSlice() ([]any, error) {
	v, err := r.JsonAny()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	out, ok := v.([]any)
	if !ok {
		return nil, errors.New("json is not an array; use JsonMap/JsonInto")
	}
	return out, nil
}

// JsonMap faz unmarshal do body JSON em um map[string]any.
func (r *Response) JsonMap() (map[string]any, error) {
	v, err := r.JsonAny()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	out, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("json is not an object; use JsonSlice/JsonInto")
	}
	return out, nil
}

// JsonInto faz unmarshal do body JSON para a struct/ponteiro informado.
func (r *Response) JsonInto(dst any) error {
	if r == nil || r.Response == nil || r.Body == nil {
		return nil
	}
	b, err := r.Json()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

// JsonAs faz unmarshal do body JSON para o tipo genérico informado.
func JsonAs[T any](r *Response) (T, error) {
	var zero T
	if r == nil || r.Response == nil || r.Body == nil {
		return zero, nil
	}
	b, err := r.Json()
	if err != nil {
		return zero, err
	}
	var out T
	if err := json.Unmarshal(b, &out); err != nil {
		return zero, err
	}
	return out, nil
}

