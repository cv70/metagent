# Metagent

Metagent is a Go-first implementation of the Agent Ontogenesis Framework (AOF):
it compiles a declarative Morphogen specification into runtime crystals and emits
a portable specialized-agent directory.

The current implementation is offline-first. It does not call an LLM by itself;
instead it generates the agent contract (`system.md`), runtime metadata
(`agent.json`), compiled crystal configs, and per-crystal payloads that another
agent runtime can load.

## Quick Start

```bash
# Create a starter Morphogen spec
metagent init -o morphogen.yaml

# Validate and compile the specialization protocol
metagent validate -f morphogen.yaml
metagent compile -f morphogen.yaml -o compiled.json

# Generate a specialized agent directory
metagent generate -f morphogen.yaml -o dist/go-review-agent

# Dry-run crystal selection for a task
metagent run --agent dist/go-review-agent/agent.json -p "review this Go concurrency patch"
```

Without an installed binary, build from source with Go and run `./metagent`.

```bash
go build -o metagent .
./metagent --help
```

## Morphogen DSL

Metagent accepts JSON or the small YAML subset emitted by `metagent init`.

```yaml
morphogen:
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
  evolution:
    inherit: true
    max_generations: 8
```

Capability types:

- `knowledge`: domain facts, rubrics, ontology fragments, reasoning hints
- `tool`: tool contracts and invocation boundaries
- `protocol`: interaction style, review workflow, dialogue policy
- `memory`: episodic memory schema and retrieval semantics
- `belief`: confidence and safety boundaries

Activation modes:

- `always`: active when the agent starts
- `contextual`: selected when task text matches id, summary, or keywords
- `manual`: held latent until an external runtime or user activates it

## Generated Layout

```text
dist/go-review-agent/
├── README.md
├── agent.json
├── compiled.json
├── system.md
└── crystals/
    ├── go-review-knowledge.json
    ├── review-protocol.json
    └── safety-belief.json
```

`system.md` is the human-readable specialized-agent contract. `agent.json` is the
runtime manifest. Files in `crystals/` are the internalized capability units.

## Architecture

See `ARCHITECTURE.md` for the full AOF design: Morphogen Compiler, Crystal Forge,
Agent Differentiator, Ontogeny Runtime, Crystal Repository, and Germline Archive.

This repository currently implements the first practical slice:

- Morphogen parser and validator
- Crystal config emitter
- Runtime crystal forge and binder metadata
- Specialized agent manifest generation
- Offline cortex-style crystal selection

## Reference

The CLI shape and project organization are intentionally inspired by
`/home/x/space/superman`: a simple top-level `main.go`, internal packages, and a
command-oriented interface.
