package shared

// Context keys for config values injected by Handle()
const KeyEndpoint = "endpoint"
const KeyAdminHomeURL = "admin_home_url"
const KeySettingAdminURL = "setting_admin_url"
const KeyFileManagerURL = "file_manager_url"

// Controller names used in the ?controller= query parameter
const (
	CONTROLLER_SETTINGS = "settings"
	CONTROLLER_MANAGER  = "manager"
)

// CatchAll is the catch-all route suffix
const CatchAll = "/*"

// Error messages
const ERROR_STORE_IS_NIL = "store cannot be nil"
const ERROR_LOGGER_IS_NIL = "logger cannot be nil"
