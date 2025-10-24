package handlers

import (
	"encoding/json"
	"net/http"
	"resco/services"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
}

// GetBOMByItemCode handles GET requests for BOM data by item code
func GetBOMByItemCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to get BOM data
	results, err := services.GetBOMByCodeParameterized(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    results,
		Count:   len(results),
		Message: "BOM data retrieved successfully",
	})
}

// GetBOMByItemCodeCN handles GET requests for BOM data with Chinese translations
func GetBOMByItemCodeCN(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to get BOM data with Chinese translations
	results, err := services.GetBOMByCodeWithTranslation(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    results,
		Count:   len(results),
		Message: "BOM data with Chinese translations retrieved successfully",
	})
}

// GetBOMByItemCodeCombined handles GET requests for BOM data with both Turkish and Chinese
func GetBOMByItemCodeCombined(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to get BOM data with both Turkish and Chinese
	results, err := services.GetBOMByCodeCombined(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    results,
		Count:   len(results),
		Message: "BOM data with Turkish and Chinese retrieved successfully",
	})
}

// GetBOMTotal handles GET requests for unique BOM codes with sequential numbers
func GetBOMTotal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to get unique codes with sequential numbers
	results, err := services.GetBOMTotal(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    results,
		Count:   len(results),
		Message: "Unique BOM codes retrieved successfully",
	})
}

// QueryHeihu handles GET requests to query the external Heihu API
func QueryHeihu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to query Heihu API
	result, err := services.QueryHeihu(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return the result from Heihu API
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// CheckProduct handles GET requests to check all BOM products against Heihu API
func CheckProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get item code from URL parameters
	vars := mux.Vars(r)
	itemCode := vars["itemCode"]

	if itemCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item code is required"})
		return
	}

	// Call the service to check products
	results, err := services.CheckProducts(itemCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    results,
		Count:   len(results),
		Message: "Product check completed successfully",
	})
}

// HealthCheck handles health check requests
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"message": "API is running",
	})
}