package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultDailyLimit = 10
	adminSessionToken = "admin-authenticated-token"
)

type AccessCode struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Label      string `json:"label"`
	IsActive   bool   `json:"is_active"`
	DailyLimit int    `json:"daily_limit"`
	UsedToday  int    `json:"used_today"`
	Remaining  int    `json:"remaining"`
	CreatedAt  string `json:"created_at"`
}

type AccessRequest struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Reason        string `json:"reason"`
	Status        string `json:"status"`
	GeneratedCode string `json:"generated_code"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// In-memory fallback store when PostgreSQL is unavailable
var (
	memAccessMu       sync.Mutex
	memAccessCodes    []AccessCode
	memAccessUsage    = make(map[string]int) // key: "codeID:YYYY-MM-DD"
	memAccessRequests []AccessRequest
	memNextCodeID     = 1
	memNextRequestID  = 1
)

func dbAvailable() bool {
	return db != nil && db.Ping() == nil
}

func normalizeAccessCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func generateRandomAccessCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	part := func(length int) string {
		b := make([]byte, length)
		for i := range b {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				b[i] = charset[0]
				continue
			}
			b[i] = charset[n.Int64()]
		}
		return string(b)
	}
	return fmt.Sprintf("DEV-%s-%s", part(4), part(4))
}

func todayKey() string {
	return time.Now().Format("2006-01-02")
}

func isAdminRequest(r *http.Request) bool {
	if r.Header.Get("X-Admin-Token") == adminSessionToken {
		return true
	}
	if r.URL.Query().Get("admin_token") == adminSessionToken {
		return true
	}
	return false
}

func getAccessCodeFromRequest(r *http.Request) string {
	if c := r.Header.Get("X-Access-Code"); strings.TrimSpace(c) != "" {
		return normalizeAccessCode(c)
	}
	if c := r.URL.Query().Get("access_code"); strings.TrimSpace(c) != "" {
		return normalizeAccessCode(c)
	}
	return ""
}

// VerifyAccessCode checks if a code exists and is active, returning its daily usage.
func VerifyAccessCode(rawCode string) (*AccessCode, error) {
	code := normalizeAccessCode(rawCode)
	if code == "" {
		return nil, errors.New("access code is required")
	}

	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		for _, ac := range memAccessCodes {
			if ac.Code == code && ac.IsActive {
				used := memAccessUsage[fmt.Sprintf("%d:%s", ac.ID, todayKey())]
				copyAc := ac
				copyAc.UsedToday = used
				rem := copyAc.DailyLimit - used
				if rem < 0 {
					rem = 0
				}
				copyAc.Remaining = rem
				return &copyAc, nil
			}
		}
		return nil, errors.New("invalid or inactive access code")
	}

	query := `
		SELECT ac.id, ac.code, ac.label, ac.is_active, ac.daily_limit, ac.created_at,
		       COALESCE(u.request_count, 0) AS used_today
		FROM access_codes ac
		LEFT JOIN access_code_usage u
		  ON u.code_id = ac.id AND u.usage_date = CURRENT_DATE
		WHERE ac.code = $1 AND ac.is_active = true
	`
	var ac AccessCode
	var createdAt time.Time
	err := db.QueryRow(query, code).Scan(
		&ac.ID, &ac.Code, &ac.Label, &ac.IsActive, &ac.DailyLimit, &createdAt, &ac.UsedToday,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid or inactive access code")
		}
		return nil, err
	}
	ac.CreatedAt = createdAt.Format(time.RFC3339)
	rem := ac.DailyLimit - ac.UsedToday
	if rem < 0 {
		rem = 0
	}
	ac.Remaining = rem
	return &ac, nil
}

// ConsumeAccessCodeQuota verifies the code and atomically increments today's usage if under daily_limit.
// Returns (accessCode, allowed, error).
func ConsumeAccessCodeQuota(rawCode string) (*AccessCode, bool, error) {
	code := normalizeAccessCode(rawCode)
	if code == "" {
		return nil, false, errors.New("access code is required")
	}

	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		for _, ac := range memAccessCodes {
			if ac.Code == code && ac.IsActive {
				key := fmt.Sprintf("%d:%s", ac.ID, todayKey())
				used := memAccessUsage[key]
				copyAc := ac
				if used >= copyAc.DailyLimit {
					copyAc.UsedToday = used
					copyAc.Remaining = 0
					return &copyAc, false, nil
				}
				used++
				memAccessUsage[key] = used
				copyAc.UsedToday = used
				copyAc.Remaining = copyAc.DailyLimit - used
				if copyAc.Remaining < 0 {
					copyAc.Remaining = 0
				}
				return &copyAc, true, nil
			}
		}
		return nil, false, errors.New("invalid or inactive access code")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	var ac AccessCode
	var createdAt time.Time
	err = tx.QueryRow(`
		SELECT id, code, label, is_active, daily_limit, created_at
		FROM access_codes
		WHERE code = $1 AND is_active = true
		FOR UPDATE
	`, code).Scan(&ac.ID, &ac.Code, &ac.Label, &ac.IsActive, &ac.DailyLimit, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, errors.New("invalid or inactive access code")
		}
		return nil, false, err
	}
	ac.CreatedAt = createdAt.Format(time.RFC3339)

	var usedToday int
	err = tx.QueryRow(`
		SELECT request_count
		FROM access_code_usage
		WHERE code_id = $1 AND usage_date = CURRENT_DATE
		FOR UPDATE
	`, ac.ID).Scan(&usedToday)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, err
	}

	if usedToday >= ac.DailyLimit {
		ac.UsedToday = usedToday
		ac.Remaining = 0
		return &ac, false, nil
	}

	err = tx.QueryRow(`
		INSERT INTO access_code_usage (code_id, usage_date, request_count)
		VALUES ($1, CURRENT_DATE, 1)
		ON CONFLICT (code_id, usage_date) DO UPDATE
		SET request_count = access_code_usage.request_count + 1
		RETURNING request_count
	`, ac.ID).Scan(&usedToday)
	if err != nil {
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	ac.UsedToday = usedToday
	ac.Remaining = ac.DailyLimit - usedToday
	if ac.Remaining < 0 {
		ac.Remaining = 0
	}
	return &ac, true, nil
}

func CreateAccessCode(label, customCode string, dailyLimit int) (*AccessCode, error) {
	code := normalizeAccessCode(customCode)
	if code == "" {
		code = generateRandomAccessCode()
	}
	if dailyLimit <= 0 {
		dailyLimit = defaultDailyLimit
	}
	label = strings.TrimSpace(label)

	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		for _, existing := range memAccessCodes {
			if existing.Code == code {
				return nil, errors.New("access code already exists")
			}
		}
		ac := AccessCode{
			ID:         memNextCodeID,
			Code:       code,
			Label:      label,
			IsActive:   true,
			DailyLimit: dailyLimit,
			UsedToday:  0,
			Remaining:  dailyLimit,
			CreatedAt:  time.Now().Format(time.RFC3339),
		}
		memNextCodeID++
		memAccessCodes = append([]AccessCode{ac}, memAccessCodes...)
		return &ac, nil
	}

	var ac AccessCode
	var createdAt time.Time
	err := db.QueryRow(`
		INSERT INTO access_codes (code, label, is_active, daily_limit)
		VALUES ($1, $2, true, $3)
		RETURNING id, code, label, is_active, daily_limit, created_at
	`, code, label, dailyLimit).Scan(&ac.ID, &ac.Code, &ac.Label, &ac.IsActive, &ac.DailyLimit, &createdAt)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, errors.New("access code already exists")
		}
		return nil, err
	}
	ac.CreatedAt = createdAt.Format(time.RFC3339)
	ac.UsedToday = 0
	ac.Remaining = ac.DailyLimit
	return &ac, nil
}

func ListAccessCodes() ([]AccessCode, error) {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		res := make([]AccessCode, 0, len(memAccessCodes))
		for _, ac := range memAccessCodes {
			used := memAccessUsage[fmt.Sprintf("%d:%s", ac.ID, todayKey())]
			copyAc := ac
			copyAc.UsedToday = used
			copyAc.Remaining = copyAc.DailyLimit - used
			if copyAc.Remaining < 0 {
				copyAc.Remaining = 0
			}
			res = append(res, copyAc)
		}
		return res, nil
	}

	rows, err := db.Query(`
		SELECT ac.id, ac.code, ac.label, ac.is_active, ac.daily_limit, ac.created_at,
		       COALESCE(u.request_count, 0) AS used_today
		FROM access_codes ac
		LEFT JOIN access_code_usage u
		  ON u.code_id = ac.id AND u.usage_date = CURRENT_DATE
		ORDER BY ac.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []AccessCode{}
	for rows.Next() {
		var ac AccessCode
		var createdAt time.Time
		if err := rows.Scan(&ac.ID, &ac.Code, &ac.Label, &ac.IsActive, &ac.DailyLimit, &createdAt, &ac.UsedToday); err != nil {
			continue
		}
		ac.CreatedAt = createdAt.Format(time.RFC3339)
		ac.Remaining = ac.DailyLimit - ac.UsedToday
		if ac.Remaining < 0 {
			ac.Remaining = 0
		}
		res = append(res, ac)
	}
	return res, nil
}

