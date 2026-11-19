package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/settingstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
// This follows the shopadmin/logadmin pattern.
type UiInterface interface {
	Store() settingstore.StoreInterface
	Logger() *slog.Logger
	Groups() []SettingGroup
	Fields() FieldRenderer

	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
