package accounting

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"pkv/api/src/domain/accounting"
	"pkv/api/src/repository/t"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) UpdateBalanceSheet(sheet *accounting.BalanceSheet, message string) error {
	re := regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4}) (\D*) (-?[\d,]+(?:\.\d{2})?) \w+ - (.+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) != 5 {
		return t.Errorf("message format is incorrect")
	}
	date, err := time.Parse("02.01.2006", matches[1])
	if err != nil {
		return err
	}
	balanceChange, err := strconv.ParseFloat(strings.Replace(matches[3], ",", "", -1), 64)
	if err != nil {
		return err
	}
	if matches[2] == "♻️" {
		balanceChange = -balanceChange
	}
	newEntry := accounting.Entry{
		Date:          date,
		BalanceChange: balanceChange,
		Notes:         matches[4],
	}
	for _, entry := range sheet.Entries {
		if entry.Date.Equal(newEntry.Date) && entry.BalanceChange == newEntry.BalanceChange && entry.Notes == newEntry.Notes {
			return nil // Same entry exists
		}
	}
	sheet.Entries = append(sheet.Entries, newEntry)
	sheet.Modified = time.Now()
	return nil
}

func (s *Service) UpdateBalanceSheet2(sheet *accounting.BalanceSheet, message string) error {
	scanner := bufio.NewScanner(strings.NewReader(message))

	dateRegex := regexp.MustCompile(`📅 (\d{2}\.\d{2}\.\d{4})`)
	transactionRegex := regexp.MustCompile(`([+\-−]\d{1,3}(?:,\d{3})*(?:\.\d{2})?) \| (.+)`)

	var currentDate time.Time

	for scanner.Scan() {
		line := scanner.Text()

		if dateMatch := dateRegex.FindStringSubmatch(line); dateMatch != nil {
			var err error
			currentDate, err = time.Parse("02.01.2006", dateMatch[1])
			if err != nil {
				return err
			}
			continue
		}

		if transactionMatch := transactionRegex.FindStringSubmatch(line); transactionMatch != nil {
			balanceChangeStr := strings.Replace(transactionMatch[1], ",", "", -1)
			balanceChangeStr = strings.Replace(balanceChangeStr, "+", "", 1)
			balanceChangeStr = strings.Replace(balanceChangeStr, "−", "-", 1)
			balanceChange, err := strconv.ParseFloat(balanceChangeStr, 64)
			if err != nil {
				return err
			}
			notes := transactionMatch[2]

			newEntry := accounting.Entry{
				Date:          currentDate,
				BalanceChange: balanceChange,
				Notes:         notes,
			}

			entryExists := false
			for _, entry := range sheet.Entries {
				if entry.Date.Equal(newEntry.Date) && entry.BalanceChange == newEntry.BalanceChange && entry.Notes == newEntry.Notes {
					entryExists = true
					break
				}
			}
			if !entryExists {
				sheet.Entries = append(sheet.Entries, newEntry)
			}
		}
	}

	return scanner.Err()
}

func (s *Service) ExportToCSV(sheet accounting.BalanceSheet) (string, error) {
	sort.Slice(sheet.Entries, func(i, j int) bool {
		return sheet.Entries[i].Date.Before(sheet.Entries[j].Date)
	})
	var csvData strings.Builder
	writer := csv.NewWriter(&csvData)
	if err := writer.Write([]string{"Date", "Balance Change", "Notes"}); err != nil {
		return "", err
	}
	var val float64
	for _, entry := range sheet.Entries {
		row := []string{
			entry.Date.Format(time.DateOnly),
			fmt.Sprintf("%.2f", entry.BalanceChange),
			entry.Notes,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
		val += entry.BalanceChange
	}
	if err := writer.Write([]string{"Total", fmt.Sprintf("%.2f", val), ""}); err != nil {
		return "", err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return csvData.String(), nil
}

func (s *Service) LoadFromJson(filename string) (accounting.BalanceSheet, error) {
	var sheet accounting.BalanceSheet
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return sheet, err
	}
	err = json.Unmarshal(jsonFile, &sheet)
	if err != nil {
		return sheet, err
	}
	return sheet, nil
}

func (s *Service) SaveToJson(sheet accounting.BalanceSheet, filename string) error {
	jsonFile, err := json.Marshal(sheet)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonFile, 0600)
	if err != nil {
		return err
	}
	return nil
}
