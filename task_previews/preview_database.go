package task_previews

import (
	"fmt"
	"strings"
)

// A backend preview's own throwaway database. The runner runs one from an official
// image beside a Task's API preview when the API needs a database; only the
// engines and versions below are ever run. app-server validates a person's
// PreviewDatabaseV1 setting against the same list.

// Preview database engines.
const (
	PreviewDatabaseEnginePostgres = "postgres"
	PreviewDatabaseEngineMySQL    = "mysql"
	PreviewDatabaseEngineMongoDB  = "mongodb"
)

// PreviewDatabaseV1.Mode values.
const (
	// PreviewDatabaseModeAutomatic detects the engine and variables from the
	// Task's code (the same as a nil setting).
	PreviewDatabaseModeAutomatic = "automatic"
	// PreviewDatabaseModeNone gives the preview no database of its own.
	PreviewDatabaseModeNone = "none"
	// PreviewDatabaseModeUse runs Engine at Version and points Variables at it,
	// with no detection.
	PreviewDatabaseModeUse = "use"
)

// PreviewDatabaseV1 is a person's setting for a service's preview database.
type PreviewDatabaseV1 struct {
	Mode string
	// Engine, Version and Variables are read only when Mode is "use". Version ""
	// means the engine's default. Variables are the names of the variables that
	// hold the database's address (one URL variable, or split variables).
	Engine    string
	Version   string
	Variables []string
}

// PreviewDatabaseEngineV1 is one allowlisted engine.
type PreviewDatabaseEngineV1 struct {
	Engine string
	// Versions are the image tags that may run; DefaultVersion is one of them.
	Versions       []string
	DefaultVersion string
	Port           int64
	// ImageRepository is the official image's repository; Image adds the tag.
	ImageRepository string
	Cpu             string
	Memory          string
}

var previewDatabaseEngines = []PreviewDatabaseEngineV1{
	{
		Engine:          PreviewDatabaseEnginePostgres,
		Versions:        []string{"13", "14", "15", "16", "17"},
		DefaultVersion:  "16",
		Port:            5432,
		ImageRepository: "public.ecr.aws/docker/library/postgres",
		Cpu:             "256",
		Memory:          "512",
	},
	{
		Engine:          PreviewDatabaseEngineMySQL,
		Versions:        []string{"8.0", "8.4"},
		DefaultVersion:  "8.4",
		Port:            3306,
		ImageRepository: "public.ecr.aws/docker/library/mysql",
		Cpu:             "256",
		Memory:          "1024",
	},
	{
		Engine:          PreviewDatabaseEngineMongoDB,
		Versions:        []string{"6.0", "7.0"},
		DefaultVersion:  "7.0",
		Port:            27017,
		ImageRepository: "public.ecr.aws/docker/library/mongo",
		Cpu:             "256",
		Memory:          "1024",
	},
}

// PreviewDatabaseEngines returns the allowlist (a copy).
func PreviewDatabaseEngines() []PreviewDatabaseEngineV1 {
	out := make([]PreviewDatabaseEngineV1, len(previewDatabaseEngines))
	for i, e := range previewDatabaseEngines {
		e.Versions = append([]string(nil), e.Versions...)
		out[i] = e
	}
	return out
}

// LookupPreviewDatabaseEngine returns the allowlisted engine, or false.
func LookupPreviewDatabaseEngine(engine string) (PreviewDatabaseEngineV1, bool) {
	for _, e := range PreviewDatabaseEngines() {
		if e.Engine == engine {
			return e, true
		}
	}
	return PreviewDatabaseEngineV1{}, false
}

// SupportsVersion reports whether version is one of the engine's versions.
func (e PreviewDatabaseEngineV1) SupportsVersion(version string) bool {
	for _, v := range e.Versions {
		if v == version {
			return true
		}
	}
	return false
}

// Image is the image reference for version, or an error when the version isn't
// allowlisted ("" means the default).
func (e PreviewDatabaseEngineV1) Image(version string) (string, error) {
	if version == "" {
		version = e.DefaultVersion
	}
	if !e.SupportsVersion(version) {
		return "", fmt.Errorf("%s %s is not a supported preview database version (supported: %s)",
			e.Engine, version, strings.Join(e.Versions, ", "))
	}
	return e.ImageRepository + ":" + version, nil
}

// DatabaseVariableKind is what a variable name says about holding a database's
// address.
type DatabaseVariableKind int

const (
	// NotDatabaseVariable: the name doesn't look like a database's.
	NotDatabaseVariable DatabaseVariableKind = iota
	// DatabaseURLVariable holds the whole address (e.g. DATABASE_URL, MONGODB_URI).
	DatabaseURLVariable
	// DatabaseSplitVariable holds one part of it (e.g. DB_HOST, PGPASSWORD).
	DatabaseSplitVariable
)

// The part of the address a split variable holds.
const (
	DatabasePartHost     = "host"
	DatabasePartPort     = "port"
	DatabasePartUser     = "user"
	DatabasePartPassword = "password"
	DatabasePartName     = "name"
)

var databaseNameTokens = []string{"DATABASE", "DB", "POSTGRES", "PG", "MYSQL", "MONGO"}

var databaseURLSuffixes = []string{"URL", "URI", "DSN"}

var databaseSplitSuffixes = []struct{ suffix, part string }{
	{"_HOST", DatabasePartHost},
	{"_PORT", DatabasePartPort},
	{"_USER", DatabasePartUser},
	{"_USERNAME", DatabasePartUser},
	{"_PASSWORD", DatabasePartPassword},
	{"_PASS", DatabasePartPassword},
	{"_NAME", DatabasePartName},
	{"_DATABASE", DatabasePartName},
	{"_DB", DatabasePartName},
}

var libpqVariables = map[string]string{
	"PGHOST":     DatabasePartHost,
	"PGPORT":     DatabasePartPort,
	"PGUSER":     DatabasePartUser,
	"PGPASSWORD": DatabasePartPassword,
	"PGDATABASE": DatabasePartName,
}

// ClassifyDatabaseVariable says whether name looks like a variable holding a
// database's address, and for a split variable which part (DatabasePart*):
//   - a URL variable contains DATABASE, DB, POSTGRES, PG, MYSQL or MONGO and ends
//     in URL, URI or DSN;
//   - a split variable is one of those followed by _HOST, _PORT, _USER, _USERNAME,
//     _PASSWORD, _PASS, _NAME, _DATABASE or _DB, or one of libpq's PGHOST, PGPORT,
//     PGUSER, PGPASSWORD, PGDATABASE.
//
// Names are compared case-sensitively, upper-case, as environment variables are
// conventionally written.
func ClassifyDatabaseVariable(name string) (DatabaseVariableKind, string) {
	if part, ok := libpqVariables[name]; ok {
		return DatabaseSplitVariable, part
	}
	for _, token := range databaseNameTokens {
		for _, s := range databaseSplitSuffixes {
			if name == token+s.suffix {
				return DatabaseSplitVariable, s.part
			}
		}
	}
	for _, suffix := range databaseURLSuffixes {
		if !strings.HasSuffix(name, suffix) {
			continue
		}
		for _, token := range databaseNameTokens {
			if strings.Contains(name, token) {
				return DatabaseURLVariable, ""
			}
		}
	}
	return NotDatabaseVariable, ""
}
