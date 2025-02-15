package main

import (
	"bytes"
	"encoding/json"
	"github.com/VMerlin/receipt-processor/src/processor"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	service := processor.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		processor.Handle(service, w, r)
	}))
	defer server.Close()

	t.Run("TestProcessReceipt", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"retailer":     "M&M Corner Market",
			"purchaseDate": "2022-03-20",
			"purchaseTime": "14:33",
			"items": []map[string]string{
				{"shortDescription": "Gatorade", "price": "2.25"},
				{"shortDescription": "Gatorade", "price": "2.25"},
				{"shortDescription": "Gatorade", "price": "2.25"},
				{"shortDescription": "Gatorade", "price": "2.25"},
			},
			"total": "9.00",
		}
		body, _ := json.Marshal(requestBody)
		resp, err := http.Post(server.URL+"/receipts/process", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.StatusCode)
		}
	})

	t.Run("TestGetPoints", func(t *testing.T) {
		receiptID := "12345"
		resp, err := http.Get(server.URL + "/receipts/" + receiptID + "/points")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.StatusCode)
		}
	})
}
