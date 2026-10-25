package task_previews

import (
	"strings"
	"testing"
)

func TestPreviewDatabaseAllowlist(t *testing.T) {
	want := map[string]struct {
		versions, def, image, cpu, memory string
		port                              int64
	}{
		PreviewDatabaseEnginePostgres: {"13,14,15,16,17", "16", "public.ecr.aws/docker/library/postgres:16", "256", "512", 5432},
		PreviewDatabaseEngineMySQL:    {"8.0,8.4", "8.4", "public.ecr.aws/docker/library/mysql:8.4", "256", "1024", 3306},
		PreviewDatabaseEngineMongoDB:  {"6.0,7.0", "7.0", "public.ecr.aws/docker/library/mongo:7.0", "256", "1024", 27017},
	}
	engines := PreviewDatabaseEngines()
	if len(engines) != len(want) {
		t.Fatalf("got %d engines; want %d", len(engines), len(want))
	}
	for _, e := range engines {
		w, ok := want[e.Engine]
		if !ok {
			t.Fatalf("unexpected engine %q", e.Engine)
		}
		if got := strings.Join(e.Versions, ","); got != w.versions {
			t.Errorf("%s versions = %s; want %s", e.Engine, got, w.versions)
		}
		if e.DefaultVersion != w.def || !e.SupportsVersion(e.DefaultVersion) {
			t.Errorf("%s default %q must be %q and in its versions", e.Engine, e.DefaultVersion, w.def)
		}
		img, err := e.Image("")
		if err != nil || img != w.image {
			t.Errorf("%s Image(\"\") = %q, %v; want %q", e.Engine, img, err, w.image)
		}
		if e.Port != w.port || e.Cpu != w.cpu || e.Memory != w.memory {
			t.Errorf("%s port/cpu/memory = %d/%s/%s", e.Engine, e.Port, e.Cpu, e.Memory)
		}
	}
	pg, _ := LookupPreviewDatabaseEngine(PreviewDatabaseEnginePostgres)
	if _, err := pg.Image("12"); err == nil {
		t.Error("postgres 12 must not be allowed")
	}
	if img, _ := pg.Image("13"); img != "public.ecr.aws/docker/library/postgres:13" {
		t.Errorf("postgres 13 image = %q", img)
	}
	for _, other := range []string{"mariadb", "redis", "sqlite", "cockroachdb", ""} {
		if _, ok := LookupPreviewDatabaseEngine(other); ok {
			t.Errorf("%q must not be allowlisted", other)
		}
	}
	// A caller mutating the returned list must not change the allowlist.
	engines[0].Versions[0] = "1"
	if PreviewDatabaseEngines()[0].Versions[0] == "1" {
		t.Error("PreviewDatabaseEngines must return a copy")
	}
}

func TestClassifyDatabaseVariable(t *testing.T) {
	cases := []struct {
		name string
		kind DatabaseVariableKind
		part string
	}{
		{"DATABASE_URL", DatabaseURLVariable, ""},
		{"POSTGRES_URL", DatabaseURLVariable, ""},
		{"MONGODB_URI", DatabaseURLVariable, ""},
		{"MYSQL_DSN", DatabaseURLVariable, ""},
		{"DIRECT_DATABASE_URL", DatabaseURLVariable, ""},
		{"PG_URL", DatabaseURLVariable, ""},
		{"DB_HOST", DatabaseSplitVariable, DatabasePartHost},
		{"DB_PORT", DatabaseSplitVariable, DatabasePartPort},
		{"DB_USER", DatabaseSplitVariable, DatabasePartUser},
		{"DB_USERNAME", DatabaseSplitVariable, DatabasePartUser},
		{"DB_PASSWORD", DatabaseSplitVariable, DatabasePartPassword},
		{"DB_PASS", DatabaseSplitVariable, DatabasePartPassword},
		{"DB_NAME", DatabaseSplitVariable, DatabasePartName},
		{"DB_DATABASE", DatabaseSplitVariable, DatabasePartName},
		{"MONGO_DB", DatabaseSplitVariable, DatabasePartName},
		{"POSTGRES_HOST", DatabaseSplitVariable, DatabasePartHost},
		{"MYSQL_PASSWORD", DatabaseSplitVariable, DatabasePartPassword},
		{"DATABASE_NAME", DatabaseSplitVariable, DatabasePartName},
		{"PGHOST", DatabaseSplitVariable, DatabasePartHost},
		{"PGPORT", DatabaseSplitVariable, DatabasePartPort},
		{"PGUSER", DatabaseSplitVariable, DatabasePartUser},
		{"PGPASSWORD", DatabaseSplitVariable, DatabasePartPassword},
		{"PGDATABASE", DatabaseSplitVariable, DatabasePartName},
		{"API_URL", NotDatabaseVariable, ""},
		{"REDIS_URL", NotDatabaseVariable, ""},
		{"PORT", NotDatabaseVariable, ""},
		{"DATABASE", NotDatabaseVariable, ""},
		{"SECRET_KEY", NotDatabaseVariable, ""},
		{"REDIS_HOST", NotDatabaseVariable, ""},
	}
	for _, c := range cases {
		kind, part := ClassifyDatabaseVariable(c.name)
		if kind != c.kind || part != c.part {
			t.Errorf("ClassifyDatabaseVariable(%q) = %d, %q; want %d, %q", c.name, kind, part, c.kind, c.part)
		}
	}
}
