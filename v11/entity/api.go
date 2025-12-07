//go:generate go run ../cmd/entity-gen
package entity

// SendPrivateMsgRequest send_private_msg API 的请求参数
// 发送私聊消息.
type SendPrivateMsgRequest struct {
	// 对方 QQ 号
	UserId int64 `json:"user_id"`
	// 要发送的内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效 | 可能的值: false
	AutoEscape bool `json:"auto_escape"`
}

// SendPrivateMsgResponse send_private_msg API 的响应数据.
type SendPrivateMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// SendGroupMsgRequest send_group_msg API 的请求参数
// 发送群消息.
type SendGroupMsgRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要发送的内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效 | 可能的值: false
	AutoEscape bool `json:"auto_escape"`
}

// SendGroupMsgResponse send_group_msg API 的响应数据.
type SendGroupMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// SendMsgRequest send_msg API 的请求参数
// 发送消息.
type SendMsgRequest struct {
	// 消息类型，支持 `private`、`group`，分别对应私聊、群组，如不传入，则根据传入的 `*_id` 参数判断
	MessageType MessageType `json:"message_type"`
	// 对方 QQ 号（消息类型为 `private` 时需要）
	UserId int64 `json:"user_id"`
	// 群号（消息类型为 `group` 时需要）
	GroupId int64 `json:"group_id"`
	// 要发送的内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效 | 可能的值: false
	AutoEscape bool `json:"auto_escape"`
}

// SendMsgResponse send_msg API 的响应数据.
type SendMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// DeleteMsgRequest delete_msg API 的请求参数
// 撤回消息.
type DeleteMsgRequest struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// DeleteMsgResponse delete_msg API 的响应数据.
type DeleteMsgResponse struct{}

