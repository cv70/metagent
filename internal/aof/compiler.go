package aof

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadMorphogen(path string) (Morphogen, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Morphogen{}, err
	}
	var m Morphogen
	if err := json.Unmarshal(data, &m); err == nil && (m.Identity.Domain != "" || len(m.Capabilities) > 0) {
		return normalizeMorphogen(m), nil
	}
	var wrapped struct {
		Morphogen Morphogen `json:"morphogen"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && (wrapped.Morphogen.Identity.Domain != "" || len(wrapped.Morphogen.Capabilities) > 0) {
		return normalizeMorphogen(wrapped.Morphogen), nil
	}
	m, err = parseYAMLSubset(string(data))
	if err != nil {
		return Morphogen{}, fmt.Errorf("parse morphogen %s: %w", path, err)
	}
	return normalizeMorphogen(m), nil
}

func CompileFile(path string) (CompileResult, error) {
	m, err := LoadMorphogen(path)
	if err != nil {
		return CompileResult{}, err
	}
	configs, err := Compile(m)
	if err != nil {
		return CompileResult{}, err
	}
	return CompileResult{Morphogen: m, Crystals: configs}, nil
}

func Compile(m Morphogen) ([]CrystalConfig, error) {
	if err := Validate(m); err != nil {
		return nil, err
	}
	out := make([]CrystalConfig, 0, len(m.Capabilities))
	for _, cap := range m.Capabilities {
		out = append(out, CrystalConfig{
			ID:     cap.ID,
			Type:   cap.Type,
			Source: cloneMap(cap.Config),
			Synthesis: SynthesisConfig{
				Strategy: synthesisStrategy(cap.Type),
				Inputs:   configKeys(cap.Config),
			},
			Activation: ActivationPolicy{Mode: cap.Activation, Condition: cap.Config["condition"]},
			Dependencies: append([]string(nil), cap.Dependencies...),
		})
	}
	return out, nil
}

func Validate(m Morphogen) error {
	var problems []string
	if m.Identity.Domain == "" {
		problems = append(problems, "identity.domain is required")
	}
	if len(m.Capabilities) == 0 {
		problems = append(problems, "at least one capability is required")
	}
	seen := map[string]bool{}
	for i, cap := range m.Capabilities {
		label := fmt.Sprintf("capabilities[%d]", i)
		if cap.ID == "" {
			problems = append(problems, label+".id is required")
		}
		if seen[cap.ID] {
			problems = append(problems, "duplicate capability id: "+cap.ID)
		}
		seen[cap.ID] = true
		switch cap.Type {
		case CrystalKnowledge, CrystalTool, CrystalProtocol, CrystalMemory, CrystalBelief:
		default:
			problems = append(problems, label+".type must be knowledge/tool/protocol/memory/belief")
		}
		switch cap.Activation {
		case ActivationAlways, ActivationContextual, ActivationManual:
		default:
			problems = append(problems, label+".activation must be always/contextual/manual")
		}
	}
	for _, cap := range m.Capabilities {
		for _, dep := range cap.Dependencies {
			if !seen[dep] {
				problems = append(problems, fmt.Sprintf("capability %s depends on unknown crystal %s", cap.ID, dep))
			}
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func normalizeMorphogen(m Morphogen) Morphogen {
	if m.Version == "" {
		m.Version = "1.0"
	}
	for i := range m.Capabilities {
		if m.Capabilities[i].Activation == "" {
			m.Capabilities[i].Activation = ActivationContextual
		}
		if m.Capabilities[i].Config == nil {
			m.Capabilities[i].Config = map[string]string{}
		}
	}
	if m.Evolution.MaxGenerations == 0 {
		m.Evolution.MaxGenerations = 8
	}
	return m
}

func synthesisStrategy(t CrystalType) string {
	switch t {
	case CrystalKnowledge:
		return "knowledge-index"
	case CrystalTool:
		return "tool-contract"
	case CrystalProtocol:
		return "interaction-policy"
	case CrystalMemory:
		return "episodic-schema"
	case CrystalBelief:
		return "belief-boundary"
	default:
		return "generic"
	}
}

func configKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func cloneMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
