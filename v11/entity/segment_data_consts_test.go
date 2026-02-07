package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentDataTypeConstants(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		got  SegmentDataType
		want string
	}{
		{name: "text", got: SegmentDataTypeText, want: "text"},
		{name: "face", got: SegmentDataTypeFace, want: "face"},
		{name: "image", got: SegmentDataTypeImage, want: "image"},
		{name: "record", got: SegmentDataTypeRecord, want: "record"},
		{name: "video", got: SegmentDataTypeVideo, want: "video"},
		{name: "at", got: SegmentDataTypeAt, want: "at"},
		{name: "rps", got: SegmentDataTypeRps, want: "rps"},
		{name: "dice", got: SegmentDataTypeDice, want: "dice"},
		{name: "shake", got: SegmentDataTypeShake, want: "shake"},
		{name: "poke", got: SegmentDataTypePoke, want: "poke"},
		{name: "anonymous", got: SegmentDataTypeAnonymous, want: "anonymous"},
		{name: "share", got: SegmentDataTypeShare, want: "share"},
		{name: "contact", got: SegmentDataTypeContact, want: "contact"},
		{name: "location", got: SegmentDataTypeLocation, want: "location"},
		{name: "music", got: SegmentDataTypeMusic, want: "music"},
		{name: "reply", got: SegmentDataTypeReply, want: "reply"},
		{name: "forward", got: SegmentDataTypeForward, want: "forward"},
		{name: "node", got: SegmentDataTypeNode, want: "node"},
		{name: "xml", got: SegmentDataTypeXml, want: "xml"},
		{name: "json", got: SegmentDataTypeJson, want: "json"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.want, string(testCase.got))
		})
	}
}

func TestImageSegmentDataTypeConstants(t *testing.T) {
	t.Parallel()

	require.Empty(t, string(ImageSegmentDataTypeCommon))
	require.Equal(t, "flash", string(ImageSegmentDataTypeFlash))
}

func TestContactSegmentDataTypeConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, "qq", string(ContactSegmentDataTypeQQ))
	require.Equal(t, "group", string(ContactSegmentDataTypeGroup))
}

func TestFlagConstants(t *testing.T) {
	t.Parallel()

	require.Equal(t, CacheFlagNo, CacheFlag(0))
	require.Equal(t, CacheFlagYes, CacheFlag(1))

	require.Equal(t, ProxyFlagNo, ProxyFlag(0))
	require.Equal(t, ProxyFlagYes, ProxyFlag(1))

	require.Equal(t, MagicFlagNo, MagicFlag(0))
	require.Equal(t, MagicFlagYes, MagicFlag(1))

	require.Equal(t, IgnoreFlagNo, IgnoreFlag(0))
	require.Equal(t, IgnoreFlagYes, IgnoreFlag(1))
}

func TestMusicTypeConstants(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		got  MusicType
		want string
	}{
		{name: "qq", got: MusicTypeQQ, want: "qq"},
		{name: "netease", got: MusicTypeNetEase, want: "163"},
		{name: "xiami", got: MusicTypeXiami, want: "xm"},
		{name: "custom", got: MusicTypeCustom, want: "custom"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.want, string(testCase.got))
		})
	}
}

func TestPokeSegmentDataTypeConstants(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		got  string
		want string
	}{
		{name: "poke", got: PokeSegmentDataTypePoke, want: "戳一戳"},              //nolint:gosmopolitan
		{name: "show_love", got: PokeSegmentDataTypeShowLove, want: "比心"},      //nolint:gosmopolitan
		{name: "like", got: PokeSegmentDataTypeLike, want: "点赞"},               //nolint:gosmopolitan
		{name: "heartbroken", got: PokeSegmentDataTypeHeartbroken, want: "心碎"}, //nolint:gosmopolitan
		{name: "six_six_six", got: PokeSegmentDataTypeSixSixSix, want: "666"},
		{name: "fang_da_zhao", got: PokeSegmentDataTypeFangDaZhao, want: "放大招"},   //nolint:gosmopolitan
		{name: "bao_bei_qiu", got: PokeSegmentDataTypeBaoBeiQiu, want: "宝贝球"},     //nolint:gosmopolitan
		{name: "rose", got: PokeSegmentDataTypeRose, want: "玫瑰花"},                 //nolint:gosmopolitan
		{name: "zhao_huan_shu", got: PokeSegmentDataTypeZhaoHuanShu, want: "召唤术"}, //nolint:gosmopolitan
		{name: "rang_ni_pi", got: PokeSegmentDataTypeRangNiPi, want: "让你皮"},       //nolint:gosmopolitan
		{name: "jie_yin", got: PokeSegmentDataTypeJieYin, want: "结印"},             //nolint:gosmopolitan
		{name: "shou_lei", got: PokeSegmentDataTypeShouLei, want: "手雷"},           //nolint:gosmopolitan
		{name: "gou_yin", got: PokeSegmentDataTypeGouYin, want: "勾引"},             //nolint:gosmopolitan
		{name: "zhua_yi_xia", got: PokeSegmentDataTypeZhuaYiXia, want: "抓一下"},     //nolint:gosmopolitan
		{name: "sui_ping", got: PokeSegmentDataTypeSuiPing, want: "碎屏"},           //nolint:gosmopolitan
		{name: "qiao_men", got: PokeSegmentDataTypeQiaoMen, want: "敲门"},           //nolint:gosmopolitan
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.want, testCase.got)
		})
	}
}

