package entity

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errBoom = errors.New("boom")

type badMarshaler struct{}

func (badMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errBoom
}

func TestActionResponseToActionRawResponseNilReceiver(t *testing.T) {
	t.Parallel()

	var response *ActionResponse[struct{}]

	raw, err := response.ToActionRawResponse()
	require.NoError(t, err)
	require.Nil(t, raw)
}

func TestActionResponseToActionRawResponseWithData(t *testing.T) {
	t.Parallel()

	type sample struct {
		Foo   string `json:"foo"`
		Count int    `json:"count"`
	}

	data := &sample{
		Foo:   "bar",
		Count: 3,
	}

	response := &ActionResponse[sample]{
		Status:  StatusOK,
		Retcode: RetcodeSuccess,
		Data:    data,
		Message: "ok",
	}

	raw, err := response.ToActionRawResponse()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Equal(t, StatusOK, raw.Status)
	require.Equal(t, RetcodeSuccess, raw.Retcode)
	require.Equal(t, "ok", raw.Message)
	require.NotNil(t, raw.Data)

	var decoded sample

	err = json.Unmarshal(raw.Data, &decoded)
	require.NoError(t, err)
	require.Equal(t, *data, decoded)
}

func TestActionResponseToActionRawResponseWithNilData(t *testing.T) {
	t.Parallel()

	response := &ActionResponse[struct{}]{
		Status:  StatusAsync,
		Retcode: RetcodeAsync,
		Data:    nil,
		Message: "pending",
	}

	raw, err := response.ToActionRawResponse()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Equal(t, StatusAsync, raw.Status)
	require.Equal(t, RetcodeAsync, raw.Retcode)
	require.Equal(t, "pending", raw.Message)
	require.Nil(t, raw.Data)
}

func TestActionResponseToActionRawResponseMarshalError(t *testing.T) {
	t.Parallel()

	data := &badMarshaler{}
	response := &ActionResponse[badMarshaler]{
		Data: data,
	}

	raw, err := response.ToActionRawResponse()
	require.Error(t, err)
	require.ErrorContains(t, err, "fail to marshal Data field")
	require.Nil(t, raw)
}

func TestActionErrorError(t *testing.T) {
	t.Parallel()

	err := &ActionError{
		UrlPath: "/send_msg",
		Status:  StatusFailed,
		Retcode: 100,
		Message: "oops",
	}
	require.Equal(t, "action /send_msg failed: status=failed retcode=100 message=oops", err.Error())

	var nilErr *ActionError
	require.Empty(t, nilErr.Error())
}
