package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
	"testing/fstest"
	"time"
)

func snapshot(names []string) fstest.MapFS {
	rows := []Shipment{}
	for i, name := range names {
		rows = append(rows, Shipment{fmt.Sprint(i), time.Unix(int64(i), 0).UTC(), name, fmt.Sprintf("%x", sha256.Sum256([]byte("ok")))})
	}
	raw, _ := json.Marshal(rows)
	return fstest.MapFS{"shipments.json": &fstest.MapFile{Data: raw}, "objects/a": &fstest.MapFile{Data: []byte("ok")}}
}

func TestComplete(t *testing.T) {
	rep, err := Evaluate(snapshot([]string{"a"}))
	if err != nil || !rep.Complete || rep.LastCoherent == nil { t.Fatalf("%+v %v", rep, err) }
}

func TestGapBreaksPrefix(t *testing.T) {
	rep, err := Evaluate(snapshot([]string{"a", "missing", "a"}))
	if err != nil || rep.Complete || rep.Checks[1].Result != "missing_document" ||
		rep.LastCoherent == nil || rep.LastCoherent.Unix() != 0 { t.Fatalf("%+v %v", rep, err) }
}

func TestWrongDigest(t *testing.T) {
	data := snapshot([]string{"a"})
	data["objects/a"].Data = []byte("wrong")
	rep, err := Evaluate(data)
	if err != nil || rep.LastCoherent != nil || rep.Checks[0].Result != "digest_mismatch" { t.Fatalf("%+v %v", rep, err) }
}

func TestInvalidName(t *testing.T) {
	if _, err := Evaluate(snapshot([]string{"../a"})); err == nil { t.Fatal("invalid name accepted") }
}
