package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// EventDispatcher 根据事件类型字段路由到对应 handler.
type EventDispatcher struct {
	handlers map[string]EventHandler
}

var _ EventRequestHandler = (*EventDispatcher)(nil)

// NewEventDispatcher 创建事件分发器.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{handlers: make(map[string]EventHandler)}
}

// Register 注册事件处理器.
// key 格式: "post_type" 或 "post_type/type" 或 "post_type/type/sub_type"
// 例如:
//   - "message" - 所有消息事件
//   - "message/private" - 所有私聊消息
//   - "message/private/friend" - 好友私聊消息
//   - "notice/group_upload" - 群文件上传
//   - "notice/notify/poke" - 群内戳一戳
//   - "request/friend" - 好友请求
//   - "meta_event/lifecycle" - 生命周期事件
func (d *EventDispatcher) Register(key string, h EventHandler) {
	d.handlers[key] = h
}

// HandleEvent 调用对应事件 handler.
func (d *EventDispatcher) HandleEvent(ctx context.Context, event entity.Event) (map[string]any, error) {
	keys := d.buildEventKeys(event)

	// 按优先级尝试匹配：最具体的优先
	for _, key := range keys {
		if h, ok := d.handlers[key]; ok {
			return h(ctx, event)
		}
	}

	// 如果没有匹配的处理器，返回 nil（204 No Content）
	return nil, ErrNoEventHandler
}

// buildEventKeys 根据事件类型字段构建可能的匹配键，按优先级从高到低排序（最具体的优先）.
func (d *EventDispatcher) buildEventKeys(event entity.Event) []string {
	postType := event.GetPostType()
	postTypeStr := string(postType)

	eventMap, err := d.parseEventToMap(event)
	if err != nil {
		// 如果解析失败，只返回 post_type
		return []string{postTypeStr}
	}

	var keys []string

	switch postType {
	case entity.EventPostTypeMessage:
		keys = d.buildMessageKeys(postTypeStr, eventMap)

	case entity.EventPostTypeNotice:
		keys = d.buildNoticeKeys(postTypeStr, eventMap)

	case entity.EventPostTypeRequest:
		keys = d.buildRequestKeys(postTypeStr, eventMap)

	case entity.EventPostTypeMetaEvent:
		keys = d.buildMetaEventKeys(postTypeStr, eventMap)

	default:
		keys = []string{postTypeStr}
	}

	return keys
}

// parseEventToMap 将事件解析为 map.
func (d *EventDispatcher) parseEventToMap(event entity.Event) (map[string]any, error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshal event failed: %w", err)
	}

	var eventMap map[string]any

	err = json.Unmarshal(eventJSON, &eventMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal event failed: %w", err)
	}

	return eventMap, nil
}

// buildMessageKeys 构建消息事件的键.
func (d *EventDispatcher) buildMessageKeys(postTypeStr string, eventMap map[string]any) []string {
	var keys []string

	if msgType, ok := eventMap["message_type"].(string); ok {
		if subType, ok := eventMap["sub_type"].(string); ok && subType != "" {
			keys = append(keys, fmt.Sprintf("%s/%s/%s", postTypeStr, msgType, subType))
		}

		keys = append(keys, fmt.Sprintf("%s/%s", postTypeStr, msgType))
	}

	keys = append(keys, postTypeStr)

	return keys
}

// buildNoticeKeys 构建通知事件的键.
func (d *EventDispatcher) buildNoticeKeys(postTypeStr string, eventMap map[string]any) []string {
	var keys []string

	if noticeType, ok := eventMap["notice_type"].(string); ok {
		if subType, ok := eventMap["sub_type"].(string); ok && subType != "" {
			keys = append(keys, fmt.Sprintf("%s/%s/%s", postTypeStr, noticeType, subType))
		}

		keys = append(keys, fmt.Sprintf("%s/%s", postTypeStr, noticeType))
	}

	keys = append(keys, postTypeStr)

	return keys
}

// buildRequestKeys 构建请求事件的键.
func (d *EventDispatcher) buildRequestKeys(postTypeStr string, eventMap map[string]any) []string {
	var keys []string

	if reqType, ok := eventMap["request_type"].(string); ok {
		if subType, ok := eventMap["sub_type"].(string); ok && subType != "" {
			keys = append(keys, fmt.Sprintf("%s/%s/%s", postTypeStr, reqType, subType))
		}

		keys = append(keys, fmt.Sprintf("%s/%s", postTypeStr, reqType))
	}

	keys = append(keys, postTypeStr)

	return keys
}

// buildMetaEventKeys 构建元事件的键.
func (d *EventDispatcher) buildMetaEventKeys(postTypeStr string, eventMap map[string]any) []string {
	var keys []string

	if metaType, ok := eventMap["meta_event_type"].(string); ok {
		if subType, ok := eventMap["sub_type"].(string); ok && subType != "" {
			keys = append(keys, fmt.Sprintf("%s/%s/%s", postTypeStr, metaType, subType))
		}

		keys = append(keys, fmt.Sprintf("%s/%s", postTypeStr, metaType))
	}

	keys = append(keys, postTypeStr)

	return keys
}
