package audact

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// KYCRequirements lists regulatory requirements for a country and number type.
// Calls GET /v1/kyc/requirements.
func (s *NumbersService) KYCRequirements(ctx context.Context, params *KYCRequirementsParams) (*ListKYCRequirementsResponse, error) {
	q := url.Values{}
	if params != nil {
		addString(q, "country", params.Country)
		addString(q, "number_type", params.NumberType)
	}
	var out ListKYCRequirementsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/kyc/requirements", authAPIKey, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListKYCDocuments lists uploaded KYC documents. Calls GET /v1/kyc/documents.
func (s *NumbersService) ListKYCDocuments(ctx context.Context) (*ListKYCDocumentsResponse, error) {
	var out ListKYCDocumentsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/kyc/documents", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// multipartBody builds a multipart/form-data body with optional text fields
// and one file part.
func multipartBody(fields [][2]string, filename string, file io.Reader) (*rawBody, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for _, f := range fields {
		if err := mw.WriteField(f[0], f[1]); err != nil {
			return nil, fmt.Errorf("audact: build multipart: %w", err)
		}
	}
	if filename == "" {
		filename = "upload"
	}
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("audact: build multipart: %w", err)
	}
	if _, err := io.Copy(fw, file); err != nil {
		return nil, fmt.Errorf("audact: read upload: %w", err)
	}
	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("audact: build multipart: %w", err)
	}
	return &rawBody{contentType: mw.FormDataContentType(), data: buf.Bytes()}, nil
}

func kycIdem(opts *KYCUploadOptions) map[string]string {
	if opts == nil {
		return nil
	}
	return idemHeader(opts.IdempotencyKey)
}

// UploadKYCDocument uploads a KYC document as multipart/form-data. opts may be nil.
// Calls POST /v1/kyc/documents.
func (s *NumbersService) UploadKYCDocument(ctx context.Context, file io.Reader, filename, docType string, opts *KYCUploadOptions) (*KYCDocument, error) {
	body, err := multipartBody([][2]string{{"doc_type", docType}}, filename, file)
	if err != nil {
		return nil, err
	}
	var out KYCDocument
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/kyc/documents", authAPIKey, nil, kycIdem(opts), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKYCDocument returns a KYC document's metadata. Calls GET /v1/kyc/documents/{id}.
func (s *NumbersService) GetKYCDocument(ctx context.Context, id string) (*KYCDocument, error) {
	var out KYCDocument
	if _, err := s.t.do(ctx, http.MethodGet, "/kyc/documents/"+url.PathEscape(id), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReplaceKYCDocument replaces a KYC document's file as multipart/form-data.
// opts may be nil. Calls PUT /v1/kyc/documents/{id}.
func (s *NumbersService) ReplaceKYCDocument(ctx context.Context, id string, file io.Reader, filename string, opts *KYCUploadOptions) (*KYCDocument, error) {
	body, err := multipartBody(nil, filename, file)
	if err != nil {
		return nil, err
	}
	var out KYCDocument
	if _, err := s.t.doHeaders(ctx, http.MethodPut, "/kyc/documents/"+url.PathEscape(id), authAPIKey, nil, kycIdem(opts), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteKYCDocument deletes a KYC document. Calls DELETE /v1/kyc/documents/{id}.
func (s *NumbersService) DeleteKYCDocument(ctx context.Context, id string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/kyc/documents/"+url.PathEscape(id), authAPIKey, nil, nil, nil)
	return err
}

// ListKYCDocumentSets lists KYC document sets. Calls GET /v1/kyc/document-sets.
func (s *NumbersService) ListKYCDocumentSets(ctx context.Context) (*ListKYCDocumentSetsResponse, error) {
	var out ListKYCDocumentSetsResponse
	if _, err := s.t.do(ctx, http.MethodGet, "/kyc/document-sets", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateKYCDocumentSet creates a KYC document set; a non-empty idempotencyKey
// is sent as the Idempotency-Key header. Calls POST /v1/kyc/document-sets.
func (s *NumbersService) CreateKYCDocumentSet(ctx context.Context, req *CreateKYCDocumentSetRequest, idempotencyKey string) (*KYCDocumentSet, error) {
	var out KYCDocumentSet
	if _, err := s.t.doHeaders(ctx, http.MethodPost, "/kyc/document-sets", authAPIKey, nil, idemHeader(idempotencyKey), orEmpty(req), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKYCDocumentSet returns a KYC document set. Calls GET /v1/kyc/document-sets/{id}.
func (s *NumbersService) GetKYCDocumentSet(ctx context.Context, id string) (*KYCDocumentSet, error) {
	var out KYCDocumentSet
	if _, err := s.t.do(ctx, http.MethodGet, "/kyc/document-sets/"+url.PathEscape(id), authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteKYCDocumentSet deletes a KYC document set. Calls DELETE /v1/kyc/document-sets/{id}.
func (s *NumbersService) DeleteKYCDocumentSet(ctx context.Context, id string) error {
	_, err := s.t.do(ctx, http.MethodDelete, "/kyc/document-sets/"+url.PathEscape(id), authAPIKey, nil, nil, nil)
	return err
}

// RotateKYCDEK rotates the tenant's document encryption key. Calls POST /v1/kyc/dek/rotate.
func (s *NumbersService) RotateKYCDEK(ctx context.Context) (*RotateKYCDEKResponse, error) {
	var out RotateKYCDEKResponse
	if _, err := s.t.do(ctx, http.MethodPost, "/kyc/dek/rotate", authAPIKey, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
