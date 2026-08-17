package aws_policy_schema

import (
	"encoding/json"
	"fmt"
)

// AddActionsToSid appends actions to the statement carrying the given Sid inside an
// IAM policy document, and returns the updated document.
//
// The document is treated as opaque JSON: only the statement matching sid is
// rewritten, and only its Action list. Every other statement is carried across
// byte-for-byte, so operators (Bool, StringNotEquals, DateLessThan, ...), NotAction,
// NotResource, Principal shapes we don't model, and any future IAM key survive the
// round-trip untouched. Do NOT replace this with a marshal of PolicyDocumentData —
// that struct is lossy and dropping a Condition silently widens a permission.
//
// If no statement carries sid, one is appended as
// {"Sid":sid,"Effect":"Allow","Action":[...],"Resource":["*"]}.
//
// Actions already present on the statement are skipped. changed reports whether
// the document actually needs to be written back; when it is false the returned
// document is the input unchanged.
func AddActionsToSid(document []byte, sid string, actions []string) (updated []byte, changed bool, err error) {
	var doc map[string]json.RawMessage
	if err = json.Unmarshal(document, &doc); err != nil {
		return nil, false, fmt.Errorf("error parsing policy document: %s", err)
	}

	statements, err := splitStatements(doc["Statement"])
	if err != nil {
		return nil, false, err
	}

	index := -1
	var existing map[string]json.RawMessage
	for i, raw := range statements {
		var statement map[string]json.RawMessage
		if err = json.Unmarshal(raw, &statement); err != nil {
			//a statement we can't even open as an object isn't ours - leave it alone
			continue
		}
		var statementSid string
		if err = json.Unmarshal(statement["Sid"], &statementSid); err != nil {
			continue
		}
		if statementSid == sid {
			index, existing = i, statement
			break
		}
	}

	var oldActions []string
	if existing != nil {
		if oldActions, err = decodeStringOrSlice(existing["Action"]); err != nil {
			return nil, false, fmt.Errorf("error parsing Action of statement %s: %s", sid, err)
		}
	}

	oldActionsSet := make(map[string]bool, len(oldActions))
	for _, action := range oldActions {
		oldActionsSet[action] = true
	}
	newActions := oldActions
	for _, action := range actions {
		if !oldActionsSet[action] {
			oldActionsSet[action] = true
			newActions = append(newActions, action)
		}
	}
	if len(newActions) == len(oldActions) {
		return document, false, nil
	}

	encodedActions, err := json.Marshal(newActions)
	if err != nil {
		return nil, false, err
	}

	if existing == nil {
		existing = map[string]json.RawMessage{
			"Sid":      json.RawMessage(mustMarshalString(sid)),
			"Effect":   json.RawMessage(`"Allow"`),
			"Resource": json.RawMessage(`["*"]`),
		}
	}
	existing["Action"] = encodedActions

	//only the matched statement is re-encoded - the rest stay as the raw bytes we read
	encodedStatement, err := json.Marshal(existing)
	if err != nil {
		return nil, false, err
	}
	if index == -1 {
		statements = append(statements, encodedStatement)
	} else {
		statements[index] = encodedStatement
	}

	if doc["Statement"], err = json.Marshal(statements); err != nil {
		return nil, false, err
	}
	updated, err = json.Marshal(doc)
	if err != nil {
		return nil, false, err
	}
	return updated, true, nil
}

// splitStatements returns the document's statements as raw JSON. IAM accepts both a
// list of statements and a single bare statement object.
func splitStatements(raw json.RawMessage) ([]json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var statements []json.RawMessage
	if err := json.Unmarshal(raw, &statements); err == nil {
		return statements, nil
	}
	var single map[string]json.RawMessage
	if err := json.Unmarshal(raw, &single); err != nil {
		return nil, fmt.Errorf("error parsing policy statements: %s", err)
	}
	return []json.RawMessage{raw}, nil
}

// decodeStringOrSlice reads an IAM field that is either a bare string
// ("Action":"s3:GetObject") or a list of them.
func decodeStringOrSlice(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var slice []string
	if err := json.Unmarshal(raw, &slice); err == nil {
		return slice, nil
	}
	var single string
	if err := json.Unmarshal(raw, &single); err != nil {
		return nil, fmt.Errorf("expected a string or list of strings, got %s", raw)
	}
	return []string{single}, nil
}

func mustMarshalString(s string) []byte {
	encoded, err := json.Marshal(s)
	if err != nil {
		//json.Marshal of a string can't fail
		panic(err)
	}
	return encoded
}
