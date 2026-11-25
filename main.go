// Revaskel vérifie le préfixe d'expéditions récupérables d'un instantané local.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"time"
)

type Shipment struct {
	ID string `json:"id"`
	At time.Time `json:"committed_at"`
	Object string `json:"object"`
	SHA string `json:"sha256"`
}

type Check struct {
	ID string `json:"id"`
	Result string `json:"result"`
}

type Report struct {
	Complete bool `json:"complete"`
	LastCoherent *time.Time `json:"last_coherent_at"`
	Checks []Check `json:"checks"`
}

var objectName = regexp.MustCompile(`^[a-zA-Z0-9_-][a-zA-Z0-9_.-]*$`)

func Evaluate(snapshot fs.FS) (Report, error) {
	raw, err := fs.ReadFile(snapshot, "shipments.json")
	if err != nil { return Report{}, err }
	var rows []Shipment
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&rows); err != nil { return Report{}, err }
	if decoder.Decode(new(any)) != io.EOF { return Report{}, errors.New("données après le manifeste") }
	if len(rows) == 0 { return Report{}, errors.New("instantané vide") }

	report := Report{Complete: true}
	seen := make(map[string]bool)
	var previous time.Time
	for _, row := range rows {
		digest, err := hex.DecodeString(row.SHA)
		if err != nil || len(digest) != sha256.Size || row.ID == "" || seen[row.ID] ||
			row.At.IsZero() || !row.At.After(previous) || !objectName.MatchString(row.Object) {
			return Report{}, fmt.Errorf("expédition invalide ou désordonnée %q", row.ID)
		}
		previous, seen[row.ID] = row.At, true
		content, err := fs.ReadFile(snapshot, path.Join("objects", row.Object))
		result := "ok"
		switch {
		case errors.Is(err, fs.ErrNotExist):
			result = "missing_document"
		case err != nil:
			return Report{}, err
		default:
			actual := sha256.Sum256(content)
			if !bytes.Equal(actual[:], digest) { result = "digest_mismatch" }
		}
		report.Checks = append(report.Checks, Check{row.ID, result})
		if result != "ok" { report.Complete = false }
		if report.Complete {
			at := row.At
			report.LastCoherent = &at
		}
	}
	return report, nil
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: revaskel <répertoire-instantané>")
		return 2
	}
	report, err := Evaluate(os.DirFS(args[0]))
	if err != nil { fmt.Fprintln(stderr, err); return 2 }
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil { fmt.Fprintln(stderr, err); return 2 }
	if !report.Complete { return 1 }
	return 0
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
