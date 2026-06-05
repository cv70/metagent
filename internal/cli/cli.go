package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai4next/metagent/internal/aof"
)

func Execute(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}
	switch args[0] {
	case "init":
		return cmdInit(args[1:])
	case "validate":
		return cmdValidate(args[1:])
	case "compile":
		return cmdCompile(args[1:])
	case "generate":
		return cmdGenerate(args[1:])
	case "run":
		return cmdRun(args[1:])
	case "help", "--help", "-h":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func cmdInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	out := fs.String("o", "morphogen.yaml", "output morphogen path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil && filepath.Dir(*out) != "." {
		return err
	}
	return os.WriteFile(*out, []byte(sampleMorphogen), 0o644)
}

func cmdValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	file := fs.String("f", "morphogen.yaml", "morphogen file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	m, err := aof.LoadMorphogen(*file)
	if err != nil {
		return err
	}
	if err := aof.Validate(m); err != nil {
		return err
	}
	fmt.Printf("ok: %s (%d capabilities)\n", m.Identity.Domain, len(m.Capabilities))
	return nil
}

func cmdCompile(args []string) error {
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	file := fs.String("f", "morphogen.yaml", "morphogen file")
	out := fs.String("o", "", "output JSON path; stdout when empty")
	if err := fs.Parse(args); err != nil {
		return err
	}
	result, err := aof.CompileFile(*file)
	if err != nil {
		return err
	}
	if *out != "" {
		return aof.WriteJSON(*out, result)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func cmdGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	file := fs.String("f", "morphogen.yaml", "morphogen file")
	out := fs.String("o", "dist/agent", "output agent directory")
	base := fs.String("base", "local-base-kernel", "base kernel reference")
	if err := fs.Parse(args); err != nil {
		return err
	}
	agent, err := aof.GenerateAgent(*file, *out, *base)
	if err != nil {
		return err
	}
	fmt.Printf("generated %s at %s with %d crystals\n", agent.ID, *out, len(agent.Crystals))
	return nil
}

func cmdRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	file := fs.String("agent", "dist/agent/agent.json", "generated agent.json")
	prompt := fs.String("p", "", "user prompt")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *prompt == "" && fs.NArg() > 0 {
		*prompt = fs.Arg(0)
	}
	if *prompt == "" {
		return fmt.Errorf("prompt is required: metagent run -p %q", "your task")
	}
	data, err := os.ReadFile(*file)
	if err != nil {
		return err
	}
	var agent aof.AgentInstance
	if err := json.Unmarshal(data, &agent); err != nil {
		return err
	}
	fmt.Print(aof.RenderOfflineRun(agent, *prompt))
	return nil
}

func printHelp() {
	fmt.Print(`metagent - Agent Ontogenesis Framework CLI

Usage:
  metagent init [-o morphogen.yaml]
  metagent validate -f morphogen.yaml
  metagent compile -f morphogen.yaml [-o compiled.json]
  metagent generate -f morphogen.yaml [-o dist/agent] [--base kernel]
  metagent run --agent dist/agent/agent.json -p "task"

The implementation is offline-first: it compiles a Morphogen specification into
runtime crystals and emits a portable specialized-agent directory.
`)
}

const sampleMorphogen = `morphogen:
  version: "1.0"
  identity:
    domain: "go-code-review"
    ontology_ref: "docs/go-review-ontology.md"
    description: "专注 Go 项目正确性、边界、并发和测试缺口的代码审查 agent。"
  capabilities:
    - id: "go-review-knowledge"
      type: "knowledge"
      activation: "always"
      dependencies: []
      config:
        summary: "Go 代码审查知识：错误处理、context 传播、并发安全、接口边界、测试覆盖。"
        keywords: "go review error context goroutine mutex race test"
    - id: "review-protocol"
      type: "protocol"
      activation: "always"
      dependencies: ["go-review-knowledge"]
      config:
        summary: "先列严重问题和文件位置，再列开放问题，最后给简短变更概述。"
        keywords: "findings severity file line risk tests"
    - id: "safety-belief"
      type: "belief"
      activation: "always"
      dependencies: []
      config:
        summary: "不能臆测未读取文件；无法验证时必须标明残余风险。"
        keywords: "uncertainty boundary verification"
  evolution:
    inherit: true
    max_generations: 8
    rules:
      - trigger: "review_completed && finding_confirmed"
        action: "commit_pattern_to_germline"
        priority: 10
`
