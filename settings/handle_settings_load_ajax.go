package settings

import (
	"net/http"

	"github.com/dracory/api"
)

// settingFieldPayload is one field in the load response: metadata plus
// the server-rendered HTML (from the per-type field template).
type settingFieldPayload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	HTML string `json:"html"`
}

// settingGroupPayload is one group in the load response.
type settingGroupPayload struct {
	ID          string                `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Fields      []settingFieldPayload `json:"fields"`
}

// handleLoadSettings returns the group schema with each field's current
// value rendered into its template HTML. The Vue app injects the HTML
// via v-html, so field types stay server-side templates.
func (u *ui) handleLoadSettings(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if u.Store() == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Setting store not available").ToString()))
		return ""
	}
	if u.Fields() == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Field renderer not available").ToString()))
		return ""
	}

	groups := []settingGroupPayload{}
	for _, group := range u.Groups() {
		payload := settingGroupPayload{
			ID:          group.ID,
			Title:       group.Title,
			Description: group.Description,
			Fields:      []settingFieldPayload{},
		}

		for _, field := range group.Fields {
			value, err := u.Store().Get(ctx, field.Name, field.DefaultValue)
			if err != nil {
				u.Logger().Error("Failed to load setting", "error", err, "key", field.Name)
				value = field.DefaultValue
			}

			fieldHTML, err := u.Fields().Render(field, value)
			if err != nil {
				u.Logger().Error("Failed to render field", "error", err, "field", field.Name)
				fieldHTML = ""
			}

			payload.Fields = append(payload.Fields, settingFieldPayload{
				Name: field.Name,
				Type: field.Type,
				HTML: fieldHTML,
			})
		}

		groups = append(groups, payload)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Settings loaded successfully", map[string]any{
		"groups": groups,
	}).ToString()))
	return ""
}
