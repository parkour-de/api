package description

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/t"
)

// TranslateDocument translates a document from one language to another.
func TranslateDocument(text, srcLang, destLang string, ctx context.Context) (string, error) {
	deeplKey := dpv.ConfigInstance.Auth.DeepLKey
	deeplUrl, err := url.Parse(dpv.ConfigInstance.Auth.DeepLUrl)
	if err != nil {
		return "", t.Errorf("invalid DeepL url: %w", err)
	}
	q := deeplUrl.Query()
	q.Add("auth_key", deeplKey)
	q.Add("text", text)
	q.Add("source_lang", srcLang)
	q.Add("target_lang", destLang)
	deeplUrl.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodPost, deeplUrl.String(), nil)
	if err != nil {
		return "", t.Errorf("creating DeepL request failed: %w", err)
	}
	req.Header.Add("User-Agent", "dpv-api")
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", t.Errorf("DeepL request failed: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		type ErrorResponse struct {
			ErrMessage string `json:"message"`
		}
		outStruct := &ErrorResponse{}
		if err := json.Unmarshal(bodyBytes, outStruct); err != nil {
			return "", t.Errorf("DeepL request failed with status %v, decoding error JSON failed: %w, error message: %v", resp.StatusCode, err, string(bodyBytes))
		}
		return "", t.Errorf("DeepL request failed with status %v: %v", resp.StatusCode, outStruct.ErrMessage)
	}
	type translation struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	}
	type TranslateResponse struct {
		Translations []translation `json:"translations"`
	}
	outStruct := &TranslateResponse{}
	if err := json.Unmarshal(bodyBytes, outStruct); err != nil {
		return "", t.Errorf("DeepL request failed with status %v, decoding response JSON failed: %w, response: %v", resp.StatusCode, err, string(bodyBytes))
	}
	return outStruct.Translations[0].Text, nil
}
