package actorpool

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lytics/grid"
)

func TestPeerOptimisticallyLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ap := New(true)
	for p := range peers {
		ap.OptimisticallyLive(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != dead {
			t.Fatal("expected live peer")
		}
		if optimisticState != live {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ap := New(true)

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != live {
			t.Fatal("expected live peer")
		}
	}
}
func TestPeerOptimisticallyLiveToLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ap := New(true)
	for p := range peers {
		ap.OptimisticallyLive(p)
	}
	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != live {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerOptimisticallyDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New(true)

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.OptimisticallyDead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != dead {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New(true)

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != dead {
			t.Fatal("expected live peer")
		}
		if optimisticState != dead {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerOptimisticallyDeadToDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New(true)

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.OptimisticallyDead(p)
	}
	for p := range peers {
		ap.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != dead {
			t.Fatal("expected dead peer")
		}
		if optimisticState != dead {
			t.Fatal("expected dead peer")
		}
	}
}

func TestDeadPeerThatWasNeverLive(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.Dead("peer0")
	state, optimisticState := ap.peerState.State("peer0")
	if state != dead {
		t.Fatal("expected dead peer")
	}
	if optimisticState != dead {
		t.Fatal("expected dead peer")
	}
}

func TestRegisterThenUnregister(t *testing.T) {
	t.Parallel()

	ap := New(true)

	ap.Live("peer0")
	ap.Register("writer0", "peer0")

	if ap.NumRegistered() != 1 {
		t.Fatal("expected 1 registered")
	}
	if ap.NumRegisteredOn("peer0") != 1 {
		t.Fatal("expected 1 registered on peer")
	}
	if !ap.IsRegistered("writer0") {
		t.Fatal("expected registered")
	}

	ap.Unregister("writer0")
	if ap.NumRegistered() != 0 {
		t.Fatal("expected 0 registered")
	}
	if ap.NumRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered on peer")
	}
	if ap.IsRegistered("writer0") {
		t.Fatal("expected not-registered")
	}
}

func TestOptimisticallyRegisterUnregister(t *testing.T) {
	t.Parallel()

	ap := New(true)

	ap.Live("peer0")
	ap.OptimisticallyRegister("writer0", "peer0")

	if ap.NumOptimisticallyRegistered() != 1 {
		t.Fatal("expected 1 optimistically registered")
	}
	if ap.NumOptimisticallyRegisteredOn("peer0") != 1 {
		t.Fatal("expected 1 registered on peer")
	}
	if !ap.IsOptimisticallyRegistered("writer0") {
		t.Fatal("expected registered")
	}

	ap.OptimisticallyUnregister("writer0")
	if ap.NumRegistered() != 0 {
		t.Fatal("expected 0 optimistically registered")
	}
	if ap.NumOptimisticallyRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered on peer")
	}
	if ap.IsOptimisticallyRegistered("writer0") {
		t.Fatal("expected not-registered")
	}
}

