package entity

// SendPrivateMsgRequest 表示 send_private_msg API 的请求参数
// 对应文档: 发送私聊消息
type SendPrivateMsgRequest struct {
	// 对方 QQ 号
	UserId int64 `json:"user_id"`
	// 要发送的内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效 | 可能的值: false
	AutoEscape bool `json:"auto_escape"`
}

// SendPrivateMsgResponse 表示 send_private_msg API 的响应数据
type SendPrivateMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// SendGroupMsgRequest 表示 send_group_msg API 的请求参数
// 对应文档: 发送群消息
type SendGroupMsgRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要发送的内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
	// 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效 | 可能的值: false
	AutoEscape bool `json:"auto_escape"`
}

// SendGroupMsgResponse 表示 send_group_msg API 的响应数据
type SendGroupMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// SendMsgRequest 表示 send_msg API 的请求参数
// 对应文档: 发送消息
type SendMsgRequest struct {
	// 消息类型，支持 `private`、`group`，分别对应私聊、群组，如不传入，则根据传入的 `*_id` 参数判断
	MessageType string `json:"message_type"`
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

// SendMsgResponse 表示 send_msg API 的响应数据
type SendMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// DeleteMsgRequest 表示 delete_msg API 的请求参数
// 对应文档: 撤回消息
type DeleteMsgRequest struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// DeleteMsgResponse 表示 delete_msg API 的响应数据
type DeleteMsgResponse struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// GetMsgRequest 表示 get_msg API 的请求参数
// 对应文档: 获取消息
type GetMsgRequest struct {
	// 消息 ID
	MessageId int64 `json:"message_id"`
}

// GetMsgResponse 表示 get_msg API 的响应数据
type GetMsgResponse struct {
	// 发送时间
	Time int64 `json:"time"`
	// 消息类型，同 [消息事件](../event/message.md)
	MessageType string `json:"message_type"`
	// 消息 ID
	MessageId int64 `json:"message_id"`
	// 消息真实 ID
	RealId int64 `json:"real_id"`
	// 发送人信息，同 [消息事件](../event/message.md)
	Sender *map[string]interface{} `json:"sender"`
	// 消息内容
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
}

// GetForwardMsgRequest 表示 get_forward_msg API 的请求参数
// 对应文档: 获取合并转发消息
type GetForwardMsgRequest struct {
	// 合并转发 ID
	Id string `json:"id"`
}

// GetForwardMsgResponse 表示 get_forward_msg API 的响应数据
type GetForwardMsgResponse struct {
	// 消息内容，使用 [消息的数组格式](../message/array.md) 表示，数组中的消息段全部为 [`node` 消息段](../message/segment.md#合并转发自定义节点)
	// 可以是字符串 (CQ 码格式) 或消息段数组
	Message *MessageValue `json:"message"`
}

// SendLikeRequest 表示 send_like API 的请求参数
// 对应文档: 发送好友赞
type SendLikeRequest struct {
	// 对方 QQ 号
	UserId int64 `json:"user_id"`
	// 赞的次数，每个好友每天最多 10 次 | 默认值: 1
	Times int64 `json:"times,omitempty"`
}

// SendLikeResponse 表示 send_like API 的响应数据
type SendLikeResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要踢的 QQ 号
	UserId int64 `json:"user_id"`
	// 拒绝此人的加群请求 | 可能的值: false
	RejectAddRequest bool `json:"reject_add_request"`
}

// SetGroupKickRequest 表示 set_group_kick API 的请求参数
// 对应文档: 群组踢人
type SetGroupKickRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要踢的 QQ 号
	UserId int64 `json:"user_id"`
	// 拒绝此人的加群请求 | 可能的值: false
	RejectAddRequest bool `json:"reject_add_request"`
}

// SetGroupKickResponse 表示 set_group_kick API 的响应数据
type SetGroupKickResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要禁言的 QQ 号
	UserId int64 `json:"user_id"`
	// 禁言时长，单位秒，0 表示取消禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupBanRequest 表示 set_group_ban API 的请求参数
