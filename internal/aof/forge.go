package aof

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func Forge(domain string, configs []CrystalConfig) ([]RuntimeCrystal, error) {
	out := make([]RuntimeCrystal, 0, len(configs))
	for _, cfg := range configs {
		state := "latent"
		if cfg.Activation.Mode == ActivationAlways {
			state = "active"
		}
		payload := cloneMap(cfg.Source)
		payload["binder"] = binderFor(cfg.Type)
		payload["summary"] = crystalSummary(cfg)
		out = append(out, RuntimeCrystal{
			ID:               cfg.ID,
			Type:             cfg.Type,
			Version:          fingerprint(cfg.ID + string(cfg.Type) + strings.Join(cfg.Synthesis.Inputs, ","))[:12],
			State:            state,
			Payload:          payload,
			ActivationPolicy: cfg.Activation,
			Meta: CrystalMeta{
				CreatedAt:        time.Now().UTC(),
				SourceMorphogen:   domain,
				Confidence:        0.80,
				SynthesisStrategy: cfg.Synthesis.Strategy,
			},
		})
	}
	return out, nil
}

func Differentiate(baseKernelRef, morphogenRef string, m Morphogen, crystals []RuntimeCrystal) AgentInstance {
	return AgentInstance{
		ID:            "agent-" + fingerprint(m.Identity.Domain+m.Version)[:10],
		Domain:        m.Identity.Domain,
		Description:   m.Identity.Description,
		BaseKernelRef: baseKernelRef,
		MorphogenRef:  morphogenRef,
		Crystals:      crystals,
		Lineage:       Lineage{Generation: 1},
	}
}

func BuildSystemPrompt(agent AgentInstance) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Specialized Agent: %s\n\n", agent.Domain)
	if agent.Description != "" {
		fmt.Fprintf(&b, "%s\n\n", agent.Description)
	}
	b.WriteString("## Operating Contract\n")
	b.WriteString("- Internalize active crystals as durable capabilities, not as external snippets.\n")
	b.WriteString("- Prefer domain-specific reasoning from knowledge and protocol crystals.\n")
	b.WriteString("- Respect belief crystals as hard capability boundaries.\n")
	b.WriteString("- Activate contextual crystals only when the user task matches their source or summary.\n")
	b.WriteString("- Ask for confirmation before crossing manual activation or uncertain tool boundaries.\n\n")
	b.WriteString("## Crystals\n")
	for _, c := range agent.Crystals {
		fmt.Fprintf(&b, "- %s (%s, %s): %s\n", c.ID, c.Type, c.State, c.Payload["summary"])
	}
	return b.String()
}

func SelectCrystals(agent AgentInstance, input string) []RuntimeCrystal {
	terms := strings.Fields(strings.ToLower(input))
	var selected []RuntimeCrystal
	for _, c := range agent.Crystals {
		if c.State == "active" || c.ActivationPolicy.Mode == ActivationAlways {
			selected = append(selected, c)
			continue
		}
		haystack := strings.ToLower(c.ID + " " + c.Payload["summary"] + " " + c.Payload["keywords"])
		for _, t := range terms {
			if len(t) > 2 && strings.Contains(haystack, t) {
				selected = append(selected, c)
				break
			}
		}
	}
	return selected
}

func RenderOfflineRun(agent AgentInstance, input string) string {
	selected := SelectCrystals(agent, input)
	var b strings.Builder
	fmt.Fprintf(&b, "Agent %s received task:\n%s\n\n", agent.Domain, input)
	b.WriteString("Activated crystals:\n")
	if len(selected) == 0 {
		b.WriteString("- none; fall back to base kernel and ask clarifying questions if needed\n")
	} else {
		for _, c := range selected {
			fmt.Fprintf(&b, "- %s (%s): %s\n", c.ID, c.Type, c.Payload["summary"])
		}
	}
	b.WriteString("\nSuggested response frame:\n")
	b.WriteString("\n1. Restate the domain-specific objective.")
	b.WriteString("\n2. Apply active crystal constraints and knowledge.")
	b.WriteString("\n3. Call tools only if the corresponding tool crystal is active.")
	b.WriteString("\n4. Surface uncertainty when belief boundaries are missing or weak.\n")
	return b.String()
}

func binderFor(t CrystalType) string {
	switch t {
	case CrystalKnowledge:
		return "context-injection+reasoning-guide"
	case CrystalTool:
		return "function-call-contract"
	case CrystalProtocol:
		return "dialogue-policy"
	case CrystalMemory:
		return "episodic-retriever"
	case CrystalBelief:
		return "meta-cognition-filter"
	default:
		return "context"
	}
}

func crystalSummary(c CrystalConfig) string {
	for _, k := range []string{"summary", "description", "content", "path", "schema"} {
		if v := c.Source[k]; v != "" {
			if len(v) > 160 {
				return v[:157] + "..."
			}
			return v
		}
	}
	return c.Synthesis.Strategy
}

func fingerprint(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
