package controller

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/lniklison/bill-splitter-go/internal/bills"
)

const maxRequestBodyBytes = 1 << 20

type HTTP struct {
	service *bills.Service
}

type replaceSharesRequest struct {
	Shares []shareInput `json:"shares"`
}

type shareInput struct {
	PersonName string `json:"personName"`
	Percentage string `json:"percentage"`
}

type allocationResponse struct {
	Bill   billResponse    `json:"bill"`
	Shares []shareResponse `json:"shares"`
}

type billResponse struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	TotalAmount string `json:"totalAmount"`
}

type shareResponse struct {
	PersonName string `json:"personName"`
	Percentage string `json:"percentage"`
	Amount     string `json:"amount"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHTTP(service *bills.Service) *HTTP {
	return &HTTP{service: service}
}

func (h *HTTP) GetShares(w http.ResponseWriter, r *http.Request) {
	billID, ok := parseBillID(w, r)
	if !ok {
		return
	}

	allocation, err := h.service.GetShares(r.Context(), billID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFor(allocation))
}

func (h *HTTP) ReplaceShares(w http.ResponseWriter, r *http.Request) {
	billID, ok := parseBillID(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request replaceSharesRequest
	if err := decoder.Decode(&request); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Request body must be valid JSON")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Request body must contain one JSON object")
		return
	}

	inputs := make([]bills.ShareInput, len(request.Shares))
	for index, share := range request.Shares {
		inputs[index] = bills.ShareInput{
			PersonName: share.PersonName,
			Percentage: share.Percentage,
		}
	}
	allocation, err := h.service.ReplaceShares(r.Context(), billID, inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFor(allocation))
}

func (h *HTTP) MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Allow", "GET, PUT")
	writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
}

func (h *HTTP) InvalidPath(w http.ResponseWriter, _ *http.Request) {
	writeAPIError(w, http.StatusBadRequest, "invalid_request", "Bill shares path is invalid")
}

func parseBillID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	billID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || billID <= 0 {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Bill ID must be a positive integer")
		return 0, false
	}
	return billID, true
}

func responseFor(allocation bills.Allocation) allocationResponse {
	response := allocationResponse{
		Bill: billResponse{
			ID:          allocation.BillID,
			Description: allocation.Description,
			TotalAmount: bills.FormatCents(allocation.TotalAmount),
		},
		Shares: make([]shareResponse, len(allocation.Shares)),
	}
	for index, share := range allocation.Shares {
		response.Shares[index] = shareResponse{
			PersonName: share.PersonName,
			Percentage: share.Percentage.String(),
			Amount:     bills.FormatCents(share.Amount),
		}
	}
	return response
}

func writeServiceError(w http.ResponseWriter, err error) {
	var domainError *bills.Error
	if errors.As(err, &domainError) {
		status := http.StatusUnprocessableEntity
		if domainError.Code == bills.ErrBillNotFound.Code {
			status = http.StatusNotFound
		}
		writeAPIError(w, status, domainError.Code, domainError.Message)
		return
	}
	log.Printf("handle bill shares request: %v", err)
	writeAPIError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