// 对应文档: 群组单人禁言
type SetGroupBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要禁言的 QQ 号
	UserId int64 `json:"user_id"`
	// 禁言时长，单位秒，0 表示取消禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupBanResponse 表示 set_group_ban API 的响应数据
type SetGroupBanResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 可选，要禁言的匿名用户对象（群消息上报的 `anonymous` 字段）
	Anonymous *map[string]interface{} `json:"anonymous"`
	// 可选，要禁言的匿名用户的 flag（需从群消息上报的数据中获得）
	AnonymousFlag string `json:"anonymous_flag"`
	// 禁言时长，单位秒，无法取消匿名用户禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupAnonymousBanRequest 表示 set_group_anonymous_ban API 的请求参数
// 对应文档: 群组匿名用户禁言
type SetGroupAnonymousBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 可选，要禁言的匿名用户对象（群消息上报的 `anonymous` 字段）
	Anonymous *map[string]interface{} `json:"anonymous"`
	// 可选，要禁言的匿名用户的 flag（需从群消息上报的数据中获得）
	AnonymousFlag string `json:"anonymous_flag"`
	// 禁言时长，单位秒，无法取消匿名用户禁言 | 可能的值: 30 * 60
	Duration int64 `json:"duration"`
}

// SetGroupAnonymousBanResponse 表示 set_group_anonymous_ban API 的响应数据
type SetGroupAnonymousBanResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否禁言 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupWholeBanRequest 表示 set_group_whole_ban API 的请求参数
// 对应文档: 群组全员禁言
type SetGroupWholeBanRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否禁言 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupWholeBanResponse 表示 set_group_whole_ban API 的响应数据
type SetGroupWholeBanResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置管理员的 QQ 号
	UserId int64 `json:"user_id"`
	// true 为设置，false 为取消 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAdminRequest 表示 set_group_admin API 的请求参数
// 对应文档: 群组设置管理员
type SetGroupAdminRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置管理员的 QQ 号
	UserId int64 `json:"user_id"`
	// true 为设置，false 为取消 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAdminResponse 表示 set_group_admin API 的响应数据
type SetGroupAdminResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否允许匿名聊天 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAnonymousRequest 表示 set_group_anonymous API 的请求参数
// 对应文档: 群组匿名
type SetGroupAnonymousRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否允许匿名聊天 | 可能的值: true
	Enable bool `json:"enable"`
}

// SetGroupAnonymousResponse 表示 set_group_anonymous API 的响应数据
type SetGroupAnonymousResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置的 QQ 号
	UserId int64 `json:"user_id"`
	// 群名片内容，不填或空字符串表示删除群名片 | 默认值: 空
	Card string `json:"card,omitempty"`
}

// SetGroupCardRequest 表示 set_group_card API 的请求参数
// 对应文档: 设置群名片（群备注）
type SetGroupCardRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置的 QQ 号
	UserId int64 `json:"user_id"`
	// 群名片内容，不填或空字符串表示删除群名片 | 默认值: 空
	Card string `json:"card,omitempty"`
}

// SetGroupCardResponse 表示 set_group_card API 的响应数据
type SetGroupCardResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 新群名
	GroupName string `json:"group_name"`
}

// SetGroupNameRequest 表示 set_group_name API 的请求参数
// 对应文档: 设置群名
type SetGroupNameRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 新群名
	GroupName string `json:"group_name"`
}

// SetGroupNameResponse 表示 set_group_name API 的响应数据
type SetGroupNameResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散 | 可能的值: false
	IsDismiss bool `json:"is_dismiss"`
}

// SetGroupLeaveRequest 表示 set_group_leave API 的请求参数
// 对应文档: 退出群组
type SetGroupLeaveRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散 | 可能的值: false
	IsDismiss bool `json:"is_dismiss"`
}

// SetGroupLeaveResponse 表示 set_group_leave API 的响应数据
type SetGroupLeaveResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要设置的 QQ 号
	UserId int64 `json:"user_id"`
	// 专属头衔，不填或空字符串表示删除专属头衔 | 默认值: 空
	SpecialTitle string `json:"special_title,omitempty"`
	// 专属头衔有效期，单位秒，-1 表示永久，不过此项似乎没有效果，可能是只有某些特殊的时间长度有效，有待测试 | 可能的值: -1
	Duration int64 `json:"duration"`
}

