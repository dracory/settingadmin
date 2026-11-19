// Package e2e contains end-to-end tests that drive the settingadmin
// panel with a real headless Chrome via chromedp (CDP). The admin is
// served in-process by httptest on a random port — no binary build
// needed.
package e2e

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dracory/settingadmin"
	"github.com/dracory/settingadmin/testutils"
	_ "modernc.org/sqlite"
)

var (
	baseURL  string
	allocCtx context.Context
)

const adminPath = "/admin/settings"

func TestMain(m *testing.M) {
	// In-process admin server backed by an in-memory SQLite DB.
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		panic("e2e: failed to init store: " + err.Error())
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := store.Set(context.Background(), "site.name", "E2E Site"); err != nil {
		panic("e2e: failed to seed store: " + err.Error())
	}

	admin, err := settingadmin.New(settingadmin.AdminOptions{
		Store:           store,
		Logger:          logger,
		AdminHomeURL:    "/admin",
		SettingAdminURL: adminPath,
		Groups: []settingadmin.SettingGroup{
			{
				ID:    "general",
				Title: "General",
				Fields: []settingadmin.Field{
					{Name: "site.name", Label: "Site Name", Type: settingadmin.FIELD_TYPE_INPUT},
					{Name: "site.maintenance", Label: "Maintenance", Type: settingadmin.FIELD_TYPE_CHECKBOX},
				},
			},
		},
		// AuthUserID intentionally nil — the panel is open.
	})
	if err != nil {
		panic("e2e: failed to create settingadmin: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle(adminPath, http.HandlerFunc(admin.Handle))
	mux.Handle(adminPath+"/", http.HandlerFunc(admin.Handle))
	srv := httptest.NewServer(mux)
	baseURL = srv.URL

	// One shared browser process for the entire test run.
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Headless)
	if bin := chromeBin(); bin != "" {
		opts = append(opts, chromedp.ExecPath(bin))
	}
	var cancel context.CancelFunc
	allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)

	// Probe: if no usable Chrome binary exists in this environment,
	// skip the e2e suite instead of failing (CI provides Chrome).
	probe, probeCancel := chromedp.NewContext(allocCtx)
	probeCtx, probeTimeout := context.WithTimeout(probe, 10*time.Second)
	if err := chromedp.Run(probeCtx); err != nil {
		browserBroken = true
	}
	probeTimeout()
	probeCancel()

	code := m.Run()

	cancel()
	srv.Close()
	os.Exit(code)
}

// browserBroken is set in TestMain when Chrome cannot start; tests
// skip instead of failing on machines without a browser.
var browserBroken bool

// newTab opens a fresh tab in the shared browser and returns its
// context. Each tab gets a 30s deadline; cleanup closes the tab.
func newTab(t *testing.T) context.Context {
	t.Helper()
	if browserBroken {
		t.Skip("skipping e2e: no usable Chrome/Chromium binary in this environment")
	}
	ctx, cancel := chromedp.NewContext(allocCtx)
	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	t.Cleanup(func() { cancelTimeout(); cancel() })
	return ctx
}

// adminURL builds a full URL under the mounted admin path. The
// controller name is passed via ?controller=.
func adminURL(controller string) string {
	if controller == "" {
		return baseURL + adminPath
	}
	return baseURL + adminPath + "?controller=" + controller
}

// chromeBin resolves a headless Chromium binary. Order: CHROME_BIN env,
// Playwright's downloaded headless shell, empty (chromedp finds Chrome
// on PATH / standard install locations).
func chromeBin() string {
	if b := os.Getenv("CHROME_BIN"); b != "" {
		return b
	}
	patterns := []string{
		// Windows
		filepath.Join(os.Getenv("LOCALAPPDATA"),
			"ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe"),
		// Linux
		filepath.Join(os.Getenv("HOME"),
			".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"),
	}
	for _, pattern := range patterns {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return ""
}
