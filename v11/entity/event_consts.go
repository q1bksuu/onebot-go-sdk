package entity

type EventPostType string // 上报类型

const (
	EventPostTypeMessage   EventPostType = "message"    // 消息事件
	EventPostTypeNotice    EventPostType = "notice"     // 通知事件
	EventPostTypeRequest   EventPostType = "request"    // 请求事件
	EventPostTypeMetaEvent EventPostType = "meta_event" // 元事件
)

type EventMessageType string // 消息类型

const (
	EventMessageTypePrivate EventMessageType = "private"
	EventMessageTypeGroup   EventMessageType = "group"
)

type EventPrivateMessageSubType string

const (
	EventPrivateMessageSubTypeFriend EventPrivateMessageSubType = "friend"
	EventPrivateMessageSubTypeGroup  EventPrivateMessageSubType = "group"
	EventPrivateMessageSubTypeOther  EventPrivateMessageSubType = "other"
)

type EventNoticeType string // 通知类型

const (
	EventNoticeTypeGroupUpload   EventNoticeType = "group_upload"
	EventNoticeTypeGroupAdmin    EventNoticeType = "group_admin"
	EventNoticeTypeGroupDecrease EventNoticeType = "group_decrease"
	EventNoticeTypeGroupIncrease EventNoticeType = "group_increase"
	EventNoticeTypeGroupBan      EventNoticeType = "group_ban"
	EventNoticeTypeFriendAdd     EventNoticeType = "friend_add"
	EventNoticeTypeGroupRecall   EventNoticeType = "group_recall"
	EventNoticeTypeFriendRecall  EventNoticeType = "friend_recall"
	EventNoticeTypeNotify        EventNoticeType = "notify"
)

type EventGroupMessageSubType string

const (
	EventGroupMessageSubTypeNormal    EventGroupMessageSubType = "normal"
	EventGroupMessageSubTypeAnonymous EventGroupMessageSubType = "anonymous"
	EventGroupMessageSubTypeNotice    EventGroupMessageSubType = "notice"
)

type EventGroupAdminChangeSubType string

const (
	EventGroupAdminChangeSubTypeSet   EventGroupAdminChangeSubType = "set"
	EventGroupAdminChangeSubTypeUnset EventGroupAdminChangeSubType = "unset"
)

type EventGroupMemberDecreaseSubType string

const (
	EventGroupMemberDecreaseSubTypeLeave  EventGroupMemberDecreaseSubType = "leave"
	EventGroupMemberDecreaseSubTypeKick   EventGroupMemberDecreaseSubType = "kick"
	EventGroupMemberDecreaseSubTypeKickMe EventGroupMemberDecreaseSubType = "kick_me"
)

type EventGroupMemberIncreaseSubType string

const (
	EventGroupMemberIncreaseSubTypeApprove EventGroupMemberIncreaseSubType = "approve"
	EventGroupMemberIncreaseSubTypeInvite  EventGroupMemberIncreaseSubType = "invite"
)

type EventGroupBanSubType string

const (
	EventGroupBanSubTypeBan     EventGroupBanSubType = "ban"
	EventGroupBanSubTypeLiftBan EventGroupBanSubType = "lift_ban"
)

type EventNoticeSubType string

const (
	EventNoticeSubTypeGroupPoke      EventNoticeSubType = "poke"
	EventNoticeSubTypeGroupLuckyKing EventNoticeSubType = "lucky_king"
	EventNoticeSubTypeGroupHonor     EventNoticeSubType = "honor"
)

type EventGroupHonorChangeHonorType string

const (
	EventGroupHonorChangeHonorTypeTalkative EventGroupHonorChangeHonorType = "talkative"
	EventGroupHonorChangeHonorTypePerformer EventGroupHonorChangeHonorType = "performer"
	EventGroupHonorChangeHonorTypeEmotion   EventGroupHonorChangeHonorType = "emotion"
)

type EventRequestType string // 请求类型

const (
	EventRequestTypeFriend EventRequestType = "friend"
	EventRequestTypeGroup  EventRequestType = "group"
)

type EventGroupRequestSubType string

const (
	EventGroupRequestSubTypeAdd    EventGroupRequestSubType = "add"
	EventGroupRequestSubTypeInvite EventGroupRequestSubType = "invite"
)

type EventMetaType string // 元事件

const (
	EventMetaTypeLifecycle EventMetaType = "lifecycle"
	EventMetaTypeHeartbeat EventMetaType = "heartbeat"
)

type EventLifecycleSubType string

const (
	EventLifecycleSubTypeEnable  EventLifecycleSubType = "enable"
	EventLifecycleSubTypeDisable EventLifecycleSubType = "disable"
	EventLifecycleSubTypeConnect EventLifecycleSubType = "connect"
)
