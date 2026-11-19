package manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/neat"
	"github.com/dracory/settingstore"
)

func (u *ui) handleLoadSettings(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	settingStore := u.Store()
	if settingStore == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Setting store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page       int         `json:"page"`
		PerPage    int         `json:"per_page"`
		SortBy     string      `json:"sort_by"`
		Sort       string      `json:"sort"`
		Key        string      `json:"key"`
		Conditions []condition `json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.Page < 0 {
		reqBody.Page = 0
	}
	if reqBody.PerPage <= 0 {
		reqBody.PerPage = 10
	}
	if reqBody.SortBy == "" {
		reqBody.SortBy = settingstore.COLUMN_SETTING_KEY
	}
	if reqBody.Sort == "" {
		reqBody.Sort = neat.SortAsc
	}

	query := settingstore.NewSettingQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortOrder(reqBody.Sort)

	if reqBody.Key != "" {
		query.SetKey(reqBody.Key)
	}

	applyConditions(query, reqBody.Conditions)

	settings, err := settingStore.SettingList(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to load settings", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to load settings").ToString()))
		return ""
	}

	total, err := settingStore.SettingCount(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to count settings", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to count settings").ToString()))
		return ""
	}

	settingList := []map[string]any{}
	for _, setting := range settings {
		settingList = append(settingList, map[string]any{
			FieldID:        setting.GetID(),
			FieldKey:       setting.GetKey(),
			FieldValue:     setting.GetValue(),
			FieldCreatedAt: setting.GetCreatedAt(),
			FieldUpdatedAt: setting.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Settings loaded successfully", map[string]any{
		FieldSettings: settingList,
		FieldTotal:    total,
		"page":        reqBody.Page,
		"per_page":    reqBody.PerPage,
	}).ToString()))
	return ""
}
