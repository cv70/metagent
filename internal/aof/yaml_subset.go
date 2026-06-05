package aof

import (
	"fmt"
	"strconv"
	"strings"
)

// parseYAMLSubset intentionally supports only the small, human-editable DSL that
// metagent emits from `init`. JSON remains the canonical lossless format.
func parseYAMLSubset(src string) (Morphogen, error) {
	var m Morphogen
	var section string
	var current *Capability
	var currentRule *EvolutionRule
	var inConfig bool

	lines := strings.Split(src, "\n")
	for i := 0; i < len(lines); i++ {
		raw := strings.TrimRight(lines[i], " \t\r")
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || line == "morphogen:" {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		if indent == 2 && strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			inConfig = false
			continue
		}
		if indent == 2 {
			k, v, ok := splitKV(line)
			if ok && k == "version" {
				m.Version = scalar(v)
			}
			continue
		}
		if section == "identity" && indent == 4 {
			k, v, ok := splitKV(line)
			if !ok {
				continue
			}
			switch k {
			case "domain":
				m.Identity.Domain = scalar(v)
			case "ontology_ref":
				m.Identity.OntologyRef = scalar(v)
			case "description":
				m.Identity.Description = scalar(v)
			}
			continue
		}
		if section == "capabilities" {
			if indent == 4 && strings.HasPrefix(line, "- ") {
				m.Capabilities = append(m.Capabilities, Capability{Config: map[string]string{}})
				current = &m.Capabilities[len(m.Capabilities)-1]
				inConfig = false
				k, v, ok := splitKV(strings.TrimSpace(strings.TrimPrefix(line, "- ")))
				if ok {
					assignCapability(current, k, v)
				}
				continue
			}
			if current == nil {
				return m, fmt.Errorf("line %d: capability field before list item", i+1)
			}
			if indent == 6 {
				k, v, ok := splitKV(line)
				if !ok {
					continue
				}
				if k == "config" {
					inConfig = true
					continue
				}
				inConfig = false
				assignCapability(current, k, v)
				continue
			}
			if inConfig && indent >= 8 {
				k, v, ok := splitKV(line)
				if ok {
					current.Config[k] = scalar(v)
				}
				continue
			}
		}
		if section == "evolution" {
			if indent == 4 {
				k, v, ok := splitKV(line)
				if !ok {
					continue
				}
				switch k {
				case "inherit":
					m.Evolution.Inherit = scalar(v) == "true"
				case "max_generations":
					m.Evolution.MaxGenerations, _ = strconv.Atoi(scalar(v))
				}
				continue
			}
			if indent == 6 && strings.HasPrefix(line, "- ") {
				m.Evolution.Rules = append(m.Evolution.Rules, EvolutionRule{})
				currentRule = &m.Evolution.Rules[len(m.Evolution.Rules)-1]
				k, v, ok := splitKV(strings.TrimSpace(strings.TrimPrefix(line, "- ")))
				if ok {
					assignRule(currentRule, k, v)
				}
				continue
			}
			if indent == 8 && currentRule != nil {
				k, v, ok := splitKV(line)
				if ok {
					assignRule(currentRule, k, v)
				}
			}
		}
	}
	return m, nil
}

func splitKV(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
}

func scalar(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'")
	return s
}

func listScalar(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" || s == "[]" {
		return nil
	}
	s = strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := scalar(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func assignCapability(c *Capability, k, v string) {
	sv := scalar(v)
	switch k {
	case "id":
		c.ID = sv
	case "type":
		c.Type = CrystalType(sv)
	case "activation":
		c.Activation = ActivationMode(sv)
	case "dependencies":
		c.Dependencies = listScalar(v)
	}
}

func assignRule(r *EvolutionRule, k, v string) {
	sv := scalar(v)
	switch k {
	case "trigger":
		r.Trigger = sv
	case "action":
		r.Action = sv
	case "priority":
		r.Priority, _ = strconv.Atoi(sv)
	}
}
