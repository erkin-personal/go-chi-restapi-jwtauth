package handlers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
)

type CustomHandler struct {
	Certificates map[string]string
	PrivateKey   *rsa.PrivateKey
}

func NewCustomHandler() *CustomHandler {
	// Generate a private key for demonstration purposes.
	// In a real-world application, use a securely-stored private key.
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return &CustomHandler{
		Certificates: map[string]string{
			"cert1": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
			"cert2": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		},
		PrivateKey: privateKey,
	}
}

func (ch *CustomHandler) PreFlight(w http.ResponseWriter, r *http.Request) {
	// Allow CORS from any origin (not recommended for production)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Set the allowed HTTP methods for CORS requests
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Set the allowed headers for CORS requests
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Set the max age for the preflight request to be cached by the browser
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

	// Respond with an empty 204 No Content status to acknowledge the preflight request
	w.WriteHeader(http.StatusNoContent)
}

func (ch *CustomHandler) TestGet(w http.ResponseWriter, r *http.Request) {
	// Define the response data structure
	type responseData struct {
		Message string `json:"message"`
	}

	// Create an instance of responseData with a custom message
	response := responseData{
		Message: "Hello from TestGet!",
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the response with a 200 OK status and the JSON-encoded responseData
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ch *CustomHandler) TestPost(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := os.ReadFile("path/to/your/file")
	if err != nil {
		// handle error
	}

	// Check if the request body is valid JSON
	var requestBody interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the response with a 200 OK status and the JSON-encoded request body
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ch *CustomHandler) GetCertificate(w http.ResponseWriter, r *http.Request) {
	// Define the request data structure
	type requestData struct {
		ID string `json:"id"`
	}

	// Decode the request JSON payload
	var request requestData
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Find the certificate by ID
	cert, ok := ch.Certificates[request.ID]
	if !ok {
		http.Error(w, "Certificate not found", http.StatusNotFound)
		return
	}

	// Define the response data structure
	type responseData struct {
		Certificate string `json:"certificate"`
	}

	// Create an instance of responseData with the fetched certificate
	response := responseData{
		Certificate: cert,
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the response with a 200 OK status and the JSON-encoded responseData
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ch *CustomHandler) CheckSignature(w http.ResponseWriter, r *http.Request) {
	// Define the request data structure
	type requestData struct {
		Certificate string `json:"certificate"`
		Signature   string `json:"signature"`
		Data        string `json:"data"`
	}

	// Decode the request JSON payload
	var request requestData
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse the certificate and extract the public key
	pubKey, err := parseCertificate(request.Certificate)
	if err != nil {
		http.Error(w, "Invalid certificate", http.StatusBadRequest)
		return
	}

	// Decode the base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(request.Signature)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	// Compute the SHA-256 hash of the data
	hash := sha256.Sum256([]byte(request.Data))

	// Verify the signature
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Define the response data structure
	type responseData struct {
		Verified bool `json:"verified"`
	}

	// Create an instance of responseData with the verification result
	response := responseData{
		Verified: true,
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the response with a 200 OK status and the JSON-encoded responseData
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ch *CustomHandler) CreateSignature(w http.ResponseWriter, r *http.Request) {
	// Define the request data structure
	type requestData struct {
		Data string `json:"data"`
	}

	// Decode the request JSON payload
	var request requestData
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Compute the SHA-256 hash of the data
	hash := sha256.Sum256([]byte(request.Data))

	// Sign the hash using the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, ch.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the signature as a base64 string
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Define the response data structure
	type responseData struct {
		Signature string `json:"signature"`
	}

	// Create an instance of responseData with the generated signature
	response := responseData{
		Signature: signatureBase64,
	}

	// Set the response Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the response with a 200 OK status and the JSON-encoded responseData
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func parseCertificate(certPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to cast public key to RSA public key")
	}

	return pubKey, nil
}
