package main

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// PasswordEntry represents a password entry in the database
type PasswordEntry struct {
	ID        *int64  `json:"id,omitempty"`
	Title     string  `json:"title"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	URL       *string `json:"url,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Category  string  `json:"category"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

// EntryRequest represents a request to create/update an entry
type EntryRequest struct {
	Title    string  `json:"title"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	URL      *string `json:"url,omitempty"`
	Notes    *string `json:"notes,omitempty"`
	Category string  `json:"category"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool    `json:"success"`
	Token   *string `json:"token,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// EntryResponse represents entry operation response
type EntryResponse struct {
	Success bool            `json:"success"`
	Entry   *PasswordEntry  `json:"entry,omitempty"`
	Entries []PasswordEntry `json:"entries,omitempty"`
	Error   *string         `json:"error,omitempty"`
}

// BackendService handles communication with the Rust backend
type BackendService struct {
	backendPath string
}

// NewBackendService creates a new backend service
func NewBackendService(backendPath string) *BackendService {
	return &BackendService{
		backendPath: backendPath,
	}
}

// executeCommand runs a backend command and returns the output
func (b *BackendService) executeCommand(args ...string) (string, error) {
	cmd := exec.Command(b.backendPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// InitializeDatabase initializes the database with a master password
func (b *BackendService) InitializeDatabase(masterPassword string) error {
	_, err := b.executeCommand("init", masterPassword)
	return err
}

// Authenticate verifies the master password
func (b *BackendService) Authenticate(masterPassword string) (*AuthResponse, error) {
	output, err := b.executeCommand("auth", masterPassword)
	if err != nil {
		return nil, err
	}

	var response AuthResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// AddEntry adds a new password entry
func (b *BackendService) AddEntry(entry EntryRequest, token string) (*EntryResponse, error) {
	entryJSON, _ := json.Marshal(entry)
	output, err := b.executeCommand("add", string(entryJSON), token)
	if err != nil {
		return nil, err
	}

	var response EntryResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAllEntries retrieves all password entries
func (b *BackendService) GetAllEntries(token string) (*EntryResponse, error) {
	output, err := b.executeCommand("list", token)
	if err != nil {
		return nil, err
	}

	var response EntryResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetEntryByID retrieves a specific entry by ID
func (b *BackendService) GetEntryByID(id int64, token string) (*EntryResponse, error) {
	output, err := b.executeCommand("get", string(rune(id)), token)
	if err != nil {
		return nil, err
	}

	var response EntryResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateEntry updates an existing entry
func (b *BackendService) UpdateEntry(entry EntryRequest, id int64, token string) (*EntryResponse, error) {
	entryJSON, _ := json.Marshal(entry)
	output, err := b.executeCommand("update", string(rune(id)), string(entryJSON), token)
	if err != nil {
		return nil, err
	}

	var response EntryResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteEntry deletes an entry by ID
func (b *BackendService) DeleteEntry(id int64, token string) error {
	_, err := b.executeCommand("delete", string(rune(id)), token)
	return err
}

// SearchEntries searches for entries matching a query
func (b *BackendService) SearchEntries(query string, token string) (*EntryResponse, error) {
	output, err := b.executeCommand("search", query, token)
	if err != nil {
		return nil, err
	}

	var response EntryResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCategories retrieves all unique categories
func (b *BackendService) GetCategories(token string) ([]string, error) {
	output, err := b.executeCommand("categories", token)
	if err != nil {
		return nil, err
	}

	var categories []string
	err = json.Unmarshal([]byte(output), &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
