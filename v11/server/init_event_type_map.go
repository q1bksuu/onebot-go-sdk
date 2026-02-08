package server

import (
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
)

// initEventTypeMap 初始化事件类型多层 Radix Tree，使用高性能前缀匹配，完全避免字符串拼接.
func initEventTypeMap() {
	// 第一层：按 post_type 分类
	eventRadixTree = util.NewRadixTreeFromMap(map[entity.EventPostType]*eventTreeValue{
		// message 类型：message.private, message.group
		entity.EventPostTypeMessage: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.MessageType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventMessageTypePrivate): {
					constructor: func() entity.Event { return &entity.PrivateMessageEvent{} }},
				string(entity.EventMessageTypeGroup): {
					constructor: func() entity.Event { return &entity.GroupMessageEvent{} }},
			}),
		},
		// notice 类型
		entity.EventPostTypeNotice: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.NoticeType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventNoticeTypeGroupUpload): {
					constructor: func() entity.Event { return &entity.GroupFileUploadEvent{} }},
				string(entity.EventNoticeTypeGroupAdmin): {
					constructor: func() entity.Event { return &entity.GroupAdminChangeEvent{} }},
				string(entity.EventNoticeTypeGroupDecrease): {
					constructor: func() entity.Event { return &entity.GroupMemberDecreaseEvent{} }},
				string(entity.EventNoticeTypeGroupIncrease): {
					constructor: func() entity.Event { return &entity.GroupMemberIncreaseEvent{} }},
				string(entity.EventNoticeTypeGroupBan): {
					constructor: func() entity.Event { return &entity.GroupBanEvent{} }},
				string(entity.EventNoticeTypeFriendAdd): {
					constructor: func() entity.Event { return &entity.FriendAddEvent{} }},
				string(entity.EventNoticeTypeGroupRecall): {
					constructor: func() entity.Event { return &entity.GroupRecallEvent{} }},
				string(entity.EventNoticeTypeFriendRecall): {
					constructor: func() entity.Event { return &entity.FriendRecallEvent{} }},
				string(entity.EventNoticeTypeNotify): {
					subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
						string(entity.EventNoticeSubTypeGroupPoke): {
							constructor: func() entity.Event { return &entity.GroupPokeEvent{} }},
						string(entity.EventNoticeSubTypeGroupLuckyKing): {
							constructor: func() entity.Event { return &entity.GroupLuckyKingEvent{} }},
						string(entity.EventNoticeSubTypeGroupHonor): {
							constructor: func() entity.Event { return &entity.GroupHonorChangeEvent{} }},
					})},
			}),
		},
		// request 类型：request.friend, request.group
		entity.EventPostTypeRequest: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.RequestType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventRequestTypeFriend): {
					constructor: func() entity.Event { return &entity.FriendRequestEvent{} }},
				string(entity.EventRequestTypeGroup): {
					constructor: func() entity.Event { return &entity.GroupRequestEvent{} }},
			}),
		},
		// meta_event 类型：meta_event.lifecycle, meta_event.heartbeat
		entity.EventPostTypeMetaEvent: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.MetaEventType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventMetaTypeLifecycle): {
					constructor: func() entity.Event { return &entity.LifecycleEvent{} }},
				string(entity.EventMetaTypeHeartbeat): {
					constructor: func() entity.Event { return &entity.HeartbeatEvent{} }},
			}),
		},
	})
}
