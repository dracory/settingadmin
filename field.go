package settingadmin

import "github.com/dracory/settingadmin/shared"

// FieldOption represents a single selectable option for select, radio,
// and checkbox fields. Alias of shared.FieldOption.
type FieldOption = shared.FieldOption

// Field type constants. Each type maps to an embedded template file
// templates/fields/field_<type>.html. Unknown types fall back to
// field_fallback.html.
const (
	FIELD_TYPE_INPUT    = shared.FIELD_TYPE_INPUT
	FIELD_TYPE_TEXTAREA = shared.FIELD_TYPE_TEXTAREA
	FIELD_TYPE_SELECT   = shared.FIELD_TYPE_SELECT
	FIELD_TYPE_CHECKBOX = shared.FIELD_TYPE_CHECKBOX
	FIELD_TYPE_NUMBER   = shared.FIELD_TYPE_NUMBER
	FIELD_TYPE_DATE     = shared.FIELD_TYPE_DATE
	FIELD_TYPE_DATETIME = shared.FIELD_TYPE_DATETIME
	FIELD_TYPE_EMAIL    = shared.FIELD_TYPE_EMAIL
	FIELD_TYPE_PASSWORD = shared.FIELD_TYPE_PASSWORD
	FIELD_TYPE_URL      = shared.FIELD_TYPE_URL
	FIELD_TYPE_TEL      = shared.FIELD_TYPE_TEL
	FIELD_TYPE_COLOR    = shared.FIELD_TYPE_COLOR
	FIELD_TYPE_RADIO    = shared.FIELD_TYPE_RADIO
	FIELD_TYPE_HIDDEN   = shared.FIELD_TYPE_HIDDEN
	FIELD_TYPE_HTMLAREA = shared.FIELD_TYPE_HTMLAREA
	FIELD_TYPE_FILE     = shared.FIELD_TYPE_FILE
	FIELD_TYPE_IMAGE    = shared.FIELD_TYPE_IMAGE
	FIELD_TYPE_RAW      = shared.FIELD_TYPE_RAW
	FIELD_TYPE_FALLBACK = shared.FIELD_TYPE_FALLBACK
)

// Field describes a single editable setting. Field.Name is the key used
// to persist the value in the settingstore. Field.DefaultValue is used
// when the store has no value for the key.
// Alias of shared.Field.
type Field = shared.Field

// SettingGroup groups related fields under a titled card in the admin UI.
// Alias of shared.SettingGroup.
type SettingGroup = shared.SettingGroup