type pokeSegmentConstructorCase struct {
	name string
	fn   func() *PokeSegmentData
	want PokeSegmentData
}

func assertPokeSegmentDataConstructors(t *testing.T, cases []pokeSegmentConstructorCase) {
	t.Helper()

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := testCase.fn()
			require.NotNil(t, got)
			require.Equal(t, testCase.want, *got)
		})
	}
}

func TestPokeSegmentDataConstructorsBasic(t *testing.T) {
	t.Parallel()

	cases := []pokeSegmentConstructorCase{
		{
			name: "poke",
			fn:   PokeSegmentDataPoke,
			want: PokeSegmentData{Type: PokeSegmentDataTypePoke, Id: 1, Name: -1},
		},
		{
			name: "show_love",
			fn:   PokeSegmentDataShowLove,
			want: PokeSegmentData{Type: PokeSegmentDataTypeShowLove, Id: 2, Name: -1},
		},
		{
			name: "like",
			fn:   PokeSegmentDataLike,
			want: PokeSegmentData{Type: PokeSegmentDataTypeLike, Id: 3, Name: -1},
		},
		{
			name: "heartbroken",
			fn:   PokeSegmentDataHeartbroken,
			want: PokeSegmentData{Type: PokeSegmentDataTypeHeartbroken, Id: 4, Name: -1},
		},
		{
			name: "six_six_six",
			fn:   PokeSegmentDataSixSixSix,
			want: PokeSegmentData{Type: PokeSegmentDataTypeSixSixSix, Id: 5, Name: -1},
		},
		{
			name: "fang_da_zhao",
			fn:   PokeSegmentDataFangDaZhao,
			want: PokeSegmentData{Type: PokeSegmentDataTypeFangDaZhao, Id: 6, Name: -1},
		},
		{
			name: "bao_bei_qiu",
			fn:   PokeSegmentDataBaoBeiQiu,
			want: PokeSegmentData{Type: PokeSegmentDataTypeBaoBeiQiu, Id: 126, Name: 2011},
		},
		{
			name: "rose",
			fn:   PokeSegmentDataRose,
			want: PokeSegmentData{Type: PokeSegmentDataTypeRose, Id: 126, Name: 2007},
		},
	}

	assertPokeSegmentDataConstructors(t, cases)
}

func TestPokeSegmentDataConstructorsExtended(t *testing.T) {
	t.Parallel()

	cases := []pokeSegmentConstructorCase{
		{
			name: "zhao_huan_shu",
			fn:   PokeSegmentDataZhaoHuanShu,
			want: PokeSegmentData{Type: PokeSegmentDataTypeZhaoHuanShu, Id: 126, Name: 2006},
		},
		{
			name: "rang_ni_pi",
			fn:   PokeSegmentDataRangNiPi,
			want: PokeSegmentData{Type: PokeSegmentDataTypeRangNiPi, Id: 126, Name: 2009},
		},
		{
			name: "jie_yin",
			fn:   PokeSegmentDataJieYin,
			want: PokeSegmentData{Type: PokeSegmentDataTypeJieYin, Id: 126, Name: 2005},
		},
		{
			name: "shou_lei",
			fn:   PokeSegmentDataShouLei,
			want: PokeSegmentData{Type: PokeSegmentDataTypeShouLei, Id: 126, Name: 2004},
		},
		{
			name: "gou_yin",
			fn:   PokeSegmentDataGouYin,
			want: PokeSegmentData{Type: PokeSegmentDataTypeGouYin, Id: 126, Name: 2003},
		},
		{
			name: "zhua_yi_xia",
			fn:   PokeSegmentDataZhuaYiXia,
			want: PokeSegmentData{Type: PokeSegmentDataTypeZhuaYiXia, Id: 126, Name: 2001},
		},
		{
			name: "sui_ping",
			fn:   PokeSegmentDataSuiPing,
			want: PokeSegmentData{Type: PokeSegmentDataTypeSuiPing, Id: 126, Name: 2002},
		},
		{
			name: "qiao_men",
			fn:   PokeSegmentDataQiaoMen,
			want: PokeSegmentData{Type: PokeSegmentDataTypeQiaoMen, Id: 126, Name: 2002},
		},
	}

	assertPokeSegmentDataConstructors(t, cases)
}
