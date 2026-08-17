package aws_policy_schema

import (
	"encoding/json"
	"reflect"
	"testing"
)

const sid = "deploymentIoAdded"

// statementByIndex pulls out statement i of a document as a generic map, so tests can
// assert on keys the schema types don't model.
func statementByIndex(t *testing.T, document []byte, i int) map[string]interface{} {
	t.Helper()
	var doc struct {
		Statement []map[string]interface{} `json:"Statement"`
	}
	if err := json.Unmarshal(document, &doc); err != nil {
		t.Fatalf("result is not valid JSON: %s\n%s", err, document)
	}
	if i >= len(doc.Statement) {
		t.Fatalf("wanted statement %d, document only has %d:\n%s", i, len(doc.Statement), document)
	}
	return doc.Statement[i]
}

func actionsOf(t *testing.T, statement map[string]interface{}) []string {
	t.Helper()
	raw, ok := statement["Action"].([]interface{})
	if !ok {
		t.Fatalf("Action is not a list: %#v", statement["Action"])
	}
	var actions []string
	for _, a := range raw {
		actions = append(actions, a.(string))
	}
	return actions
}

// Every one of these documents is legal IAM that the PolicyDocumentData struct either
// rejects outright or quietly mangles. What we assert is that the statement we don't
// own comes back out identical to what went in.
func TestAddActionsToSid_PreservesForeignStatements(t *testing.T) {
	tests := []struct {
		name    string
		foreign string
	}{
		{
			name:    "bare string Action",
			foreign: `{"Sid":"custom","Effect":"Allow","Action":"s3:GetObject","Resource":["*"]}`,
		},
		{
			name:    "bare string Resource",
			foreign: `{"Sid":"custom","Effect":"Allow","Action":["s3:GetObject"],"Resource":"arn:aws:s3:::bucket/*"}`,
		},
		{
			name:    "unmodeled Condition operator",
			foreign: `{"Sid":"custom","Effect":"Deny","Action":["*"],"Resource":["*"],"Condition":{"Bool":{"aws:SecureTransport":"false"}}}`,
		},
		{
			name:    "Condition with DateLessThan and StringNotEquals",
			foreign: `{"Sid":"custom","Effect":"Allow","Action":["s3:*"],"Resource":["*"],"Condition":{"DateLessThan":{"aws:CurrentTime":"2030-01-01T00:00:00Z"},"StringNotEquals":{"aws:PrincipalTag/team":"ops"}}}`,
		},
		{
			name:    "NotAction",
			foreign: `{"Sid":"custom","Effect":"Deny","NotAction":["s3:Get*"],"Resource":["*"]}`,
		},
		{
			name:    "NotResource",
			foreign: `{"Sid":"custom","Effect":"Deny","Action":["s3:*"],"NotResource":["arn:aws:s3:::public/*"]}`,
		},
		{
			name:    "Principal with AWS and Federated",
			foreign: `{"Sid":"custom","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::123456789012:root","Federated":"cognito-identity.amazonaws.com"},"Action":["sts:AssumeRole"],"Resource":["*"]}`,
		},
		{
			name:    "Principal wildcard string",
			foreign: `{"Sid":"custom","Effect":"Allow","Principal":"*","Action":["sts:AssumeRole"],"Resource":["*"]}`,
		},
		{
			name:    "statement with no Sid",
			foreign: `{"Effect":"Allow","Action":"ecr:GetAuthorizationToken","Resource":"*"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := []byte(`{"Version":"2012-10-17","Statement":[` + test.foreign +
				`,{"Sid":"` + sid + `","Effect":"Allow","Action":["ecs:DescribeTasks"],"Resource":["*"]}]}`)

			updated, changed, err := AddActionsToSid(document, sid, []string{"secretsmanager:GetSecretValue"})
			if err != nil {
				t.Fatalf("AddActionsToSid: %s", err)
			}
			if !changed {
				t.Fatal("changed = false, want true - a new action was added")
			}

			var want, got interface{}
			if err = json.Unmarshal([]byte(test.foreign), &want); err != nil {
				t.Fatalf("bad test fixture: %s", err)
			}
			gotRaw, err := json.Marshal(statementByIndex(t, updated, 0))
			if err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal(gotRaw, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("foreign statement was altered by the round-trip\n want: %s\n got:  %s", test.foreign, gotRaw)
			}

			//and our own statement did get the action
			ours := actionsOf(t, statementByIndex(t, updated, 1))
			if !reflect.DeepEqual(ours, []string{"ecs:DescribeTasks", "secretsmanager:GetSecretValue"}) {
				t.Errorf("our actions = %v", ours)
			}
		})
	}
}

func TestAddActionsToSid_NoMatchingStatementAppendsOne(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":[{"Sid":"custom","Effect":"Deny","NotAction":["s3:Get*"],"Resource":["*"]}]}`)

	updated, changed, err := AddActionsToSid(document, sid, []string{"secretsmanager:GetSecretValue"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}

	added := statementByIndex(t, updated, 1)
	if added["Sid"] != sid || added["Effect"] != "Allow" {
		t.Errorf("appended statement = %#v", added)
	}
	if got := actionsOf(t, added); !reflect.DeepEqual(got, []string{"secretsmanager:GetSecretValue"}) {
		t.Errorf("appended actions = %v", got)
	}
	if got := statementByIndex(t, updated, 0); got["NotAction"] == nil {
		t.Errorf("NotAction was dropped from the untouched statement: %#v", got)
	}
}

// Our own statement is legal IAM too - it may have been hand-edited into a bare string.
func TestAddActionsToSid_OurStatementWithBareStringAction(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":[{"Sid":"` + sid + `","Effect":"Allow","Action":"ecs:DescribeTasks","Resource":["*"]}]}`)

	updated, changed, err := AddActionsToSid(document, sid, []string{"secretsmanager:GetSecretValue"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}
	got := actionsOf(t, statementByIndex(t, updated, 0))
	if !reflect.DeepEqual(got, []string{"ecs:DescribeTasks", "secretsmanager:GetSecretValue"}) {
		t.Errorf("actions = %v", got)
	}
}

// A statement of ours that carries a Condition must keep it - this is the silent
// permission-widening case.
func TestAddActionsToSid_KeepsConditionOnOurStatement(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":[{"Sid":"` + sid + `","Effect":"Allow","Action":["ecs:DescribeTasks"],"Resource":["*"],"Condition":{"Bool":{"aws:SecureTransport":"true"}}}]}`)

	updated, _, err := AddActionsToSid(document, sid, []string{"secretsmanager:GetSecretValue"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	condition, ok := statementByIndex(t, updated, 0)["Condition"].(map[string]interface{})
	if !ok || condition["Bool"] == nil {
		t.Errorf("Condition was dropped or emptied: %#v", condition)
	}
}

func TestAddActionsToSid_NoNewActionsIsANoOp(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":[{"Sid":"` + sid + `","Effect":"Allow","Action":["ecs:DescribeTasks"],"Resource":["*"]}]}`)

	updated, changed, err := AddActionsToSid(document, sid, []string{"ecs:DescribeTasks"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	if changed {
		t.Error("changed = true, want false - nothing new to add")
	}
	if string(updated) != string(document) {
		t.Errorf("document was rewritten on a no-op:\n%s", updated)
	}
}

func TestAddActionsToSid_DedupsAndKeepsOrder(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":[{"Sid":"` + sid + `","Effect":"Allow","Action":["a:one"],"Resource":["*"]}]}`)

	updated, _, err := AddActionsToSid(document, sid, []string{"a:one", "b:two", "b:two", "c:three"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	got := actionsOf(t, statementByIndex(t, updated, 0))
	if !reflect.DeepEqual(got, []string{"a:one", "b:two", "c:three"}) {
		t.Errorf("actions = %v", got)
	}
}

// IAM accepts a single statement object where a list is usual.
func TestAddActionsToSid_SingleStatementObject(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Statement":{"Sid":"custom","Effect":"Allow","Action":"s3:GetObject","Resource":"*"}}`)

	updated, changed, err := AddActionsToSid(document, sid, []string{"secretsmanager:GetSecretValue"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}
	if got := statementByIndex(t, updated, 0); got["Action"] != "s3:GetObject" {
		t.Errorf("original statement altered: %#v", got)
	}
	if got := actionsOf(t, statementByIndex(t, updated, 1)); !reflect.DeepEqual(got, []string{"secretsmanager:GetSecretValue"}) {
		t.Errorf("appended actions = %v", got)
	}
}

func TestAddActionsToSid_PreservesUnknownTopLevelKeys(t *testing.T) {
	document := []byte(`{"Version":"2012-10-17","Id":"policy-id","Statement":[{"Sid":"` + sid + `","Effect":"Allow","Action":["a:one"],"Resource":["*"]}]}`)

	updated, _, err := AddActionsToSid(document, sid, []string{"b:two"})
	if err != nil {
		t.Fatalf("AddActionsToSid: %s", err)
	}
	var doc map[string]interface{}
	if err = json.Unmarshal(updated, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["Id"] != "policy-id" || doc["Version"] != "2012-10-17" {
		t.Errorf("top level keys lost: %#v", doc)
	}
}

func TestAddActionsToSid_InvalidDocument(t *testing.T) {
	if _, _, err := AddActionsToSid([]byte(`not json`), sid, []string{"a:one"}); err == nil {
		t.Error("expected an error for a non-JSON document")
	}
	//an Action we can't read at all is worth failing on rather than overwriting
	document := []byte(`{"Statement":[{"Sid":"` + sid + `","Action":{"nope":true}}]}`)
	if _, _, err := AddActionsToSid(document, sid, []string{"a:one"}); err == nil {
		t.Error("expected an error for an unreadable Action")
	}
}
