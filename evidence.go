package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type observation struct {
	ID          string            `json:"id"`
	Time        int64             `json:"time_s"`
	Uncertainty int64             `json:"uncertainty_s"`
	Kind        string            `json:"kind"`
	Subject     string            `json:"subject"`
	Facts       map[string]any    `json:"facts"`
}

func exportEvidence(args []string, out, errOut io.Writer) int {
	options := flag.NewFlagSet("evidence", flag.ContinueOnError)
	options.SetOutput(errOut)
	caseID := options.String("case-id", "", "identifiant commun de l’exercice")
	timeS := options.Int64("observed-at", 0, "seconde de l’observation sur l’horloge de l’exercice")
	errorS := options.Int64("uncertainty-s", 0, "incertitude en secondes")
	if err := options.Parse(args); err != nil {
		return 2
	}
	hasTime := false
	options.Visit(func(f *flag.Flag) {
		if f.Name == "observed-at" {
			hasTime = true
		}
	})
	if options.NArg() != 1 || strings.TrimSpace(*caseID) == "" || *errorS < 0 || !hasTime {
		fmt.Fprintln(errOut, "usage: revaskel evidence --case-id ID --observed-at S [--uncertainty-s N] INSTANTANÉ")
		return 2
	}
	report, err := Evaluate(os.DirFS(options.Arg(0)))
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	events := make([]observation, 0, len(report.Checks))
	for i, check := range report.Checks {
		events = append(events, observation{
			ID: fmt.Sprintf("revaskel-%d", i+1), Time: *timeS, Uncertainty: *errorS,
			Kind: "document", Subject: check.ID,
			Facts: map[string]any{"expected": check.expected, "observed": check.observed},
		})
	}
	packet := map[string]any{"schema": "evidence.v1", "source": "revaskel", "case_id": *caseID, "events": events}
	if err := json.NewEncoder(out).Encode(packet); err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	return 0
}
