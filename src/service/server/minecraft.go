package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type whitelistEntry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (s *Service) AddUsernameToWhitelist(username string, ctx context.Context) error {
	matched, err := regexp.MatchString(`[A-Za-z0-9_]{3,16}`, username)
	if err != nil {
		return fmt.Errorf("could not validate minecraft username: %w", err)
	}
	if !matched {
		return errors.New("the provided username is not valid in minecraft")
	}

	cmd := exec.Command("journalctl", "-u", "minecraft-server.service", "--since", "10min ago")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("could not obtain minecraft server logs: %w", err)
	}

	pattern := fmt.Sprintf(`UUID of player %s is ([a-f0-9\-]{36})`, regexp.QuoteMeta(username))
	uuidRegex := regexp.MustCompile(pattern)

	submatches := uuidRegex.FindStringSubmatch(string(output))
	if submatches == nil || len(submatches) < 2 {
		return errors.New("make sure the user has tried to connect within the last 10 minutes")
	}
	uuid := submatches[1]

	const whitelistFile = "/var/lib/minecraft/server/whitelist.json"
	data, err := os.ReadFile(whitelistFile)
	if err != nil {
		return fmt.Errorf("could not open minecraft server whitelist: %w", err)
	}

	var whitelist []whitelistEntry
	err = json.Unmarshal(data, &whitelist)
	if err != nil {
		return fmt.Errorf("could not read minecraft server whitelist: %w", err)
	}

	for _, entry := range whitelist {
		if entry.Name == username {
			return fmt.Errorf("user is already whitelisted: %w", err)
		}
	}

	newEntry := whitelistEntry{
		UUID: uuid,
		Name: username,
	}

	whitelist = append(whitelist, newEntry)
	updatedData, err := json.Marshal(whitelist)
	if err != nil {
		return fmt.Errorf("could not prepare updated minecraft server whitelist: %w", err)
	}

	err = os.WriteFile(whitelistFile, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("could not write minecraft server whitelist: %w", err)
	}

	return nil
}
