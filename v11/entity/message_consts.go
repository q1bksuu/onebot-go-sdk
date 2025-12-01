package entity

type MessageValueType string

const (
	MessageValueTypeString MessageValueType = "string"
	MessageValueTypeArray  MessageValueType = "array"
)

type MessageSegmentType string

const (
	MessageSegmentTypeText      MessageSegmentType = "text"
	MessageSegmentTypeFace      MessageSegmentType = "face"
	MessageSegmentTypeImage     MessageSegmentType = "image"
	MessageSegmentTypeRecord    MessageSegmentType = "record"
	MessageSegmentTypeVideo     MessageSegmentType = "video"
	MessageSegmentTypeAt        MessageSegmentType = "at"
	MessageSegmentTypeRps       MessageSegmentType = "rps"
	MessageSegmentTypeDice      MessageSegmentType = "dice"
	MessageSegmentTypeShake     MessageSegmentType = "shake"
	MessageSegmentTypePoke      MessageSegmentType = "poke"
	MessageSegmentTypeAnonymous MessageSegmentType = "anonymous"
	MessageSegmentTypeShare     MessageSegmentType = "share"
	MessageSegmentTypeContact   MessageSegmentType = "contact"
	MessageSegmentTypeLocation  MessageSegmentType = "location"
	MessageSegmentTypeMusic     MessageSegmentType = "music"
	MessageSegmentTypeReply     MessageSegmentType = "reply"
	MessageSegmentTypeForward   MessageSegmentType = "forward"
	MessageSegmentTypeNode      MessageSegmentType = "node"
	MessageSegmentTypeXml       MessageSegmentType = "xml"
	MessageSegmentTypeJson      MessageSegmentType = "json"
)

type ImageSegmentType string

const ImageSegmentTypeFlash ImageSegmentType = "flash"

type ContactSegmentType string

const (
	ContactSegmentTypeQQ    ContactSegmentType = "qq"
	ContactSegmentTypeGroup ContactSegmentType = "group"
)

// PokeSegment enums.

func PokeSegmentPoke() *PokeSegment        { return &PokeSegment{"戳一戳", 1, -1} }     // 戳一戳
func PokeSegmentShowLove() *PokeSegment    { return &PokeSegment{"比心", 2, -1} }      // 比心
func PokeSegmentLike() *PokeSegment        { return &PokeSegment{"点赞", 3, -1} }      // 点赞
func PokeSegmentHeartbroken() *PokeSegment { return &PokeSegment{"心碎", 4, -1} }      // 心碎
func PokeSegmentSixSixSix() *PokeSegment   { return &PokeSegment{"666", 5, -1} }     // 666
func PokeSegmentFangDaZhao() *PokeSegment  { return &PokeSegment{"放大招", 6, -1} }     // 放大招
func PokeSegmentBaoBeiQiu() *PokeSegment   { return &PokeSegment{"宝贝球", 126, 2011} } // 宝贝球 (SVIP)
func PokeSegmentRose() *PokeSegment        { return &PokeSegment{"玫瑰花", 126, 2007} } // 玫瑰花 (SVIP)
func PokeSegmentZhaoHuanShu() *PokeSegment { return &PokeSegment{"召唤术", 126, 2006} } // 召唤术 (SVIP)
func PokeSegmentRangNiPi() *PokeSegment    { return &PokeSegment{"让你皮", 126, 2009} } // 让你皮 (SVIP)
func PokeSegmentJieYin() *PokeSegment      { return &PokeSegment{"结印", 126, 2005} }  // 结印 (SVIP)
func PokeSegmentShouLei() *PokeSegment     { return &PokeSegment{"手雷", 126, 2004} }  // 手雷 (SVIP)
func PokeSegmentGouYin() *PokeSegment      { return &PokeSegment{"勾引", 126, 2003} }  // 勾引
func PokeSegmentZhuaYiXia() *PokeSegment   { return &PokeSegment{"抓一下", 126, 2001} } // 抓一下 (SVIP)
func PokeSegmentSuiPing() *PokeSegment     { return &PokeSegment{"碎屏", 126, 2002} }  // 碎屏 (SVIP)
func PokeSegmentQiaoMen() *PokeSegment     { return &PokeSegment{"敲门", 126, 2002} }  // 敲门 (SVIP)
