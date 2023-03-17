package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://h5.tu.qq.com/"

func RandomUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}
	uuid[8] = 0x80 | (uuid[8] & 0x3F)
	uuid[6] = 0x40 | (uuid[6] & 0x0F)
	return hex.EncodeToString(uuid)
}

func Jadianime(image []byte) ([]byte, error) {
	imageData := base64.StdEncoding.EncodeToString(image)
	data := map[string]interface{}{
		"busiId": "different_dimension_me_img_entry",
		"extra": map[string]interface{}{
			"face_rects": []interface{}{},
			"version":    2,
			"platform":   "web",
			"data_report": map[string]interface{}{
				"parent_trace_id": RandomUUID(),
				"root_channel":    "",
				"level":           0,
			},
		},
		"images": []string{imageData},
	}
	dataStr, _ := json.Marshal(data)
	encodedData := url.QueryEscape(string(dataStr))
	encodedData = strings.Replace(encodedData, "%7B%7D", "{}", -1)
	encodedData = strings.Replace(encodedData, "+", "%20", -1)
	encodedData = strings.Replace(encodedData, "*", "%2A", -1)
	encodedData = strings.Replace(encodedData, "%7E", "~", -1)
	signStr := baseURL + encodedData + "HQ31X02e"
	h := md5.New()
	h.Write([]byte(signStr))
	sign := hex.EncodeToString(h.Sum(nil))
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://h5.tu.qq.com")
	headers.Set("Referer", "https://h5.tu.qq.com/")
	headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	headers.Set("x-sign-value", sign)
	body := bytes.NewReader(dataStr)
	req, _ := http.NewRequest("POST", "https://ai.tu.qq.com/trpc.shadow_cv.ai_processor_cgi.AIProcessorCgi/Process", body)
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody := bytes.Buffer{}
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}
