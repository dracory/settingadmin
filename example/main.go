// Example settingadmin server.
//
// Boots an HTTP server with an in-memory SQLite database and mounts the
// settingadmin panel at /admin/settings. No external services required.
//
// Run from the settingadmin module root:
//
//	go run ./example
//
// Then open http://localhost:8080/ in your browser.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dracory/settingadmin"
	"github.com/dracory/settingstore"
	_ "modernc.org/sqlite"
)

const (
	addr     = ":8080"
	dbFile   = ":memory:"
	adminURL = "/admin/settings"
	homeURL  = "/admin"
	dbDriver = "sqlite"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := sql.Open(dbDriver, dbFile+"?parseTime=true")
	if err != nil {
		logger.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	store, err := settingstore.NewStore(settingstore.NewStoreOptions{
		DB:                 db,
		SettingTableName:   "setting",
		AutomigrateEnabled: true,
	})
	if err != nil {
		logger.Error("failed to create settingstore", "err", err)
		os.Exit(1)
	}

	seedSettings(store, logger)

	admin, err := settingadmin.New(settingadmin.AdminOptions{
		Store:           store,
		Logger:          logger,
		SettingAdminURL: adminURL,
		AdminHomeURL:    homeURL,
		Groups: []settingadmin.SettingGroup{
			{
				ID:          "general",
				Title:       "General",
				Description: "Site-wide settings",
				Fields: []settingadmin.Field{
					{Name: "site.name", Label: "Site Name", Type: settingadmin.FIELD_TYPE_INPUT, DefaultValue: "My Site"},
					{Name: "site.description", Label: "Site Description", Type: settingadmin.FIELD_TYPE_TEXTAREA, Help: "Shown in search results"},
					{Name: "site.maintenance", Label: "Maintenance Mode", Type: settingadmin.FIELD_TYPE_CHECKBOX, Placeholder: "Take the site offline"},
				},
			},
			{
				ID:          "shop",
				Title:       "Shop",
				Description: "Storefront settings",
				Fields: []settingadmin.Field{
					{Name: "shop.currency", Label: "Currency", Type: settingadmin.FIELD_TYPE_SELECT, DefaultValue: "USD", Options: []settingadmin.FieldOption{
						{Value: "USD", Label: "US Dollar"},
						{Value: "EUR", Label: "Euro"},
						{Value: "GBP", Label: "British Pound"},
					}},
					{Name: "shop.per_page", Label: "Products Per Page", Type: settingadmin.FIELD_TYPE_NUMBER, DefaultValue: "20"},
					{Name: "shop.support_email", Label: "Support Email", Type: settingadmin.FIELD_TYPE_EMAIL},
				},
			},
		},
	})
	if err != nil {
		logger.Error("failed to create settingadmin", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// settingadmin.AdminInterface exposes Handle(w, r) but not ServeHTTP,
	// so wrap it in an http.HandlerFunc.
	mux.Handle(adminURL, http.HandlerFunc(admin.Handle))
	mux.Handle(adminURL+"/", http.HandlerFunc(admin.Handle))

	// Landing page that links into the admin.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<html><body style="font-family:sans-serif;padding:40px">
			<h1>settingadmin example</h1>
			<ul>
				<li><a href="%s">Settings (grouped editor)</a></li>
				<li><a href="%s?controller=manager">Settings Manager (raw editor)</a></li>
			</ul>
		</body></html>`, adminURL, adminURL)
	})

	logger.Info("serving", "addr", "http://localhost"+addr, "admin", adminURL)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}

// seedSettings stores a few values so the editors have data to show.
func seedSettings(store settingstore.StoreInterface, logger *slog.Logger) {
	ctx := context.Background()
	seeds := map[string]string{
		"site.name":       "Example Site",
		"shop.currency":   "EUR",
		"shop.per_page":   "24",
		"feature.enabled": "true",
	}
	for key, value := range seeds {
		if err := store.Set(ctx, key, value); err != nil {
			logger.Error("seed failed", "key", key, "err", err)
		}
	}
}
