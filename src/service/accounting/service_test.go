package accounting

import (
	"pkv/api/src/domain"
	"pkv/api/src/domain/accounting"
	"strings"
	"testing"
	"time"
)

func TestService_UpdateBalanceSheet(t *testing.T) {
	s := NewService()
	sheet := accounting.BalanceSheet{
		Entity: domain.Entity{
			Key:      "123",
			Created:  time.Now(),
			Modified: time.Now(),
		},
	}

	// Example messages
	messages := []string{
		"13.01.2024 📩 12345.67 EUR - Create account",
		"14.01.2024 📤 -5.00 EUR - Server Fee",
		"15.01.2024 ♻️ 5.00 EUR - Maintenance Fee",
		"17.01.2024 ♻️ -5.00 EUR - Refunded maintenance fee",
	}

	// Update BalanceSheet with messages
	for _, msg := range messages {
		err := s.UpdateBalanceSheet(&sheet, msg)
		if err != nil {
			t.Errorf("Error updating BalanceSheet: %v", err)
		}
	}

	// Export to CSV
	csvData, err := s.ExportToCSV(sheet)
	if err != nil {
		t.Errorf("Error exporting to CSV: %v", err)
	}
	// assert csv Data headers are correct
	csvLines := strings.Split(csvData, "\n")
	if len(csvLines) != 6 {
		t.Errorf("Expected 6 lines in CSV, got %d", len(csvLines))
	}
	if csvLines[0] != "Date,Balance Change,Notes" {
		t.Errorf("Wrong CSV header, got %s", csvLines[0])
	}
	if csvLines[1] != "2024-01-13,12345.67,Create account" {
		t.Errorf("Wrong CSV entry, got %s", csvLines[1])
	}
	if !strings.Contains(csvLines[2], "-5.00") {
		t.Errorf("Expected negative value, got %s", csvLines[2])
	}
	if !strings.Contains(csvLines[3], "-5.00") {
		t.Errorf("Expected negative value, got %s", csvLines[3])
	}
	if !strings.Contains(csvLines[4], "5.00") {
		t.Errorf("Expected positive value, got %s", csvLines[4])
	}
	if csvLines[5] != "" {
		t.Errorf("Expected empty line, got %s", csvLines[4])
	}
}
