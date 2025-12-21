package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		filename    = flag.String("file", os.Getenv("GOFILE"), "Go source file to process")
		typeName    = flag.String("type", "", "Struct type name to process (optional, process all if empty)")
		outputRaw   = flag.String("output", "", "Output file (default: {filename}_setter_getter.go)")
		constsFiles = flag.String(
			"consts",
			"",
			"Additional consts files to scan (comma-separated). Auto-scan *_consts.go by default")
		noAutoScan = flag.Bool("no-auto-scan", false, "Disable auto-scanning of *_consts.go files")
	)

	flag.Parse()

	if *filename == "" {
		log.Fatal(`Usage: entity-gen -file=<file.go> [-type=<TypeName>] [-output=<output.go>]
[-consts=<file1,file2>] [-no-auto-scan]

Options:
  -file          Go source file to process
  -type          Struct type name(s) to process (optional, process all if empty)
  -output        Output file (default: {filename}_setter_getter.go)
  -consts        Additional consts files to scan (comma-separated).
  -no-auto-scan  Disable auto-scanning of *_consts.go files`)
	}

	err := runGenerator(*filename, *typeName, *outputRaw, *constsFiles, *noAutoScan)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Generated setter/getter methods in output file")
}

func runGenerator(filename, typeName, outputRaw, constsFilesStr string, noAutoScan bool) error {
	// 确定输出文件路径
	outputFile := outputRaw
	if outputFile == "" {
		outputFile = generateOutputFilename(filename)
	}

	// 收集要扫描的 consts 文件
	scanFiles := collectScanFiles(filename, constsFilesStr, noAutoScan)

	// 解析AST并生成方法
	generator, err := NewGenerator(filename, scanFiles)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	typeList := parseTypeNames(typeName)

	methods, err := generator.Generate(typeList)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	err = os.WriteFile(outputFile, []byte(methods), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func collectScanFiles(filename, constsFilesStr string, noAutoScan bool) []string {
	var scanFiles []string

	if !noAutoScan {
		autoFiles, err := findConstsFiles(filepath.Dir(filename))
		if err == nil {
			scanFiles = append(scanFiles, autoFiles...)
		}
	}

	if constsFilesStr != "" {
		files := strings.Split(constsFilesStr, ",")
		for _, f := range files {
			f = strings.TrimSpace(f)
			if f != "" {
				scanFiles = append(scanFiles, f)
			}
		}
	}

	return scanFiles
}

func parseTypeNames(typeName string) []string {
	var typeList []string

	typeNames := strings.Split(typeName, ",")
	for _, t := range typeNames {
		t = strings.TrimSpace(t)
		if t != "" {
			typeList = append(typeList, t)
		}
	}

	return typeList
}

func generateOutputFilename(original string) string {
	dir := filepath.Dir(original)
	base := filepath.Base(original)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	return filepath.Join(dir, name+"_setter_getter.go")
}

// findConstsFiles 自动扫描目录中的 *_consts.go 文件.
func findConstsFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	var constFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// 匹配 *_consts.go 文件
		if strings.HasSuffix(name, "_consts.go") && !strings.HasSuffix(name, "_setter_getter.go") {
			constFiles = append(constFiles, filepath.Join(dir, name))
		}
	}

	return constFiles, nil
}