type BaseUser struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex SexType `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
	// 备注名
	Remark string `json:"remark"`
}

// GetMsgRequest get_msg API 的请求参数
// 获取消息.
type GetMsgRequest struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// GetMsgResponse get_msg API 的响应数据.
type GetMsgResponse struct {
	// 发送时间
	Time int64 `json:"time"`
	// 消息类型，同 [消息事件](../event/message.md)
	MessageType MessageType `json:"message_type"`
	// 消息 ID
	MessageId int64 `json:"message_id"`
	// 消息真实 ID
	RealId int64 `json:"real_id"`
	// 发送人信息，同 [消息事件](../event/message.md)
	Sender *BaseUser `json:"sender"`
	// 消息内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
}

// GetForwardMsgRequest get_forward_msg API 的请求参数
// 获取合并转发消息.
type GetForwardMsgRequest struct {
	// 合并转发 ID
	Id string `json:"id"`
}

// GetForwardMsgResponse get_forward_msg API 的响应数据.
type GetForwardMsgResponse struct {
	// 消息内容，使用 [消息的数组格式](../message/array.md) 表示，数组中的消息段全部为 [`node` 消息段](../message/segment.md#合并转发自定义节点)
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
}

// SendLikeRequest send_like API 的请求参数
// 发送好友赞.
type SendLikeRequest struct {
	// 对方 QQ 号
	UserId int64 `json:"user_id"`
	// 赞的次数，每个好友每天最多 10 次 | 默认值: 1
	Times int64 `json:"times,omitempty"`
}

// SendLikeResponse send_like API 的响应数据.
type SendLikeResponse struct{}

// SetGroupKickRequest set_group_kick API 的请求参数
// 群组踢人.
type SetGroupKickRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要踢的 QQ 号
	UserId int64 `json:"user_id"`
	// 拒绝此人的加群请求 | 可能的值: false
	RejectAddRequest bool `json:"reject_add_request"`
}

// SetGroupKickResponse set_group_kick API 的响应数据.
type SetGroupKickResponse struct{}

// SetGroupBanRequest set_group_ban API 的请求参数
// 群组单人禁言.
type SetGroupBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要禁言的 QQ 号
	UserId int64 `json:"user_id"`
	// 禁言时长，单位秒，0 表示取消禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupBanResponse set_group_ban API 的响应数据.
type SetGroupBanResponse struct{}

// SetGroupAnonymousBanRequest set_group_anonymous_ban API 的请求参数
// 群组匿名用户禁言.
type SetGroupAnonymousBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 可选，要禁言的匿名用户对象（群消息上报的 `anonymous` 字段）
	Anonymous *GroupAnonymousUser `json:"anonymous,omitempty"`
	// 可选，要禁言的匿名用户的 flag（需从群消息上报的数据中获得）
	AnonymousFlag string `json:"anonymous_flag,omitempty"`
	// 禁言时长，单位秒，无法取消匿名用户禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupAnonymousBanResponse set_group_anonymous_ban API 的响应数据.
type SetGroupAnonymousBanResponse struct{}

// SetGroupWholeBanRequest set_group_whole_ban API 的请求参数
// 群组全员禁言.
type SetGroupWholeBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否禁言 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupWholeBanResponse set_group_whole_ban API 的响应数据.
type SetGroupWholeBanResponse struct{}

// SetGroupAdminRequest set_group_admin API 的请求参数
// 群组设置管理员.
type SetGroupAdminRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置管理员的 QQ 号
	UserId int64 `json:"user_id"`
	// true 为设置，false 为取消 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAdminResponse set_group_admin API 的响应数据.
type SetGroupAdminResponse struct{}

// SetGroupAnonymousRequest set_group_anonymous API 的请求参数
// 群组匿名.
type SetGroupAnonymousRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否允许匿名聊天 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAnonymousResponse set_group_anonymous API 的响应数据.
type SetGroupAnonymousResponse struct{}

// SetGroupCardRequest set_group_card API 的请求参数
// 设置群名片（群备注）.
type SetGroupCardRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置的 QQ 号
	UserId int64 `json:"user_id"`
	// 群名片内容，不填或空字符串表示删除群名片 | 默认值: 空
	Card string `json:"card,omitempty"`
}

// SetGroupCardResponse set_group_card API 的响应数据.
type SetGroupCardResponse struct{}

// SetGroupNameRequest set_group_name API 的请求参数
// 设置群名.
type SetGroupNameRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 新群名
	GroupName string `json:"group_name"`
}

// SetGroupNameResponse set_group_name API 的响应数据.
type SetGroupNameResponse struct{}

// SetGroupLeaveRequest set_group_leave API 的请求参数
// 退出群组.
type SetGroupLeaveRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散 | 可能的值: false
	IsDismiss bool `json:"is_dismiss"`
}

// SetGroupLeaveResponse set_group_leave API 的响应数据.
type SetGroupLeaveResponse struct{}

// SetGroupSpecialTitleRequest set_group_special_title API 的请求参数
// 设置群组专属头衔.
type SetGroupSpecialTitleRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置的 QQ 号
	UserId int64 `json:"user_id"`
	// 专属头衔，不填或空字符串表示删除专属头衔 | 默认值: 空
	SpecialTitle string `json:"special_title,omitempty"`
	// 专属头衔有效期，单位秒，-1 表示永久，不过此项似乎没有效果，可能是只有某些特殊的时间长度有效，有待测试 | 可能的值: -1
	Duration int64 `json:"duration"`
}

// SetGroupSpecialTitleResponse set_group_special_title API 的响应数据.
type SetGroupSpecialTitleResponse struct{}

// SetFriendAddRequestRequest set_friend_add_request API 的请求参数
// 处理加好友请求.
type SetFriendAddRequestRequest struct {
	// 加好友请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// 是否同意请求 | 可能的值: true
	Approve bool `json:"approve"`
	// 添加后的好友备注（仅在同意时有效） | 默认值: 空
	Remark string `json:"remark,omitempty"`
}

// SetFriendAddRequestResponse set_friend_add_request API 的响应数据.
type SetFriendAddRequestResponse struct{}

// SetGroupAddRequestRequest set_group_add_request API 的请求参数
// 处理加群请求／邀请.
type SetGroupAddRequestRequest struct {
	// 加群请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// `add` 或 `invite`，请求类型（需要和上报消息中的 `sub_type` 字段相符）
	SubType SetGroupAddRequestSubType `json:"sub_type"`
	// 是否同意请求／邀请 | 可能的值: true
	Approve bool `json:"approve"`
	// 拒绝理由（仅在拒绝时有效） | 默认值: 空
	Reason string `json:"reason,omitempty"`
}

// SetGroupAddRequestResponse set_group_add_request API 的响应数据.
type SetGroupAddRequestResponse struct{}

// GetLoginInfoRequest get_login_info API 的请求参数
// 获取登录号信息.
type GetLoginInfoRequest struct{}

// GetLoginInfoResponse get_login_info API 的响应数据.
type GetLoginInfoResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// QQ 昵称
	Nickname string `json:"nickname"`
}

// GetStrangerInfoRequest get_stranger_info API 的请求参数
// 获取陌生人信息.
type GetStrangerInfoRequest struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetStrangerInfoResponse get_stranger_info API 的响应数据.
type GetStrangerInfoResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex SexType `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
}

