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
	if len(csvLines) != 7 {
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
	if csvLines[5] != "Total,12340.67," {
		t.Errorf("Expected total, got %s", csvLines[5])
	}
	if csvLines[6] != "" {
		t.Errorf("Expected empty line, got %s", csvLines[4])
	}
}

func TestService_UpdateBalanceSheet2(t *testing.T) {
	s := NewService()
	sheet := accounting.BalanceSheet{
		Entity: domain.Entity{
			Key:      "123",
			Created:  time.Now(),
			Modified: time.Now(),
		},
	}

	messages := `🗂 Here is your report:

📅 15.01.2024
00:00:00
+5.00 | Earnings
📅 14.01.2024
00:00:00
+3.99 | More earnings
-8.00 | Monthly service charge
📅 13.01.2024
00:00:00
📅 12.01.2024
00:00:00
+1,234.56 | admin_add
−1,234.56 | Investment
`

	err := s.UpdateBalanceSheet2(&sheet, messages)
	if err != nil {
		t.Errorf("Error updating BalanceSheet: %v", err)
	}

	// Export to CSV
	csvData, err := s.ExportToCSV(sheet)
	if err != nil {
		t.Errorf("Error exporting to CSV: %v", err)
	}
	// assert csv Data headers are correct
	csvLines := strings.Split(csvData, "\n")
	if len(csvLines) != 8 {
		t.Errorf("Expected 7 lines in CSV, got %d", len(csvLines))
	}
	if csvLines[0] != "Date,Balance Change,Notes" {
		t.Errorf("Wrong CSV header, got %s", csvLines[0])
	}
	if csvLines[1] != "2024-01-12,1234.56,admin_add" {
		t.Errorf("Wrong CSV entry, got %s", csvLines[1])
	}
	if !strings.Contains(csvLines[2], "-1234.56") {
		t.Errorf("Expected negative value, got %s", csvLines[2])
	}
	if !strings.Contains(csvLines[3], "3.99") {
		t.Errorf("Expected positive value, got %s", csvLines[3])
	}
	if !strings.Contains(csvLines[4], "-8.00") {
		t.Errorf("Expected negative value, got %s", csvLines[4])
	}
	if csvLines[6] != "Total,0.99," {
		t.Errorf("Expected total, got %s", csvLines[5])
	}
	if csvLines[7] != "" {
		t.Errorf("Expected empty line, got %s", csvLines[4])
	}
}
