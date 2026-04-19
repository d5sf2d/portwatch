package notify_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func baseMsg() notify.Message {
	return notify.Message{
		Level: notify.LevelWarn,
		Title: "port opened",
		Body:  "port 8080 is now open",
		Tags:  []string{"tcp"},
	}
}

func TestLogNotifier_Send(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)
	if err := n.Send(baseMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "port opened") {
		t.Errorf("expected title in output, got: %s", out)
	}
	if !strings.Contains(out, string(notify.LevelWarn)) {
		t.Errorf("expected level in output, got: %s", out)
	}
}

func TestLogNotifier_NilWriterDefaultsToStdout(t *testing.T) {
	n := notify.NewLogNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestWebhookNotifier_Send(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL)
	if err := n.Send(baseMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["title"] != "port opened" {
		t.Errorf("unexpected payload: %v", received)
	}
}

func TestWebhookNotifier_BadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL)
	if err := n.Send(baseMsg()); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestMulti_Send(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	m := notify.NewMulti(notify.NewLogNotifier(&buf1), notify.NewLogNotifier(&buf2))
	if err := m.Send(baseMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both notifiers to receive message")
	}
}
