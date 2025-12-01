package entity

type PrivateMessageEventPostType string

const (
	PrivateMessageEventPostTypeMessage PrivateMessageEventPostType = "message"
)

type PrivateMessageEventMessageType string

const (
	PrivateMessageEventMessageTypePrivate PrivateMessageEventMessageType = "private"
)

type PrivateMessageEventSubType string

const (
	PrivateMessageEventSubTypeFriend PrivateMessageEventSubType = "friend"
	PrivateMessageEventSubTypeGroup  PrivateMessageEventSubType = "group"
	PrivateMessageEventSubTypeOther  PrivateMessageEventSubType = "other"
)

type GroupMessageEventPostType string

const (
	GroupMessageEventPostTypeMessage GroupMessageEventPostType = "message"
)

type GroupMessageEventMessageType string

const (
	GroupMessageEventMessageTypeGroup GroupMessageEventMessageType = "group"
)

type GroupMessageEventSubType string

const (
	GroupMessageEventSubTypeNormal    GroupMessageEventSubType = "normal"
	GroupMessageEventSubTypeAnonymous GroupMessageEventSubType = "anonymous"
	GroupMessageEventSubTypeNotice    GroupMessageEventSubType = "notice"
)

type GroupFileUploadEventPostType string

const (
	GroupFileUploadEventPostTypeNotice GroupFileUploadEventPostType = "notice"
)

type GroupFileUploadEventNoticeType string

const (
	GroupFileUploadEventNoticeTypeGroupUpload GroupFileUploadEventNoticeType = "group_upload"
)

type GroupAdminChangeEventPostType string

const (
	GroupAdminChangeEventPostTypeNotice GroupAdminChangeEventPostType = "notice"
)

type GroupAdminChangeEventNoticeType string

const (
	GroupAdminChangeEventNoticeTypeGroupAdmin GroupAdminChangeEventNoticeType = "group_admin"
)

type GroupAdminChangeEventSubType string

const (
	GroupAdminChangeEventSubTypeSet   GroupAdminChangeEventSubType = "set"
	GroupAdminChangeEventSubTypeUnset GroupAdminChangeEventSubType = "unset"
)

type GroupMemberDecreaseEventPostType string

const (
	GroupMemberDecreaseEventPostTypeNotice GroupMemberDecreaseEventPostType = "notice"
)

type GroupMemberDecreaseEventNoticeType string

const (
	GroupMemberDecreaseEventNoticeTypeGroupDecrease GroupMemberDecreaseEventNoticeType = "group_decrease"
)

type GroupMemberDecreaseEventSubType string

const (
	GroupMemberDecreaseEventSubTypeLeave  GroupMemberDecreaseEventSubType = "leave"
	GroupMemberDecreaseEventSubTypeKick   GroupMemberDecreaseEventSubType = "kick"
	GroupMemberDecreaseEventSubTypeKickMe GroupMemberDecreaseEventSubType = "kick_me"
)

type GroupMemberIncreaseEventPostType string

const (
	GroupMemberIncreaseEventPostTypeNotice GroupMemberIncreaseEventPostType = "notice"
)

type GroupMemberIncreaseEventNoticeType string

const (
	GroupMemberIncreaseEventNoticeTypeGroupIncrease GroupMemberIncreaseEventNoticeType = "group_increase"
)

type GroupMemberIncreaseEventSubType string

const (
	GroupMemberIncreaseEventSubTypeApprove GroupMemberIncreaseEventSubType = "approve"
	GroupMemberIncreaseEventSubTypeInvite  GroupMemberIncreaseEventSubType = "invite"
)

type GroupBanEventPostType string

const (
	GroupBanEventPostTypeNotice GroupBanEventPostType = "notice"
)

type GroupBanEventNoticeType string

const (
	GroupBanEventNoticeTypeGroupBan GroupBanEventNoticeType = "group_ban"
)

type GroupBanEventSubType string

const (
	GroupBanEventSubTypeBan     GroupBanEventSubType = "ban"
	GroupBanEventSubTypeLiftBan GroupBanEventSubType = "lift_ban"
)

type FriendAddEventPostType string

const (
	FriendAddEventPostTypeNotice FriendAddEventPostType = "notice"
)

