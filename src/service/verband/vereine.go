package verband

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pkv/api/src/repository/dpv"
	"sort"
	"strings"
)

type answer struct {
	Id         int    `json:"id"`
	QuestionId int    `json:"questionId"`
	Text       string `json:"text"`
}

type answerList []answer

func (al answerList) findByQuestionId(questionId int) answer {
	for _, e := range al {
		if e.QuestionId == questionId {
			return e
		}
	}
	return answer{}
}

type meta struct {
	Status     string `json:"status"`
	Statuscode int    `json:"statuscode"`
	Message    string `json:"message"`
}

type submission struct {
	Id      int `json:"id"`
	FormId  int `json:"formId"`
	Answers answerList
}

type question struct {
	Id     int    `json:"id"`
	FormId int    `json:"formId"`
	Text   string `json:"text"`
}

type data struct {
	Submissions []submission `json:"submissions"`
	Questions   []question   `json:"questions"`
}

type ocs struct {
	Meta meta `json:"meta"`
	Data data `json:"data"`
}

type nextcloudResponse struct {
	OCS ocs `json:"ocs"`
}

type Verein struct {
	Bundesland string `json:"bundesland"`
	Stadt      string `json:"stadt"`
	Name       string `json:"name"`
	Webseite   string `json:"webseite"`
}

func (v Verein) SortKey() string {
	return v.Bundesland + "|" + v.Stadt + "|" + v.Name
}

func (s *Service) GetVereine(ctx context.Context) ([]Verein, error) {
	url := dpv.ConfigInstance.Nextcloud.URL + "ocs/v2.php/apps/forms/api/v2.4/submissions/" + dpv.ConfigInstance.Nextcloud.FormID
	user := dpv.ConfigInstance.Nextcloud.User
	pass := dpv.ConfigInstance.Nextcloud.Pass

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.SetBasicAuth(user, pass)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get vereine: %w", err)
	}

	var response nextcloudResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	if response.OCS.Meta.Status != "ok" {
		return nil, fmt.Errorf("could not get vereine: %s", response.OCS.Meta.Message)
	}

	vereine := s.ExtractVereineList(response)
	s.sortVereine(vereine)

	return vereine, nil
}

func normalizeURL(inputURL string) (string, error) {
	inputURL = strings.TrimSpace(inputURL)
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	if u.Host == "" && u.Scheme == "" {
		return normalizeURL("https://" + inputURL)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String(), nil
}

func (s *Service) ExtractVereineList(response nextcloudResponse) []Verein {
	ocsData := response.OCS.Data

	var vereine []Verein

	for _, answer := range ocsData.Submissions {
		if strings.Contains(answer.Answers.findByQuestionId(16).Text, "Ja") {
			normalizedURL, _ := normalizeURL(answer.Answers.findByQuestionId(6).Text)
			vereine = append(vereine, Verein{
				Bundesland: strings.TrimSpace(answer.Answers.findByQuestionId(17).Text),
				Stadt:      strings.TrimSpace(answer.Answers.findByQuestionId(12).Text),
				Name:       strings.TrimSpace(answer.Answers.findByQuestionId(13).Text),
				Webseite:   normalizedURL,
			})
		}
	}
	return vereine
}

func (s *Service) sortVereine(vereine []Verein) {
	sort.Slice(vereine, func(i, j int) bool {
		return vereine[i].SortKey() < vereine[j].SortKey()
	})
}
