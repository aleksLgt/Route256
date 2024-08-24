package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"route256/cart/internal/domain"
)

type (
	AddItemRequest struct {
		Count int `json:"count"`
	}

	CartCheckoutRequest struct {
		User int64 `json:"user"`
	}

	GetCartItemsResponse struct {
		Items      []domain.ListItem `json:"items"`
		TotalPrice int               `json:"total_price"`
	}

	CartCheckoutResponse struct {
		OrderID int `json:"orderID"`
	}
)

func TestCartItemDelete(t *testing.T) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8082/user/31337/cart/1076963", http.NoBody)

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Invalid response status: %d", resp.StatusCode)
	}
}

func TestCartItemsClear(t *testing.T) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8082/user/31337/cart/", http.NoBody)

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Invalid response status: %d", resp.StatusCode)
	}
}

func TestCartList(t *testing.T) {
	client := http.Client{}

	addItemRequest := AddItemRequest{
		Count: 1,
	}
	data, err := json.Marshal(addItemRequest)

	if err != nil {
		t.Fatalf("failed to encode request %v", err)
	}

	addReq, err := http.NewRequest(http.MethodPost, "http://localhost:8082/user/31337/cart/1076963", bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	addResp, err := client.Do(addReq)
	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer addResp.Body.Close()

	getReq, err := http.NewRequest(http.MethodGet, "http://localhost:8082/cart/31337/list/", http.NoBody)
	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	getResp, err := client.Do(getReq)

	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer getResp.Body.Close()

	response := &GetCartItemsResponse{}
	err = json.NewDecoder(getResp.Body).Decode(response)

	if err != nil {
		t.Fatalf("failed to decode error response")
	}

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status: %d", getResp.StatusCode)
	}

	if len(response.Items) != 1 {
		t.Fatalf("The number of products must be equal to 1")
	}
}

func TestAddItem(t *testing.T) {
	client := http.Client{}

	addItemRequest := AddItemRequest{
		Count: 1,
	}
	data, err := json.Marshal(addItemRequest)

	if err != nil {
		t.Fatalf("failed to encode request %v", err)
	}

	addReq, err := http.NewRequest(http.MethodPost, "http://localhost:8082/user/31337/cart/1076963", bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	addResp, err := client.Do(addReq)
	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer addResp.Body.Close()

	if addResp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status: %d", addResp.StatusCode)
	}
}

func TestCheckoutCart(t *testing.T) {
	client := http.Client{}

	addItemRequest := AddItemRequest{
		Count: 1,
	}
	data, err := json.Marshal(addItemRequest)

	if err != nil {
		t.Fatalf("failed to encode request %v", err)
	}

	addReq, err := http.NewRequest(http.MethodPost, "http://localhost:8082/user/31337/cart/1076963", bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	addResp, err := client.Do(addReq)
	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer addResp.Body.Close()

	checkoutCartRequest := CartCheckoutRequest{
		User: 31337,
	}
	data, err = json.Marshal(checkoutCartRequest)

	if err != nil {
		t.Fatalf("failed to encode request %v", err)
	}

	checkoutReq, err := http.NewRequest(http.MethodPost, "http://localhost:8082/cart/checkout", bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("Error creating the request: %v", err)
	}

	checkoutResp, err := client.Do(checkoutReq)
	if err != nil {
		t.Fatalf("Error when executing the request: %v", err)
	}

	defer checkoutResp.Body.Close()

	response := &CartCheckoutResponse{}
	err = json.NewDecoder(checkoutResp.Body).Decode(response)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if checkoutResp.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response status: %d", checkoutResp.StatusCode)
	}

	if response.OrderID != 1 {
		t.Fatalf("The orderID must be equal to 1")
	}
}
