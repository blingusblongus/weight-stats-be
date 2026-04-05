package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func initDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS measurements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			measured_at TEXT UNIQUE NOT NULL,
			weight REAL NOT NULL,
			bmi REAL,
			body_fat REAL,
			fat_free_weight REAL,
			subcutaneous_fat REAL,
			visceral_fat INTEGER,
			body_water REAL,
			muscle_mass REAL,
			skeletal_muscles REAL,
			bone_mass REAL,
			protein REAL,
			bmr INTEGER,
			metabolic_age INTEGER
		)
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	return db
}

func insertMeasurements(db *sql.DB, measurements []Measurement) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO measurements (
			measured_at, weight, bmi, body_fat, fat_free_weight,
			subcutaneous_fat, visceral_fat, body_water, muscle_mass,
			skeletal_muscles, bone_mass, protein, bmr, metabolic_age
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	var inserted int64
	for _, m := range measurements {
		res, err := stmt.Exec(
			m.MeasuredAt, m.Weight, m.BMI, m.BodyFat, m.FatFreeWeight,
			m.SubcutaneousFat, m.VisceralFat, m.BodyWater, m.MuscleMass,
			m.SkeletalMuscles, m.BoneMass, m.Protein, m.BMR, m.MetabolicAge,
		)
		if err != nil {
			return 0, fmt.Errorf("exec: %w", err)
		}
		n, _ := res.RowsAffected()
		inserted += n
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return inserted, nil
}

func queryMeasurements(db *sql.DB, start, end string) ([]Measurement, error) {
	query := `SELECT measured_at, weight, bmi, body_fat, fat_free_weight,
		subcutaneous_fat, visceral_fat, body_water, muscle_mass,
		skeletal_muscles, bone_mass, protein, bmr, metabolic_age
		FROM measurements`

	var args []any
	if start != "" && end != "" {
		query += " WHERE measured_at >= ? AND measured_at <= ?"
		args = append(args, start, end+"T23:59:59")
	} else if start != "" {
		query += " WHERE measured_at >= ?"
		args = append(args, start)
	} else if end != "" {
		query += " WHERE measured_at <= ?"
		args = append(args, end+"T23:59:59")
	}

	query += " ORDER BY measured_at ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var results []Measurement
	for rows.Next() {
		var m Measurement
		err := rows.Scan(
			&m.MeasuredAt, &m.Weight, &m.BMI, &m.BodyFat, &m.FatFreeWeight,
			&m.SubcutaneousFat, &m.VisceralFat, &m.BodyWater, &m.MuscleMass,
			&m.SkeletalMuscles, &m.BoneMass, &m.Protein, &m.BMR, &m.MetabolicAge,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		results = append(results, m)
	}
	return results, rows.Err()
}
