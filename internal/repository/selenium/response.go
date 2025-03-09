package selenium

import "encoding/json"

// GeneralResponse defines the common structure returned by the Selenium API.
type GeneralResponse struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Status  string          `json:"status"`
}

// ParseGeneralResponse unmarshals the full response and then the dynamic data into out.
func ParseGeneralResponse(body []byte, out interface{}) (string, error) {
	var gr GeneralResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return "", err
	}
	if out != nil {
		if err := json.Unmarshal(gr.Data, out); err != nil {
			return gr.Message, err
		}
	}
	return gr.Message, nil
}
