package verband

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"net/textproto"
	"pkv/api/src/domain/verband"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxLengthTextArea = 100000
	maxLengthInput    = 100
	smtpServer        = "localhost"
	smtpPort          = "25"
)

var validAGs = map[string]string{
	"bildung":         "Bildung, Forschung und Wissenschaft",
	"bjoern":          "Björn's VIP-Club",
	"design":          "Logo & Corporate Design",
	"finanzen":        "Finanzen",
	"it":              "IT",
	"lizenzen":        "Lizenzen und Ausbildung",
	"oeffentlichkeit": "Öffentlichkeitsarbeit",
	"parkourparks":    "Parkour-Parks",
	"satzung":         "Satzung",
	"wettkampf":       "Wettkampf",
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
	if match, err := regexp.Match("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", []byte(data.Email)); !match || err != nil {
		return fmt.Errorf("invalid email - email must pass this spec: https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address - %w", err)
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
<p>Ich freue mich auf eure Antwort, oder schreibt mir unter <a href="mailto:%s">%s</a>.</p>
<p>Bis bald!</p>
<p>%s</p>
</body>
</html>`,
		prettyAG, time.Now().Format("Monday, 2 January 2006 at 15:04"), data.Name,
		indent(data.Kompetenzen), indent(data.Fragen), data.Email, data.Email, data.Name)

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

func encodeBase64Header(header string) string {
	return "=?utf-8?B?" + base64.StdEncoding.EncodeToString([]byte(header)) + "?="
}

func (s *Service) SendMail(from, to, subject, body string) error {
	if err := validateLine(from); err != nil {
		return fmt.Errorf("invalid from: %w", err)
	}
	if err := validateLine(to); err != nil {
		return fmt.Errorf("invalid to: %w", err)
	}
	if err := validateLine(subject); err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}
	header := textproto.MIMEHeader{}
	header.Set("MIME-Version:", "1.0")
	header.Set("Date", time.Now().Format(time.RFC1123Z))
	header.Set("From", "noreply@8bj.de")
	header.Set("To", to)
	header.Set("Reply-To", from)
	header.Set("Subject", encodeBase64Header(subject))
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Content-Transfer-Encoding", "quoted-printable")

	var buf bytes.Buffer
	writer := quotedprintable.NewWriter(&buf)
	_, err := writer.Write([]byte(body))
	if err != nil {
		return err
	}
	writer.Close()

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v[0])
	}
	msg += "\r\n" + buf.String()

	println(msg)

	conn, err := net.Dial("tcp", smtpServer+":"+smtpPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return err
	}
	defer c.Quit()

	if err := c.Mail("noreply@8bj.de"); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
