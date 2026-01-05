package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"trip-planner/models"
	"trip-planner/utils"
)

type AIService struct {
	apiKey string
	model  string
}

func NewAIService(apiKey, model string) *AIService {
	return &AIService{
		apiKey: apiKey,
		model:  model,
	}
}

func (s *AIService) GenerateTrip(
	ctx context.Context,
	req models.TripRequest,
	places any,
	weather any,
) (map[string]any, error) {

	prompt := buildPrompt(req, places, weather)

	// ✅ CORRECT Responses API payload (2025 schema)
	payload := map[string]any{
		"model": s.model, // e.g. gpt-5.2
		"input": prompt,

		// ✅ JSON enforcement (THIS IS THE RIGHT WAY)
		"text": map[string]any{
			"format": map[string]any{
				"type": "json_object",
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.apiKey,
		"Content-Type":  "application/json",
	}

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var raw struct {
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := utils.PostJSON(
		ctx,
		"https://api.openai.com/v1/responses",
		payload,
		&raw,
		headers,
	); err != nil {
		return nil, err
	}

	// ✅ Collect ALL text parts
	var text string
	for _, out := range raw.Output {
		for _, c := range out.Content {
			if c.Text != "" {
				text += c.Text
			}
		}
	}

	if text == "" {
		return nil, fmt.Errorf("AI returned empty output")
	}

	// ✅ Parse guaranteed JSON
	var obj map[string]any
	if err := json.Unmarshal([]byte(text), &obj); err != nil {
		return nil, fmt.Errorf("AI returned invalid JSON: %w\nRAW:\n%s", err, text)
	}

	return obj, nil
}

func buildPrompt(req models.TripRequest, places any, weather any) string {
	bReq, _ := json.MarshalIndent(req, "", "  ")

	return fmt.Sprintf(`
You are a Sri Lanka trip planner. Return ONLY valid JSON (no markdown, no extra text).

JSON structure:
{
  "summary":"string",
  "route":["City1","City2"],
  "total_budget":{"currency":"LKR","low":0,"mid":0,"high":0,"notes":"..."},
  "tips":["..."],
  "warnings":["..."],
  "days":[
    {
      "day_number":1,
      "date":"YYYY-MM-DD",
      "base_city":"string",
      "theme":"string",
      "items":[
        {
          "time_block":"08:00-10:30",
          "title":"...",
          "description":"...",
          "location":"...",
          "travel_mode":"car",
          "travel_mins":30
        }
      ],
      "meals":[
        {"meal_type":"breakfast","suggestion":"...","area":"..."}
      ],
      "hotel_area":"string",
      "cost_range":{"currency":"LKR","low":0,"mid":0,"high":0,"notes":"..."}
    }
  ]
}

Rules:
- Realistic Sri Lanka travel times.
- Prefer sensible routes (Colombo → Kandy → Ella → Yala → Mirissa).
- Family-safe and practical.
- Currency must be LKR.
- Use places context to pick REAL attractions.

User request:
%s

PLACES (Google raw JSON):
%v

WEATHER (raw JSON):
%v
`, string(bReq), places, weather)
}
