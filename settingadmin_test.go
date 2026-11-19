package settingadmin_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/settingadmin"
	"github.com/dracory/settingadmin/testutils"
	_ "modernc.org/sqlite"
)

func logger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func TestNewRequiresStore(t *testing.T) {
	_, err := settingadmin.New(settingadmin.AdminOptions{Logger: logger()})
	if err != settingadmin.ErrStoreRequired {
		t.Fatalf("expected ErrStoreRequired, got %v", err)
	}
}

func TestNewRequiresLogger(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}
	_, err = settingadmin.New(settingadmin.AdminOptions{Store: store})
	if err != settingadmin.ErrLoggerRequired {
		t.Fatalf("expected ErrLoggerRequired, got %v", err)
	}
}

func TestHandleRedirectsUnauthenticated(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	admin, err := settingadmin.New(settingadmin.AdminOptions{
		Store:      store,
		Logger:     logger(),
		AuthUserID: func(r *http.Request) string { return "" },
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	admin.Handle(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin" {
		t.Fatalf("expected redirect to /admin, got %s", loc)
	}
}

func TestHandleRendersSettingsPage(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	admin, err := settingadmin.New(settingadmin.AdminOptions{
		Store:  store,
		Logger: logger(),
		Groups: []settingadmin.SettingGroup{
			{ID: "g1", Title: "Group One", Fields: []settingadmin.Field{
				{Name: "k1", Label: "Key One", Type: settingadmin.FIELD_TYPE_INPUT},
			}},
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	admin.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `id="app"`) || !strings.Contains(body, "load-settings") {
		t.Fatal("expected Vue app shell with load-settings endpoint in settings page")
	}
}

func TestHandleManagerController(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	admin, err := settingadmin.New(settingadmin.AdminOptions{
		Store:  store,
		Logger: logger(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings?controller=manager", nil)
	admin.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Settings. Manager") {
		t.Fatal("expected manager page header")
	}
}
