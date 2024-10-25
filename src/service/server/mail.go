package server

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"os"
	"os/exec"
	"pkv/api/src/repository/security"
	"pkv/api/src/repository/t"
	"strings"
)

func (s *Service) ChangeMailPassword(email string, oldpassword string, newpassword string, ctx context.Context) error {
	if email == "" {
		return t.Errorf("email is empty")
	}
	emailSafe := strings.ReplaceAll(email, "@", ".")
	illegalCharacters := "#%&{}\\<>*/?$!'\":;,+`|=[] "
	if emailSafe[0] == '.' || emailSafe[len(emailSafe)-1] == '.' {
		return t.Errorf("email is invalid")
	}
	for _, char := range illegalCharacters {
		if strings.ContainsRune(emailSafe, char) {
			return t.Errorf("email is not supported")
		}
	}

	filename := emailSafe + ".key"
	filepath := "/var/config/mail/" + filename

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return t.Errorf("email does not exist")
		}
		return t.Errorf("could not obtain the password")
	}

	err = bcrypt.CompareHashAndPassword(data, []byte(oldpassword))
	if err != nil {
		return t.Errorf("the old password is incorrect")
	}

	// Check password length
	if len(newpassword) < 8 {
		return t.Errorf("password too short (minimum 8 chars)")
	}

	// Check password strength
	if !security.IsStrongPassword(newpassword) {
		return t.Errorf("password too weak (contains only numbers)")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), 6)
	if err != nil {
		return t.Errorf("cannot create the new password")
	}

	err = os.WriteFile(filepath, newHashedPassword, 0644)
	if err != nil {
		return t.Errorf("cannot save the new password")
	}

	services := []string{"dovecot2.service", "postfix-setup.service", "postfix.service"}
	for _, service := range services {
		if err := restartService(service); err != nil {
			return t.Errorf("the password has been changed successfully, but the mail server could not be restarted - you may still have to use the old password, or you can try restarting it again by typing in your new password in all three password fields: %w", err)
		}
	}

	return nil
}

func restartService(serviceName string) error {
	cmd := exec.Command("/run/wrappers/bin/doas", "-u", "root", "/run/current-system/sw/bin/systemctl", "restart", serviceName)
	err := cmd.Run()
	if err != nil {
		return t.Errorf("failed to restart %s: %w", serviceName, err)
	}
	return nil
}
