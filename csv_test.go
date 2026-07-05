package main

import (
	"encoding/csv"
	"strings"
	"testing"
)

func splitRecord(t *testing.T, line string) []string {
	t.Helper()
	reader := strings.NewReader(line)
	rec, err := csv.NewReader(reader).Read()
	if err != nil {
		t.Fatalf("splitRecord helper failed: %v", err)
	}

	return rec
}

func TestParseRecord(t *testing.T) {
	cases := []string{
		"6/2/2026 7:08 AM,207.5 lb,29.7 ,25.7 %,154.1 lb,22.1 %,12 ,53.5 %,146.2 lb,47.9 %,7.7 lb,16.9 %,2003 kcal,37",
		"6/5/2026 7:28 AM,207.7lb,29.7 ,25.7 %,154.2lb,22.1%,12,53.5%,146.4lb,47.9 %,7.7lb,16.9 %,2004kcal,37",
	}

	for i := range cases {
		rec := splitRecord(t, cases[i])

		_, err := parseRecord(rec)
		if err != nil {
			t.Fatalf("parseRecord failed (case %d): %v", i, err)
		}
	}
}
