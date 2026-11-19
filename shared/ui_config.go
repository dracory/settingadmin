package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/settingstore"
)

// UiConfig holds the dependencies passed to subcontroller UI factories.
// The Layout function uses an anonymous struct for options to match
// shopadmin/logadmin exactly, allowing consumers to reuse their
// existing layout function for settingadmin.
type UiConfig struct {
	Store  settingstore.StoreInterface
	Logger *slog.Logger
	Groups []SettingGroup
	Fields FieldRenderer
	Layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