func UpdateAccessCodeStatus(id int, isActive bool) error {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		for i := range memAccessCodes {
			if memAccessCodes[i].ID == id {
				memAccessCodes[i].IsActive = isActive
				return nil
			}
		}
		return errors.New("access code not found")
	}

	res, err := db.Exec(`UPDATE access_codes SET is_active = $1 WHERE id = $2`, isActive, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("access code not found")
	}
	return nil
}

func DeleteAccessCode(id int) error {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		filtered := make([]AccessCode, 0, len(memAccessCodes))
		for _, ac := range memAccessCodes {
			if ac.ID != id {
				filtered = append(filtered, ac)
			}
		}
		memAccessCodes = filtered
		return nil
	}

	_, err := db.Exec(`DELETE FROM access_codes WHERE id = $1`, id)
	return err
}

func CreateAccessRequest(name, email, reason string) (*AccessRequest, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	reason = strings.TrimSpace(reason)
	if name == "" {
		return nil, errors.New("name is required")
	}

	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		now := time.Now().Format(time.RFC3339)
		req := AccessRequest{
			ID:        memNextRequestID,
			Name:      name,
			Email:     email,
			Reason:    reason,
			Status:    "pending",
			CreatedAt: now,
			UpdatedAt: now,
		}
		memNextRequestID++
		memAccessRequests = append([]AccessRequest{req}, memAccessRequests...)
		return &req, nil
	}

	var req AccessRequest
	var createdAt, updatedAt time.Time
	err := db.QueryRow(`
		INSERT INTO access_requests (name, email, reason, status)
		VALUES ($1, $2, $3, 'pending')
		RETURNING id, name, email, reason, status, generated_code, created_at, updated_at
	`, name, email, reason).Scan(
		&req.ID, &req.Name, &req.Email, &req.Reason, &req.Status, &req.GeneratedCode, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	req.CreatedAt = createdAt.Format(time.RFC3339)
	req.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &req, nil
}

