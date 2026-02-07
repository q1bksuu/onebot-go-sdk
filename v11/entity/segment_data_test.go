package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentDataSegmentType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		data SegmentData
		want SegmentDataType
	}{
		{name: "text", data: &TextSegmentData{}, want: SegmentDataTypeText},
		{name: "face", data: &FaceSegmentData{}, want: SegmentDataTypeFace},
		{name: "image", data: &ImageSegmentData{}, want: SegmentDataTypeImage},
		{name: "record", data: &RecordSegmentData{}, want: SegmentDataTypeRecord},
		{name: "video", data: &VideoSegmentData{}, want: SegmentDataTypeVideo},
		{name: "at", data: &AtSegmentData{}, want: SegmentDataTypeAt},
		{name: "rps", data: &RpsSegmentData{}, want: SegmentDataTypeRps},
		{name: "dice", data: &DiceSegmentData{}, want: SegmentDataTypeDice},
		{name: "shake", data: &ShakeSegmentData{}, want: SegmentDataTypeShake},
		{name: "poke", data: &PokeSegmentData{}, want: SegmentDataTypePoke},
		{name: "anonymous", data: &AnonymousSegmentData{}, want: SegmentDataTypeAnonymous},
		{name: "share", data: &ShareSegmentData{}, want: SegmentDataTypeShare},
		{name: "contact", data: &ContactSegmentData{}, want: SegmentDataTypeContact},
		{name: "location", data: &LocationSegmentData{}, want: SegmentDataTypeLocation},
		{name: "music", data: &MusicSegmentData{}, want: SegmentDataTypeMusic},
		{name: "reply", data: &ReplySegmentData{}, want: SegmentDataTypeReply},
		{name: "forward", data: &ForwardSegmentData{}, want: SegmentDataTypeForward},
		{name: "node", data: &NodeSegmentData{}, want: SegmentDataTypeNode},
		{name: "xml", data: &XmlSegmentData{}, want: SegmentDataTypeXml},
		{name: "json", data: &JsonSegmentData{}, want: SegmentDataTypeJson},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, tc.data.SegmentType())
		})
	}
}

func TestNewSegment(t *testing.T) {
	t.Parallel()

	data := &TextSegmentData{Text: "hello"}
	segment := NewSegment(data)

	require.NotNil(t, segment)
	require.Equal(t, SegmentDataTypeText, segment.Type)
	got, ok := segment.Data.(*TextSegmentData)
	require.True(t, ok)
	require.Same(t, data, got)
}

func TestSegmentMarshalJSON(t *testing.T) {
	t.Parallel()

	segment := NewSegment(&TextSegmentData{Text: "hi"})
	data, err := json.Marshal(segment)
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"text","data":{"text":"hi"}}`, string(data))
}

func TestPokeSegmentDataMarshalJSON(t *testing.T) {
	t.Parallel()

	data := &PokeSegmentData{
		Type: "poke",
		Id:   12,
		Name: 34,
	}
	payload, err := json.Marshal(data)
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"poke","id":"12","name":"34"}`, string(payload))
}
