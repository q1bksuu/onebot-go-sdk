package entity

import (
	"bytes"
	"encoding/json"
)

type MessageSegmentData interface {
	SegmentType() MessageSegmentType
}

// MessageSegment 表示一个通用消息段
// 所有具体的消息段类型都应该能转换为此类型
type MessageSegment struct {
	Type MessageSegmentType `json:"type"`
	Data MessageSegmentData `json:"data"`
}

func NewMessageSegment(data MessageSegmentData) *MessageSegment {
	return &MessageSegment{
		Type: data.SegmentType(),
		Data: data,
	}
}

// MessageValue 表示 OneBot message 字段的值
// 可以是纯文本字符串或消息段数组
type MessageValue struct {
	// 如果 Type 为 "string"，则使用 StringValue
	// 如果 Type 为 "array"，则使用 ArrayValue
	Type        MessageValueType  `json:"-"`
	StringValue string            `json:"-"`
	ArrayValue  []*MessageSegment `json:"-"`
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
// 用于在反序列化时自动选择正确的类型
func (m *MessageValue) UnmarshalJSON(data []byte) error {
	// 首先尝试作为字符串解析
	var str string
	if bytes.HasPrefix(data, []byte{'"'}) {
		if err := json.Unmarshal(data, &str); err == nil {
			m.Type = MessageValueTypeString
			m.StringValue = str
			return nil
		}
	}

	// 然后尝试作为数组解析
	var arr []*MessageSegment
	if bytes.HasPrefix(data, []byte{'['}) {
		if err := json.Unmarshal(data, &arr); err == nil {
			m.Type = MessageValueTypeArray
			m.ArrayValue = arr
			return nil
		}
	}

	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
// 用于在序列化时正确地输出值
func (m *MessageValue) MarshalJSON() ([]byte, error) {
	if m.Type == MessageValueTypeString {
		return json.Marshal(m.StringValue)
	}
	if m.Type == MessageValueTypeArray {
		return json.Marshal(m.ArrayValue)
	}
	return []byte{'n', 'u', 'l', 'l'}, nil
}

// TextSegment 纯文本
// 消息段类型: text
// 支持发送、支持接收
type TextSegment struct {
	// 纯文本内容
	Text string `json:"text,omitempty"`
}

func (s *TextSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeText
}

// FaceSegment QQ 表情
// 消息段类型: face
// 支持发送、支持接收
type FaceSegment struct {
	// QQ 表情 ID | 可能的值: 见 [QQ 表情 ID 表](https://github.com/richardchien/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8)
	Id string `json:"id,omitempty"`
}

func (s *FaceSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeFace
}

// ImageSegment 图片
// 消息段类型: image
// 支持发送、支持接收
type ImageSegment struct {
	// 图片文件名
	File string `json:"file,omitempty"`
	// 图片类型，`flash` 表示闪照，无此参数表示普通图片 | 可能的值: flash
	Type ImageSegmentType `json:"type,omitempty"`
	// 图片 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *int `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *int `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间，默认不超时
	Timeout int64 `json:"timeout,omitempty"`
}

func (s *ImageSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeImage
}

// RecordSegment 语音
// 消息段类型: record
// 支持发送、支持接收
type RecordSegment struct {
	// 语音文件名
	File string `json:"file,omitempty"`
	// 发送时可选，默认 `0`，设置为 `1` 表示变声 | 可能的值: `0` `1`
	Magic *int `json:"magic,omitempty"`
	// 语音 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *int `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *int `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
	Timeout int64 `json:"timeout,omitempty"`
}

func (s *RecordSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeRecord
}

// VideoSegment 短视频
// 消息段类型: video
// 支持发送、支持接收
type VideoSegment struct {
	// 视频文件名
	File string `json:"file,omitempty"`
	// 视频 URL
	Url string `json:"url,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	Cache *int `json:"cache,omitempty"`
	// 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	Proxy *int `json:"proxy,omitempty"`
	// 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
	Timeout int64 `json:"timeout,omitempty"`
}

func (s *VideoSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeVideo
}

// AtSegment @某人
// 消息段类型: at
// 支持发送、支持接收
type AtSegment struct {
	// @的 QQ 号，`all` 表示全体成员 | 可能的值: QQ 号, all
	Qq string `json:"qq,omitempty"`
}

func (s *AtSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeAt
}

// RpsSegment 猜拳魔法表情
// 消息段类型: rps
// 支持发送、支持接收
type RpsSegment struct {
	// 无参数
}

func (s *RpsSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeRps
}

// DiceSegment 掷骰子魔法表情
// 消息段类型: dice
// 支持发送、支持接收
type DiceSegment struct {
	// 无参数
}

func (s *DiceSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeDice
}

// ShakeSegment 窗口抖动（戳一戳）
// 消息段类型: shake
// 支持发送
type ShakeSegment struct {
	// 无参数
}

func (s *ShakeSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeShake
}

// PokeSegment 戳一戳
// 消息段类型: poke
// 支持发送、支持接收
type PokeSegment struct {
	// 类型 | 可能的值: 见 [Mirai 的 PokeMessage 类](https://github.com/mamoe/mirai/blob/f5eefae7ecee84d18a66afce3f89b89fe1584b78/mirai-core/src/commonMain/kotlin/net.mamoe.mirai/message/data/HummerMessage.kt#L49)
	Type string `json:"type,omitempty"`
	// ID | 可能的值: 同上
	Id int64 `json:"id,omitempty"`
	// 表情名
	Name int64 `json:"name,omitempty"`
}

func (s *PokeSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypePoke
}

// AnonymousSegment 匿名发消息
// 消息段类型: anonymous
// 支持发送
type AnonymousSegment struct {
	// 可选，表示无法匿名时是否继续发送 | 可能的值: `0`, `1`
	Ignore *int `json:"ignore,omitempty"`
}

func (s *AnonymousSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeAnonymous
}

// ShareSegment 链接分享
// 消息段类型: share
// 支持发送、支持接收
type ShareSegment struct {
	// URL
	Url string `json:"url,omitempty"`
	// 标题
	Title string `json:"title,omitempty"`
	// 发送时可选，内容描述
	Content *string `json:"content,omitempty"`
	// 发送时可选，图片 URL
	Image *string `json:"image,omitempty"`
}

func (s *ShareSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeShare
}

// ContactSegment 推荐好友
// 消息段类型: contact
// 支持发送、支持接收
type ContactSegment struct {
	// 推荐好友 | 可能的值: qq, group
	Type ContactSegmentType `json:"type,omitempty"`
	// 被推荐人的 QQ 号/群号
	Id string `json:"id,omitempty"`
}

func (s *ContactSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeContact
}

// LocationSegment 位置
// 消息段类型: location
// 支持发送、支持接收
type LocationSegment struct {
	// 纬度
	Lat string `json:"lat,omitempty"`
	// 经度
	Lon string `json:"lon,omitempty"`
	// 发送时可选，标题
	Title *string `json:"title,omitempty"`
	// 发送时可选，内容描述
	Content *string `json:"content,omitempty"`
}

func (s *LocationSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeLocation
}

// MusicSegment 音乐自定义分享
// 消息段类型: music
// 支持发送
type MusicSegment struct {
	// 表示使用 QQ 音乐、网易云音乐、虾米音乐或音乐自定义分享
	Type string `json:"type,omitempty"`

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

func (s *MusicSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeMusic
}

// ReplySegment 回复
// 消息段类型: reply
// 支持发送、支持接收
type ReplySegment struct {
	// 回复时引用的消息 ID
	Id string `json:"id,omitempty"`
}

func (s *ReplySegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeReply
}

// ForwardSegment 合并转发
// 消息段类型: forward
// 支持接收
type ForwardSegment struct {
	// 合并转发 ID，需通过 get_forward_msg-获取合并转发消息 获取具体内容
	Id string `json:"id,omitempty"`
}

func (s *ForwardSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeForward
}

// NodeSegmentId 合并转发节点
// 消息段类型: node
// 支持发送
type NodeSegmentId struct {
	// 转发的消息 ID
	Id string `json:"id,omitempty"`
}

func (s *NodeSegmentId) SegmentType() MessageSegmentType {
	return MessageSegmentTypeNode
}

// NodeSegment 合并转发自定义节点
// 消息段类型: node
// 支持发送、支持接收
type NodeSegment struct {
	// 发送者 QQ 号
	UserId string `json:"user_id,omitempty"`
	// 发送者昵称
	Nickname string `json:"nickname,omitempty"`
	// 消息内容，支持发送消息时的 `message` 数据类型
	Content *MessageValue `json:"content,omitempty"`
}

func (s *NodeSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeNode
}

// XmlSegment XML 消息
// 消息段类型: xml
// 支持发送、支持接收
type XmlSegment struct {
	// XML 内容
	Data string `json:"data,omitempty"`
}

func (s *XmlSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeXml
}

// JsonSegment JSON 消息
// 消息段类型: json
// 支持发送、支持接收
type JsonSegment struct {
	// JSON 内容
	Data json.RawMessage `json:"data,omitempty"`
}

func (s *JsonSegment) SegmentType() MessageSegmentType {
	return MessageSegmentTypeJson
}
