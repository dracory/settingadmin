package manager

import (
	"strings"

	"github.com/dracory/settingstore"
)

const (
	FieldID        = "id"
	FieldKey       = "key"
	FieldValue     = "value"
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldSettings  = "settings"
	FieldTotal     = "total"
	FieldSettingID = "setting_id"
)

// Condition field constants for the multi-filter system.
const (
	CondFieldKey      = "key"
	CondFieldID       = "id"
	CondFieldDateFrom = "date_from"
	CondFieldDateTo   = "date_to"
)

// condition is a single filter condition sent from the frontend.
type condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// applyConditions maps filter conditions to the setting query.
// Note: settingstore's query has no KeyLike filter (unlike shopstore),
// so key conditions use exact match via SetKey.
func applyConditions(query settingstore.SettingQueryInterface, conds []condition) {
	for _, c := range conds {
		val := strings.TrimSpace(c.Value)
		if val == "" {
			continue
		}
		switch c.Field {
		case CondFieldKey:
			query.SetKey(val)
		case CondFieldID:
			query.SetID(val)
		case CondFieldDateFrom:
			query.SetCreatedAtGte(val + " 00:00:00")
		case CondFieldDateTo:
			query.SetCreatedAtLte(val + " 23:59:59")
		}
	}
}