// GetFriendListRequest get_friend_list API 的请求参数
// 获取好友列表.
type GetFriendListRequest struct{}

// GetFriendListResponse get_friend_list API 的响应数据.
type GetFriendListResponse []*GetFriendListResponseItem

type GetFriendListResponseItem struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 备注名
	Remark string `json:"remark"`
}

// GetGroupInfoRequest get_group_info API 的请求参数
// 获取群信息.
type GetGroupInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupInfoResponse get_group_info API 的响应数据.
type GetGroupInfoResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 群名称
	GroupName string `json:"group_name"`
	// 成员数
	MemberCount int64 `json:"member_count"`
	// 最大成员数（群容量）
	MaxMemberCount int64 `json:"max_member_count"`
}

// GetGroupListRequest get_group_list API 的请求参数
// 获取群列表.
type GetGroupListRequest struct{}

// GetGroupListResponse get_group_list API 的响应数据.
type GetGroupListResponse []GetGroupInfoResponse

// GetGroupMemberInfoRequest get_group_member_info API 的请求参数
// 获取群成员信息.
type GetGroupMemberInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupMemberInfoResponse get_group_member_info API 的响应数据.
type GetGroupMemberInfoResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 群名片／备注
	Card string `json:"card"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex SexType `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
	// 地区
	Area string `json:"area"`
	// 加群时间戳
	JoinTime int64 `json:"join_time"`
	// 最后发言时间戳
	LastSentTime int64 `json:"last_sent_time"`
	// 成员等级
	Level string `json:"level"`
	// 角色，`owner` 或 `admin` 或 `member`
	Role GroupMemberRoleType `json:"role"`
	// 是否不良记录成员
	Unfriendly bool `json:"unfriendly"`
	// 专属头衔
	Title string `json:"title"`
	// 专属头衔过期时间戳
	TitleExpireTime int64 `json:"title_expire_time"`
	// 是否允许修改群名片
	CardChangeable bool `json:"card_changeable"`
}

// GetGroupMemberListRequest get_group_member_list API 的请求参数
// 获取群成员列表.
type GetGroupMemberListRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
}

// GetGroupMemberListResponse get_group_member_list API 的响应数据
// 但对于同一个群组的同一个成员，获取列表时和获取单独的成员信息时，某些字段可能有所不同，例如 `area`、`title` 等字段在获取列表时无法获得，具体应以单独的成员信息为准。
type GetGroupMemberListResponse []GetGroupMemberInfoResponse

// GetGroupHonorInfoRequest get_group_honor_info API 的请求参数
// 获取群荣誉信息.
type GetGroupHonorInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要获取的群荣誉类型，可传入 `talkative` `performer` `legend` `strong_newbie` `emotion` 以分别获取单个类型的群荣誉数据，或传入 `all` 获取所有数据
	Type GroupHonorType `json:"type"`
}

// GetGroupHonorInfoResponse get_group_honor_info API 的响应数据.
type GetGroupHonorInfoResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 当前龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	CurrentTalkative *GroupHonorInfoCurrentTalkative `json:"current_talkative"`
	// 历史龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	TalkativeList []*GroupHonorInfoListItem `json:"talkative_list"`
	// 群聊之火，仅 `type` 为 `performer` 或 `all` 时有数据
	PerformerList []*GroupHonorInfoListItem `json:"performer_list"`
	// 群聊炽焰，仅 `type` 为 `legend` 或 `all` 时有数据
	LegendList []*GroupHonorInfoListItem `json:"legend_list"`
	// 冒尖小春笋，仅 `type` 为 `strong_newbie` 或 `all` 时有数据
	StrongNewbieList []*GroupHonorInfoListItem `json:"strong_newbie_list"`
	// 快乐之源，仅 `type` 为 `emotion` 或 `all` 时有数据
	EmotionList []*GroupHonorInfoListItem `json:"emotion_list"`
}

type GroupHonorInfoCurrentTalkative struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 头像 URL
	Avatar string `json:"avatar"`
	// 持续天数
	DayCount int32 `json:"day_count"`
}

type GroupHonorInfoListItem struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 头像 URL
	Avatar string `json:"avatar"`
	// 荣誉描述
	Description string `json:"description"`
}

