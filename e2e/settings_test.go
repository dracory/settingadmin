package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestSettingsPageLoads verifies the grouped settings editor renders
// the configured group card and its fields.
func TestSettingsPageLoads(t *testing.T) {
	ctx := newTab(t)

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("settings")),
		chromedp.WaitVisible("#group-general", chromedp.ByID),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("settings page failed: %v", err)
	}
	if !strings.Contains(body, "General") {
		t.Fatal("expected 'General' group card on settings page")
	}
	if !strings.Contains(body, "site.name") {
		t.Fatal("expected 'site.name' field on settings page")
	}
}

// TestManagerPageLoads verifies the raw settings manager mounts its Vue
// app and lists seeded settings.
func TestManagerPageLoads(t *testing.T) {
	ctx := newTab(t)

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("manager")),
		chromedp.Poll(`!!document.querySelector('#app .list-group')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("manager page failed: %v", err)
	}
	if !strings.Contains(body, "site.name") {
		t.Fatal("expected 'site.name' in manager list")
	}
}
