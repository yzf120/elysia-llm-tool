package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yzf120/elysia-llm-tool/config"
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

const (
	volcHost    = "open.volcengineapi.com"
	volcService = "ark"
	volcAction  = "GetInferenceUsage"
	volcVersion = "2024-01-01"
)

// InferenceUsageService 推理用量查询服务
type InferenceUsageService struct {
	cfg *config.Config
}

// NewInferenceUsageService 创建推理用量查询服务
func NewInferenceUsageService() *InferenceUsageService {
	return &InferenceUsageService{
		cfg: config.GetConfig(),
	}
}

// volcRequestBody 火山API请求体
type volcRequestBody struct {
	QueryInterval    string          `json:"QueryInterval"`
	StartTime        string          `json:"StartTime"`
	EndTime          string          `json:"EndTime"`
	Filters          []volcFilter    `json:"Filters,omitempty"`
	ProjectName      string          `json:"ProjectName,omitempty"`
	ShowWindowDetail bool            `json:"ShowWindowDetail,omitempty"`
}

type volcFilter struct {
	Key       string   `json:"Key"`
	Values    []string `json:"Values,omitempty"`
	ValueLike string   `json:"ValueLike,omitempty"`
}

// volcResponse 火山API响应
type volcResponse struct {
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
	} `json:"ResponseMetadata"`
	Result struct {
		Fields    []volcField  `json:"Fields"`
		Data      [][]string   `json:"Data"`
		DataCount int          `json:"DataCount"`
	} `json:"Result"`
}

type volcField struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
}

// GetInferenceUsage 查询推理用量
func (s *InferenceUsageService) GetInferenceUsage(req *pb.GetInferenceUsageRequest) (*pb.GetInferenceUsageResponse, error) {
	if s.cfg.VolcAccessKeyID == "" || s.cfg.VolcSecretAccessKey == "" {
		return nil, fmt.Errorf("火山开放平台 AK/SK 未配置，请设置 VOLC_ACCESS_KEY_ID 和 VOLC_SECRET_ACCESS_KEY")
	}

	// 构造请求体
	body := volcRequestBody{
		QueryInterval: req.QueryInterval,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	}
	if req.ProjectName != "" {
		body.ProjectName = req.ProjectName
	}
	for _, f := range req.Filters {
		vf := volcFilter{Key: f.Key}
		if len(f.Values) > 0 {
			vf.Values = f.Values
		}
		if f.ValueLike != "" {
			vf.ValueLike = f.ValueLike
		}
		body.Filters = append(body.Filters, vf)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 发起签名请求
	respData, err := s.signedRequest(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("调用火山API失败: %v", err)
	}

	// 解析响应
	var volcResp volcResponse
	if err := json.Unmarshal(respData, &volcResp); err != nil {
		return nil, fmt.Errorf("解析火山API响应失败: %v", err)
	}

	if volcResp.ResponseMetadata.Error != nil {
		return nil, fmt.Errorf("火山API返回错误: [%s] %s",
			volcResp.ResponseMetadata.Error.Code,
			volcResp.ResponseMetadata.Error.Message)
	}

	// 转换为 proto 响应
	resp := &pb.GetInferenceUsageResponse{
		DataCount: int32(volcResp.Result.DataCount),
	}

	for _, f := range volcResp.Result.Fields {
		resp.Fields = append(resp.Fields, &pb.UsageField{
			Name: f.Name,
			Type: f.Type,
		})
	}

	for _, row := range volcResp.Result.Data {
		resp.Data = append(resp.Data, &pb.UsageDataRow{
			Values: row,
		})
	}

	log.Printf("[InferenceUsage] 查询成功，返回 %d 条数据", volcResp.Result.DataCount)
	return resp, nil
}

// signedRequest 对火山开放平台发起 HMAC-SHA256 签名请求
func (s *InferenceUsageService) signedRequest(body []byte) ([]byte, error) {
	now := time.Now().UTC()
	xDate := now.Format("20060102T150405Z")
	shortDate := now.Format("20060102")

	// 计算 body hash
	bodyHash := sha256Hex(body)

	// 构造 canonical request
	queryString := fmt.Sprintf("Action=%s&Version=%s", volcAction, volcVersion)
	canonicalHeaders := fmt.Sprintf("content-type:application/json; charset=UTF-8\nhost:%s\nx-content-sha256:%s\nx-date:%s\n",
		volcHost, bodyHash, xDate)
	signedHeaders := "content-type;host;x-content-sha256;x-date"

	canonicalRequest := fmt.Sprintf("POST\n/\n%s\n%s\n%s\n%s",
		queryString, canonicalHeaders, signedHeaders, bodyHash)

	// 构造 string to sign
	credentialScope := fmt.Sprintf("%s/%s/%s/request", shortDate, s.cfg.VolcRegion, volcService)
	stringToSign := fmt.Sprintf("HMAC-SHA256\n%s\n%s\n%s",
		xDate, credentialScope, sha256Hex([]byte(canonicalRequest)))

	// 计算签名
	kDate := hmacSHA256([]byte(shortDate), []byte(s.cfg.VolcSecretAccessKey))
	kRegion := hmacSHA256([]byte(s.cfg.VolcRegion), kDate)
	kService := hmacSHA256([]byte(volcService), kRegion)
	kSigning := hmacSHA256([]byte("request"), kService)
	signature := hex.EncodeToString(hmacSHA256([]byte(stringToSign), kSigning))

	// 构造 Authorization header
	authorization := fmt.Sprintf("HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.cfg.VolcAccessKeyID, credentialScope, signedHeaders, signature)

	// 发起 HTTP 请求
	url := fmt.Sprintf("https://%s/?%s", volcHost, queryString)
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	httpReq.Header.Set("Host", volcHost)
	httpReq.Header.Set("X-Date", xDate)
	httpReq.Header.Set("X-Content-Sha256", bodyHash)
	httpReq.Header.Set("Authorization", authorization)

	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	if httpResp.StatusCode != 200 {
		log.Printf("[InferenceUsage] HTTP %d, body: %s", httpResp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// sha256Hex 计算 SHA256 哈希的十六进制字符串
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// hmacSHA256 计算 HMAC-SHA256
func hmacSHA256(data []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
