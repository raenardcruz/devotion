package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetInMemoryAccessState() {
	memAccessMu.Lock()
	defer memAccessMu.Unlock()
	memAccessCodes = nil
	memAccessUsage = make(map[string]int)
	memAccessRequests = nil
	memNextCodeID = 1
	memNextRequestID = 1
}

func TestAccessCodeDailyLimitAndAdminBypass(t *testing.T) {
	resetInMemoryAccessState()

	ac, err := CreateAccessCode("Test Guest", "", 10)
	if err != nil {
		t.Fatalf("failed to create access code: %v", err)
	}
	if ac.DailyLimit != 10 {
		t.Fatalf("expected daily limit 10, got %d", ac.DailyLimit)
	}

	protected := accessControlMiddleware(true, true, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// 1. Missing code should return 403
	req := httptest.NewRequest(http.MethodGet, "/context", nil)
	rr := httptest.NewRecorder()
	protected.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 without code, got %d", rr.Code)
	}

	// 2. First 10 requests with valid code should succeed
	for i := 1; i <= 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/context", nil)
		req.Header.Set("X-Access-Code", ac.Code)
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d (%s)", i, rr.Code, rr.Body.String())
		}
	}

	// 3. 11th request should return 429 Too Many Requests
	req11 := httptest.NewRequest(http.MethodGet, "/context", nil)
	req11.Header.Set("X-Access-Code", ac.Code)
	rr11 := httptest.NewRecorder()
	protected.ServeHTTP(rr11, req11)
	if rr11.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on 11th request, got %d", rr11.Code)
	}

	// 4. Admin token should bypass both access code and daily limit
	for i := 0; i < 15; i++ {
		adminReq := httptest.NewRequest(http.MethodGet, "/context", nil)
		adminReq.Header.Set("X-Admin-Token", adminSessionToken)
		adminRR := httptest.NewRecorder()
		protected.ServeHTTP(adminRR, adminReq)
		if adminRR.Code != http.StatusOK {
			t.Fatalf("admin request %d: expected 200, got %d", i, adminRR.Code)
		}
	}

	// 5. Revoking code should return 403 even if quota check isn't the blocker
	ac2, err := CreateAccessCode("Revokable", "TEST-CODE-123", 10)
	if err != nil {
		t.Fatalf("failed to create second code: %v", err)
	}
	if err := UpdateAccessCodeStatus(ac2.ID, false); err != nil {
		t.Fatalf("failed to revoke code: %v", err)
	}
	reqRevoked := httptest.NewRequest(http.MethodGet, "/context", nil)
	reqRevoked.Header.Set("X-Access-Code", ac2.Code)
	rrRevoked := httptest.NewRecorder()
	protected.ServeHTTP(rrRevoked, reqRevoked)
	if rrRevoked.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for revoked code, got %d", rrRevoked.Code)
	}
}

func TestAccessRequestAndApprovalFlow(t *testing.T) {
	resetInMemoryAccessState()

	// Guest submits access request
	body, _ := json.Marshal(map[string]string{
		"name":   "Jane Doe",
		"email":  "jane@example.com",
		"reason": "Parish bible study",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/access/request", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	requestAccessHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 from requestAccessHandler, got %d: %s", rr.Code, rr.Body.String())
	}

	reqs, err := ListAccessRequests()
	if err != nil || len(reqs) != 1 {
		t.Fatalf("expected 1 access request, got %d (err: %v)", len(reqs), err)
	}

	// Admin approves request
	approveBody, _ := json.Marshal(map[string]interface{}{
		"id": reqs[0].ID,
	})
	appReq := httptest.NewRequest(http.MethodPost, "/api/admin/access-requests/approve", bytes.NewReader(approveBody))
	appReq.Header.Set("X-Admin-Token", adminSessionToken)
	appRR := httptest.NewRecorder()
	adminApproveAccessRequestHandler(appRR, appReq)
	if appRR.Code != http.StatusOK {
		t.Fatalf("expected 200 from approve handler, got %d: %s", appRR.Code, appRR.Body.String())
	}

	var resp struct {
		Success bool          `json:"success"`
		Request AccessRequest `json:"request"`
		Code    AccessCode    `json:"code"`
	}
	if err := json.Unmarshal(appRR.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal approve response: %v", err)
	}
	if resp.Code.Code == "" || resp.Request.Status != "approved" {
		t.Fatalf("expected approved request with generated code, got %+v", resp)
	}

	// Verify the generated code works via verifyAccessHandler
	verifyBody, _ := json.Marshal(map[string]string{"code": resp.Code.Code})
	vReq := httptest.NewRequest(http.MethodPost, "/api/access/verify", bytes.NewReader(verifyBody))
	vRR := httptest.NewRecorder()
	verifyAccessHandler(vRR, vReq)
	if vRR.Code != http.StatusOK {
		t.Fatalf("expected 200 from verifyAccessHandler, got %d: %s", vRR.Code, vRR.Body.String())
	}
}
