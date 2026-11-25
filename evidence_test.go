package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestEvidenceReflectsSnapshot(t *testing.T) {
	var out, log bytes.Buffer
	code := run([]string{"evidence", "--observed-at", "8", "--case-id", "exercise", "examples/partial"}, &out, &log)
	var got struct {
		Schema string `json:"schema"`
		CaseID string `json:"case_id"`
		Events []observation `json:"events"`
	}
	if code != 0 || json.Unmarshal(out.Bytes(), &got) != nil || len(got.Events) != 2 || got.Schema != "evidence.v1" || got.CaseID != "exercise" {
		t.Fatalf("code %d, sortie %s, erreur %s", code, out.String(), log.String())
	}
	if got.Events[0].Facts["expected"] != got.Events[0].Facts["observed"] || got.Events[1].Facts["observed"] != nil || got.Events[1].ID != "revaskel-2" {
		t.Fatalf("observations inattendues: %+v", got.Events)
	}
}

func TestEvidenceRequiresObservationTime(t *testing.T) {
	var out, log bytes.Buffer
	if run([]string{"evidence", "--case-id", "exercise", "examples/partial"}, &out, &log) != 2 || out.Len() != 0 {
		t.Fatal("instant manquant accepté")
	}
}
