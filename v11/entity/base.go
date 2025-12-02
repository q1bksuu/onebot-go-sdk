package entity

type StatusMeta struct {
	// 当前 QQ 在线，`null` 表示无法查询到在线状态
	Online bool `json:"online"`
	// 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	Good bool `json:"good"`
	// TODO 其他状态信息，视 OneBot 实现而定
}

type GroupAnonymousUser struct {
	// 匿名用户 ID
	Id int64 `json:"id"`
	// 匿名用户名称
	Name string `json:"name"`
	// 匿名用户 flag，在调用禁言 API 时需要传入
	Flag string `json:"flag"`
}
