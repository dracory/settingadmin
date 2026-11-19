package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/settingstore"
)

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs can embed this to get the Store(),
// Logger(), Groups(), Fields(), and Layout() methods for free,
// following the shopadmin/logadmin pattern.
type UiBase struct {
	StoreField  settingstore.StoreInterface
	LoggerField *slog.Logger
	GroupsField []SettingGroup
	FieldsField FieldRenderer
	LayoutField func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

// Store returns the setting store
func (u UiBase) Store() settingstore.StoreInterface { return u.StoreField }

// Logger returns the logger
func (u UiBase) Logger() *slog.Logger { return u.LoggerField }

// Groups returns the declarative setting groups schema
func (u UiBase) Groups() []SettingGroup { return u.GroupsField }

// Fields returns the field renderer registry
func (u UiBase) Fields() FieldRenderer { return u.FieldsField }

// Layout renders the page layout
func (u UiBase) Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	return u.LayoutField(w, r, webpageTitle, webpageHtml, options)
}

// NewUiBase creates a UiBase from a UiConfig
func NewUiBase(config UiConfig) UiBase {
	return UiBase{
		StoreField:  config.Store,
		LoggerField: config.Logger,
		GroupsField: config.Groups,
		FieldsField: config.Fields,
		LayoutField: config.Layout,
	}
}
