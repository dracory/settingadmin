package settings

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/settingadmin/shared"
)

const (
	actionLoadSettings = "load-settings"
	actionSaveGroup    = "save-group"
)

// UiInterface defines the grouped settings controller's UI interface
type UiInterface interface {
	shared.UiInterface
	Settings(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new grouped settings controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// Settings handles the grouped settings controller requests
func (u *ui) Settings(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the settings request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadSettings:
		return u.handleLoadSettings(w, r)
	case actionSaveGroup:
		return u.handleGroupSave(w, r)
	default:
		return u.renderPage(w, r)
	}
}
