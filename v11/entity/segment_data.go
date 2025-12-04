package entity

// SegmentData 表示 OneBot 消息的数据片段
// SegmentDataType 和数据片段唯一对应
type SegmentData interface {
	SegmentType() SegmentDataType
}

// Segment 表示一个通用消息段
// 所有具体的消息段类型都应该能转换为此类型
type Segment struct {
	Type SegmentDataType `json:"type"`
	Data SegmentData     `json:"data"`
}

func NewSegment(data SegmentData) *Segment {
	return &Segment{
		Type: data.SegmentType(),
		Data: data,
	}
}

// TextSegmentData 纯文本
// 消息段类型: text
// 支持发送、支持接收
type TextSegmentData struct {
	// 纯文本内容
	Text string `json:"text,omitempty"`
}

func (s *TextSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeText
}

// FaceSegmentData QQ 表情
// 消息段类型: face
// 支持发送、支持接收
type FaceSegmentData struct {
	// QQ 表情 ID | 可能的值: 见 [QQ 表情 ID 表](https://github.com/richardchien/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8)
	Id string `json:"id,omitempty"`
}

func (s *FaceSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeFace
}

// ImageSegmentData 图片
// 消息段类型: image
// 支持发送、支持接收
type ImageSegmentData struct {
	// 图片文件名
	File string `json:"file,omitempty"`
	// 图片类型，`flash` 表示闪照，无此参数表示普通图片 | 可能的值: flash
	Type ImageSegmentDataType `json:"type,omitempty"`
	// 图片 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *CacheFlag `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *ProxyFlag `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间，默认不超时
	Timeout *int64 `json:"timeout,omitempty"`
}

func (s *ImageSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeImage
}

// RecordSegmentData 语音
// 消息段类型: record
// 支持发送、支持接收
type RecordSegmentData struct {
	// 语音文件名
	File string `json:"file,omitempty"`
	// 发送时可选，默认 `0`，设置为 `1` 表示变声 | 可能的值: `0` `1`
	Magic *MagicFlag `json:"magic,omitempty"`
	// 语音 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *CacheFlag `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *ProxyFlag `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
	Timeout *int64 `json:"timeout,omitempty"`
}

func (s *RecordSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeRecord
}

// VideoSegmentData 短视频
// 消息段类型: video
// 支持发送、支持接收
type VideoSegmentData struct {
	// 视频文件名
	File string `json:"file,omitempty"`
	// 视频 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *CacheFlag `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *ProxyFlag `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
	Timeout *int64 `json:"timeout,omitempty"`
}

func (s *VideoSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeVideo
}

// AtSegmentData @某人
// 消息段类型: at
// 支持发送、支持接收
type AtSegmentData struct {
	// @的 QQ 号，`all` 表示全体成员 | 可能的值: QQ 号, all
	QQ string `json:"qq,omitempty"`
}

func (s *AtSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeAt
}

// RpsSegmentData 猜拳魔法表情
// 消息段类型: rps
// 支持发送、支持接收
type RpsSegmentData struct{}

func (s *RpsSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeRps
}

// DiceSegmentData 掷骰子魔法表情
// 消息段类型: dice
// 支持发送、支持接收
type DiceSegmentData struct{}

func (s *DiceSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeDice
}

// ShakeSegmentData 窗口抖动（戳一戳）
// 消息段类型: shake
// 支持发送
type ShakeSegmentData struct{}

func (s *ShakeSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeShake
}

// PokeSegmentData 戳一戳
// 消息段类型: poke
// 支持发送、支持接收
type PokeSegmentData struct {
	// 类型 | 可能的值: 见 segment_consts.go 中的 PokeSegmentData enums.
	Type string `json:"type,omitempty"`
	// ID | 可能的值: 同上
	Id int64 `json:"id,string,omitempty"`
	// 表情名
	Name int64 `json:"name,string,omitempty"`
}

func (s *PokeSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypePoke
}

// AnonymousSegmentData 匿名发消息
// 消息段类型: anonymous
// 支持发送
type AnonymousSegmentData struct {
	// 可选，表示无法匿名时是否继续发送 | 可能的值: `0`, `1`
	Ignore *IgnoreFlag `json:"ignore,omitempty"`
}

