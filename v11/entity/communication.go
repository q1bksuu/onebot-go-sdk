//go:generate go run ../cmd/entity-gen
package entity

import (
	"encoding/json"
	"fmt"
)

// ActionRequest 表示上报到传输层的动作请求.
type ActionRequest struct {
	// Action 为路径段，例如 send_private_msg
	Action string `json:"action"`
	// Params 保存解析后的参数映射，可能来自 query/form/JSON
	Params map[string]any `json:"params,omitempty"`
}

// ActionRawResponse 表示传输层返回给调用方的标准响应.
type ActionRawResponse struct {
	Status  ActionResponseStatus  `json:"status"`            // ok | async | failed
	Retcode ActionResponseRetcode `json:"retcode"`           // 0=成功,1=异步,其他=失败
	Data    json.RawMessage       `json:"data,omitempty"`    // 返回数据，可能为 null
	Message string                `json:"message,omitempty"` // 可选，人类可读的错误信息
}

// ActionResponse 表示传输层返回给调用方的标准响应，Data 字段已解码.
type ActionResponse[T any] struct {
	Status  ActionResponseStatus  `json:"status"`            // ok | async | failed
	Retcode ActionResponseRetcode `json:"retcode"`           // 0=成功,1=异步,其他=失败
	Data    *T                    `json:"data,omitempty"`    // 返回数据，可能为 null
	Message string                `json:"message,omitempty"` // 可选，人类可读的错误信息
}

func (r *ActionResponse[T]) ToActionRawResponse() (*ActionRawResponse, error) {
	if r == nil {
		//nolint:nilnil
		return nil, nil
	}

	var rawData json.RawMessage

	if r.Data != nil {
		b, err := json.Marshal(r.Data)
		if err != nil {
			return nil, fmt.Errorf("fail to marshal Data field: %w", err)
		}

		rawData = b
	}

	return &ActionRawResponse{
		Status:  r.Status,
		Retcode: r.Retcode,
		Data:    rawData,
		Message: r.Message,
	}, nil
}

type ActionError struct {
	UrlPath string
	Status  ActionResponseStatus
	Retcode ActionResponseRetcode
	Message string
}

func (e *ActionError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("action %s failed: status=%s retcode=%d message=%s", e.UrlPath, e.Status, e.Retcode, e.Message)
}
