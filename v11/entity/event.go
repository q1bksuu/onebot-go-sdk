//go:generate go run ../cmd/entity-gen
package entity

// PrivateMessageEvent 私聊消息
// 事件类型: private
// 子类型: friend,group,other.
type PrivateMessageEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: message
	PostType PrivateMessageEventPostType `json:"post_type"`
	// 消息类型 | 可能的值: private
	MessageType PrivateMessageEventMessageType `json:"message_type"`
	// 消息子类型，如果是好友则是 `friend`，如果是群临时会话则是 `group` | 可能的值: friend, group, other
	SubType PrivateMessageEventSubType `json:"sub_type"`
	// 消息 ID
	MessageId int64 `json:"message_id"`
	// 发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 消息内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 原始消息内容
	RawMessage string `json:"raw_message"`
	// 字体
	Font int64 `json:"font"`
	// 发送人信息
	Sender *PrivateMessageEventSender `json:"sender"`
}

type PrivateMessageEventSender struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex SexType `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
}

// GroupMessageEvent 群消息
// 事件类型: group
// 子类型: normal,anonymous,notice.
type GroupMessageEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: message
	PostType GroupMessageEventPostType `json:"post_type"`
	// 消息类型 | 可能的值: group
	MessageType GroupMessageEventMessageType `json:"message_type"`
	// 消息子类型，正常消息是 `normal`，匿名消息是 `anonymous`，系统提示（如「管理员已禁止群内匿名聊天」）是 `notice` | 可能的值: normal, anonymous, notice
	SubType GroupMessageEventSubType `json:"sub_type"`
	// 消息 ID
	MessageId int64 `json:"message_id"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 匿名信息，如果不是匿名消息则为 null
	Anonymous *GroupAnonymousUser `json:"anonymous"`
	// 消息内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 原始消息内容
	RawMessage string `json:"raw_message"`
	// 字体
	Font int64 `json:"font"`
	// 发送人信息
	Sender *GroupMessageEventSender `json:"sender"`
}

type GroupMessageEventSender struct {
	// 发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 群名片／备注
	Card string `json:"card"`
	// 性别，male 或 female 或 unknown
	Sex SexType `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
	// 地区
	Area string `json:"area"`
	// 成员等级
	Level string `json:"level"`
	// 角色，owner 或 admin 或 member
	Role GroupMemberRoleType `json:"role"`
	// 专属头衔
	Title string `json:"title"`
}

// GroupFileUploadEvent 群文件上传
// 事件类型: group_upload.
type GroupFileUploadEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupFileUploadEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_upload
	NoticeType GroupFileUploadEventNoticeType `json:"notice_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 文件信息
	File *GroupFileUploadEventFile `json:"file"`
}

type GroupFileUploadEventFile struct {
	// 文件 ID
	Id string `json:"id"`
	// 文件名
	Name string `json:"name"`
	// 文件大小（字节数）
	Size int64 `json:"size"`
	// busid（目前不清楚有什么作用）
	BusId int64 `json:"busid"`
}

// GroupAdminChangeEvent 群管理员变动
// 事件类型: group_admin
// 子类型: set,unset.
type GroupAdminChangeEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupAdminChangeEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_admin
	NoticeType GroupAdminChangeEventNoticeType `json:"notice_type"`
	// 事件子类型，分别表示设置和取消管理员 | 可能的值: set, unset
	SubType GroupAdminChangeEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 管理员 QQ 号
	UserId int64 `json:"user_id"`
}

// GroupMemberDecreaseEvent 群成员减少
// 事件类型: group_decrease
// 子类型: leave,kick,kick_me.
type GroupMemberDecreaseEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupMemberDecreaseEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_decrease
	NoticeType GroupMemberDecreaseEventNoticeType `json:"notice_type"`
	// 事件子类型，分别表示主动退群、成员被踢、登录号被踢 | 可能的值: leave, kick, kick_me
	SubType GroupMemberDecreaseEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 操作者 QQ 号（如果是主动退群，则和 `user_id` 相同）
	OperatorId int64 `json:"operator_id"`
	// 离开者 QQ 号
	UserId int64 `json:"user_id"`
}

// GroupMemberIncreaseEvent 群成员增加
// 事件类型: group_increase
// 子类型: approve,invite.
type GroupMemberIncreaseEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupMemberIncreaseEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_increase
	NoticeType GroupMemberIncreaseEventNoticeType `json:"notice_type"`
	// 事件子类型，分别表示管理员已同意入群、管理员邀请入群 | 可能的值: approve, invite
	SubType GroupMemberIncreaseEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 操作者 QQ 号
	OperatorId int64 `json:"operator_id"`
	// 加入者 QQ 号
	UserId int64 `json:"user_id"`
}

// GroupBanEvent 群禁言
// 事件类型: group_ban
// 子类型: ban,lift_ban.
type GroupBanEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupBanEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_ban
	NoticeType GroupBanEventNoticeType `json:"notice_type"`
	// 事件子类型，分别表示禁言、解除禁言 | 可能的值: ban, lift_ban
	SubType GroupBanEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 操作者 QQ 号
	OperatorId int64 `json:"operator_id"`
	// 被禁言 QQ 号
	UserId int64 `json:"user_id"`
	// 禁言时长，单位秒
	Duration int64 `json:"duration"`
}

