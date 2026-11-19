package settingadmin_test

import (
	"strings"
	"testing"

	"github.com/dracory/settingadmin"
	"github.com/dracory/settingadmin/shared"
)

func newRegistry(t *testing.T) *settingadmin.FieldRegistry {
	t.Helper()
	r, err := settingadmin.NewFieldRegistry(nil)
	if err != nil {
		t.Fatalf("NewFieldRegistry: %v", err)
	}
	return r
}

func TestFieldRegistryRendersEmbeddedTypes(t *testing.T) {
	r := newRegistry(t)

	for _, fieldType := range []string{
		"input", "textarea", "select", "checkbox", "number", "date",
		"datetime", "email", "password", "url", "tel", "color",
		"radio", "hidden", "htmlarea", "file", "image", "raw",
	} {
		if !r.HasType(fieldType) {
			t.Errorf("missing embedded template for type %q", fieldType)
		}
	}
}

func TestFieldRegistryRenderInput(t *testing.T) {
	r := newRegistry(t)

	html, err := r.Render(shared.Field{
		Name:  "site.name",
		Label: "Site Name",
		Type:  settingadmin.FIELD_TYPE_INPUT,
	}, "hello")
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(html, `name="site.name"`) {
		t.Error("expected input with name site.name")
	}
	if !strings.Contains(html, `value="hello"`) {
		t.Error("expected value attribute")
	}
}

func TestFieldRegistryRenderSelect(t *testing.T) {
	r := newRegistry(t)

	html, err := r.Render(shared.Field{
		Name: "currency",
		Type: settingadmin.FIELD_TYPE_SELECT,
		Options: []shared.FieldOption{
			{Value: "USD", Label: "US Dollar"},
			{Value: "EUR", Label: "Euro"},
		},
	}, "EUR")
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(html, `value="EUR" selected`) {
		t.Error("expected EUR option selected")
	}
}

func TestFieldRegistryUnknownTypeFallsBack(t *testing.T) {
	r := newRegistry(t)

	html, err := r.Render(shared.Field{
		Name: "weird",
		Type: "nonexistent",
	}, "v")
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(html, "unknown type") {
		t.Error("expected fallback template to mark unknown type")
	}
}

func TestFieldRegistryOverride(t *testing.T) {
	r, err := settingadmin.NewFieldRegistry(map[string]string{
		"input": `<div class="custom">{{.Value}}</div>`,
		"mine":  `<b>{{.Field.Name}}</b>`,
	})
	if err != nil {
		t.Fatalf("NewFieldRegistry: %v", err)
	}

	html, err := r.Render(shared.Field{Name: "x", Type: "input"}, "ov")
	if err != nil || !strings.Contains(html, `class="custom"`) {
		t.Errorf("override not applied: %v %q", err, html)
	}

	html, err = r.Render(shared.Field{Name: "y", Type: "mine"}, "")
	if err != nil || !strings.Contains(html, "<b>y</b>") {
		t.Errorf("custom type not registered: %v %q", err, html)
	}
}
