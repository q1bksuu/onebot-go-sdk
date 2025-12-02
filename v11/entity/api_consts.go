package entity

type MessageType string

const (
	MessageTypePrivate MessageType = "private"
	MessageTypeGroup   MessageType = "group"
)

type SexType string

const (
	SexTypeMale    SexType = "male"
	SexTypeFemale  SexType = "female"
	SexTypeUnknown SexType = "unknown"
)

type GroupMemberRoleType string

const (
	GroupMemberRoleTypeOwner   GroupMemberRoleType = "owner"
	GroupMemberRoleTypeAdmin   GroupMemberRoleType = "admin"
	GroupMemberRoleTypeMember  GroupMemberRoleType = "member"
	GroupMemberRoleTypeUnknown GroupMemberRoleType = "unknown"
)

type SetGroupAddRequestSubType string

const (
	SetGroupAddRequestSubTypeAdd    SetGroupAddRequestSubType = "add"
	SetGroupAddRequestSubTypeInvite SetGroupAddRequestSubType = "invite"
)

type GroupHonorType string

const (
	GroupHonorTypeTalkative    GroupHonorType = "talkative"
	GroupHonorTypePerformer    GroupHonorType = "performer"
	GroupHonorTypeLegend       GroupHonorType = "legend"
	GroupHonorTypeStrongNewbie GroupHonorType = "strong_newbie"
	GroupHonorTypeEmotion      GroupHonorType = "emotion"
	GroupHonorTypeAll          GroupHonorType = "all"
)

type GetRecordOutputFormat string

const (
	GetRecordOutputFormatMP3  GetRecordOutputFormat = "mp3"
	GetRecordOutputFormatAMR  GetRecordOutputFormat = "amr"
	GetRecordOutputFormatWMA  GetRecordOutputFormat = "wma"
	GetRecordOutputFormatM4A  GetRecordOutputFormat = "m4a"
	GetRecordOutputFormatSPX  GetRecordOutputFormat = "spx"
	GetRecordOutputFormatOGG  GetRecordOutputFormat = "ogg"
	GetRecordOutputFormatWAV  GetRecordOutputFormat = "wav"
	GetRecordOutputFormatFLAC GetRecordOutputFormat = "flac"
)