// FriendAddEvent 好友添加
// 事件类型: friend_add.
type FriendAddEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType FriendAddEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: friend_add
	NoticeType FriendAddEventNoticeType `json:"notice_type"`
	// 新添加好友 QQ 号
	UserId int64 `json:"user_id"`
}

// GroupRecallEvent 群消息撤回
// 事件类型: group_recall.
type GroupRecallEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupRecallEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: group_recall
	NoticeType GroupRecallEventNoticeType `json:"notice_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 消息发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 操作者 QQ 号
	OperatorId int64 `json:"operator_id"`
	// 被撤回的消息 ID
	MessageId int64 `json:"message_id"`
}

// FriendRecallEvent 好友消息撤回
// 事件类型: friend_recall.
type FriendRecallEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType FriendRecallEventPostType `json:"post_type"`
	// 通知类型 | 可能的值: friend_recall
	NoticeType FriendRecallEventNoticeType `json:"notice_type"`
	// 好友 QQ 号
	UserId int64 `json:"user_id"`
	// 被撤回的消息 ID
	MessageId int64 `json:"message_id"`
}

// GroupPokeEvent 群内戳一戳
// 事件类型: notify
// 子类型: poke.
type GroupPokeEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupPokeEventPostType `json:"post_type"`
	// 消息类型 | 可能的值: notify
	NoticeType GroupPokeEventNoticeType `json:"notice_type"`
	// 提示类型 | 可能的值: poke
	SubType GroupPokeEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 被戳者 QQ 号
	TargetId int64 `json:"target_id"`
}

// GroupLuckyKingEvent 群红包运气王
// 事件类型: notify
// 子类型: lucky_king.
type GroupLuckyKingEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupLuckyKingEventPostType `json:"post_type"`
	// 消息类型 | 可能的值: notify
	NoticeType GroupLuckyKingEventNoticeType `json:"notice_type"`
	// 提示类型 | 可能的值: lucky_king
	SubType GroupLuckyKingEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 红包发送者 QQ 号
	UserId int64 `json:"user_id"`
	// 运气王 QQ 号
	TargetId int64 `json:"target_id"`
}

// GroupHonorChangeEvent 群成员荣誉变更
// 事件类型: notify
// 子类型: honor.
type GroupHonorChangeEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: notice
	PostType GroupHonorChangeEventPostType `json:"post_type"`
	// 消息类型 | 可能的值: notify
	NoticeType GroupHonorChangeEventNoticeType `json:"notice_type"`
	// 提示类型 | 可能的值: honor
	SubType GroupHonorChangeEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 荣誉类型，分别表示龙王、群聊之火、快乐源泉 | 可能的值: talkative, performer, emotion
	HonorType GroupHonorChangeEventHonorType `json:"honor_type"`
	// 成员 QQ 号
	UserId int64 `json:"user_id"`
}

// FriendRequestEvent 加好友请求
// 事件类型: friend.
type FriendRequestEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: request
	PostType FriendRequestEventPostType `json:"post_type"`
	// 请求类型 | 可能的值: friend
	RequestType FriendRequestEventRequestType `json:"request_type"`
	// 发送请求的 QQ 号
	UserId int64 `json:"user_id"`
	// 验证信息
	Comment string `json:"comment"`
	// 请求 flag，在调用处理请求的 API 时需要传入
	Flag string `json:"flag"`
}

// GroupRequestEvent 加群请求／邀请
// 事件类型: group
// 子类型: add,invite.
type GroupRequestEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: request
	PostType GroupRequestEventPostType `json:"post_type"`
	// 请求类型 | 可能的值: group
	RequestType GroupRequestEventRequestType `json:"request_type"`
	// 请求子类型，分别表示加群请求、邀请登录号入群 | 可能的值: add, invite
	SubType GroupRequestEventSubType `json:"sub_type"`
	// 群号
	GroupId int64 `json:"group_id"`
	// 发送请求的 QQ 号
	UserId int64 `json:"user_id"`
	// 验证信息
	Comment string `json:"comment"`
	// 请求 flag，在调用处理请求的 API 时需要传入
	Flag string `json:"flag"`
}

// LifecycleEvent 生命周期
// 事件类型: lifecycle
// 子类型: enable,disable,connect.
type LifecycleEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: meta_event
	PostType LifecycleEventPostType `json:"post_type"`
	// 元事件类型 | 可能的值: lifecycle
	MetaEventType LifecycleEventMetaEventType `json:"meta_event_type"`
	// 事件子类型，分别表示 OneBot 启用、停用、WebSocket 连接成功 | 可能的值: enable, disable, connect
	SubType LifecycleEventSubType `json:"sub_type"`
}

// HeartbeatEvent 心跳
// 事件类型: heartbeat.
type HeartbeatEvent struct {
	// 事件发生的时间戳
	Time int64 `json:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int64 `json:"self_id"`
	// 上报类型 | 可能的值: meta_event
	PostType HeartbeatEventPostType `json:"post_type"`
	// 元事件类型 | 可能的值: heartbeat
	MetaEventType HeartbeatEventMetaEventType `json:"meta_event_type"`
	// 状态信息
	Status *StatusMeta `json:"status"`
	// 到下次心跳的间隔，单位毫秒
	Interval int64 `json:"interval"`
}
