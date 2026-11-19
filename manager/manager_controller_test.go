package manager_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/settingstore"
)

func TestRenderPageShowsManager(t *testing.T) {
	env := newTestEnv(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings?controller=manager", nil)
	env.ui.Manager(rec, req)

	if !strings.Contains(rec.Body.String(), "Settings. Manager") {
		t.Error("expected manager page header")
	}
}

func TestLoadSettingsReturnsSeeded(t *testing.T) {
	env := newTestEnv(t)
	if err := env.store.Set(context.Background(), "alpha", "1"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resp := postAction(t, env, "load-settings", nil, map[string]any{
		"page": 0, "per_page": 10,
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}
	settings, ok := resp.Data["settings"].([]any)
	if !ok || len(settings) != 1 {
		t.Fatalf("expected 1 setting, got %v", resp.Data["settings"])
	}
}

func TestSaveSettingCreates(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-setting", nil, map[string]any{
		"key": "beta", "value": "2",
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	got, _ := env.store.Get(context.Background(), "beta", "")
	if got != "2" {
		t.Fatalf("expected beta=2, got %q", got)
	}
}

func TestSaveSettingRejectsEmptyKey(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-setting", nil, map[string]any{
		"key": "  ", "value": "2",
	})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %v", resp)
	}
}

func TestDeleteSettingRemoves(t *testing.T) {
	env := newTestEnv(t)
	if err := env.store.Set(context.Background(), "gone", "x"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	setting, err := env.store.SettingFindByKey(context.Background(), "gone")
	if err != nil {
		t.Fatalf("find: %v", err)
	}

	resp := postAction(t, env, "delete-setting", nil, map[string]any{
		"setting_id": setting.GetID(),
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	found, _ := env.store.SettingFindByKey(context.Background(), "gone")
	if found != nil {
		t.Fatal("expected setting deleted")
	}
}

func TestDeleteSelectedRemovesAll(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	for _, key := range []string{"d1", "d2"} {
		if err := env.store.Set(ctx, key, "x"); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	list, err := env.store.SettingList(ctx, settingstore.NewSettingQuery())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	ids := []string{}
	for _, s := range list {
		ids = append(ids, s.GetID())
	}

	resp := postAction(t, env, "delete-selected", nil, map[string]any{
		"bulk_setting_ids": ids,
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %v", resp)
	}

	count, _ := env.store.SettingCount(ctx, settingstore.NewSettingQuery())
	if count != 0 {
		t.Fatalf("expected 0 settings, got %d", count)
	}
}

func TestLoadSettingsKeyFilter(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	_ = env.store.Set(ctx, "wanted", "1")
	_ = env.store.Set(ctx, "other", "2")

	resp := postAction(t, env, "load-settings", nil, map[string]any{
		"page": 0, "per_page": 10, "key": "wanted",
	})
	settings, ok := resp.Data["settings"].([]any)
	if !ok || len(settings) != 1 {
		t.Fatalf("expected 1 filtered setting, got %v", resp.Data["settings"])
	}
	first := settings[0].(map[string]any)
	if first["key"] != "wanted" {
		t.Fatalf("expected key=wanted, got %v", first["key"])
	}
}
