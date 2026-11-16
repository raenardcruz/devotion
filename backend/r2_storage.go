package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type R2Config struct {
	AccountID       string `json:"r2_account_id"`
	Endpoint        string `json:"r2_endpoint"`
	Bucket          string `json:"r2_bucket"`
	AccessKeyID     string `json:"r2_access_key_id"`
	SecretAccessKey string `json:"r2_secret_access_key"`
	PublicURL       string `json:"r2_public_url"`
}

var r2HTTPClient = &http.Client{Timeout: 45 * time.Second}

func ResolveR2Config(override *R2Config) R2Config {
	settings := GetSettings()
	cfg := R2Config{
		AccountID:       strings.TrimSpace(settings.R2AccountID),
		Endpoint:        strings.TrimSpace(settings.R2Endpoint),
		Bucket:          strings.TrimSpace(settings.R2Bucket),
		AccessKeyID:     strings.TrimSpace(settings.R2AccessKeyID),
		SecretAccessKey: strings.TrimSpace(settings.R2SecretAccessKey),
		PublicURL:       strings.TrimSpace(settings.R2PublicURL),
	}

	if override != nil {
		if v := strings.TrimSpace(override.AccountID); v != "" {
			cfg.AccountID = v
		}
		if v := strings.TrimSpace(override.Endpoint); v != "" {
			cfg.Endpoint = v
		}
		if v := strings.TrimSpace(override.Bucket); v != "" {
			cfg.Bucket = v
		}
		if v := strings.TrimSpace(override.AccessKeyID); v != "" {
			cfg.AccessKeyID = v
		}
		if v := strings.TrimSpace(override.SecretAccessKey); v != "" && !isMaskedKey(v) {
			cfg.SecretAccessKey = v
		}
		if v := strings.TrimSpace(override.PublicURL); v != "" {
			cfg.PublicURL = v
		}
	}

	if cfg.AccountID == "" {
		cfg.AccountID = strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID"))
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = strings.TrimSpace(os.Getenv("R2_ENDPOINT"))
	}
	if cfg.Bucket == "" {
		cfg.Bucket = strings.TrimSpace(os.Getenv("R2_BUCKET"))
	}
	if cfg.AccessKeyID == "" {
		cfg.AccessKeyID = strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID"))
	}
	if cfg.SecretAccessKey == "" {
		cfg.SecretAccessKey = strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY"))
	}
	if cfg.PublicURL == "" {
		cfg.PublicURL = strings.TrimSpace(os.Getenv("R2_PUBLIC_URL"))
	}

	return cfg
}

func (c R2Config) ResolvedEndpoint() string {
	ep := strings.TrimRight(strings.TrimSpace(c.Endpoint), "/")
	if ep != "" {
		return ep
	}
	if acct := strings.TrimSpace(c.AccountID); acct != "" {
		return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", acct)
	}
	return ""
}

func (c R2Config) IsConfigured() bool {
	return c.ResolvedEndpoint() != "" &&
		strings.TrimSpace(c.Bucket) != "" &&
		strings.TrimSpace(c.AccessKeyID) != "" &&
		strings.TrimSpace(c.SecretAccessKey) != ""
}

func (c R2Config) PublicAssetURL(objectKey string) string {
	cleanKey := strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if pub := strings.TrimRight(strings.TrimSpace(c.PublicURL), "/"); pub != "" {
		return fmt.Sprintf("%s/%s", pub, awsPercentEncode(cleanKey, false))
	}
	return fmt.Sprintf("/api/manga/asset?key=%s", url.QueryEscape(cleanKey))
}

func (c R2Config) objectRequestURL(objectKey string, rawQuery string) (string, error) {
	ep := c.ResolvedEndpoint()
	if ep == "" {
		return "", fmt.Errorf("R2 endpoint or account ID is not configured")
	}
	u, err := url.Parse(ep)
	if err != nil {
		return "", fmt.Errorf("invalid R2 endpoint URL: %w", err)
	}

	basePath := strings.TrimRight(u.Path, "/")
	bucketSeg := awsPercentEncode(strings.TrimSpace(c.Bucket), true)
	cleanKey := strings.TrimLeft(strings.TrimSpace(objectKey), "/")

	if cleanKey != "" {
		u.Path = fmt.Sprintf("%s/%s/%s", basePath, bucketSeg, awsPercentEncode(cleanKey, false))
	} else {
		u.Path = fmt.Sprintf("%s/%s", basePath, bucketSeg)
	}
	u.RawPath = u.Path
	u.RawQuery = rawQuery
	return u.String(), nil
}

