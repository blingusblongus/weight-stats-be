package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type Measurement struct {
	MeasuredAt      string  `json:"measured_at"`
	Weight          float64 `json:"weight"`
	BMI             float64 `json:"bmi"`
	BodyFat         float64 `json:"body_fat"`
	FatFreeWeight   float64 `json:"fat_free_weight"`
	SubcutaneousFat float64 `json:"subcutaneous_fat"`
	VisceralFat     int     `json:"visceral_fat"`
	BodyWater       float64 `json:"body_water"`
	MuscleMass      float64 `json:"muscle_mass"`
	SkeletalMuscles float64 `json:"skeletal_muscles"`
	BoneMass        float64 `json:"bone_mass"`
	Protein         float64 `json:"protein"`
	BMR             int     `json:"bmr"`
	MetabolicAge    int     `json:"metabolic_age"`
}

const csvTimestampLayout = "1/2/2006 3:04 PM"

func parseCSV(r io.Reader) ([]Measurement, error) {
	reader := csv.NewReader(r)

	// Skip header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if len(header) < 14 {
		return nil, fmt.Errorf("expected at least 14 columns, got %d", len(header))
	}

	var measurements []Measurement
	lineNum := 1
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("csv line %d: skipping, read error: %v", lineNum, err)
			continue
		}
		if len(record) < 14 {
			log.Printf("csv line %d: skipping, only %d columns", lineNum, len(record))
			continue
		}

		m, err := parseRecord(record)
		if err != nil {
			log.Printf("csv line %d: skipping, parse error: %v", lineNum, err)
			continue
		}
		measurements = append(measurements, m)
	}

	return measurements, nil
}

func parseRecord(record []string) (Measurement, error) {
	var m Measurement

	t, err := time.Parse(csvTimestampLayout, strings.TrimSpace(record[0]))
	if err != nil {
		return m, fmt.Errorf("timestamp %q: %w", record[0], err)
	}
	m.MeasuredAt = t.Format("2006-01-02T15:04:05")

	m.Weight, err = parseFloat(record[1])
	if err != nil {
		return m, fmt.Errorf("weight: %w", err)
	}

	m.BMI, _ = parseFloat(record[2])
	m.BodyFat, _ = parseFloat(record[3])
	m.FatFreeWeight, _ = parseFloat(record[4])
	m.SubcutaneousFat, _ = parseFloat(record[5])
	m.VisceralFat, _ = parseInt(record[6])
	m.BodyWater, _ = parseFloat(record[7])
	m.MuscleMass, _ = parseFloat(record[8])
	m.SkeletalMuscles, _ = parseFloat(record[9])
	m.BoneMass, _ = parseFloat(record[10])
	m.Protein, _ = parseFloat(record[11])
	m.BMR, _ = parseInt(record[12])
	m.MetabolicAge, _ = parseInt(record[13])

	return m, nil
}

func stripUnit(s string) string {
	s = strings.TrimSpace(s)
	for _, suffix := range []string{" lb", " %", " kcal"} {
		s = strings.TrimSuffix(s, suffix)
	}
	return strings.TrimSpace(s)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(stripUnit(s), 64)
}

func parseInt(s string) (int, error) {
	v, err := strconv.Atoi(stripUnit(s))
	return v, err
}
