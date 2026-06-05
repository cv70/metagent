package aof

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompileYAMLSubset(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "morphogen.yaml")
	if err := os.WriteFile(path, []byte(`morphogen:
  version: "1.0"
  identity:
    domain: "support-agent"
    ontology_ref: "ontology.md"
    description: "Support specialist"
  capabilities:
    - id: "faq"
      type: "knowledge"
      activation: "always"
      dependencies: []
      config:
        summary: "Product FAQ"
        keywords: "refund billing"
    - id: "boundary"
      type: "belief"
      activation: "always"
      dependencies: ["faq"]
      config:
        summary: "Escalate legal questions"
  evolution:
    inherit: true
    max_generations: 3
`), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := CompileFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if result.Morphogen.Identity.Domain != "support-agent" {
		t.Fatalf("domain = %q", result.Morphogen.Identity.Domain)
	}
	if len(result.Crystals) != 2 {
		t.Fatalf("crystal count = %d", len(result.Crystals))
	}
	if result.Crystals[1].Dependencies[0] != "faq" {
		t.Fatalf("dependency not parsed: %#v", result.Crystals[1].Dependencies)
	}
}

func TestValidateUnknownDependency(t *testing.T) {
	_, err := Compile(Morphogen{
		Identity: Identity{Domain: "bad"},
		Capabilities: []Capability{{
			ID:           "a",
			Type:         CrystalKnowledge,
			Activation:   ActivationAlways,
			Dependencies: []string{"missing"},
		}},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
