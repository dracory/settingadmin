# settingadmin

Standalone admin panel for [`dracory/settingstore`](https://github.com/dracory/settingstore),
following the same architecture as `shopadmin`/`logadmin`/`blogadmin`
(folder-per-controller, embeddable, layout-agnostic).

## Features

- **Grouped settings editor** — define groups and fields in code; each
  group renders as a card with its own Save button. Values persist in
  the settingstore keyed by field name.
- **Field templates** — every field type renders through its own
  embedded template (`templates/fields/field_<type>.html`), so adding
  or restyling a type means touching one file.
- **Raw key/value manager** — paginated, filterable table for arbitrary
  settings with save/delete (Vue 3 frontend, `?controller=manager`).
- **Template overrides** — `AdminOptions.FieldTemplates` lets the host
  app override existing types or register new ones.

## Quick start

```bash
go get github.com/dracory/settingadmin
```

```go
store, _ := settingstore.NewStore(settingstore.NewStoreOptions{
    DB: db, SettingTableName: "setting", AutomigrateEnabled: true,
})

admin, _ := settingadmin.New(settingadmin.AdminOptions{
    Store:  store,
    Logger: logger,
    Groups: []settingadmin.SettingGroup{
        {
            ID:    "general",
            Title: "General",
            Fields: []settingadmin.Field{
                {Name: "site.name", Label: "Site Name", Type: settingadmin.FIELD_TYPE_INPUT},
                {Name: "site.maintenance", Label: "Maintenance", Type: settingadmin.FIELD_TYPE_CHECKBOX},
                {Name: "shop.currency", Label: "Currency", Type: settingadmin.FIELD_TYPE_SELECT,
                    Options: []settingadmin.FieldOption{{Value: "USD", Label: "US Dollar"}}},
            },
        },
    },
})

mux.Handle("/admin/settings", http.HandlerFunc(admin.Handle))
mux.Handle("/admin/settings/", http.HandlerFunc(admin.Handle))
```

- `/admin/settings` → grouped editor (`?controller=settings`)
- `/admin/settings?controller=manager` → raw key/value manager

## Field types

`input`, `textarea`, `select`, `checkbox`, `number`, `date`, `datetime`,
`email`, `password`, `url`, `tel`, `color`, `radio`, `hidden`,
`htmlarea`, `file`, `image`, `raw` — each backed by
`templates/fields/field_<type>.html`. Unknown types render via
`field_fallback.html`.

To add/override a type:

```go
settingadmin.New(settingadmin.AdminOptions{
    ...
    FieldTemplates: map[string]string{
        "markdown": `<textarea name="{{.Field.Name}}">{{.Value}}</textarea>`,
    },
})
```

## Example

```bash
go run ./example   # serves http://localhost:8080, admin at /admin/settings
```

## Testing

```bash
go test ./...      # unit tests (in-memory SQLite)
go test ./e2e      # chromedp e2e (skips when no Chrome available)
```

## License

AGPL-3.0. Commercial license available — see LICENSE_COMMERCIAL.txt.