/*
func TestRelocateOne(t *testing.T) {
	ap := New(true)

	ap.SetRequired(grid.NewActorStart("writer0"))
	ap.SetRequired(grid.NewActorStart("writer1"))

	ap.Live("peer0")
	ap.Live("peer1")

	ap.Register("writer0", "peer0")
	ap.Register("writer1", "peer0")

	plan := ap.Relocate()
	if len(plan.Relocations) != 1 {
		t.Fatal("expected one relocation")
	}
}

func TestRelocateZero(t *testing.T) {
	ap := New(true)

	ap.SetRequired(grid.NewActorStart("writer0"))
	ap.SetRequired(grid.NewActorStart("writer1"))

	ap.Live("peer0")
	ap.Live("peer1")

	ap.Register("writer0", "peer0")
	ap.Register("writer1", "peer1")

	plan := ap.Relocate()
	if len(plan.Relocations) != 0 {
		t.Fatal("expected zero relocations")
	}
}

func TestRelocateOddPeersEvenActors(t *testing.T) {
	ap := New(true)

	for _, w := range []string{
		"writer0",
		"writer1",
		"writer2",
		"writer3",
	} {
		def := grid.NewActorStart(w)
		def.Type = "writer"
		ap.SetRequired(def)
	}
	for _, p := range []string{
		"peer0",
		"peer1",
		"peer2",
	} {
		ap.Live(p)
	}

	p1, _ := ap.ByHash("writer0")
	ap.Register("writer0", p1)
	p2, _ := ap.ByHash("writer1")
	ap.Register("writer1", p2)
	p3, _ := ap.ByHash("writer2")
	ap.Register("writer2", p3)
	p4, _ := ap.ByHash("writer3")
	ap.Register("writer3", p4)

	plan := ap.Relocate()
	if len(plan.Relocations) != 0 {
		t.Fatalf("expected zero relocations: got %v", plan.Relocations)
	}
}

func TestRelocateEvenPeersOddActors(t *testing.T) {
	ap := New(true)

	for _, w := range []string{
		"writer0",
		"writer1",
		"writer2",
	} {
		def := grid.NewActorStart(w)
		def.Type = "writer"
		ap.SetRequired(def)
	}
	for i := 0; i < 42; i++ {
		p := fmt.Sprintf("peer%s", i)
		ap.Live(p)
	}

	pickBadPeer := func(name string) string {
		goodpeer, _ := ap.ByHash(name)
		for peer, _ := range ap.peers {
			if peer != goodpeer {
				return peer
			}
		}
		t.Fatalf("unreachable???")
		return ""
	}
	//mix up the assignments so we get a relocation plan
	ap.Register("writer0", pickBadPeer("writer0"))
	ap.Register("writer1", pickBadPeer("writer1"))
	ap.Register("writer2", pickBadPeer("writer2"))

	plan := ap.Relocate()
	if len(plan.Relocations) != 3 {
		t.Fatalf("expected 3 relocations: got %v", plan.Relocations)
	}
}

func TestRelocateWithDeadPeer(t *testing.T) {
	ap := New(true)

	for _, p := range []string{
		"peer0",
		"peer1",
		"peer2",
	} {
		ap.Live(p)
	}

	for _, w := range []string{
		"writer0",
		"writer1",
		"writer2",
		"writer3",
		"writer4",
		"writer5",
		"writer6",
		"writer7",
	} {
		def := grid.NewActorStart(w)
		def.Type = "writer"
		ap.SetRequired(def)

		p, err := ap.ByHash(def.Name)
		if err != nil {
			t.Fatal(err)
		}

		ap.Register(w, p)
	}

	plan := ap.Relocate()
	if len(plan.Relocations) != 0 {
		t.Fatalf("expected zero relocations[1]: got %v", plan.Relocations)
	}

	t.Logf("locations: ")
	for w, pi := range ap.peers {
		t.Logf("   %v --> %v ", w, pi.registered)
	}

	peerToKill, _ := ap.peers["peer1"]
	ap.Unregister("writer2") //writer2 should be on peer0 at this point.
	ap.Dead(peerToKill.name)
	expectedRelocatons := peerToKill.NumActors() //killing peer1 should free up :  [writer0 writer3 writer6]
	expectedRelocatons += 2                      // plus changing the size is going to cause [writer5, writer4] to relocate

	if len(ap.Missing()) != 1 {
		t.Fatal("expected writer2 to be missing")
	}
	for _, def := range ap.Missing() {
		if def.Name != "writer2" {
			t.Fatal("expected writer2 to be missing")
		}
	}

	plan = ap.Relocate()
	if len(plan.Relocations) != expectedRelocatons {
		d, _ := json.Marshal(ap.Status())
		t.Logf("%v", string(d))
		t.Logf("peerToKill:%v actorToUnReg:writer2 expectedRelocatons:%v", peerToKill, expectedRelocatons)
		t.Fatalf("expected %v relocations[3]: got %v", expectedRelocatons, plan.Relocations)
	}

	//move the actors over to their new peers
	for _, a := range plan.Relocations {
		ap.Unregister(a)
		p, err := ap.ByHash(a)
		if err != nil {
			t.Fatal(err)
		}
		ap.Register(a, p)
	}

	t.Logf("locations: ")
	for w, pi := range ap.peers {
		t.Logf("   %v --> %v ", w, pi.registered)
	}

	ap.Live(peerToKill.name)
	plan = ap.Relocate()
	//we expect the same number of relocations because the actors should move back to their
	// original peers after peer1 rejoins.
	if len(plan.Relocations) != expectedRelocatons {
		t.Logf("peerToKill:%v expectedRelocatons:%v peer_selectors:%v", peerToKill, expectedRelocatons, ap.selector.peers)
		t.Fatalf("expected %v relocations[4]: got %v", expectedRelocatons, plan.Relocations)
	}
}

func TestRegisterRemoveOptimisticallyRegistered(t *testing.T) {
	ap := New(true)

	ap.Live("peer0")
	ap.OptimisticallyRegister("writer0", "peer0")

	if 0 != ap.NumRegistered() {
		t.Fatal("expected 0 registered")
	}
	if 1 != ap.NumOptimisticallyRegistered() {
		t.Fatal("expected 1 optimistically registered")
	}

	ap.Register("writer0", "peer0")
	if 1 != ap.NumRegistered() {
		t.Fatal("expected 1 registered")
	}
	if 0 != ap.NumOptimisticallyRegistered() {
		t.Fatal("expected 0 optimistically registered")
	}
}
*/

