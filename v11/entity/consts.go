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
