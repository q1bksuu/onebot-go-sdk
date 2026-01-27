package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func main() {
	var (
		configPath                 = flag.String("config", "cmd/event-bindings-gen/config.yaml", "YAML config path")
		serverEventsRegisterOutput = flag.String(
			"http-server-events-register-output", "",
			"output go file for server event bindings",
		)
	)

	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		exitErr("load config failed", err)
	}

	err = validateConfig(cfg)
	if err != nil {
		exitErr("validate config failed", err)
	}

	if serverEventsRegisterOutput != nil && strings.TrimSpace(*serverEventsRegisterOutput) != "" {
		outputPath := strings.TrimSpace(*serverEventsRegisterOutput)

		err = os.MkdirAll(filepath.Dir(outputPath), 0o750)
		if err != nil {
			exitErr("mkdir output dir failed", err)
		}

		code, err := render(cfg, serverEventsRegisterTpl)
		if err != nil {
			exitErr("render failed", err)
		}

		err = os.WriteFile(outputPath, code, 0o600)
		if err != nil {
			exitErr("write output failed", err)
		}
	}
}

func loadConfig(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	// 这是一个代码生成工具，文件路径来自命令行参数，这是预期的行为
	//nolint:gosec // 路径来自命令行参数，是安全的
	fileBytes, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if strings.TrimSpace(cfg.CombinedService.Name) == "" {
		return fmt.Errorf("%w: combined service name is empty", errInvalidCombinedService)
	}

	if len(cfg.Groups) == 0 {
		return fmt.Errorf("%w", errNoGroupsConfigured)
	}

	for gi, g := range cfg.Groups {
		err := validateGroup(g, gi)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateGroup(group Group, groupIndex int) error {
	if strings.TrimSpace(group.Name) == "" {
		return fmt.Errorf("%w: group[%d] name is empty", errInvalidGroupName, groupIndex)
	}

	if strings.TrimSpace(group.ServiceName) == "" {
		return fmt.Errorf("%w: group[%d] service_name is empty", errInvalidGroupName, groupIndex)
	}

	for ei, e := range group.Events {
		err := validateEvent(e, group.Name, ei)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateEvent(event Event, groupName string, eventIndex int) error {
	if strings.TrimSpace(event.Method) == "" {
		return fmt.Errorf("%w: group[%s] event[%d] method is empty", errInvalidEventMethod, groupName, eventIndex)
	}

	if strings.TrimSpace(event.Key) == "" {
		return fmt.Errorf("%w: group[%s] event[%s] key is empty", errInvalidEventKey, groupName, event.Method)
	}

	if strings.TrimSpace(event.Type) == "" {
		return fmt.Errorf("%w: group[%s] event[%s] type is empty", errInvalidEventType, groupName, event.Method)
	}

	return nil
}

func render(data *Config, t *template.Template) ([]byte, error) {
	var buf bytes.Buffer

	err := t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format source: %w", err)
	}

	return formattedCode, nil
}

func exitErr(stage string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", stage, err)

	os.Exit(1)
}

var (
	errNoGroupsConfigured     = errors.New("no groups configured")
	errInvalidGroupName       = errors.New("invalid group name")
	errInvalidCombinedService = errors.New("invalid combined service")
	errInvalidEventMethod     = errors.New("invalid event method")
	errInvalidEventKey        = errors.New("invalid event key")
	errInvalidEventType       = errors.New("invalid event type")
)

// serverEventsRegisterTpl 是代码生成模板，作为全局变量是合理的
//
//nolint:gochecknoglobals // 模板作为全局变量是合理的
var (
	//go:embed templates/server_events_register.gohtml
	serverEventsRegisterTplText string
	serverEventsRegisterTpl     = template.Must(template.New("server-events").Parse(serverEventsRegisterTplText))
)