func ListAccessRequests() ([]AccessRequest, error) {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		copied := make([]AccessRequest, len(memAccessRequests))
		copy(copied, memAccessRequests)
		return copied, nil
	}

	rows, err := db.Query(`
		SELECT id, name, email, reason, status, generated_code, created_at, updated_at
		FROM access_requests
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []AccessRequest{}
	for rows.Next() {
		var req AccessRequest
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&req.ID, &req.Name, &req.Email, &req.Reason, &req.Status, &req.GeneratedCode, &createdAt, &updatedAt); err != nil {
			continue
		}
		req.CreatedAt = createdAt.Format(time.RFC3339)
		req.UpdatedAt = updatedAt.Format(time.RFC3339)
		res = append(res, req)
	}
	return res, nil
}

func ApproveAccessRequest(id int, customCode string, dailyLimit int) (*AccessRequest, *AccessCode, error) {
	if !dbAvailable() {
		memAccessMu.Lock()
		var targetIdx = -1
		for i := range memAccessRequests {
			if memAccessRequests[i].ID == id {
				targetIdx = i
				break
			}
		}
		if targetIdx == -1 {
			memAccessMu.Unlock()
			return nil, nil, errors.New("access request not found")
		}
		guestLabel := memAccessRequests[targetIdx].Name
		if memAccessRequests[targetIdx].Email != "" {
			guestLabel = fmt.Sprintf("%s (%s)", memAccessRequests[targetIdx].Name, memAccessRequests[targetIdx].Email)
		}
		memAccessMu.Unlock()

		ac, err := CreateAccessCode(guestLabel, customCode, dailyLimit)
		if err != nil {
			return nil, nil, err
		}

		memAccessMu.Lock()
		memAccessRequests[targetIdx].Status = "approved"
		memAccessRequests[targetIdx].GeneratedCode = ac.Code
		memAccessRequests[targetIdx].UpdatedAt = time.Now().Format(time.RFC3339)
		updatedReq := memAccessRequests[targetIdx]
		memAccessMu.Unlock()
		return &updatedReq, ac, nil
	}

	var req AccessRequest
	var createdAt, updatedAt time.Time
	err := db.QueryRow(`
		SELECT id, name, email, reason, status, generated_code, created_at, updated_at
		FROM access_requests
		WHERE id = $1
	`, id).Scan(&req.ID, &req.Name, &req.Email, &req.Reason, &req.Status, &req.GeneratedCode, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("access request not found")
		}
		return nil, nil, err
	}

	guestLabel := req.Name
	if req.Email != "" {
		guestLabel = fmt.Sprintf("%s (%s)", req.Name, req.Email)
	}

	ac, err := CreateAccessCode(guestLabel, customCode, dailyLimit)
	if err != nil {
		return nil, nil, err
	}

	err = db.QueryRow(`
		UPDATE access_requests
		SET status = 'approved', generated_code = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, email, reason, status, generated_code, created_at, updated_at
	`, ac.Code, id).Scan(&req.ID, &req.Name, &req.Email, &req.Reason, &req.Status, &req.GeneratedCode, &createdAt, &updatedAt)
	if err != nil {
		return nil, nil, err
	}
	req.CreatedAt = createdAt.Format(time.RFC3339)
	req.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &req, ac, nil
}

func RejectAccessRequest(id int) error {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		for i := range memAccessRequests {
			if memAccessRequests[i].ID == id {
				memAccessRequests[i].Status = "rejected"
				memAccessRequests[i].UpdatedAt = time.Now().Format(time.RFC3339)
				return nil
			}
		}
		return errors.New("access request not found")
	}

	res, err := db.Exec(`UPDATE access_requests SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("access request not found")
	}
	return nil
}

