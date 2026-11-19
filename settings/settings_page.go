package settings

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
	//go:embed settings.html
	settingsHTML string

	//go:embed settings.js
	settingsJS string
)

// renderPage renders the grouped settings editor shell. The Vue app
// fetches the schema + values via the load-settings action and renders
// each group's fields (server-rendered HTML injected via v-html).
func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Setting store is not initialized")
	}
	if u.Fields() == nil {
		return shared.ErrorAlert("Field renderer is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Settings", URL: shared.URLR(r, shared.CONTROLLER_SETTINGS, nil)},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadSettings": linksHelper.Settings(map[string]string{"action": actionLoadSettings}),
		"saveGroup":    linksHelper.Settings(map[string]string{"action": actionSaveGroup}),
	}
	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.settingAdminUrls = %s;", string(urlsJSON)))

	managerLink := hb.Hyperlink().
		Class("btn btn-outline-secondary btn-sm").
		Href(linksHelper.Manager(nil)).
		HTML("Manage Raw Settings")

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Settings", managerLink)).
		Child(urlsScript).
		Child(hb.Raw(settingsHTML)).
		Child(hb.Script(settingsJS))

	return u.Layout(w, r, "Settings", content.ToHTML(), struct {
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
