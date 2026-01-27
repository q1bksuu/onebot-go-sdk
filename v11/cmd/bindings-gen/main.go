package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func main() {
	var (
		configPath                      = flag.String("config", "cmd/bindings-gen/config.yaml", "YAML config path")
		httpServerActionsRegisterOutput = flag.String(
			"http-server-actions-register-output",
			"",
			"output go file for server bindings",
		)
		httpClientActionsOutput = flag.String("http-client-actions-output", "", "output go file for http client")
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

	log.Printf("using config: %s\n", *configPath)

	for _, genInfo := range []struct {
		template *template.Template
		output   *string
		desc     string
	}{
		{
			template: httpServerActionsRegisterTpl,
			output:   httpServerActionsRegisterOutput,
			desc:     "http server actions register output",
		},
		{
			template: httpClientActionsTpl,
			output:   httpClientActionsOutput,
			desc:     "http client actions output",
		},
	} {
		if genInfo.output != nil && strings.TrimSpace(*genInfo.output) != "" {
			outputPath := strings.TrimSpace(*genInfo.output)
			log.Printf("start generating %s -> %s\n", genInfo.desc, outputPath)

			err = os.MkdirAll(filepath.Dir(outputPath), 0o750)
			if err != nil {
				exitErr(genInfo.desc+": mkdir output dir failed", err)
			}

			code, err := render(cfg, genInfo.template)
			if err != nil {
				exitErr(genInfo.desc+": render failed", err)
			}

			err = os.WriteFile(outputPath, code, 0o600)
			if err != nil {
				exitErr(genInfo.desc+": write server output failed", err)
			}

			log.Printf("generate finished %s (%d bytes) -> %s\n", genInfo.desc, len(code), outputPath)
		}
	}
}

func loadConfig(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	fileBytes, err := os.ReadFile(abs) //nolint:gosec // config path is from flag parameter
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

	for ai, a := range group.Actions {
		err := validateAction(a, group.Name, ai)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateAction(action Action, groupName string, actionIndex int) error {
	if strings.TrimSpace(action.Method) == "" {
		return fmt.Errorf("%w: group[%s] action[%d] method is empty", errInvalidActionMethod, groupName, actionIndex)
	}

	if strings.TrimSpace(action.Request) == "" || strings.TrimSpace(action.Response) == "" {
		return fmt.Errorf(
			"%w: group[%s] action[%s] request/response is empty",
			errEmptyRequestOrResponse, groupName, action.Method,
		)
	}

	err := validateHTTPMethod(action.HTTPMethod, groupName, action.Method)
	if err != nil {
		return err
	}

	return nil
}

func validateHTTPMethod(httpMethod, groupName, actionMethod string) error {
	if m := strings.TrimSpace(httpMethod); m != "" {
		up := strings.ToUpper(m)
		if up != http.MethodGet && up != http.MethodPost {
			return fmt.Errorf(
				"%w: group[%s] action[%s] http_method must be GET or POST",
				errInvalidHTTPMethod,
				groupName,
				actionMethod,
			)
		}
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
	errInvalidActionMethod    = errors.New("invalid action method")
	errEmptyRequestOrResponse = errors.New("empty request or response")
	errInvalidHTTPMethod      = errors.New("invalid HTTP method")
)

//nolint:gochecknoglobals // templates are global for generator initialization
var (
	//go:embed templates/http_server_actions_register.gohtml
	httpServerActionsRegisterTplText string
	//go:embed templates/http_client_actions.gohtml
	httpClientActionsTplText string

	httpServerActionsRegisterTpl = template.Must(template.New("server").Parse(httpServerActionsRegisterTplText))
	httpClientActionsTpl         = template.Must(template.New("client-actions").Parse(httpClientActionsTplText))
)