func DeleteAccessRequest(id int) error {
	if !dbAvailable() {
		memAccessMu.Lock()
		defer memAccessMu.Unlock()
		filtered := make([]AccessRequest, 0, len(memAccessRequests))
		for _, req := range memAccessRequests {
			if req.ID != id {
				filtered = append(filtered, req)
			}
		}
		memAccessRequests = filtered
		return nil
	}

	_, err := db.Exec(`DELETE FROM access_requests WHERE id = $1`, id)
	return err
}

// accessControlMiddleware restricts Bible Study & Magisterium Chat to admins or guests with a valid access code.
// When consumeQuota is true, each guest request consumes 1 of their daily_limit (default 10) requests.
func accessControlMiddleware(consumeQuota bool, blockIfExhausted bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAdminRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		code := getAccessCodeFromRequest(r)
		if code == "" {
			respondWithError(w, http.StatusForbidden, "Access code required for this feature.")
			return
		}

		if consumeQuota {
			ac, allowed, err := ConsumeAccessCodeQuota(code)
			if err != nil {
				respondWithError(w, http.StatusForbidden, "Invalid or inactive access code.")
				return
			}
			if !allowed {
				respondWithError(w, http.StatusTooManyRequests, fmt.Sprintf("Daily limit of %d requests reached for this access code. Please try again tomorrow.", ac.DailyLimit))
				return
			}
			w.Header().Set("X-Access-Remaining", strconv.Itoa(ac.Remaining))
			w.Header().Set("X-Access-Limit", strconv.Itoa(ac.DailyLimit))
			next.ServeHTTP(w, r)
			return
		}

		ac, err := VerifyAccessCode(code)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Invalid or inactive access code.")
			return
		}
		if blockIfExhausted && ac.Remaining <= 0 {
			respondWithError(w, http.StatusTooManyRequests, fmt.Sprintf("Daily limit of %d requests reached for this access code. Please try again tomorrow.", ac.DailyLimit))
			return
		}
		w.Header().Set("X-Access-Remaining", strconv.Itoa(ac.Remaining))
		w.Header().Set("X-Access-Limit", strconv.Itoa(ac.DailyLimit))
		next.ServeHTTP(w, r)
	})
}

func verifyAccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if isAdminRequest(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":       true,
			"is_admin":    true,
			"daily_limit": -1,
			"used_today":  0,
			"remaining":   -1,
		})
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ac, err := VerifyAccessCode(req.Code)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid or inactive access code")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":       true,
		"is_admin":    false,
		"code":        ac.Code,
		"label":       ac.Label,
		"daily_limit": ac.DailyLimit,
		"used_today":  ac.UsedToday,
		"remaining":   ac.Remaining,
	})
}

func requestAccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	created, err := CreateAccessRequest(req.Name, req.Email, req.Reason)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Your access request has been submitted. Once approved by an administrator, you can enter your access code here.",
		"request": created,
	})
}

func adminAccessCodesHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		codes, err := ListAccessCodes()
		if err != nil {
			log.Printf("[adminAccessCodesHandler] List error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to load access codes")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"codes": codes,
		})

	case http.MethodPost:
		var req struct {
			Label      string `json:"label"`
			Code       string `json:"code"`
			DailyLimit int    `json:"daily_limit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		ac, err := CreateAccessCode(req.Label, req.Code, req.DailyLimit)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"code":    ac,
		})

	case http.MethodPut:
		var req struct {
			ID       int  `json:"id"`
			IsActive bool `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID <= 0 {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		if err := UpdateAccessCodeStatus(req.ID, req.IsActive); err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			respondWithError(w, http.StatusBadRequest, "Valid id query parameter required")
			return
		}
		if err := DeleteAccessCode(id); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete access code")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func adminAccessRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		reqs, err := ListAccessRequests()
		if err != nil {
			log.Printf("[adminAccessRequestsHandler] List error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to load access requests")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"requests": reqs,
		})

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			respondWithError(w, http.StatusBadRequest, "Valid id query parameter required")
			return
		}
		if err := DeleteAccessRequest(id); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete access request")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func adminApproveAccessRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID         int    `json:"id"`
		Code       string `json:"code"`
		DailyLimit int    `json:"daily_limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedReq, ac, err := ApproveAccessRequest(req.ID, req.Code, req.DailyLimit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"request": updatedReq,
		"code":    ac,
	})
}

func adminRejectAccessRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := RejectAccessRequest(req.ID); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
