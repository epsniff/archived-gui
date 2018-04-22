package grid

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/lytics/grid/testetcd"
)

type busyActor struct {
	ready  chan bool
	server *Server
}

func (a *busyActor) Act(c context.Context) {
	name, err := ContextActorName(c)
	if err != nil {
		return
	}

	mailbox, err := NewMailbox(a.server, name, 0)
	if err != nil {
		return
	}
	defer mailbox.Close()

	// Don't bother listening
	// to the mailbox, too
	// busy.
	a.ready <- true
	<-c.Done()
}

type echoActor struct {
	ready  chan bool
	server *Server
}

func (a *echoActor) Act(c context.Context) {
	name, err := ContextActorName(c)
	if err != nil {
		return
	}

	mailbox, err := NewMailbox(a.server, name, 1)
	if err != nil {
		return
	}
	defer mailbox.Close()

	a.ready <- true
	for {
		select {
		case <-c.Done():
			return
		case req := <-mailbox.C:
			req.Respond(req.Msg())
		}
	}
}

func init() {
	Register(EchoMsg{})
}

func TestNewClient(t *testing.T) {
	etcd := testetcd.StartAndConnect(t)

	client, err := NewClient(etcd, ClientCfg{Namespace: newNamespace()})
	if err != nil {
		t.Fatal(err)
	}
	client.Close()
}

func TestNewClientWithNilEtcd(t *testing.T) {
	_, err := NewClient(nil, ClientCfg{Namespace: newNamespace()})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClientClose(t *testing.T) {
	// Start etcd.
	etcd := testetcd.StartAndConnect(t)

	// Create client.
	client, err := NewClient(etcd, ClientCfg{Namespace: newNamespace()})
	if err != nil {
		t.Fatal(err)
	}
	client.cs = newClientStats()

	// The type clientAndConn checks if it is nil
	// in its close method, and returns an error.
	client.clientsAndConns["mock"] = nil
	err = client.Close()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClientRequestWithUnregisteredMailbox(t *testing.T) {
	const timeout = 2 * time.Second

	// Bootstrap.
	etcd, server, client := bootstrapClientTest(t)
	defer etcd.Close()
	defer server.Stop()
	defer client.Close()

	// Set client stats.
	client.cs = newClientStats()

	// Send a request to some random name.
	res, err := client.Request(timeout, "mock", NewActorStart("mock"))
	if err != ErrUnregisteredMailbox {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatal(res)
	}

	if v := client.cs.counters[numErrUnregisteredMailbox]; v == 0 {
		t.Fatal("expected non-zero error count")
	}
}

func TestClientRequestWithUnknownMailbox(t *testing.T) {
	const timeout = 2 * time.Second

	// Bootstrap.
	etcd, server, client := bootstrapClientTest(t)
	defer etcd.Close()
	defer server.Stop()
	defer client.Close()

	// Set client stats.
	client.cs = newClientStats()

	// Place a bogus entry in etcd with
	// a matching name.
	timeoutC, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := server.registry.Register(timeoutC, client.cfg.Namespace+".mailbox.mock")
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	// Send a request to some random name.
	res, err := client.Request(timeout, "mock", NewActorStart("mock"))
	if err == nil {
		t.Fatal("expected error")
	}
	if res != nil {
		t.Fatal(res)
	}
	if !strings.Contains(err.Error(), ErrUnknownMailbox.Error()) {
		t.Fatal(err)
	}

	if v := client.cs.counters[numErrUnknownMailbox]; v == 0 {
		t.Fatal("expected non-zero error count")
	}
}

func TestClientWithRunningReceiver(t *testing.T) {
	const timeout = 2 * time.Second
	expected := &EchoMsg{"testing 1, 2, 3"}

	// Bootstrap.
	etcd, server, client := bootstrapClientTest(t)
	defer etcd.Close()
	defer server.Stop()
	defer client.Close()

	// Set client stats.
	client.cs = newClientStats()

	// Create echo actor.
	a := &echoActor{ready: make(chan bool)}

	// Set grid definition.
	server.RegisterDef("echo", func(_ []byte) (Actor, error) { return a, nil })

	// Set server on echo actor.
	a.server = server

	// Discover some peers.
	peers, err := client.Query(timeout, Peers)
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 1 {
		t.Fatal("expected 1 peer")
	}

	// Start the echo actor on the first peer.
	res, err := client.Request(timeout, peers[0].Name(), NewActorStart("echo"))
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response")
	}

	// Wait for echo actor to start.
	<-a.ready

	// Make a request to echo actor.
	res, err = client.Request(timeout, "echo", expected)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response")
	}

	// Expect the same string back as a response.
	switch res := res.(type) {
	case *EchoMsg:
		if res.Msg != expected.Msg {
			t.Fatalf("expected: %v, received: %v", expected, res)
		}
	default:
		t.Fatalf("expected type: string, received type: %T", res)
	}

	if v := client.cs.counters[numErrConnectionUnavailable]; v != 0 {
		t.Fatal("expected zero error count")
	}
}

