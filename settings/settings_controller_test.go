package settings_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderPageShowsAppShell(t *testing.T) {
	env := newTestEnv(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	env.ui.Settings(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, `id="app"`) {
		t.Error("expected Vue app container")
	}
	if !strings.Contains(body, "settingAdminUrls") {
		t.Error("expected injected URLs")
	}
	if !strings.Contains(body, "load-settings") {
		t.Error("expected load-settings action URL")
	}
}

func TestLoadSettingsReturnsGroupsWithRenderedFields(t *testing.T) {
	env := newTestEnv(t)
	if err := env.store.Set(context.Background(), "site.name", "StoredValue"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resp := postAction(t, env, "load-settings", map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	groups, ok := resp.Data["groups"].([]any)
	if !ok || len(groups) != 1 {
		t.Fatalf("expected 1 group, got %v", resp.Data["groups"])
	}
	group := groups[0].(map[string]any)
	if group["id"] != "general" || group["title"] != "General" {
		t.Fatalf("unexpected group payload: %v", group)
	}

	fields := group["fields"].([]any)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %v", fields)
	}
	first := fields[0].(map[string]any)
	if first["name"] != "site.name" {
		t.Fatalf("expected field name site.name, got %v", first["name"])
	}
	if !strings.Contains(first["html"].(string), `value="StoredValue"`) {
		t.Error("expected stored value rendered in field html")
	}
}

func TestSaveGroupPersistsFields(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-group", map[string]any{
		"group_id": "general",
		"values": map[string]string{
			"site.name":        "NewName",
			"site.maintenance": "true",
		},
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	got, err := env.store.Get(context.Background(), "site.name", "")
	if err != nil || got != "NewName" {
		t.Fatalf("expected site.name=NewName, got %q err=%v", got, err)
	}
}

func TestSaveGroupCheckboxDefaultsFalse(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-group", map[string]any{
		"group_id": "general",
		"values":   map[string]string{"site.name": "X"},
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	got, _ := env.store.Get(context.Background(), "site.maintenance", "")
	if got != "false" {
		t.Fatalf("expected unchecked checkbox saved as false, got %q", got)
	}
}

func TestSaveGroupUnknownGroup(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-group", map[string]any{
		"group_id": "nope",
		"values":   map[string]string{},
	})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %v", resp)
	}
	if !strings.Contains(resp.Message, "Setting group not found") {
		t.Errorf("expected group-not-found message, got %q", resp.Message)
	}
}
