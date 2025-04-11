package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nedpals/supabase-go"
)

// SupabaseClientExtended extends the supabase.Client with additional functionality
type SupabaseClientExtended struct {
	*supabase.Client
	apiKey string
	baseURL string
}

// Functions provides access to Supabase Edge Functions
type Functions struct {
	client *SupabaseClientExtended
}

// CreateClientExtended creates an extended Supabase client
func CreateClientExtended(supabaseURL, supabaseKey string) *SupabaseClientExtended {
	client := supabase.CreateClient(supabaseURL, supabaseKey)
	return &SupabaseClientExtended{
		Client: client,
		apiKey: supabaseKey,
		baseURL: supabaseURL,
	}
}

// Functions returns the Functions API
func (c *SupabaseClientExtended) Functions() *Functions {
	return &Functions{
		client: c,
	}
}

// Invoke calls a Supabase Edge Function or RPC function
func (f *Functions) Invoke(functionName string, body interface{}, result interface{}) error {
	// Use direct HTTP request approach for all functions
	// The supabase-go library doesn't have a direct way to call RPC functions
	// Convert the request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create the request URL - handle URL parsing more robustly
	// Avoid direct parsing of URLs with hostnames containing hyphens
	baseURLStr := f.client.baseURL
	// Ensure the baseURL doesn't have a trailing slash
	if baseURLStr[len(baseURLStr)-1] == '/' {
		baseURLStr = baseURLStr[:len(baseURLStr)-1]
	}
	
	// Construct the full URL as a string first
	requestURL := fmt.Sprintf("%s/rest/v1/rpc/%s", baseURLStr, functionName)
	
	// Now parse the complete URL
	baseURL, err := url.Parse(requestURL)
	if err != nil {
		return err
	}
	
	// Create the request
	req, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", f.client.apiKey)
	req.Header.Set("Authorization", "Bearer "+f.client.apiKey)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("function call failed with status code: %d", resp.StatusCode)
	}

	// Decode the response
	return json.NewDecoder(resp.Body).Decode(result)
}