func (s *AnonymousSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeAnonymous
}

// ShareSegmentData 链接分享
// 消息段类型: share
// 支持发送、支持接收
type ShareSegmentData struct {
	// URL
	Url string `json:"url,omitempty"`
	// 标题
	Title string `json:"title,omitempty"`
	// 发送时可选，内容描述
	Content *string `json:"content,omitempty"`
	// 发送时可选，图片 URL
	Image *string `json:"image,omitempty"`
}

func (s *ShareSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeShare
}

// ContactSegmentData 推荐好友
// 消息段类型: contact
// 支持发送、支持接收
type ContactSegmentData struct {
	// 推荐好友 | 可能的值: qq, group
	Type ContactSegmentDataType `json:"type,omitempty"`
	// 被推荐人的 QQ 号/群号
	Id string `json:"id,omitempty"`
}

func (s *ContactSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeContact
}

// LocationSegmentData 位置
// 消息段类型: location
// 支持发送、支持接收
type LocationSegmentData struct {
	// 纬度
	Lat string `json:"lat,omitempty"`
	// 经度
	Lon string `json:"lon,omitempty"`
	// 发送时可选，标题
	Title *string `json:"title,omitempty"`
	// 发送时可选，内容描述
	Content *string `json:"content,omitempty"`
}

func (s *LocationSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeLocation
}

// MusicSegmentData 音乐自定义分享
// 消息段类型: music
// 支持发送
type MusicSegmentData struct {
	// 表示使用 QQ 音乐、网易云音乐、虾米音乐或音乐自定义分享 | 可能的值: qq, 163, xm, custom
	Type MusicType `json:"type,omitempty"`

	// 歌曲 ID (使用 QQ 音乐、网易云音乐、虾米音乐)
	Id string `json:"id,omitempty"`

	// 音乐自定义分享
	// 点击后跳转目标 URL
	Url string `json:"url,omitempty"`
	// 音乐 URL
	Audio string `json:"audio,omitempty"`
	// 标题
	Title string `json:"title,omitempty"`
	// 发送时可选，内容描述
	Content *string `json:"content,omitempty"`
	// 发送时可选，图片 URL
	Image *string `json:"image,omitempty"`
}

func (s *MusicSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeMusic
}

// ReplySegmentData 回复
// 消息段类型: reply
// 支持发送、支持接收
type ReplySegmentData struct {
	// 回复时引用的消息 ID
	Id string `json:"id,omitempty"`
}

func (s *ReplySegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeReply
}

// ForwardSegmentData 合并转发
// 消息段类型: forward
// 支持接收
type ForwardSegmentData struct {
	// 合并转发 ID，需通过 get_forward_msg-获取合并转发消息 获取具体内容
	Id string `json:"id,omitempty"`
}

func (s *ForwardSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeForward
}

// NodeSegmentData
// Id 不为空时:
// 合并转发节点
//
//	消息段类型: node
//	支持发送
//
// Id 为空时:
// 合并转发自定义节点
// 消息段类型: node
// 支持发送、支持接收
type NodeSegmentData struct {
	// 转发的消息 ID
	Id string `json:"id,omitempty"`

	// 发送者 QQ 号
	UserId string `json:"user_id,omitempty"`
	// 发送者昵称
	Nickname string `json:"nickname,omitempty"`
	// 消息内容，支持发送消息时的 `message` 数据类型
	Content *MessageValue `json:"content,omitempty"`
}

func (s *NodeSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeNode
}

// XmlSegmentData XML 消息
// 消息段类型: xml
// 支持发送、支持接收
type XmlSegmentData struct {
	// XML 内容
	Data string `json:"data,omitempty"`
}

func (s *XmlSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeXml
}

// JsonSegmentData JSON 消息
// 消息段类型: json
// 支持发送、支持接收
type JsonSegmentData struct {
	// JSON 内容
	Data string `json:"data,omitempty"`
}

func (s *JsonSegmentData) SegmentType() SegmentDataType {
	return SegmentDataTypeJson
}
