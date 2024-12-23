package util

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"payment-gateway/internal/models/request"
	"strings"
)

func DecodeRequest(r *http.Request, request *request.Transaction) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(request)
	case "text/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	case "application/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	default:
		return fmt.Errorf("unsupported content type")
	}
}

func SendEncodedResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	// get type of content
	ct := w.Header().Get("Content-Type")

	// set default content type to JSON
	if ct == "" || strings.Contains(ct, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response in JSON", http.StatusInternalServerError)
		}
	} else if strings.Contains(ct, "application/soap+xml") {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(statusCode)

		if err := xml.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response in SOAP", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Content-Type is not supported.", http.StatusUnsupportedMediaType)
	}
}
