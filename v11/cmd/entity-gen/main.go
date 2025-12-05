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
		// 对应 go generate 的 $GOFILE, $GOPACKAGE 等变量
		filename    = flag.String("file", os.Getenv("GOFILE"), "Go source file to process")
		typeName    = flag.String("type", "", "Struct type name to process (optional, process all if empty)")
		outputRaw   = flag.String("output", "", "Output file (default: {filename}_setter_getter.go)")
		constsFiles = flag.String("consts", "", "Additional consts files to scan (comma-separated). Auto-scan *_consts.go by default")
		noAutoScan  = flag.Bool("no-auto-scan", false, "Disable auto-scanning of *_consts.go files")
	)
	flag.Parse()

	if *filename == "" {
		fmt.Println(`Usage: entity-gen -file=<file.go> [-type=<TypeName>] [-output=<output.go>] [-consts=<file1,file2>] [-no-auto-scan]

Or use with 'go generate':

  //go:generate entity-gen -type=MyType
  //go:generate entity-gen -type=MyType,AnotherType
  //go:generate entity-gen (generate all types)

Options:
  -file          Go source file to process
  -type          Struct type name(s) to process (optional, comma-separated)
  -output        Output file (default: {filename}_setter_getter.go)
  -consts        Additional consts files to scan (comma-separated)
  -no-auto-scan  Disable auto-scanning of *_consts.go files in the same directory`)
		os.Exit(1)
	}

	// 确定输出文件路径
	var outputFile string
	if *outputRaw != "" {
		outputFile = *outputRaw
	} else {
		outputFile = generateOutputFilename(*filename)
	}

	// 收集要扫描的 consts 文件
	var scanFiles []string
	if !*noAutoScan {
		// 自动扫描同目录的 *_consts.go 文件
		autoFiles, _ := findConstsFiles(filepath.Dir(*filename))
		scanFiles = append(scanFiles, autoFiles...)
	}
	if *constsFiles != "" {
		// 添加手动指定的 consts 文件
		files := strings.Split(*constsFiles, ",")
		for _, f := range files {
			f = strings.TrimSpace(f)
			if f != "" {
				scanFiles = append(scanFiles, f)
			}
		}
	}

	// 解析AST并生成方法
	generator, err := NewGenerator(*filename, scanFiles)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}
	var methods string

	var typeList []string
	// 一次运行生成多个类型
	typeNames := strings.Split(*typeName, ",")
	for _, t := range typeNames {
		t = strings.TrimSpace(t)
		if t != "" {
			typeList = append(typeList, t)
		}
	}
	methods, genErr := generator.Generate(typeList)
	if genErr != nil {
		log.Fatalf("Generation failed: %v", genErr)
	}

	if writeErr := os.WriteFile(outputFile, []byte(methods), 0o644); writeErr != nil {
		log.Fatalf("Failed to write output file: %v", writeErr)
	}

	if *typeName != "" {
		fmt.Printf("✓ Generated setter/getter methods for %s in %s\n", *typeName, outputFile)
	} else {
		fmt.Printf("✓ Generated setter/getter methods in %s\n", outputFile)
	}
}

func generateOutputFilename(original string) string {
	dir := filepath.Dir(original)
	base := filepath.Base(original)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(dir, name+"_setter_getter.go")
}

// findConstsFiles 自动扫描目录中的 *_consts.go 文件
func findConstsFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
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
