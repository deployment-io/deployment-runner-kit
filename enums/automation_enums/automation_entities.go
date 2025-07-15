package automation_enums

type Entity uint

const (
	WebService Entity = iota + 1
	StaticSite
	PrivateService
	Environment
	Database
	Repository

	EntityMax //add others before EntityMax
)

var entityToString = map[Entity]string{
	WebService:     "Web service automation",
	StaticSite:     "Static site automation",
	PrivateService: "Private service automation",
	Environment:    "Environment automation",
	Database:       "Database automation",
	Repository:     "Repository automation",
}

func (e Entity) String() string {
	return entityToString[e]
}

// Define a map of tool types for each Entity
var entityToTools = map[Entity][]ToolType{
	WebService:     {GetCPUMemoryUsage},
	PrivateService: {GetCPUMemoryUsage},
	Repository:     {QueryCode},
	//StaticSite:     {"Hugo", "Jekyll"},
	//Environment:    {"Terraform", "Ansible"},
	//Database:       {"PostgreSQL", "MySQL"},
}

// ToolTypes function returns the tool types that an entity can use
func (e Entity) ToolTypes() []ToolType {
	return entityToTools[e]
}

func (e Entity) IsServiceType() bool {
	return e == PrivateService || e == WebService || e == StaticSite
}

func (e Entity) IsBackendServiceType() bool {
	return e == PrivateService || e == WebService
}

func (e Entity) IsDatabaseType() bool {
	return e == Database
}
