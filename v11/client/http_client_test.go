package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func jsonRespOk(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func newTestClient(t *testing.T, f rtFunc, opts ...Option) *HTTPClient {
	t.Helper()

	hc := &http.Client{Transport: f}
	c, err := NewHTTPClient("http://example", append(opts, WithHTTPClient(hc))...)
	require.NoError(t, err, "create client")

	return c
}

func TestNewHTTPClient_EmptyBaseURL_Error(t *testing.T) {
	t.Parallel()

	_, err := NewHTTPClient("")
	require.Error(t, err, "expected error for empty baseURL")
}

func TestHTTPClient_do_GetQueryMerge(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodGet, r.Method)
		qs := r.URL.Query()
		require.Equal(t, "bar", qs.Get("foo"))
		tags := qs["tags"]
		assert.ElementsMatch(t, []string{"a", "b"}, tags)
		require.Equal(t, "1", qs.Get("extra"))

		return jsonRespOk(`{"status":"ok","retcode":0,"data":{},"message":""}`), nil
	})

	req := struct {
		Foo  string   `json:"foo"`
		Tags []string `json:"tags"`
	}{
		Foo:  "bar",
		Tags: []string{"a", "b"},
	}

	_, err := client.do(context.Background(), "/get", http.MethodGet, req, WithQuery("extra", "1"))
	require.NoError(t, err, "do GET")
}

func TestHTTPClient_do_PostHeadersBody(t *testing.T) {
	t.Parallel()

	token := "test-token"
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))
		require.Equal(t, "yes", r.Header.Get("X-Extra"))
		require.Equal(t, "ok", r.URL.Query().Get("q"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "alice", body["name"])

		return jsonRespOk(`{"status":"ok","retcode":0,"data":{"result":"ok"},"message":""}`), nil
	}, WithAccessToken(token))

	req := map[string]string{"name": "alice"}
	resp, err := client.do(
		context.Background(), "/post", http.MethodPost, req,
		WithHeader("X-Extra", "yes"), WithQuery("q", "ok"),
	)
	require.NoError(t, err, "do POST")
	require.Equal(t, entity.StatusOK, resp.Status)
	require.Equal(t, entity.RetcodeSuccess, resp.Retcode)
}

func TestHTTPClient_do_StatusNot2xx(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Body:       io.NopCloser(strings.NewReader("oops")),
			Header:     make(http.Header),
		}, nil
	})

	_, err := client.do(context.Background(), "/err", http.MethodGet, struct{}{})
	require.Error(t, err)
	require.ErrorContains(t, err, errHTTPStatus.Error()+": 500")
}

func TestHTTPClient_do_NonZeroRetcode(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ *http.Request) (*http.Response, error) {
		return jsonRespOk(`{"status":"ok","retcode":114,"data":{},"message":"bad"}`), nil
	})

	_, err := client.do(context.Background(), "/retcode", http.MethodPost, struct{}{})
	require.Error(t, err)

	var ae *entity.ActionError
	require.ErrorAs(t, err, &ae)
	require.Equal(t, entity.ActionResponseRetcode(114), ae.Retcode)
	require.Equal(t, "bad", ae.Message)
}

func TestHTTPClient_do_DecodeError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("not-json")),
			Header:     make(http.Header),
		}, nil
	})

	_, err := client.do(context.Background(), "/bad-json", http.MethodGet, struct{}{})
	require.Error(t, err)
	require.ErrorContains(t, err, "decode action response")
}

func TestHTTPClient_SendPrivateMsg_Success(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPost, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.InEpsilon(t, float64(123), body["user_id"], 1e-9)

		return jsonRespOk(`{"status":"ok","retcode":0,"data":{"message_id":321},"message":""}`), nil
	})

	req := &entity.SendPrivateMsgRequest{UserId: 123}
	resp, err := client.SendPrivateMsg(context.Background(), req, WithMethod(http.MethodPost))
	require.NoError(t, err, "SendPrivateMsg error")
	require.Equal(t, int64(321), resp.Data.MessageId)
}
