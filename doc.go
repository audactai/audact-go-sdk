// Package audact is a hand-written, standard-library-only client for the Audact API.
//
// Quickstart:
//
//	import audact "github.com/audactai/audact-go-sdk"
//
//	c, err := audact.NewClient("") // reads AUDACT_API_KEY
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use the service accessors, for example c.Agents, c.Calls or c.Evidence.
//	_ = c
//
// # Base URL
//
// The base URL is resolved in this order: WithBaseURL, the AUDACT_BASE_URL
// environment variable, then DefaultBaseURL. A trailing slash and a trailing
// /v1 are stripped; every request path is prefixed with /v1.
//
// # API key
//
// An empty apiKey falls back to the AUDACT_API_KEY environment variable.
// A missing key is not an error because some operations are public.
package audact
