package settings

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/settingadmin/shared"
)

// handleGroupSave persists all schema fields of one group. The Vue app
// posts a JSON body of {group_id, values: {name: value}} collected from
// the rendered inputs. Checkbox fields missing from values are saved
// as "false". Raw and readonly fields are never written.
func (u *ui) handleGroupSave(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if u.Store() == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Setting store not available").ToString()))
		return ""
	}

	var reqBody struct {
		GroupID string            `json:"group_id"`
		Values  map[string]string `json:"values"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	group, found := u.findGroup(reqBody.GroupID)
	if !found {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Setting group not found: " + reqBody.GroupID).ToString()))
		return ""
	}

	for _, field := range group.Fields {
		if field.Type == shared.FIELD_TYPE_RAW || field.Readonly {
			continue
		}

		value, ok := reqBody.Values[field.Name]
		if !ok {
			if field.Type == shared.FIELD_TYPE_CHECKBOX {
				value = "false"
			} else {
				continue
			}
		}

		if err := u.Store().Set(ctx, field.Name, value); err != nil {
			u.Logger().Error("Failed to save setting", "error", err, "key", field.Name)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(api.Error("Failed to save setting: " + field.Name).ToString()))
			return ""
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Settings saved successfully").ToString()))
	return ""
}

// findGroup returns the group with the given ID.
func (u *ui) findGroup(groupID string) (shared.SettingGroup, bool) {
	for _, group := range u.Groups() {
		if group.ID == groupID {
			return group, true
		}
	}
	return shared.SettingGroup{}, false
}
