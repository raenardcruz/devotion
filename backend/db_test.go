package main

import (
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestParseConnDetails(t *testing.T) {
	t.Run("URL format", func(t *testing.T) {
		connStr := "postgres://postgres:postgres@localhost:5432/devotion?sslmode=disable"
		targetDB, maintenanceConnStr, err := parseConnDetails(connStr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if targetDB != "devotion" {
			t.Errorf("expected targetDB 'devotion', got %q", targetDB)
		}
		expectedMaint := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		if maintenanceConnStr != expectedMaint {
			t.Errorf("expected maintenanceConnStr %q, got %q", expectedMaint, maintenanceConnStr)
		}
	})

	t.Run("Key-value format", func(t *testing.T) {
		connStr := "host=localhost port=5432 user=postgres password=postgres dbname=devotion sslmode=disable"
		targetDB, maintenanceConnStr, err := parseConnDetails(connStr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if targetDB != "devotion" {
			t.Errorf("expected targetDB 'devotion', got %q", targetDB)
		}
		expectedMaint := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
		if maintenanceConnStr != expectedMaint {
			t.Errorf("expected maintenanceConnStr %q, got %q", expectedMaint, maintenanceConnStr)
		}
	})
}

func TestIsDBNotExistError(t *testing.T) {
	if isDBNotExistError(nil) {
		t.Errorf("expected false for nil error")
	}

	pqErr := &pq.Error{Code: "3D000", Message: `database "devotion" does not exist`}
	if !isDBNotExistError(pqErr) {
		t.Errorf("expected true for pq.Error 3D000")
	}

	genericErr := fmt.Errorf("pq: database \"devotion\" does not exist")
	if !isDBNotExistError(genericErr) {
		t.Errorf("expected true for string matching does not exist")
	}
}
