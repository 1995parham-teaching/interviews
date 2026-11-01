package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultsAreValid(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestFlagsOverrideTheFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "relay.json")
	os.WriteFile(file, []byte(`{
		"id": "from-file",
		"backends": [{"addr": "10.0.0.1:80", "weight": 2}],
		"tombstone_ttl": "1h",
		"retries": 3
	}`), 0o644)

	c, err := Load([]string{"-config", file, "-retries", "0", "-listen", "127.0.0.1:9999"})
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != "from-file" || c.Retries != 0 || c.Listen != "127.0.0.1:9999" {
		t.Fatalf("got %+v", c)
	}
	if c.TombstoneTTL.D() != time.Hour || c.Backends[0].Weight != 2 {
		t.Fatalf("file values lost: %+v", c)
	}
	if c.Fall != 3 {
		t.Fatalf("default lost: fall %d", c.Fall)
	}
}

func TestIDDefaultsToListen(t *testing.T) {
	c, err := Load([]string{"-listen", "127.0.0.1:8005"})
	if err != nil || c.ID != "127.0.0.1:8005" {
		t.Fatalf("got %q, %v", c.ID, err)
	}
}

func TestParseBackends(t *testing.T) {
	bs, err := ParseBackends("10.0.0.1:80, 10.0.0.2:80=3,")
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 2 || bs[0].Weight != 1 || bs[1].Weight != 3 {
		t.Fatalf("got %+v", bs)
	}
	if _, err := ParseBackends("10.0.0.1:80=x"); err == nil {
		t.Fatal("accepted a bad weight")
	}
}

func TestValidateRejects(t *testing.T) {
	for name, args := range map[string][]string{
		"bad backend":    {"-backends", "nohost"},
		"weight too big": {"-backends", "10.0.0.1:80=11"},
		"zero interval":  {"-gossip-interval", "0s"},
		"bad cidr":       {"-trusted-proxies", "10.0.0.0/33"},
		"negative retry": {"-retries", "-1"},
	} {
		if _, err := Load(args); err == nil {
			t.Errorf("%s: accepted %v", name, args)
		}
	}
}
