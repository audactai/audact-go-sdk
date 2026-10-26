package audact

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func callPIIPath(callID, suffix string) string {
	return "/calls/" + url.PathEscape(callID) + "/pii/" + suffix
}

// Shred erases an identity across all calls; irreversible, so pass an
// idempotencyKey to make retries safe. Calls POST /v1/pii/shred.
func (s *PIIService) Shred(ctx context.Context, req *PIIShredRequest, idempotencyKey string) (*PIIShredResponse, error) {
	var out PIIShredResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/pii/shred", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListShreds returns the PII-free erasure log. Calls GET /v1/pii/shreds.
func (s *PIIService) ListShreds(ctx context.Context, params *ListPIIShredsParams) (*PIIShredsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "cursor", params.Cursor)
		addInt(q, "limit", params.Limit)
	}
	var out PIIShredsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/pii/shreds", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListCallIdentities lists identities linked to a call. Calls GET /v1/calls/{call_id}/pii/identities.
func (s *PIIService) ListCallIdentities(ctx context.Context, callID string) (*PIIIdentitiesResponse, error) {
	var out PIIIdentitiesResponse
	if _, err := s.t.do(ctx, http.MethodGet, callPIIPath(callID, "identities"), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LinkCallIdentity links an identity to a call after the fact. Calls POST /v1/calls/{call_id}/pii/identities.
func (s *PIIService) LinkCallIdentity(ctx context.Context, callID string, req *PIILinkIdentityRequest) (*PIILinkIdentityResponse, error) {
	var out PIILinkIdentityResponse
	if _, err := s.t.do(ctx, http.MethodPost, callPIIPath(callID, "identities"), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ShredCall erases the listed identities of a call (all when none are listed);
// irreversible, so pass an idempotencyKey. Calls POST /v1/calls/{call_id}/pii/shred.
func (s *PIIService) ShredCall(ctx context.Context, callID string, req *PIICallShredRequest, idempotencyKey string) (*PIIShredResponse, error) {
	var out PIIShredResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPost, callPIIPath(callID, "shred"), authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Decrypt decrypts tokens in Text or Texts. It deliberately takes no
// idempotency key so plaintext is never cached server-side. Calls
// POST /v1/calls/{call_id}/pii/decrypt.
func (s *PIIService) Decrypt(ctx context.Context, callID string, req *PIIDecryptRequest) (*PIIDecryptResponse, error) {
	var out PIIDecryptResponse
	if _, err := s.t.do(ctx, http.MethodPost, callPIIPath(callID, "decrypt"), authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListErasedKeysets reports which of keysetIDs were crypto-erased for a call;
// keysetIDs must not be empty (no request is made otherwise). Calls
// GET /v1/calls/{call_id}/pii/erased-keysets with one keyset_id per id.
func (s *PIIService) ListErasedKeysets(ctx context.Context, callID string, keysetIDs []string) (*PIIErasedKeysetsResponse, error) {
	if len(keysetIDs) == 0 {
		return nil, errors.New("audact: list erased keysets: keysetIDs must contain at least one id")
	}
	q := url.Values{}
	for _, id := range keysetIDs {
		q.Add("keyset_id", id)
	}
	var out PIIErasedKeysetsResponse
	if _, err := s.t.do(ctx, http.MethodGet, callPIIPath(callID, "erased-keysets"), authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListExemptions lists PII exemptions. Calls GET /v1/pii/exemptions.
func (s *PIIService) ListExemptions(ctx context.Context) (*PIIExemptionsResponse, error) {
	var out PIIExemptionsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/pii/exemptions", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddExemption adds a PII exemption. Calls POST /v1/pii/exemptions.
func (s *PIIService) AddExemption(ctx context.Context, req *PIIExemptionRequest) (*PIIExemption, error) {
	var out PIIExemption
	if _, err := s.t.do(ctx, http.MethodPost, "/pii/exemptions", authAPIKey, nil, orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteExemption deletes an exemption using the legacy request-body form.
// Prefer DeleteExemptionByHash, since proxies may drop DELETE bodies. Calls
// DELETE /v1/pii/exemptions.
func (s *PIIService) DeleteExemption(ctx context.Context, req *PIIExemptionRequest) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/pii/exemptions", authAPIKey, nil, orEmpty(req), nil)
	return err
}

// DeleteExemptionByHash deletes an exemption by type and client hash. Calls
// DELETE /v1/pii/exemptions/{exempt_type}/{identifier_hash}.
func (s *PIIService) DeleteExemptionByHash(ctx context.Context, exemptType PIIExemptType, identifierHash string) error {
	path := fmt.Sprintf("/pii/exemptions/%s/%s", url.PathEscape(string(exemptType)), url.PathEscape(identifierHash))
	_, err := s.t.do(ctx, http.MethodDelete, path, authAPIKey, nil, nil, nil)
	return err
}

// Settings returns the PII encryption switch state. Calls GET /v1/pii/settings.
func (s *PIIService) Settings(ctx context.Context) (*PIISettingsResponse, error) {
	var out PIISettingsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/pii/settings", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSettings flips the PII encryption switch (owner role); a non-empty
// idempotencyKey is sent as Idempotency-Key. Calls PUT /v1/pii/settings.
func (s *PIIService) UpdateSettings(ctx context.Context, req *PIISettingsPutRequest, idempotencyKey string) (*PIISettingsPutResponse, error) {
	var out PIISettingsPutResponse
	if _, err := s.t.doHeaders(ctx, http.MethodPut, "/pii/settings", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