func TestClientWithErrConnectionIsUnregistered(t *testing.T) {
	const timeout = 2 * time.Second
	expected := &EchoMsg{"testing 1, 2, 3"}

	// Bootstrap.
	etcd, server, client := bootstrapClientTest(t)
	defer etcd.Close()
	defer client.Close()

	// Set client stats.
	client.cs = newClientStats()

	// Create echo actor.
	a := &echoActor{ready: make(chan bool)}

	// Set grid definition.
	server.RegisterDef("echo", func(_ []byte) (Actor, error) { return a, nil })

	// Set server on echo actor.
	a.server = server

	// Discover some peers.
	peers, err := client.Query(timeout, Peers)
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 1 {
		t.Fatal("expected 1 peer")
	}

	// Start the echo actor on the first peer.
	res, err := client.Request(timeout, peers[0].Name(), NewActorStart("echo"))
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response")
	}

	// Wait for echo actor to start.
	<-a.ready

	// Make a request to echo actor.
	res, err = client.Request(timeout, "echo", expected)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response")
	}

	// Stop the server.
	server.Stop()

	// Wait for the gRPC to be really stopped.
	time.Sleep(timeout)

	// Make the request again.
	res, err = client.Request(timeout, "echo", expected)
	if err == nil {
		t.Fatal("expected error")
	}
	if res != nil {
		t.Fatal(res)
	}
	if !strings.Contains(err.Error(), "unregistered mailbox") {
		t.Fatal(err)
	}

	if v := client.cs.counters[numErrUnregisteredMailbox]; v == 0 {
		t.Fatal("expected non-zero error count")
	}
}

func TestClientWithBusyReceiver(t *testing.T) {
	const timeout = 2 * time.Second
	expected := &EchoMsg{"testing 1, 2, 3"}

	// Bootstrap.
	etcd, server, client := bootstrapClientTest(t)
	defer etcd.Close()
	defer server.Stop()
	defer client.Close()

	// Set client stats.
	client.cs = newClientStats()

	// Create busy actor.
	a := &busyActor{ready: make(chan bool)}

	server.RegisterDef("busy", func(_ []byte) (Actor, error) { return a, nil })

	// Set server on busy actor.
	a.server = server

	// Discover some peers.
	peers, err := client.Query(timeout, Peers)
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 1 {
		t.Fatal("expected 1 peer")
	}

	// Start the busy actor on the first peer.
	res, err := client.Request(timeout, peers[0].Name(), NewActorStart("busy"))
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response")
	}

	// Wait for busy actor to start.
	<-a.ready

	// Make a request to busy actor.
	res, err = client.Request(timeout, "busy", expected)
	if err == nil {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatal("expected response")
	}
	if !strings.Contains(err.Error(), ErrReceiverBusy.Error()) {
		t.Fatal(err)
	}
}

func TestClientStats(t *testing.T) {
	cs := newClientStats()
	cs.Inc(numGetWireClient)
	cs.Inc(numDeleteAddress)
	if cs.counters[numGetWireClient] != 1 {
		t.Fatal("expected count of 1")
	}
	if cs.counters[numDeleteAddress] != 1 {
		t.Fatal("expected count of 1")
	}
	switch cs.String() {
	case "numGetWireClient:1, numDeleteAddress:1":
	case "numDeleteAddress:1, numGetWireClient:1":
	default:
		t.Fatal("expected string: 'numGetWireClient:1'")
	}
}

func TestNilClientStats(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("expected no panic")
		}
	}()
	var cs *clientStats
	cs.Inc(numGetWireClient)
}

func bootstrapClientTest(t *testing.T) (*clientv3.Client, *Server, *Client) {
	// Namespace for test.
	namespace := newNamespace()

	// Start etcd.
	etcd := testetcd.StartAndConnect(t)

	// Logger for actors.
	logger := log.New(os.Stderr, namespace+": ", log.LstdFlags)

	// Create the server.
	server, err := NewServer(etcd, ServerCfg{Namespace: namespace, Logger: logger})
	if err != nil {
		t.Fatal(err)
	}

	// Create the listener on a random port.
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	// Start the server in the background.
	done := make(chan error, 1)
	go func() {
		err = server.Serve(lis)
		if err != nil {
			done <- err
		}
	}()
	time.Sleep(2 * time.Second)

	// Create a grid client.
	client, err := NewClient(etcd, ClientCfg{Namespace: namespace, Logger: logger})
	if err != nil {
		t.Fatal(err)
	}

	return etcd, server, client
}