// SetGroupSpecialTitleRequest 表示 set_group_special_title API 的请求参数
// 对应文档: 设置群组专属头衔
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

// SetGroupSpecialTitleResponse 表示 set_group_special_title API 的响应数据
type SetGroupSpecialTitleResponse struct {
	// 加好友请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// 是否同意请求 | 可能的值: true
	Approve bool `json:"approve"`
	// 添加后的好友备注（仅在同意时有效） | 默认值: 空
	Remark string `json:"remark,omitempty"`
}

// SetFriendAddRequestRequest 表示 set_friend_add_request API 的请求参数
// 对应文档: 处理加好友请求
type SetFriendAddRequestRequest struct {
	// 加好友请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// 是否同意请求 | 可能的值: true
	Approve bool `json:"approve"`
	// 添加后的好友备注（仅在同意时有效） | 默认值: 空
	Remark string `json:"remark,omitempty"`
}

// SetFriendAddRequestResponse 表示 set_friend_add_request API 的响应数据
type SetFriendAddRequestResponse struct {
	// 加群请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// `add` 或 `invite`，请求类型（需要和上报消息中的 `sub_type` 字段相符）
	SubType string `json:"sub_type"`
	// 是否同意请求／邀请 | 可能的值: true
	Approve bool `json:"approve"`
	// 拒绝理由（仅在拒绝时有效） | 默认值: 空
	Reason string `json:"reason,omitempty"`
}

// SetGroupAddRequestRequest 表示 set_group_add_request API 的请求参数
// 对应文档: 处理加群请求／邀请
type SetGroupAddRequestRequest struct {
	// 加群请求的 flag（需从上报的数据中获得）
	Flag string `json:"flag"`
	// `add` 或 `invite`，请求类型（需要和上报消息中的 `sub_type` 字段相符）
	SubType string `json:"sub_type"`
	// 是否同意请求／邀请 | 可能的值: true
	Approve bool `json:"approve"`
	// 拒绝理由（仅在拒绝时有效） | 默认值: 空
	Reason string `json:"reason,omitempty"`
}

// SetGroupAddRequestResponse 表示 set_group_add_request API 的响应数据
type SetGroupAddRequestResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// QQ 昵称
	Nickname string `json:"nickname"`
}

// GetLoginInfoRequest 表示 get_login_info API 的请求参数
// 对应文档: 获取登录号信息
type GetLoginInfoRequest struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// QQ 昵称
	Nickname string `json:"nickname"`
}

// GetLoginInfoResponse 表示 get_login_info API 的响应数据
type GetLoginInfoResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// QQ 昵称
	Nickname string `json:"nickname"`
}

// GetStrangerInfoRequest 表示 get_stranger_info API 的请求参数
// 对应文档: 获取陌生人信息
type GetStrangerInfoRequest struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetStrangerInfoResponse 表示 get_stranger_info API 的响应数据
type GetStrangerInfoResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex string `json:"sex"`
	// 年龄
	Age int64 `json:"age"`
}

// GetFriendListRequest 表示 get_friend_list API 的请求参数
// 对应文档: 获取好友列表
type GetFriendListRequest struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 备注名
	Remark string `json:"remark"`
}

