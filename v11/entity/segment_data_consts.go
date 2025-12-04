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

// CacheFlag 图片/语音/视频缓存标志位
type CacheFlag int

const (
	CacheFlagNo  CacheFlag = 0 // 不使用已缓存的文件
	CacheFlagYes CacheFlag = 1 // 使用已缓存的文件（默认）
)

// ProxyFlag 图片/语音/视频代理下载标志位
type ProxyFlag int

const (
	ProxyFlagNo  ProxyFlag = 0 // 不通过代理下载文件
	ProxyFlagYes ProxyFlag = 1 // 通过代理下载文件（默认）
)

// MagicFlag 语音变声标志位
type MagicFlag int

const (
	MagicFlagNo  MagicFlag = 0 // 不变声
	MagicFlagYes MagicFlag = 1 // 变声
)

// IgnoreFlag 匿名发消息标志位
type IgnoreFlag int

const (
	IgnoreFlagNo  IgnoreFlag = 0 // 无法匿名时不继续发送
	IgnoreFlagYes IgnoreFlag = 1 // 无法匿名时继续发送
)

// MusicType 音乐分享类型
type MusicType string

const (
	MusicTypeQQ      MusicType = "qq"     // QQ 音乐
	MusicTypeNetEase MusicType = "163"    // 网易云音乐
	MusicTypeXiami   MusicType = "xm"     // 虾米音乐
	MusicTypeCustom  MusicType = "custom" // 音乐自定义分享
)

// PokeSegmentData enums. from: https://github.com/mamoe/mirai/blob/f5eefae7ecee84d18a66afce3f89b89fe1584b78/mirai-core/src/commonMain/kotlin/net.mamoe.mirai/message/data/HummerMessage.kt#L49

func PokeSegmentDataPoke() *PokeSegmentData        { return &PokeSegmentData{"戳一戳", 1, -1} }     // 戳一戳
func PokeSegmentDataShowLove() *PokeSegmentData    { return &PokeSegmentData{"比心", 2, -1} }      // 比心
func PokeSegmentDataLike() *PokeSegmentData        { return &PokeSegmentData{"点赞", 3, -1} }      // 点赞
func PokeSegmentDataHeartbroken() *PokeSegmentData { return &PokeSegmentData{"心碎", 4, -1} }      // 心碎
func PokeSegmentDataSixSixSix() *PokeSegmentData   { return &PokeSegmentData{"666", 5, -1} }     // 666
func PokeSegmentDataFangDaZhao() *PokeSegmentData  { return &PokeSegmentData{"放大招", 6, -1} }     // 放大招
func PokeSegmentDataBaoBeiQiu() *PokeSegmentData   { return &PokeSegmentData{"宝贝球", 126, 2011} } // 宝贝球 (SVIP)
func PokeSegmentDataRose() *PokeSegmentData        { return &PokeSegmentData{"玫瑰花", 126, 2007} } // 玫瑰花 (SVIP)
func PokeSegmentDataZhaoHuanShu() *PokeSegmentData { return &PokeSegmentData{"召唤术", 126, 2006} } // 召唤术 (SVIP)
func PokeSegmentDataRangNiPi() *PokeSegmentData    { return &PokeSegmentData{"让你皮", 126, 2009} } // 让你皮 (SVIP)
func PokeSegmentDataJieYin() *PokeSegmentData      { return &PokeSegmentData{"结印", 126, 2005} }  // 结印 (SVIP)
func PokeSegmentDataShouLei() *PokeSegmentData     { return &PokeSegmentData{"手雷", 126, 2004} }  // 手雷 (SVIP)
func PokeSegmentDataGouYin() *PokeSegmentData      { return &PokeSegmentData{"勾引", 126, 2003} }  // 勾引
func PokeSegmentDataZhuaYiXia() *PokeSegmentData   { return &PokeSegmentData{"抓一下", 126, 2001} } // 抓一下 (SVIP)
func PokeSegmentDataSuiPing() *PokeSegmentData     { return &PokeSegmentData{"碎屏", 126, 2002} }  // 碎屏 (SVIP)
func PokeSegmentDataQiaoMen() *PokeSegmentData     { return &PokeSegmentData{"敲门", 126, 2002} }  // 敲门 (SVIP)