// GetCookiesRequest get_cookies API 的请求参数
// 获取 Cookies.
type GetCookiesRequest struct {
	// 需要获取 cookies 的域名 | 默认值: 空
	Domain string `json:"domain,omitempty"`
}

// GetCookiesResponse get_cookies API 的响应数据.
type GetCookiesResponse struct {
	// Cookies
	Cookies string `json:"cookies"`
}

// GetCsrfTokenRequest get_csrf_token API 的请求参数
// 获取 CSRF Token.
type GetCsrfTokenRequest struct{}

// GetCsrfTokenResponse get_csrf_token API 的响应数据.
type GetCsrfTokenResponse struct {
	// CSRF Token
	Token int64 `json:"token"`
}

// GetCredentialsRequest get_credentials API 的请求参数
// 获取 QQ 相关接口凭证.
type GetCredentialsRequest struct {
	// 需要获取 cookies 的域名 | 默认值: 空
	Domain string `json:"domain,omitempty"`
}

// GetCredentialsResponse get_credentials API 的响应数据.
type GetCredentialsResponse struct {
	// Cookies
	Cookies string `json:"cookies"`
	// CSRF Token
	CsrfToken int64 `json:"csrf_token"`
}

// GetRecordRequest get_record API 的请求参数
// 获取语音.
type GetRecordRequest struct {
	// 收到的语音文件名（消息段的 `file` 参数），如 `0B38145AA44505000B38145AA4450500.silk`
	File string `json:"file"`
	// 要转换到的格式，目前支持 `mp3`、`amr`、`wma`、`m4a`、`spx`、`ogg`、`wav`、`flac`
	OutFormat GetRecordOutputFormat `json:"out_format"`
}

// GetRecordResponse get_record API 的响应数据.
type GetRecordResponse struct {
	// 转换后的语音文件路径，如 `/home/somebody/cqhttp/data/record/0B38145AA44505000B38145AA4450500.mp3`
	File string `json:"file"`
}

// GetImageRequest get_image API 的请求参数
// 获取图片.
type GetImageRequest struct {
	// 收到的图片文件名（消息段的 `file` 参数），如 `6B4DE3DFD1BD271E3297859D41C530F5.jpg`
	File string `json:"file"`
}

// GetImageResponse get_image API 的响应数据.
type GetImageResponse struct {
	// 下载后的图片文件路径，如 `/home/somebody/cqhttp/data/image/6B4DE3DFD1BD271E3297859D41C530F5.jpg`
	File string `json:"file"`
}

// CanSendImageRequest can_send_image API 的请求参数
// 检查是否可以发送图片.
type CanSendImageRequest struct{}

// CanSendImageResponse can_send_image API 的响应数据.
type CanSendImageResponse struct {
	// 是或否
	Yes bool `json:"yes"`
}

// CanSendRecordRequest can_send_record API 的请求参数
// 检查是否可以发送语音.
type CanSendRecordRequest struct{}

// CanSendRecordResponse can_send_record API 的响应数据.
type CanSendRecordResponse struct {
	// 是或否
	Yes bool `json:"yes"`
}

// GetStatusRequest get_status API 的请求参数
// 获取运行状态.
type GetStatusRequest struct{}

// GetStatusResponse get_status API 的响应数据.
type GetStatusResponse = StatusMeta

// GetVersionInfoRequest get_version_info API 的请求参数
// 获取版本信息.
type GetVersionInfoRequest struct{}

// GetVersionInfoResponse get_version_info API 的响应数据.
type GetVersionInfoResponse struct {
	// 应用标识，如 `mirai-native`
	AppName string `json:"app_name"`
	// 应用版本，如 `1.2.3`
	AppVersion string `json:"app_version"`
	// OneBot 标准版本，如 `v11`
	ProtocolVersion string `json:"protocol_version"`
	// TODO 其他状态信息，视 OneBot 实现而定
}

// SetRestartRequest set_restart API 的请求参数
// 重启 OneBot 实现.
type SetRestartRequest struct {
	// 要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为 2000 左右 | 可能的值: 0
	Delay int64 `json:"delay"`
}

// SetRestartResponse set_restart API 的响应数据.
type SetRestartResponse struct{}

// CleanCacheRequest clean_cache API 的请求参数
// 清理缓存.
type CleanCacheRequest struct{}

// CleanCacheResponse clean_cache API 的响应数据.
type CleanCacheResponse struct{}
