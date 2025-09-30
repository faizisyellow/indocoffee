package uploadthing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/faizisyellow/indocoffee/internal/uploader"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/google/uuid"
)

type Uploadthing struct {
	apiKey        string
	presignUrl    string
	poolUploadUrl string
	acl           string
	slug          string
	actor         string
	metaUrl       string
	callbackUrl   string
	deleteUrl     string
	appId         string
}

type RegisterResponse struct {
	FileKey   string `json:"key"`
	UploadURL string `json:"url"`
}

type Bucket struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Region          string `json:"region"`
	Provider        string `json:"provider"`
	ProviderOwnerId string `json:"-"`
	CreatedAt       string `json:"createdAt"`
}

type FileData struct {
	AppId        string `json:"-"`
	FileKey      string `json:"fileKey"`
	FileUrl      string `json:"fileUrl"`
	FileType     string `json:"fileType"`
	CallbackSlug string `json:"callbackSlug"`
	CallbackUrl  string `json:"callbackUrl"`
	Metadata     any    `json:"-"`
	Filename     string `json:"fileName"`
	FileSize     string `json:"fileSize"`
	CustomId     string `json:"customId"`
	Acl          string `json:"acl"`
}

type ResponsePoolUpload struct {
	Status string `json:"status"`
	FileData
	ApiKey string `json:"-"`
	UfsUrl string `json:"ufsUrl"`
}

const (
	contentDisposition = "inline"
	expiresIn          = 300
)

func New(apiKey, presign, poolUpload, acl, slg, act, mturl, cllbckurl, dltUrl, appId string) *Uploadthing {

	return &Uploadthing{
		apiKey:        apiKey,
		presignUrl:    presign,
		poolUploadUrl: poolUpload,
		acl:           acl,
		slug:          slg,
		actor:         act,
		metaUrl:       mturl,
		callbackUrl:   cllbckurl,
		deleteUrl:     dltUrl,
		appId:         appId,
	}
}

func (u *Uploadthing) UploadFile(_ context.Context, file uploader.FileInput) (string, error) {
	registerResult, err := u.Register(file.Name, file.MimeType, int(file.Size))
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 20 * time.Second}

	metaPayload := map[string]any{
		"fileKeys": []string{registerResult.FileKey},
		"metadata": map[string]string{
			"uploadedBy": u.actor,
		},
		"callbackUrl":     u.callbackUrl,
		"callbackSlug":    u.slug,
		"awaitServerData": false,
		"isDev":           false,
	}

	metaJSON, err := json.Marshal(metaPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata payload: %w", err)
	}

	metaReq, err := http.NewRequest("POST", u.metaUrl, bytes.NewBuffer(metaJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create metadata request: %w", err)
	}
	metaReq.Header.Set("Content-Type", "application/json")
	metaReq.Header.Set("x-uploadthing-api-key", u.apiKey)

	metaResp, err := client.Do(metaReq)
	if err != nil {
		return "", fmt.Errorf("metadata registration failed: %w", err)
	}
	defer metaResp.Body.Close()

	if metaResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(metaResp.Body)
		return "", fmt.Errorf("metadata error: %s", string(b))
	}

	body, contentType, err := uploader.CreateMultipartBody(file)
	if err != nil {
		return "", err
	}

	putReq, err := http.NewRequest("PUT", registerResult.UploadURL, body)
	if err != nil {
		return "", err
	}
	putReq.Header.Set("Content-Type", contentType)

	putResp, err := client.Do(putReq)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(putResp.Body)
		return "", fmt.Errorf("upload error: %s", string(b))
	}

	checkUpload, err := u.PoolUpload(registerResult.FileKey)
	if err != nil {
		return "", err
	}

	if checkUpload.Status != "done" {
		return "", errors.New("upload failed: upload not done")
	}

	return u.GetUrls(registerResult.FileKey), nil
}

func (u *Uploadthing) DeleteFile(_ context.Context, filekey string) error {
	payload := map[string]any{
		"fileKeys": []string{filekey},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata payload: %w", err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("POST", u.deleteUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-uploadthing-api-key", u.apiKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("delete error: %s", string(b))
	}

	return nil
}

func (u *Uploadthing) Register(filename, filetype string, filesize int) (*RegisterResponse, error) {

	uid := utils.UUID{Plaintoken: uuid.New().String()}

	registerReq := struct {
		Filename           string `json:"fileName"`
		Filesize           int    `json:"fileSize"`
		Slug               string `json:"slug"`
		Filetype           string `json:"fileType"`
		CustomId           string `json:"customId"`
		ContentDisposition string `json:"contentDisposition"`
		Acl                string `json:"acl"`
		ExpiresIn          int    `json:"expiresIn"`
	}{
		Filename:           filename,
		Filesize:           filesize,
		Slug:               u.slug,
		Filetype:           filetype,
		CustomId:           uid.Generate(),
		ContentDisposition: contentDisposition,
		Acl:                u.acl,
		ExpiresIn:          expiresIn,
	}

	body, _ := json.Marshal(registerReq)
	req, err := http.NewRequest("POST", u.presignUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-uploadthing-api-key", u.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register failed: %s", string(b))
	}

	var regRes RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&regRes); err != nil {
		return nil, err
	}

	if regRes.FileKey == "" && regRes.UploadURL == "" {
		return nil, errors.New("error response should be valid")
	}

	return &regRes, nil
}

func (u *Uploadthing) PoolUpload(filekey string) (ResponsePoolUpload, error) {

	// build URL
	endpoint, err := url.Parse(fmt.Sprintf("%s/%s", u.poolUploadUrl, filekey))
	if err != nil {
		return ResponsePoolUpload{}, fmt.Errorf("invalid url: %w", err)
	}

	// build request
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return ResponsePoolUpload{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-uploadthing-api-key", u.apiKey)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponsePoolUpload{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponsePoolUpload{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// decode JSON response
	var result ResponsePoolUpload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResponsePoolUpload{}, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (u *Uploadthing) GetUrls(filekey string) string {
	if filekey == "" {
		return ""
	}

	urlBuilder := strings.Builder{}
	urlBuilder.WriteString("https://")
	urlBuilder.WriteString(u.appId)
	urlBuilder.WriteString(".ufs.sh/f/")
	urlBuilder.WriteString(filekey)

	return urlBuilder.String()
}

func GetFileKey(url string) string {
	if url == "" {
		return ""
	}

	_, filekey, found := strings.Cut(url, ".ufs.sh/f/")
	if !found {
		return ""
	}

	return filekey
}
