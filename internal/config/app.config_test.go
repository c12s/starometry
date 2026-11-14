package config

import "testing"

func TestNatsAddressPrefersNatsAddress(t *testing.T) {
	t.Setenv("NATS_URL", "nats")
	t.Setenv("NATS_PORT", "4222")
	if got := natsAddress(); got != "nats:4222" {
		t.Fatalf("pair: got %q", got)
	}
	t.Setenv("NATS_ADDRESS", "192.168.122.1:4222")
	if got := natsAddress(); got != "192.168.122.1:4222" {
		t.Fatalf("NATS_ADDRESS: got %q", got)
	}
}

func TestUnsetSourceHasNoAddress(t *testing.T) {
	t.Setenv("CADVISOR_URL", "")
	t.Setenv("CADVISOR_PORT", "8081")
	if got := hostPort("CADVISOR_URL", "CADVISOR_PORT"); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestNodeIDIsTrimmed(t *testing.T) {
	c := &AppConfig{}
	c.SetNodeID("3077ffc7-aaaa\n")
	if c.GetNodeID() != "3077ffc7-aaaa" {
		t.Fatalf("got %q", c.GetNodeID())
	}
}
