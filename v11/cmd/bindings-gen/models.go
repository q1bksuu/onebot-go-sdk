package main

// Config 配置根.
type Config struct {
	Groups          []Group         `yaml:"groups"`
	CombinedService CombinedService `yaml:"combined_service"`
}

type CombinedService struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
}

// Group 表示一组业务接口，可生成独立 Service.
type Group struct {
	Name        string   `yaml:"name"`
	ServiceName string   `yaml:"service_name"`
	ServiceDesc string   `yaml:"service_desc"`
	Actions     []Action `yaml:"actions"`
}

// Action 定义单个 action 的生成规则.
type Action struct {
	// Method 生成的 Service 方法名（必填）.
	Method string `yaml:"method"`
	// Action 协议动作名（必填）.
	Action string `yaml:"action"`
	// Desc 方法描述（必填）.
	Desc string `yaml:"desc"`
	// Request / Response 类型（必填），需包含包名，例如 entity.SendPrivateMsgRequest.
	Request  string `yaml:"request"`
	Response string `yaml:"response"`
	// HTTPMethod 可选，默认 POST，可指定 GET/POST.
	HTTPMethod string `yaml:"http_method"`
	// Path 可选，默认 "/{action}".
	Path string `yaml:"path"`
}
