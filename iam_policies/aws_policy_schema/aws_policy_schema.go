package aws_policy_schema

type Principal struct {
	Service []string `json:"Service,omitempty"`
}

type Condition struct {
	ArnLike      map[string]string `json:"ArnLike,omitempty"`
	StringEquals map[string]string `json:"StringEquals,omitempty"`
	StringLike   map[string]string `json:"StringLike,omitempty"`
}

type Statement struct {
	Sid       string     `json:"Sid,omitempty"`
	Effect    string     `json:"Effect,omitempty"`
	Principal *Principal `json:"Principal,omitempty"`
	Action    []string   `json:"Action,omitempty"`
	Resource  []string   `json:"Resource,omitempty"`
	Condition *Condition `json:"Condition,omitempty"`
}

type PolicyDocumentData struct {
	Version   string      `json:"Version,omitempty"`
	Statement []Statement `json:"Statement,omitempty"`
}
