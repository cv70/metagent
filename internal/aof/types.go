package aof

import "time"

// CrystalType classifies the kind of capability being internalized by an agent.
type CrystalType string

const (
	CrystalKnowledge CrystalType = "knowledge"
	CrystalTool      CrystalType = "tool"
	CrystalProtocol  CrystalType = "protocol"
	CrystalMemory    CrystalType = "memory"
	CrystalBelief    CrystalType = "belief"
)

// ActivationMode describes when a crystal should be active.
type ActivationMode string

const (
	ActivationAlways     ActivationMode = "always"
	ActivationContextual ActivationMode = "contextual"
	ActivationManual     ActivationMode = "manual"
)

type Morphogen struct {
	Version      string       `json:"version"`
	Identity     Identity     `json:"identity"`
	Capabilities []Capability `json:"capabilities"`
	Evolution    Evolution    `json:"evolution"`
}

type Identity struct {
	Domain      string `json:"domain"`
	OntologyRef string `json:"ontology_ref"`
	Description string `json:"description"`
}

type Capability struct {
	ID           string            `json:"id"`
	Type         CrystalType       `json:"type"`
	Config       map[string]string `json:"config"`
	Activation   ActivationMode    `json:"activation"`
	Dependencies []string          `json:"dependencies"`
}

type Evolution struct {
	Rules          []EvolutionRule `json:"rules"`
	Inherit        bool            `json:"inherit"`
	MaxGenerations int             `json:"max_generations"`
}

type EvolutionRule struct {
	Trigger  string `json:"trigger"`
	Action   string `json:"action"`
	Priority int    `json:"priority"`
}

type CrystalConfig struct {
	ID           string            `json:"id"`
	Type         CrystalType       `json:"type"`
	Source       map[string]string `json:"source"`
	Synthesis    SynthesisConfig   `json:"synthesis"`
	Activation   ActivationPolicy  `json:"activation"`
	Dependencies []string          `json:"dependencies"`
}

type SynthesisConfig struct {
	Strategy string   `json:"strategy"`
	Inputs   []string `json:"inputs"`
}

type ActivationPolicy struct {
	Mode      ActivationMode `json:"mode"`
	Condition string         `json:"condition,omitempty"`
}

type RuntimeCrystal struct {
	ID               string            `json:"id"`
	Type             CrystalType       `json:"type"`
	Version          string            `json:"version"`
	State            string            `json:"state"`
	Payload          map[string]string `json:"payload"`
	ActivationPolicy ActivationPolicy  `json:"activation_policy"`
	Meta             CrystalMeta       `json:"meta"`
}

type CrystalMeta struct {
	CreatedAt        time.Time `json:"created_at"`
	SourceMorphogen   string    `json:"source_morphogen"`
	Confidence        float64   `json:"confidence"`
	UsageCount        int       `json:"usage_count"`
	LastUsed          time.Time `json:"last_used,omitempty"`
	SynthesisStrategy string    `json:"synthesis_strategy"`
}

type AgentInstance struct {
	ID            string           `json:"id"`
	Domain        string           `json:"domain"`
	Description   string           `json:"description"`
	BaseKernelRef string           `json:"base_kernel_ref"`
	MorphogenRef  string           `json:"morphogen_ref"`
	Crystals      []RuntimeCrystal `json:"crystals"`
	Lineage       Lineage          `json:"lineage"`
}

type Lineage struct {
	Generation     int      `json:"generation"`
	ParentInstance string   `json:"parent_instance,omitempty"`
	Mutations       []string `json:"mutations,omitempty"`
}

type CompileResult struct {
	Morphogen Morphogen       `json:"morphogen"`
	Crystals  []CrystalConfig `json:"crystals"`
}
