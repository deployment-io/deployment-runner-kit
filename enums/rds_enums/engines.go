package rds_enums

type Engine uint

const (
	Postgres Engine = 1
	MySQL    Engine = 2
	MariaDB  Engine = 3
)

func (e Engine) String() string {
	switch e {
	case Postgres:
		return "postgres"
	case MySQL:
		return "mysql"
	case MariaDB:
		return "mariadb"
	}
	return "unknown"
}

func (e Engine) Name() string {
	switch e {
	case Postgres:
		return "PostgreSQL"
	case MySQL:
		return "MySQL"
	case MariaDB:
		return "MariaDB"
	}
	return "unknown"
}

func (e Engine) EnablePerformanceInsights() bool {
	if e == Postgres {
		return true
	}
	return false
}

func (e Engine) GetEngineLifecycleSupport() *string {
	disabled := "open-source-rds-extended-support-disabled"
	if e == Postgres || e == MySQL {
		return &disabled
	}
	return nil
}

func (e Engine) GetLicenseModel() *string {
	if e == Postgres {
		s := "postgresql-license"
		return &s
	}
	if e == MySQL || e == MariaDB {
		s := "general-public-license"
		return &s
	}
	return nil
}
