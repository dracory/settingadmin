package shared

// FieldOption represents a single selectable option for select, radio,
// and checkbox fields.
type FieldOption struct {
	Value string
	Label string
}

// Field type constants. Each type maps to an embedded template file
// templates/fields/field_<type>.html. Unknown types fall back to
// field_fallback.html.
const (
	FIELD_TYPE_INPUT    = "input"
	FIELD_TYPE_TEXTAREA = "textarea"
	FIELD_TYPE_SELECT   = "select"
	FIELD_TYPE_CHECKBOX = "checkbox"
	FIELD_TYPE_NUMBER   = "number"
	FIELD_TYPE_DATE     = "date"
	FIELD_TYPE_DATETIME = "datetime"
	FIELD_TYPE_EMAIL    = "email"
	FIELD_TYPE_PASSWORD = "password"
	FIELD_TYPE_URL      = "url"
	FIELD_TYPE_TEL      = "tel"
	FIELD_TYPE_COLOR    = "color"
	FIELD_TYPE_RADIO    = "radio"
	FIELD_TYPE_HIDDEN   = "hidden"
	FIELD_TYPE_HTMLAREA = "htmlarea"
	FIELD_TYPE_FILE     = "file"
	FIELD_TYPE_IMAGE    = "image"
	FIELD_TYPE_RAW      = "raw"
	FIELD_TYPE_FALLBACK = "fallback"
)

// Field describes a single editable setting. Field.Name is the key used
// to persist the value in the settingstore. Field.DefaultValue is used
// when the store has no value for the key.
type Field struct {
	Name         string
	Label        string
	Type         string
	DefaultValue string
	Placeholder  string
	Help         string
	Options      []FieldOption
	Readonly     bool
	Required     bool
}

// SettingGroup groups related fields under a titled card in the admin UI.
type SettingGroup struct {
	ID          string
	Title       string
	Description string
	Fields      []Field
}

// FieldRenderer renders a single field to HTML (implemented by
// settingadmin.FieldRegistry).
type FieldRenderer interface {
	Render(field Field, value string) (string, error)
}
