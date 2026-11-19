// Package settingadmin provides a standalone settings admin interface
// following the folder-per-controller pattern. Each controller is in its
// own subfolder and handles its own views and AJAX data.
//
// It supports two controllers:
//   - "settings" — a grouped editor driven by a declarative schema
//     (AdminOptions.Groups); each field type renders through an embedded
//     template file (templates/fields/field_<type>.html).
//   - "manager" — a raw key/value table for arbitrary settings.
//
// This module is modeled on github.com/dracory/logadmin and
// github.com/dracory/shopadmin.
package settingadmin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/settingadmin/manager"
	"github.com/dracory/settingadmin/settings"
	"github.com/dracory/settingadmin/shared"
	"github.com/dracory/settingstore"
)

// AdminOptions contains all dependencies and configuration for the
// setting admin.
type AdminOptions struct {
	// Store is the settingstore.StoreInterface (required)
	Store settingstore.StoreInterface

	// Logger is required
	Logger *slog.Logger

	// Groups is the declarative schema for the grouped settings editor.
	// Each group's fields are rendered via the FieldRegistry templates.
	// Optional — when empty, the "manager" controller becomes default.
	Groups []SettingGroup

	// FieldTemplates overrides or extends the embedded field templates.
	// Key = field type (e.g. "input", "mytype"), value = template text.
	FieldTemplates map[string]string

	// FuncLayout is an optional function to render the admin interface
	// inside your own layout (branding, menus, etc.). If nil, a default
	// bare-bones HTML page is used (Bootstrap + Vue CDN).
	FuncLayout func(w http.ResponseWriter, r *http.Request, title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string

	// AdminHomeURL is the URL for the admin home page (default: "/admin")
	AdminHomeURL string

	// SettingAdminURL is the base URL for setting admin
	// (default: "/admin/settings")
	SettingAdminURL string

	// FileManagerURL is the URL for the file manager (optional)
	FileManagerURL string

	// AuthUserID returns the authenticated user ID from the request.
	// If it returns "", the user is treated as unauthenticated.
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the setting admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	store      settingstore.StoreInterface
	logger     *slog.Logger
	groups     []SettingGroup
	fields     *FieldRegistry
	funcLayout func(w http.ResponseWriter, r *http.Request, title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	adminHomeURL    string
	settingAdminURL string
	fileManagerURL  string
	authUserID      func(r *http.Request) string
	routes          map[string]func(w http.ResponseWriter, r *http.Request)
}

// New creates a new setting admin instance.
// Returns ErrStoreRequired if Store is nil, ErrLoggerRequired if Logger
// is nil.
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Store == nil {
		return nil, ErrStoreRequired
	}
	if opts.Logger == nil {
		return nil, ErrLoggerRequired
	}

	fields, err := NewFieldRegistry(opts.FieldTemplates)
	if err != nil {
		return nil, err
	}

	if opts.AdminHomeURL == "" {
		opts.AdminHomeURL = "/admin"
	}
	if opts.SettingAdminURL == "" {
		opts.SettingAdminURL = "/admin/settings"
	}

	a := &admin{
		store:           opts.Store,
		logger:          opts.Logger,
		groups:          opts.Groups,
		fields:          fields,
		funcLayout:      opts.FuncLayout,
		adminHomeURL:    opts.AdminHomeURL,
		settingAdminURL: opts.SettingAdminURL,
		fileManagerURL:  opts.FileManagerURL,
		authUserID:      opts.AuthUserID,
	}

	// Build routes once at construction time
	a.routes = a.buildRoutes()

	return a, nil
}

// Handle processes all setting admin requests.
// Config values are injected into the request context (following the
// logadmin/shopadmin pattern). Route lookup is map-based via the
// ?controller= query parameter.
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if a.authUserID != nil && a.authUserID(r) == "" {
		http.Redirect(w, r, a.adminHomeURL, http.StatusSeeOther)
		return
	}

	// Inject config into request context (like logadmin/shopadmin)
	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)
	ctx = context.WithValue(ctx, shared.KeySettingAdminURL, a.settingAdminURL)
	ctx = context.WithValue(ctx, shared.KeyFileManagerURL, a.fileManagerURL)
	r = r.WithContext(ctx)

	// Map-based route lookup
	defaultController := shared.CONTROLLER_SETTINGS
	if len(a.groups) == 0 {
		defaultController = shared.CONTROLLER_MANAGER
	}

	controller := req.GetStringTrimmed(r, "controller")
	if controller == "" {
		controller = defaultController
	}

	handler, ok := a.routes[controller]
	if !ok {
		handler = a.routes[defaultController]
	}

	handler(w, r)
}

// buildRoutes creates the handler dispatch map once at construction time.
func (a *admin) buildRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	uiConfig := shared.UiConfig{
		Store:  a.store,
		Logger: a.logger,
		Groups: a.groups,
		Fields: a.fields,
		Layout: a.render,
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_SETTINGS: func(w http.ResponseWriter, r *http.Request) {
			settings.UI(uiConfig).Settings(w, r)
		},
		shared.CONTROLLER_MANAGER: func(w http.ResponseWriter, r *http.Request) {
			manager.UI(uiConfig).Manager(w, r)
		},
	}
}

// render wraps content in the layout. If FuncLayout is provided and
// returns non-empty, it is used; otherwise the default shared.Layout
// is used (following the logadmin/shopadmin pattern).
func (a *admin) render(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	// If a custom layout is provided, try it first
	if a.funcLayout != nil {
		custom := a.funcLayout(w, r, webpageTitle, webpageHtml, options)
		if custom != "" {
			if w != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(custom))
				return ""
			}
			return custom
		}
	}

	webpage := shared.Layout(w, r, webpageTitle, webpageHtml, options)

	if w != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(webpage))
		return ""
	}

	return webpage
}
