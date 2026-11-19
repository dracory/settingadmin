package shared

import "net/http"

// Links provides URL helpers for settingadmin controllers.
// The base URL is read from request context (injected by Handle()),
// not hardcoded. This follows the shopadmin/logadmin pattern.
type Links struct {
	baseURL string
}

// NewLinks creates a Links helper with the given base URL.
// If baseURL is empty, defaults to "/admin/settings".
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/settings"
	}
	return &Links{baseURL: baseURL}
}

// NewLinksFromRequest creates a Links helper using the setting admin URL
// from the request context.
func NewLinksFromRequest(r *http.Request) *Links {
	return NewLinks(SettingAdminURL(r))
}

// Settings builds the URL for the grouped settings controller
func (l *Links) Settings(params map[string]string) string {
	return l.url(CONTROLLER_SETTINGS, params)
}

// Manager builds the URL for the raw key/value manager controller
func (l *Links) Manager(params map[string]string) string {
	return l.url(CONTROLLER_MANAGER, params)
}

// url builds a URL for the given controller. The params map is copied
// before mutation (does not modify caller's map).
func (l *Links) url(controller string, params map[string]string) string {
	return URL(l.baseURL, controller, params)
}
