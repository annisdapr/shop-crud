package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ItemResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	Stock int       `json:"stock"`
}

type ItemClient interface {
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*ItemResponse, error)
}

type itemClient struct {
	baseURL string
	client  *http.Client
}

func NewItemClient(baseURL string) ItemClient {
	return &itemClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport), 
		},
	}
}

func (c *itemClient) GetItemByID(ctx context.Context, itemID uuid.UUID) (*ItemResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/items/%s", c.baseURL, itemID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("item not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get item: %s", resp.Status)
	}

	var item ItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}
