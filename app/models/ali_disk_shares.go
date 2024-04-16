package models

import (
	"bytes"
	"encoding/json"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type AliDiskShares struct {
	orm.Model
	ShareId  string `gorm:"column:share_id"`
	Password string `gorm:"column:password"`
	SyncFlag string `gorm:"column:flag"`
	orm.SoftDeletes
}

func (t AliDiskShares) RandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (t *AliDiskShares) GetToken(ip string) string {
	var cacheName string
	if ip != "" {
		cacheName = "aliyunpansharetoken:" + ip + ":" + t.ShareId
	} else {
		cacheName = "aliyunpansharetoken:" + t.ShareId
	}
	value := facades.Cache().Get(cacheName, func() any {
		//https://api.aliyundrive.com/v2/share_link/get_share_token
		url := "https://api.aliyundrive.com/v2/share_link/get_share_token"
		// 准备请求体数据
		requestData := map[string]interface{}{
			"share_id":  t.ShareId,
			"share_pwd": t.Password,
		}
		// 将请求体数据编码为JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			panic(err) // 处理JSON编码错误
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := (&http.Client{}).Do(req)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		response := &GetShareTokenResponse{}
		if err := json.Unmarshal(body, response); err != nil {
			panic(err)
		}
		facades.Cache().Add(cacheName, body, time.Duration(response.ExpiresIn-300)*time.Second)
		return facades.Cache().Get(cacheName)
	})
	response1 := &GetShareTokenResponse{}
	if err := json.Unmarshal([]byte(value.(string)), response1); err != nil {
		panic(err)
	}
	return response1.ShareToken
}

func (t *AliDiskShares) GetFileListByFileId(fileId string, marker string) []AliDiskShareFile {
	url := "https://api.aliyundrive.com/adrive/v2/file/list_by_share"
	// 准备请求体数据
	requestData := map[string]interface{}{
		"limit":           100,
		"order_by":        "name",
		"order_direction": "DESC",
		"parent_file_id":  fileId,
		"share_id":        t.ShareId,
		"marker":          marker,
	}
	// 将请求体数据编码为JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(err) // 处理JSON编码错误
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Share-Token", t.GetToken(""))
	resp, err := (&http.Client{}).Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	response := &ListByShareResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	facades.Log().Info(string(body))
	facades.Log().Info(response)
	for _, item := range response.Items {
		facades.Log().Info(item.DriveId)
		facades.Log().Info(item.DomainId)
	}
	return response.Items
}

type ListByShareResponse struct {
	Items      []AliDiskShareFile `json:"items"`
	NextMarker string             `json:"next_marker"`
}

type AliDiskShareFile struct {
	DriveId      string `gorm:"column:drive_id" json:"drive_id"`
	DomainId     string `gorm:"column:domain_id" json:"domain_id"`
	FileId       string `gorm:"column:file_id" json:"file_id"`
	ShareId      string `gorm:"column:share_id" json:"share_id"`
	Name         string `gorm:"column:name" json:"name"`
	Type         string `gorm:"column:type" json:"type"`
	ParentFileId string `gorm:"column:parent_file_id" json:"parent_file_id"`
	Json         string `gorm:"column:json"`
	orm.Model
	orm.SoftDeletes
}

type GetShareTokenResponse struct {
	ExpireTime string `json:"expire_time"`
	ExpiresIn  int    `json:"expires_in"`
	ShareToken string `json:"share_token"`
}
