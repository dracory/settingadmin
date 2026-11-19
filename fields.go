package settingadmin

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"github.com/dracory/settingadmin/shared"
)

//go:embed templates/fields/field_*.html
var fieldTemplatesFS embed.FS

// fieldView is the data passed to each field template.
type fieldView struct {
	Field     shared.Field
	Value     string
	ValueHTML template.HTML
}

// FieldRegistry renders setting fields using one template file per field
// type. Templates are embedded under templates/fields/field_<type>.html.
// Host applications can override existing types or register new ones via
// AdminOptions.FieldTemplates.
type FieldRegistry struct {
	templates map[string]*template.Template
}

// NewFieldRegistry parses the embedded field templates and merges any
// overrides (key = field type, value = template text).
func NewFieldRegistry(overrides map[string]string) (*FieldRegistry, error) {
	registry := &FieldRegistry{templates: map[string]*template.Template{}}

	entries, err := fs.ReadDir(fieldTemplatesFS, "templates/fields")
	if err != nil {
		return nil, fmt.Errorf("field templates: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "field_") || !strings.HasSuffix(name, ".html") {
			continue
		}
		fieldType := strings.TrimSuffix(strings.TrimPrefix(name, "field_"), ".html")
		content, err := fs.ReadFile(fieldTemplatesFS, "templates/fields/"+name)
		if err != nil {
			return nil, fmt.Errorf("field template %s: %w", name, err)
		}
		tpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("field template %s: %w", name, err)
		}
		registry.templates[fieldType] = tpl
	}

	for fieldType, content := range overrides {
		tpl, err := template.New("field_" + fieldType + ".html").Parse(content)
		if err != nil {
			return nil, fmt.Errorf("field template override %q: %w", fieldType, err)
		}
		registry.templates[fieldType] = tpl
	}

	return registry, nil
}

// HasType reports whether a template exists for the given field type.
func (r *FieldRegistry) HasType(fieldType string) bool {
	_, ok := r.templates[fieldType]
	return ok
}

// Types returns the registered field type names.
func (r *FieldRegistry) Types() []string {
	types := make([]string, 0, len(r.templates))
	for t := range r.templates {
		types = append(types, t)
	}
	return types
}

// Render renders a single field with the given current value.
// Unknown types fall back to the fallback template, or a plain text
// input template when no fallback is registered.
func (r *FieldRegistry) Render(field shared.Field, value string) (string, error) {
	fieldType := field.Type
	if fieldType == "" {
		fieldType = FIELD_TYPE_INPUT
	}

	tpl, ok := r.templates[fieldType]
	if !ok {
		tpl, ok = r.templates[FIELD_TYPE_FALLBACK]
	}
	if !ok {
		return "", fmt.Errorf("no template for field type %q and no fallback registered", field.Type)
	}

	view := fieldView{Field: field, Value: value}
	if fieldType == FIELD_TYPE_RAW {
		view.ValueHTML = template.HTML(value) // #nosec G203 -- raw type is explicitly unescaped by design
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, view); err != nil {
		return "", fmt.Errorf("render field %q: %w", field.Name, err)
	}
	return buf.String(), nil
}
