package config

import (
	"strings"
	"testing"
	"time"

	"github.com/rwstauner/ynetd/procman"
	"github.com/rwstauner/ynetd/service"
)

func TestMakeServices(t *testing.T) {
	pm := procman.New()
	services, err := MakeServices(Config{
		Services: []service.Config{
			{
				Proxy: map[string]string{
					":4000": "localhost:4001",
				},
				Command: []string{"foo", "bar"},
				Timeout: "3s",
			},
		},
	}, pm)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(services) != 1 {
		t.Errorf("expected 1 service, got %d", len(services))
	}

	svc := services[0]
	if svc.Command.String() != "[foo bar] (running: false)" {
		t.Errorf("unexpected command: %s", svc.Command)
	}
	if svc.Timeout != 3*time.Second {
		t.Errorf("unexpected timeout: %s", svc.Timeout)
	}
}

func TestMakeServicesError(t *testing.T) {
	pm := procman.New()
	_, err := MakeServices(Config{
		Services: []service.Config{
			{
				Timeout: "3",
			},
		},
	}, pm)

	if err == nil {
		t.Errorf("expected error, got none")
	} else if !strings.Contains(err.Error(), "duration") {
		t.Errorf("unexpected error: %s", err)
	}
}