type FriendAddEventNoticeType string

const (
	FriendAddEventNoticeTypeFriendAdd FriendAddEventNoticeType = "friend_add"
)

type GroupRecallEventPostType string

const (
	GroupRecallEventPostTypeNotice GroupRecallEventPostType = "notice"
)

type GroupRecallEventNoticeType string

const (
	GroupRecallEventNoticeTypeGroupRecall GroupRecallEventNoticeType = "group_recall"
)

type FriendRecallEventPostType string

const (
	FriendRecallEventPostTypeNotice FriendRecallEventPostType = "notice"
)

type FriendRecallEventNoticeType string

const (
	FriendRecallEventNoticeTypeFriendRecall FriendRecallEventNoticeType = "friend_recall"
)

type GroupPokeEventPostType string

const (
	GroupPokeEventPostTypeNotice GroupPokeEventPostType = "notice"
)

type GroupPokeEventNoticeType string

const (
	GroupPokeEventNoticeTypeNotify GroupPokeEventNoticeType = "notify"
)

type GroupPokeEventSubType string

const (
	GroupPokeEventSubTypePoke GroupPokeEventSubType = "poke"
)

type GroupLuckyKingEventPostType string

const (
	GroupLuckyKingEventPostTypeNotice GroupLuckyKingEventPostType = "notice"
)

type GroupLuckyKingEventNoticeType string

const (
	GroupLuckyKingEventNoticeTypeNotify GroupLuckyKingEventNoticeType = "notify"
)

type GroupLuckyKingEventSubType string

const (
	GroupLuckyKingEventSubTypeLuckyKing GroupLuckyKingEventSubType = "lucky_king"
)

type GroupHonorChangeEventPostType string

const (
	GroupHonorChangeEventPostTypeNotice GroupHonorChangeEventPostType = "notice"
)

type GroupHonorChangeEventNoticeType string

const (
	GroupHonorChangeEventNoticeTypeNotify GroupHonorChangeEventNoticeType = "notify"
)

type GroupHonorChangeEventSubType string

const (
	GroupHonorChangeEventSubTypeHonor GroupHonorChangeEventSubType = "honor"
)

type GroupHonorChangeEventHonorType string

const (
	GroupHonorChangeEventHonorTypeTalkative GroupHonorChangeEventHonorType = "talkative"
	GroupHonorChangeEventHonorTypePerformer GroupHonorChangeEventHonorType = "performer"
	GroupHonorChangeEventHonorTypeEmotion   GroupHonorChangeEventHonorType = "emotion"
)

type FriendRequestEventPostType string

const (
	FriendRequestEventPostTypeRequest FriendRequestEventPostType = "request"
)

type FriendRequestEventRequestType string

const (
	FriendRequestEventRequestTypeFriend FriendRequestEventRequestType = "friend"
)

type GroupRequestEventPostType string

const (
	GroupRequestEventPostTypeRequest GroupRequestEventPostType = "request"
)

type GroupRequestEventRequestType string

const (
	GroupRequestEventRequestTypeGroup GroupRequestEventRequestType = "group"
)

type GroupRequestEventSubType string

const (
	GroupRequestEventSubTypeAdd    GroupRequestEventSubType = "add"
	GroupRequestEventSubTypeInvite GroupRequestEventSubType = "invite"
)

type LifecycleEventPostType string

const (
	LifecycleEventPostTypeMetaEvent LifecycleEventPostType = "meta_event"
)

type LifecycleEventMetaEventType string

const (
	LifecycleEventMetaEventTypeLifecycle LifecycleEventMetaEventType = "lifecycle"
)

type LifecycleEventSubType string

const (
	LifecycleEventSubTypeEnable  LifecycleEventSubType = "enable"
	LifecycleEventSubTypeDisable LifecycleEventSubType = "disable"
	LifecycleEventSubTypeConnect LifecycleEventSubType = "connect"
)

type HeartbeatEventPostType string

const (
	HeartbeatEventPostTypeMetaEvent HeartbeatEventPostType = "meta_event"
)

type HeartbeatEventMetaEventType string

const (
	HeartbeatEventMetaEventTypeHeartbeat HeartbeatEventMetaEventType = "heartbeat"
)
