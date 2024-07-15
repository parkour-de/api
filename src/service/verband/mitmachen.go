package verband

import (
	"bytes"
	"fmt"
	"html"
	"mime/quotedprintable"
	"net/textproto"
	"os/exec"
	"pkv/api/src/domain/verband"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxLengthTextArea = 100000
	maxLengthInput    = 100
	smtpServer        = "localhost"
	smtpPort          = "2525"
)

var validAGs = map[string]string{
	"bildung":         "Bildung, Forschung und Wissenschaft",
	"design":          "Logo & Corporate Design",
	"finanzen":        "Finanzen",
	"it":              "IT",
	"lizenzen":        "Lizenzen und Ausbildung",
	"oeffentlichkeit": "Öffentlichkeitsarbeit",
	"satzung":         "Satzung",
	"wettkampf":       "Wettkampf",
	"bjoern":          "Björn's VIP-Club",
}

func trimAndSanitize(input string, maxLength int) (string, error) {
	input = strings.TrimSpace(input)
	if utf8.RuneCountInString(input) > maxLength {
		return "", fmt.Errorf("maximum field length exceeded - maximum length is %d chars, %d given", maxLength, utf8.RuneCountInString(input))
	}
	return html.EscapeString(input), nil
}

func indent(content string) string {
	content = strings.ReplaceAll(content, "\n", "<br>")
	return fmt.Sprintf(`<p style="border-left:.3em solid #888; padding-left: .3em; margin-left: .3em;">%s</p>`, content)
}

func (s *Service) Mitmachen(data verband.MitmachenRequest) error {
	var err error
	if data.Name, err = trimAndSanitize(data.Name, maxLengthInput); err != nil {
		return fmt.Errorf("invalid name - %w", err)
	}
	if data.Email, err = trimAndSanitize(data.Email, maxLengthInput); err != nil {
		return fmt.Errorf("invalid email - %w", err)
	}
	if data.Kompetenzen, err = trimAndSanitize(data.Kompetenzen, maxLengthTextArea); err != nil {
		return fmt.Errorf("invalid kompetenzen - %w", err)
	}
	if data.Fragen, err = trimAndSanitize(data.Fragen, maxLengthTextArea); err != nil {
		return fmt.Errorf("invalid fragen - %w", err)
	}

	if data.Name == "" {
		data.Name = "ein Interessent"
	}
	if data.Email == "" {
		data.Email = "noreply@8bj.de"
	}

	prettyAG, valid := validAGs[data.AG]
	if !valid {
		return fmt.Errorf("invalid AG provided")
	}

	subject := fmt.Sprintf("[ANFRAGE] %s möchte bei %s mitmachen", data.Name, prettyAG)
	to := fmt.Sprintf("%s@parkour-deutschland.de", data.AG)
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<style>
    body { font-family: Arial, sans-serif; }
</style>
</head>
<body>
<p>Liebes %s Kommittee,</p>
<p>am %s wurde folgende Anfrage an euch gestellt:</p>
<p>Ich bin <b>%s</b> und möchte bei euch mitmachen.</p>
<p>Ich bringe folgende Kompetenzen, Erfahrungen und Referenzen mit:</p>
%s
<p>Ich möchte euch außerdem sagen:</p>
%s
<p>Ich freue mich auf eure Antwort, oder schreibt mir unter %s.</p>
<p>Bis bald!</p>
<p>%s</p>
</body>
</html>`,
		prettyAG, time.Now().Format("Monday, 2 January 2006 at 15:04"), data.Name,
		indent(data.Kompetenzen), indent(data.Fragen), data.Email, data.Name)

	if err := s.SendMail(data.Email, to, subject, body); err != nil {
		return err
	}
	return nil
}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return fmt.Errorf("smtp: A line must not contain CR or LF")
	}
	return nil
}
func (s *Service) SendMail(from, to, subject, body string) error {
	if err := validateLine(from); err != nil {
		return err
	}
	if err := validateLine(to); err != nil {
		return err
	}
	if err := validateLine(subject); err != nil {
		return err
	}
	header := textproto.MIMEHeader{}
	header.Set("From", "noreply@8bj.de")
	header.Set("To", to)
	header.Set("Reply-To", from)
	header.Set("Subject", subject)
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Content-Transfer-Encoding", "quoted-printable")

	var buf bytes.Buffer
	writer := quotedprintable.NewWriter(&buf)
	writer.Write([]byte(body))
	writer.Close()

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v[0])
	}
	msg += "\r\n" + buf.String()

	println(msg)

	cmd := exec.Command("sendmail", "-F", "Parkour Deutschland", "-f", "noreply@8bj.de", "-i", "--", to)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(msg))
	}()

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
