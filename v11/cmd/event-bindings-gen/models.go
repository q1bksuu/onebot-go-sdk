package main

type Config struct {
	Groups          []Group         `yaml:"groups"`
	CombinedService CombinedService `yaml:"combined_service"`
}

type CombinedService struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
}

type Group struct {
	Name        string  `yaml:"name"`
	ServiceName string  `yaml:"service_name"`
	ServiceDesc string  `yaml:"service_desc"`
	Events      []Event `yaml:"events"`
}

type Event struct {
	Method string `yaml:"method"`
	Key    string `yaml:"key"`
	Desc   string `yaml:"desc"`
	Type   string `yaml:"type"`
}
