package placement

func NewRelocationPlan(actorType string, total, average int) *RelocationPlan {
	return &RelocationPlan{
		ActorType:   actorType,
		Relocations: []string{}, //actor name. aka the mailbox
	}
}

type RelocationPlan struct {
	ActorType   string
	Relocations []string //actor name. aka the mailbox
}
