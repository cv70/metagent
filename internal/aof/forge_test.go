package aof

import "testing"

func TestForgeAndSelectCrystals(t *testing.T) {
	m := Morphogen{Identity: Identity{Domain: "review", Description: "Review agent"}}
	configs := []CrystalConfig{
		{ID: "always", Type: CrystalKnowledge, Source: map[string]string{"summary": "general review"}, Activation: ActivationPolicy{Mode: ActivationAlways}},
		{ID: "race", Type: CrystalKnowledge, Source: map[string]string{"summary": "goroutine race detector", "keywords": "goroutine race"}, Activation: ActivationPolicy{Mode: ActivationContextual}},
	}
	crystals, err := Forge(m.Identity.Domain, configs)
	if err != nil {
		t.Fatal(err)
	}
	agent := Differentiate("base", "spec", m, crystals)
	selected := SelectCrystals(agent, "please inspect goroutine behavior")
	if len(selected) != 2 {
		t.Fatalf("selected %d crystals", len(selected))
	}
	prompt := BuildSystemPrompt(agent)
	if prompt == "" {
		t.Fatal("empty prompt")
	}
}
