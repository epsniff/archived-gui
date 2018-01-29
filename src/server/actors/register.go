package actors

func RegisterActorsDefs() {
	// Define how actors are created.
	server.RegisterDef("leader", func(_ []byte) (grid.Actor, error) { return &LeaderActor{client: client}, nil })
}