func TestSetUnsetRequired(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.SetRequired(grid.NewActorStart("writer"))
	if !ap.IsRequired("writer") {
		t.Fatal("expected required")
	}
	ap.UnsetRequired("writer")
	if ap.IsRequired("writer") {
		t.Fatal("expected not-required")
	}
}

func TestSetMissingRequired(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.SetRequired(grid.NewActorStart("writer"))
	if !ap.IsRequired("writer") {
		t.Fatal("expected required")
	}

	starts := ap.Missing()
	if len(starts) != 1 {
		t.Fatalf("expected 1 missing, got:%v", len(starts))
	}

	ap.UnsetRequired("writer")
	if ap.IsRequired("writer") {
		t.Fatal("expected not-required")
	}

	starts = ap.Missing()
	if len(starts) != 0 {
		t.Fatalf("expected 0 missing, got:%v", len(starts))
	}

	for i := 0; i < 100; i++ {
		actorDef := grid.NewActorStart(fmt.Sprintf("writer-%d", i))
		id := actorDef.Name
		err := ap.SetRequired(actorDef)
		if err != nil {
			t.Fatalf("fot an setRequired error for actor:%+v err:%v", actorDef, err)
		}
		if !ap.IsRequired(id) {
			t.Fatalf("expected required for actor:%v ap.required:%v", id, ap.required)
		}
	}

	starts = ap.Missing()
	if len(starts) != 100 {
		t.Fatalf("expected 100 missing, got:%v", len(starts))
	}

	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("writer-%d", i)
		ap.UnsetRequired(id)
		if ap.IsRequired(id) {
			t.Fatalf("expected not-required")
		}
	}

	starts = ap.Missing()
	if len(starts) != 0 {
		t.Fatalf("expected 0 missing, got:%v", len(starts))
	}

}

func TestNoActorsOnNewPeer(t *testing.T) {
	t.Parallel()
	ap := New(true)

	if ap.NumRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered")
	}

	ap.Live("peer0")
	if ap.NumRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered")
	}
}

func TestNumRegisteredOn(t *testing.T) {
	t.Parallel()
	ap := New(true)

	if ap.NumRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered")
	}

	ap.Live("peer0")
	ap.Register("writer0", "peer0")
	if ap.NumRegisteredOn("peer0") != 1 {
		t.Fatal("expected 1 registered")
	}
}

func TestNumOptimisticallyRegisteredOn(t *testing.T) {
	t.Parallel()
	ap := New(true)

	if ap.NumOptimisticallyRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered")
	}

	ap.Live("peer0")
	ap.OptimisticallyRegister("writer0", "peer0")
	if ap.NumOptimisticallyRegisteredOn("peer0") != 1 {
		t.Fatal("expected 1 registered")
	}
}

