package shared

import "net/http"

// SettingAdminURL returns the setting admin base URL from request context
func SettingAdminURL(r *http.Request) string {
	value, ok := r.Context().Value(KeySettingAdminURL).(string)
	if !ok {
		return ""
	}
	return value
}
