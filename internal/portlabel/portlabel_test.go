package portlabel_test

import (
	"testing"

	"github.com/user/portwatch/internal/portlabel"
)

func TestLabel_WellKnownSSH(t *testing.T) {
	l := portlabel.New()
	lbl := l.Label(22)
	if lbl.Service != "ssh" {
		t.Errorf("expected service ssh, got %s", lbl.Service)
	}
	if lbl.Severity != portlabel.SeverityHigh {
		t.Errorf("expected severity high, got %s", lbl.Severity)
	}
}

func TestLabel_Telnet_IsCritical(t *testing.T) {
	l := portlabel.New()
	lbl := l.Label(23)
	if lbl.Severity != portlabel.SeverityCritical {
		t.Errorf("expected critical for telnet, got %s", lbl.Severity)
	}
}

func TestLabel_DatabasePorts_AreCritical(t *testing.T) {
	l := portlabel.New()
	for _, port := range []int{3306, 5432} {
		lbl := l.Label(port)
		if lbl.Severity != portlabel.SeverityCritical {
			t.Errorf("port %d: expected critical, got %s", port, lbl.Severity)
		}
	}
}

func TestLabel_DynamicPort_IsInfo(t *testing.T) {
	l := portlabel.New()
	lbl := l.Label(55000)
	if lbl.Severity != portlabel.SeverityInfo {
		t.Errorf("expected info for dynamic port, got %s", lbl.Severity)
	}
}

func TestLabel_SystemPort_IsMedium(t *testing.T) {
	l := portlabel.New()
	lbl := l.Label(999)
	if lbl.Severity != portlabel.SeverityMedium {
		t.Errorf("expected medium for unknown system port, got %s", lbl.Severity)
	}
}

func TestAdd_OverridesDefault(t *testing.T) {
	l := portlabel.New()
	l.Add(8080, "internal-api", portlabel.SeverityHigh, "Internal API gateway")
	lbl := l.Label(8080)
	if lbl.Service != "internal-api" {
		t.Errorf("expected internal-api, got %s", lbl.Service)
	}
	if lbl.Severity != portlabel.SeverityHigh {
		t.Errorf("expected high, got %s", lbl.Severity)
	}
}

func TestAdd_DoesNotAffectOtherPorts(t *testing.T) {
	l := portlabel.New()
	l.Add(9999, "custom", portlabel.SeverityCritical, "custom service")
	lbl := l.Label(80)
	if lbl.Service != "http" {
		t.Errorf("port 80 should still be http, got %s", lbl.Service)
	}
}

func TestLabel_String_Format(t *testing.T) {
	lbl := portlabel.Label{Port: 443, Service: "https", Severity: portlabel.SeverityLow}
	s := lbl.String()
	if s != "port=443 service=https severity=low" {
		t.Errorf("unexpected string: %s", s)
	}
}