// GetFriendListResponse 表示 get_friend_list API 的响应数据
type GetFriendListResponse struct {
	// QQ 号
	UserId int64 `json:"user_id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 备注名
	Remark string `json:"remark"`
}

// GetGroupInfoRequest 表示 get_group_info API 的请求参数
// 对应文档: 获取群信息
type GetGroupInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupInfoResponse 表示 get_group_info API 的响应数据
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

// GetGroupListRequest 表示 get_group_list API 的请求参数
// 对应文档: 获取群列表
type GetGroupListRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupListResponse 表示 get_group_list API 的响应数据
type GetGroupListResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupMemberInfoRequest 表示 get_group_member_info API 的请求参数
// 对应文档: 获取群成员信息
type GetGroupMemberInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// QQ 号
	UserId int64 `json:"user_id"`
	// 是否不使用缓存（使用缓存可能更新不及时，但响应更快） | 可能的值: false
	NoCache bool `json:"no_cache"`
}

// GetGroupMemberInfoResponse 表示 get_group_member_info API 的响应数据
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
	Sex string `json:"sex"`
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
	Role string `json:"role"`
	// 是否不良记录成员
	Unfriendly bool `json:"unfriendly"`
	// 专属头衔
	Title string `json:"title"`
	// 专属头衔过期时间戳
	TitleExpireTime int64 `json:"title_expire_time"`
	// 是否允许修改群名片
	CardChangeable bool `json:"card_changeable"`
}

// GetGroupMemberListRequest 表示 get_group_member_list API 的请求参数
// 对应文档: 获取群成员列表
type GetGroupMemberListRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
}

// GetGroupMemberListResponse 表示 get_group_member_list API 的响应数据
type GetGroupMemberListResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要获取的群荣誉类型，可传入 `talkative` `performer` `legend` `strong_newbie` `emotion` 以分别获取单个类型的群荣誉数据，或传入 `all` 获取所有数据
	Type string `json:"type"`
}

// GetGroupHonorInfoRequest 表示 get_group_honor_info API 的请求参数
// 对应文档: 获取群荣誉信息
type GetGroupHonorInfoRequest struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 要获取的群荣誉类型，可传入 `talkative` `performer` `legend` `strong_newbie` `emotion` 以分别获取单个类型的群荣誉数据，或传入 `all` 获取所有数据
	Type string `json:"type"`
}

// GetGroupHonorInfoResponse 表示 get_group_honor_info API 的响应数据
type GetGroupHonorInfoResponse struct {
	// 群号
	GroupId int64 `json:"group_id"`
	// 当前龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	CurrentTalkative *map[string]interface{} `json:"current_talkative"`
	// 历史龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	TalkativeList *[]interface{} `json:"talkative_list"`
	// 群聊之火，仅 `type` 为 `performer` 或 `all` 时有数据
	PerformerList *[]interface{} `json:"performer_list"`
	// 群聊炽焰，仅 `type` 为 `legend` 或 `all` 时有数据
	LegendList *[]interface{} `json:"legend_list"`
	// 冒尖小春笋，仅 `type` 为 `strong_newbie` 或 `all` 时有数据
	StrongNewbieList *[]interface{} `json:"strong_newbie_list"`
	// 快乐之源，仅 `type` 为 `emotion` 或 `all` 时有数据
	EmotionList *[]interface{} `json:"emotion_list"`
}

// GetCookiesRequest 表示 get_cookies API 的请求参数
// 对应文档: 获取 Cookies
type GetCookiesRequest struct {
	// 需要获取 cookies 的域名 | 默认值: 空
	Domain string `json:"domain,omitempty"`
}

// GetCookiesResponse 表示 get_cookies API 的响应数据
type GetCookiesResponse struct {
	// Cookies
	Cookies string `json:"cookies"`
}

// GetCsrfTokenRequest 表示 get_csrf_token API 的请求参数
// 对应文档: 获取 CSRF Token
type GetCsrfTokenRequest struct {
	// CSRF Token
	Token int64 `json:"token"`
}

// GetCsrfTokenResponse 表示 get_csrf_token API 的响应数据
type GetCsrfTokenResponse struct {
	// CSRF Token
	Token int64 `json:"token"`
}

// GetCredentialsRequest 表示 get_credentials API 的请求参数
// 对应文档: 获取 QQ 相关接口凭证
type GetCredentialsRequest struct {
	// 需要获取 cookies 的域名 | 默认值: 空
	Domain string `json:"domain,omitempty"`
}

// GetCredentialsResponse 表示 get_credentials API 的响应数据
type GetCredentialsResponse struct {
	// Cookies
	Cookies string `json:"cookies"`
	// CSRF Token
	CsrfToken int64 `json:"csrf_token"`
}

// GetRecordRequest 表示 get_record API 的请求参数
// 对应文档: 获取语音
type GetRecordRequest struct {
	// 收到的语音文件名（消息段的 `file` 参数），如 `0B38145AA44505000B38145AA4450500.silk`
	File string `json:"file"`
	// 要转换到的格式，目前支持 `mp3`、`amr`、`wma`、`m4a`、`spx`、`ogg`、`wav`、`flac`
	OutFormat string `json:"out_format"`
}

// GetRecordResponse 表示 get_record API 的响应数据
type GetRecordResponse struct {
	// 转换后的语音文件路径，如 `/home/somebody/cqhttp/data/record/0B38145AA44505000B38145AA4450500.mp3`
	File string `json:"file"`
}

// GetImageRequest 表示 get_image API 的请求参数
// 对应文档: 获取图片
type GetImageRequest struct {
	// 收到的图片文件名（消息段的 `file` 参数），如 `6B4DE3DFD1BD271E3297859D41C530F5.jpg`
	File string `json:"file"`
}

// GetImageResponse 表示 get_image API 的响应数据
type GetImageResponse struct {
	// 下载后的图片文件路径，如 `/home/somebody/cqhttp/data/image/6B4DE3DFD1BD271E3297859D41C530F5.jpg`
	File string `json:"file"`
}

// CanSendImageRequest 表示 can_send_image API 的请求参数
// 对应文档: 检查是否可以发送图片
type CanSendImageRequest struct {
	// 是或否
	Yes bool `json:"yes"`
}

// CanSendImageResponse 表示 can_send_image API 的响应数据
type CanSendImageResponse struct {
	// 是或否
	Yes bool `json:"yes"`
}

// CanSendRecordRequest 表示 can_send_record API 的请求参数
// 对应文档: 检查是否可以发送语音
type CanSendRecordRequest struct {
	// 是或否
	Yes bool `json:"yes"`
}

// CanSendRecordResponse 表示 can_send_record API 的响应数据
type CanSendRecordResponse struct {
	// 是或否
	Yes bool `json:"yes"`
}

// GetStatusRequest 表示 get_status API 的请求参数
// 对应文档: 获取运行状态
type GetStatusRequest struct {
	// 当前 QQ 在线，`null` 表示无法查询到在线状态
	Online bool `json:"online"`
	// 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	Good bool `json:"good"`
}

// GetStatusResponse 表示 get_status API 的响应数据
type GetStatusResponse struct {
	// 当前 QQ 在线，`null` 表示无法查询到在线状态
	Online bool `json:"online"`
	// 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	Good bool `json:"good"`
}

// GetVersionInfoRequest 表示 get_version_info API 的请求参数
// 对应文档: 获取版本信息
type GetVersionInfoRequest struct {
	// 应用标识，如 `mirai-native`
	AppName string `json:"app_name"`
	// 应用版本，如 `1.2.3`
	AppVersion string `json:"app_version"`
	// OneBot 标准版本，如 `v11`
	ProtocolVersion string `json:"protocol_version"`
}

// GetVersionInfoResponse 表示 get_version_info API 的响应数据
type GetVersionInfoResponse struct {
	// 应用标识，如 `mirai-native`
	AppName string `json:"app_name"`
	// 应用版本，如 `1.2.3`
	AppVersion string `json:"app_version"`
	// OneBot 标准版本，如 `v11`
	ProtocolVersion string `json:"protocol_version"`
}

// SetRestartRequest 表示 set_restart API 的请求参数
// 对应文档: 重启 OneBot 实现
type SetRestartRequest struct {
	// 要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为 2000 左右 | 可能的值: 0
	Delay int64 `json:"delay"`
}

// SetRestartResponse 表示 set_restart API 的响应数据
type SetRestartResponse struct {
	// 无响应数据
}

// CleanCacheRequest 表示 clean_cache API 的请求参数
// 对应文档: 清理缓存
type CleanCacheRequest struct {
}

// CleanCacheResponse 表示 clean_cache API 的响应数据
type CleanCacheResponse struct {
	// 无响应数据
}
