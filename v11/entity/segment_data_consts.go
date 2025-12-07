package entity

type SegmentDataType string

const (
	SegmentDataTypeText      SegmentDataType = "text"
	SegmentDataTypeFace      SegmentDataType = "face"
	SegmentDataTypeImage     SegmentDataType = "image"
	SegmentDataTypeRecord    SegmentDataType = "record"
	SegmentDataTypeVideo     SegmentDataType = "video"
	SegmentDataTypeAt        SegmentDataType = "at"
	SegmentDataTypeRps       SegmentDataType = "rps"
	SegmentDataTypeDice      SegmentDataType = "dice"
	SegmentDataTypeShake     SegmentDataType = "shake"
	SegmentDataTypePoke      SegmentDataType = "poke"
	SegmentDataTypeAnonymous SegmentDataType = "anonymous"
	SegmentDataTypeShare     SegmentDataType = "share"
	SegmentDataTypeContact   SegmentDataType = "contact"
	SegmentDataTypeLocation  SegmentDataType = "location"
	SegmentDataTypeMusic     SegmentDataType = "music"
	SegmentDataTypeReply     SegmentDataType = "reply"
	SegmentDataTypeForward   SegmentDataType = "forward"
	SegmentDataTypeNode      SegmentDataType = "node"
	SegmentDataTypeXml       SegmentDataType = "xml"
	SegmentDataTypeJson      SegmentDataType = "json"
)

type ImageSegmentDataType string

const (
	ImageSegmentDataTypeCommon ImageSegmentDataType = ""
	ImageSegmentDataTypeFlash  ImageSegmentDataType = "flash"
)

type ContactSegmentDataType string

const (
	ContactSegmentDataTypeQQ    ContactSegmentDataType = "qq"
	ContactSegmentDataTypeGroup ContactSegmentDataType = "group"
)

// CacheFlag 图片/语音/视频缓存标志位.
type CacheFlag int

const (
	CacheFlagNo  CacheFlag = 0 // 不使用已缓存的文件
	CacheFlagYes CacheFlag = 1 // 使用已缓存的文件（默认）
)

// ProxyFlag 图片/语音/视频代理下载标志位.
type ProxyFlag int

const (
	ProxyFlagNo  ProxyFlag = 0 // 不通过代理下载文件
	ProxyFlagYes ProxyFlag = 1 // 通过代理下载文件（默认）
)

// MagicFlag 语音变声标志位.
type MagicFlag int

const (
	MagicFlagNo  MagicFlag = 0 // 不变声
	MagicFlagYes MagicFlag = 1 // 变声
)

// IgnoreFlag 匿名发消息标志位.
type IgnoreFlag int

const (
	IgnoreFlagNo  IgnoreFlag = 0 // 无法匿名时不继续发送
	IgnoreFlagYes IgnoreFlag = 1 // 无法匿名时继续发送
)

// MusicType 音乐分享类型.
type MusicType string

const (
	MusicTypeQQ      MusicType = "qq"     // QQ 音乐
	MusicTypeNetEase MusicType = "163"    // 网易云音乐
	MusicTypeXiami   MusicType = "xm"     // 虾米音乐
	MusicTypeCustom  MusicType = "custom" // 音乐自定义分享
)

//nolint:lll
// PokeSegmentData enums.
// from: https://github.com/mamoe/mirai/blob/f5eefae7ecee84d18a66afce3f89b89fe1584b78/mirai-core/src/commonMain/kotlin/net.mamoe.mirai/message/data/HummerMessage.kt#L49

const (
	PokeSegmentDataTypePoke        = "戳一戳" //nolint:gosmopolitan
	PokeSegmentDataTypeShowLove    = "比心"  //nolint:gosmopolitan
	PokeSegmentDataTypeLike        = "点赞"  //nolint:gosmopolitan
	PokeSegmentDataTypeHeartbroken = "心碎"  //nolint:gosmopolitan
	PokeSegmentDataTypeSixSixSix   = "666"
	PokeSegmentDataTypeFangDaZhao  = "放大招" //nolint:gosmopolitan
	PokeSegmentDataTypeBaoBeiQiu   = "宝贝球" //nolint:gosmopolitan
	PokeSegmentDataTypeRose        = "玫瑰花" //nolint:gosmopolitan
	PokeSegmentDataTypeZhaoHuanShu = "召唤术" //nolint:gosmopolitan
	PokeSegmentDataTypeRangNiPi    = "让你皮" //nolint:gosmopolitan
	PokeSegmentDataTypeJieYin      = "结印"  //nolint:gosmopolitan
	PokeSegmentDataTypeShouLei     = "手雷"  //nolint:gosmopolitan
	PokeSegmentDataTypeGouYin      = "勾引"  //nolint:gosmopolitan
	PokeSegmentDataTypeZhuaYiXia   = "抓一下" //nolint:gosmopolitan
	PokeSegmentDataTypeSuiPing     = "碎屏"  //nolint:gosmopolitan
	PokeSegmentDataTypeQiaoMen     = "敲门"  //nolint:gosmopolitan
)

// PokeSegmentDataPoke 戳一戳.
func PokeSegmentDataPoke() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypePoke, 1, -1}
}

// PokeSegmentDataShowLove 比心.
func PokeSegmentDataShowLove() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeShowLove, 2, -1}
}

// PokeSegmentDataLike 点赞.
func PokeSegmentDataLike() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeLike, 3, -1}
}

// PokeSegmentDataHeartbroken 心碎.
func PokeSegmentDataHeartbroken() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeHeartbroken, 4, -1}
}

// PokeSegmentDataSixSixSix 666.
func PokeSegmentDataSixSixSix() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeSixSixSix, 5, -1}
}

// PokeSegmentDataFangDaZhao 放大招.
func PokeSegmentDataFangDaZhao() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeFangDaZhao, 6, -1}
}

// PokeSegmentDataBaoBeiQiu 宝贝球 (SVIP).
func PokeSegmentDataBaoBeiQiu() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeBaoBeiQiu, 126, 2011}
}

// PokeSegmentDataRose 玫瑰花 (SVIP).
func PokeSegmentDataRose() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeRose, 126, 2007}
}

// PokeSegmentDataZhaoHuanShu 召唤术 (SVIP).
func PokeSegmentDataZhaoHuanShu() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeZhaoHuanShu, 126, 2006}
}

// PokeSegmentDataRangNiPi 让你皮 (SVIP).
func PokeSegmentDataRangNiPi() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeRangNiPi, 126, 2009}
}

// PokeSegmentDataJieYin 结印 (SVIP).
func PokeSegmentDataJieYin() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeJieYin, 126, 2005}
}

// PokeSegmentDataShouLei 手雷 (SVIP).
func PokeSegmentDataShouLei() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeShouLei, 126, 2004}
}

// PokeSegmentDataGouYin 勾引.
func PokeSegmentDataGouYin() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeGouYin, 126, 2003}
}

// PokeSegmentDataZhuaYiXia 抓一下 (SVIP).
func PokeSegmentDataZhuaYiXia() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeZhuaYiXia, 126, 2001}
}

// PokeSegmentDataSuiPing 碎屏 (SVIP).
func PokeSegmentDataSuiPing() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeSuiPing, 126, 2002}
}

// PokeSegmentDataQiaoMen 敲门 (SVIP).
func PokeSegmentDataQiaoMen() *PokeSegmentData {
	return &PokeSegmentData{PokeSegmentDataTypeQiaoMen, 126, 2002}
}