func (c R2Config) PutObject(ctx context.Context, objectKey string, data []byte, contentType string) error {
	if !c.IsConfigured() {
		return fmt.Errorf("R2 storage is not configured")
	}

	reqURL, err := c.objectRequestURL(objectKey, "")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if strings.TrimSpace(contentType) != "" {
		req.Header.Set("Content-Type", contentType)
	}

	signR2Request(req, data, c, time.Now().UTC())

	resp, err := r2HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("R2 PutObject request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("R2 PutObject failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func (c R2Config) GetObject(ctx context.Context, objectKey string) ([]byte, string, error) {
	if !c.IsConfigured() {
		return nil, "", fmt.Errorf("R2 storage is not configured")
	}

	reqURL, err := c.objectRequestURL(objectKey, "")
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, "", err
	}

	signR2Request(req, nil, c, time.Now().UTC())

	resp, err := r2HTTPClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("R2 GetObject request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, "", fmt.Errorf("object not found in R2: %s", objectKey)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, "", fmt.Errorf("R2 GetObject failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	return data, ct, nil
}

func (c R2Config) DeleteObject(ctx context.Context, objectKey string) error {
	if !c.IsConfigured() {
		return fmt.Errorf("R2 storage is not configured")
	}

	reqURL, err := c.objectRequestURL(objectKey, "")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}

	signR2Request(req, nil, c, time.Now().UTC())

	resp, err := r2HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("R2 DeleteObject request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusNoContent || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	return fmt.Errorf("R2 DeleteObject failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

func (c R2Config) DeleteObjects(ctx context.Context, keys []string) error {
	var errs []string
	for _, k := range uniqueStrings(keys) {
		if err := c.DeleteObject(ctx, k); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed deleting some R2 objects: %s", strings.Join(errs, "; "))
	}
	return nil
}

type r2ListBucketResult struct {
	KeyCount int `xml:"KeyCount"`
	Contents []struct {
		Key string `xml:"Key"`
	} `xml:"Contents"`
}

func (c R2Config) ListObjectKeys(ctx context.Context, prefix string) ([]string, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("R2 storage is not configured")
	}

	vals := url.Values{}
	vals.Set("list-type", "2")
	cleanPrefix := strings.TrimLeft(strings.TrimSpace(prefix), "/")
	if cleanPrefix != "" {
		vals.Set("prefix", cleanPrefix)
	}

	reqURL, err := c.objectRequestURL("", vals.Encode())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	signR2Request(req, nil, c, time.Now().UTC())

	resp, err := r2HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("R2 ListObjectKeys request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("R2 ListObjectKeys failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed r2ListBucketResult
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse R2 ListBucketResult: %w", err)
	}

	keys := make([]string, 0, len(parsed.Contents))
	for _, item := range parsed.Contents {
		if k := strings.TrimSpace(item.Key); k != "" {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (c R2Config) TestConnection(ctx context.Context) (int, error) {
	if !c.IsConfigured() {
		return 0, fmt.Errorf("R2 connection config is incomplete (Account ID/Endpoint, Bucket, Access Key ID, and Secret Access Key are required)")
	}

	reqURL, err := c.objectRequestURL("", "list-type=2&max-keys=1")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	signR2Request(req, nil, c, time.Now().UTC())

	resp, err := r2HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("R2 connection failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("R2 returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed r2ListBucketResult
	_ = xml.Unmarshal(body, &parsed)
	return parsed.KeyCount, nil
}

func signR2Request(req *http.Request, payload []byte, cfg R2Config, now time.Time) {
	amzDate := now.UTC().Format("20060102T150405Z")
	dateStamp := now.UTC().Format("20060102")
	region := "auto"
	service := "s3"

	sum := sha256.Sum256(payload)
	payloadHash := hex.EncodeToString(sum[:])

	host := req.URL.Host
	req.Host = host
	req.Header.Set("Host", host)
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)

	canonicalURI := req.URL.EscapedPath()
	if canonicalURI == "" {
		canonicalURI = "/"
	}

	canonicalQuery := buildCanonicalQuery(req.URL.RawQuery)
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", host, payloadHash, amzDate)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	crSum := sha256.Sum256([]byte(canonicalRequest))
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hex.EncodeToString(crSum[:]),
	}, "\n")

	kDate := hmacSHA256([]byte("AWS4"+cfg.SecretAccessKey), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, stringToSign))

	authHeader := fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		cfg.AccessKeyID,
		credentialScope,
		signedHeaders,
		signature,
	)
	req.Header.Set("Authorization", authHeader)
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func buildCanonicalQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	parsed, err := url.ParseQuery(rawQuery)
	if err != nil {
		return rawQuery
	}
	type kv struct {
		k string
		v string
	}
	var pairs []kv
	for k, vals := range parsed {
		ek := awsPercentEncode(k, true)
		if len(vals) == 0 {
			pairs = append(pairs, kv{ek, ""})
		} else {
			for _, v := range vals {
				pairs = append(pairs, kv{ek, awsPercentEncode(v, true)})
			}
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].k == pairs[j].k {
			return pairs[i].v < pairs[j].v
		}
		return pairs[i].k < pairs[j].k
	})
	parts := make([]string, len(pairs))
	for i, p := range pairs {
		parts[i] = p.k + "=" + p.v
	}
	return strings.Join(parts, "&")
}

func awsPercentEncode(s string, encodeSlash bool) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		isUnreserved := (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~'
		if isUnreserved || (c == '/' && !encodeSlash) {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
