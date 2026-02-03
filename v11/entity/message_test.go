package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type badSegmentData struct {
	F func()
}

func (b *badSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeText
}

func TestMessageValueUnmarshalJSONString(t *testing.T) {
	t.Parallel()

	var value MessageValue

	err := json.Unmarshal([]byte(`"hello"`), &value)
	require.NoError(t, err)
	require.Equal(t, MessageValueTypeString, value.Type)
	require.Equal(t, "hello", value.StringValue)
	require.Nil(t, value.ArrayValue)
}

func TestMessageValueUnmarshalJSONArray(t *testing.T) {
	t.Parallel()

	var value MessageValue

	err := json.Unmarshal([]byte(`[{"type":"text"},{"type":"face"}]`), &value)
	require.NoError(t, err)
	require.Equal(t, MessageValueTypeArray, value.Type)
	require.Len(t, value.ArrayValue, 2)
	require.Equal(t, SegmentDataTypeText, value.ArrayValue[0].Type)
	require.Equal(t, SegmentDataTypeFace, value.ArrayValue[1].Type)
	require.Nil(t, value.ArrayValue[0].Data)
}

func TestMessageValueMarshalJSONString(t *testing.T) {
	t.Parallel()

	value := MessageValue{
		Type:        MessageValueTypeString,
		StringValue: "world",
	}
	data, err := json.Marshal(&value)
	require.NoError(t, err)
	require.Equal(t, `"world"`, string(data))
}

func TestMessageValueMarshalJSONArray(t *testing.T) {
	t.Parallel()

	value := MessageValue{
		Type: MessageValueTypeArray,
		ArrayValue: []*Segment{
			{
				Type: SegmentDataTypeText,
				Data: nil,
			},
		},
	}
	data, err := json.Marshal(&value)
	require.NoError(t, err)
	require.JSONEq(t, `[{"type":"text","data":null}]`, string(data))
}

func TestMessageValueMarshalJSONNull(t *testing.T) {
	t.Parallel()

	var value MessageValue

	data, err := json.Marshal(&value)
	require.NoError(t, err)
	require.Equal(t, `null`, string(data))
}

func TestMessageValueMarshalJSONArrayError(t *testing.T) {
	t.Parallel()

	value := MessageValue{
		Type: MessageValueTypeArray,
		ArrayValue: []*Segment{
			{
				Type: SegmentDataTypeText,
				Data: &badSegmentData{F: func() {}},
			},
		},
	}
	_, err := json.Marshal(&value)
	require.Error(t, err)
	require.ErrorContains(t, err, "marshal array value")
}
