package manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/settingadmin/shared"
)

var (
	//go:embed manager.html
	managerHTML string

	//go:embed manager.js
	managerJS string
)

// renderPage renders the raw key/value setting manager page.
func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Setting store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Settings", URL: shared.URLR(r, shared.CONTROLLER_SETTINGS, nil)},
		{Name: "Manager", URL: shared.URLR(r, shared.CONTROLLER_MANAGER, nil)},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadSettings":           linksHelper.Manager(map[string]string{"action": actionLoadSettings}),
		"saveSetting":            linksHelper.Manager(map[string]string{"action": actionSettingSave}),
		"deleteSetting":          linksHelper.Manager(map[string]string{"action": actionSettingDelete}),
		"deleteSelectedSettings": linksHelper.Manager(map[string]string{"action": actionSettingDeleteSelected}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.settingManagerUrls = %s;", string(urlsJSON)))

	settingsLink := hb.Hyperlink().
		Class("btn btn-outline-secondary btn-sm").
		Href(linksHelper.Settings(nil)).
		HTML("Back to Settings")

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Settings. Manager", settingsLink)).
		Child(urlsScript).
		Child(hb.Raw(managerHTML)).
		Child(hb.Script(managerJS))

	return u.Layout(w, r, "Settings Manager", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Sweetalert2_10(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