func TestRegisterWithZeroPeers(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.Register("writer0", "peer0")
	if ap.NumRegistered() != 1 {
		t.Fatal("expected 1 registered")
	}
}

func TestOptimisticallyRegisterWithZeroPeers(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.OptimisticallyRegister("writer0", "peer0")
	if ap.NumOptimisticallyRegistered() != 1 {
		t.Fatal("expected 1 registered")
	}
}

func TestRegisterTwiceWithoutUnregister(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.Live("peer0")
	ap.Live("peer1")

	ap.Register("writer0", "peer0")
	ap.Register("writer0", "peer1")

	if ap.NumRegisteredOn("peer0") != 0 {
		t.Fatal("expected 0 registered on peer")
	}
	if ap.NumRegisteredOn("peer1") != 1 {
		t.Fatal("expected 1 registered on peer")
	}
}

func TestOptimisticallyRegisterTwiceWithoutUnregister(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.Live("peer0")
	ap.Live("peer1")

	ap.OptimisticallyRegister("writer0", "peer0")
	ap.OptimisticallyRegister("writer0", "peer1")
	p0 := ap.NumOptimisticallyRegisteredOn("peer0")
	p1 := ap.NumOptimisticallyRegisteredOn("peer1")

	if p0 != 0 {
		l := []string{}
		for _, p := range ap.peerState.peers {
			l = append(l, fmt.Sprintf("%v\n", p))
		}
		t.Fatalf("expected 0 on peer0: got: peer0:%v peer1:%v peers:\n%+v",
			p0, p1, strings.Join(l, "\n"))
	}
	if ap.NumOptimisticallyRegisteredOn("peer1") != 1 {
		t.Fatalf("expected 1 on peer1: got: peer0:%v peer1:%v", p0, p1)
	}
}

func TestUnregisterWithOptimisticallyRegistered(t *testing.T) {
	t.Parallel()
	ap := New(true)

	ap.Live("peer0")
	ap.OptimisticallyRegister("writer0", "peer0")
	if ap.NumOptimisticallyRegisteredOn("peer0") != 1 {
		t.Fatal("expected 1 registered on peer")
	}
	ap.Unregister("writer0")
	op0 := ap.NumOptimisticallyRegisteredOn("peer0")
	if op0 != 0 {
		l := []string{}
		for _, p := range ap.peerState.peers {
			l = append(l, fmt.Sprintf("%v\n", p))
		}
		t.Fatalf("expected 0 on peer0: got: peer0:%v peers:\n%+v",
			op0, strings.Join(l, "\n"))
	}
}

func TestIsValidName(t *testing.T) {
	t.Parallel()
	if isValidName("") {
		t.Fatal("empty string is not valid name")
	}
	if !isValidName("actor-foo") {
		t.Fatal("non empty string is valid name")
	}
}

func TestEmptyActorStart(t *testing.T) {
	t.Parallel()
	ap := New(true)

	if ErrInvalidName != ap.SetRequired(grid.NewActorStart("")) {
		t.Fatal("expected invalid name error")
	}
}

func TestEmptyPeerLiveDead(t *testing.T) {
	t.Parallel()
	ap := New(true)

	checkNoEntry := func() {
		if _, ok := ap.peerState.peers[""]; ok {
			t.Fatal("expected no peer info")
		}
	}

	ap.Live("")
	checkNoEntry()

	ap.Dead("")
	checkNoEntry()

	ap.OptimisticallyLive("")
	checkNoEntry()

	ap.OptimisticallyDead("")
	checkNoEntry()
}

func TestEmptyActorRegisterUnregister(t *testing.T) {
	t.Parallel()
	ap := New(true)

	const validPeer = "peer-1"

	checkNoEntry := func() {
		if _, ok := ap.peerState.peers[validPeer]; ok {
			t.Fatal("expected no peer info")
		}
	}

	ap.Register("", validPeer)
	checkNoEntry()

	ap.Unregister("")
	checkNoEntry()

	ap.OptimisticallyRegister("", validPeer)
	checkNoEntry()

	ap.OptimisticallyUnregister("")
	checkNoEntry()
}
