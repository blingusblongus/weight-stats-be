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
		"7/18/2026, 5:26 AM,197.4lb,28.3 ,23.6 %,150.8lb,20.4%,11,55.1%,143.1lb,49.3 %,7.5lb,17.4 %,1941kcal,37",
	}

	for i := range cases {
		rec := splitRecord(t, cases[i])

		_, err := parseRecord(rec)
		if err != nil {
			t.Fatalf("parseRecord failed (case %d): %v", i, err)
		}
	}
}

func TestParseCSVNewFormat(t *testing.T) {
	// The new export format has 14 header columns but 15 fields per data row
	// because of the unquoted comma in the timestamp.
	data := `Time,Weight,BMI,Body Fat,Fat-Free Body Weight,Subcutaneous Fat,Visceral Fat,Body Water,Muscle Mass,Skeletal Muscles,Bone Mass,Protein,BMR,Metabolic Age
7/18/2026, 5:26 AM,197.4lb,28.3 ,23.6 %,150.8lb,20.4%,11,55.1%,143.1lb,49.3 %,7.5lb,17.4 %,1941kcal,37
7/19/2026, 10:56 AM,200.3lb,28.7 ,24.2 %,151.8lb,20.9%,11,54.6%,144.0lb,48.9 %,7.6lb,17.2 %,1959kcal,37
`

	measurements, err := parseCSV(strings.NewReader(data))
	if err != nil {
		t.Fatalf("parseCSV failed: %v", err)
	}
	if len(measurements) != 2 {
		t.Fatalf("expected 2 measurements, got %d", len(measurements))
	}
	if measurements[0].MeasuredAt != "2026-07-18T05:26:00" {
		t.Errorf("wrong timestamp: %s", measurements[0].MeasuredAt)
	}
	if measurements[0].Weight != 197.4 {
		t.Errorf("wrong weight: %v", measurements[0].Weight)
	}
}
