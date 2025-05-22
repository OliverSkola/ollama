package wandb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type WandbClient struct {
	apiKey   string
	entity   string
	project  string
	runID    string
	baseURL  string
	client   *http.Client
}

func NewWandbClient() *WandbClient {
	return &WandbClient{
		apiKey:   os.Getenv("O_WANDB_API_KEY"),
		entity:   os.Getenv("O_WANDB_ENTITY"),
		project:  os.Getenv("O_WANDB_PROJECT"),
		baseURL:  "https://api.wandb.ai",
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WandbClient) StartRun(runName string) error {
	url := fmt.Sprintf("%s/v1/%s/%s/runs", w.baseURL, w.entity, w.project)
	body := map[string]interface{}{
		"name": runName,
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+w.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if id, ok := result["id"].(string); ok {
		w.runID = id
	} else {
		return fmt.Errorf("could not parse run ID")
	}
	return nil
}

func (w *WandbClient) LogMetrics(step int, metrics map[string]float64) error {
	url := fmt.Sprintf("%s/v1/%s/%s/runs/%s/metrics", w.baseURL, w.entity, w.project, w.runID)

	body := map[string]interface{}{
		"step":    step,
		"metrics": metrics,
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+w.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
