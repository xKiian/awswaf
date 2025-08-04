package captcha

import (
	"awswaf/internal/aws"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"log"
	"time"
)

type WafCaptcha struct {
	session       tlsclient.HttpClient
	gokuProps     aws.GokuProps
	host          string
	hostToken     string
	domain        string
	userAgent     string
	existingToken string
}

func NewAwsWafCaptcha(
	waf *aws.Waf, gokuProps aws.GokuProps, host string, existingToken string,
) *WafCaptcha {
	return &WafCaptcha{
		session:       waf.Session,
		gokuProps:     gokuProps,
		host:          host,
		domain:        waf.Domain,
		userAgent:     waf.UserAgent,
		existingToken: existingToken,
		hostToken:     waf.Host,
	}
}

func (c *WafCaptcha) GetCaptcha() (ProblemRes ProblemResponse, err error) {
	url := fmt.Sprintf("https://%s/problem?kind=visual&domain=%s&locale=en", c.host, c.domain)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = http.Header{
		"sec-ch-ua-platform": {"\"Windows\""},
		"user-agent":         {c.userAgent},
		"sec-ch-ua":          {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://huggingface.co"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://huggingface.co/"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey:  {"sec-ch-ua-platform", "user-agent", "sec-ch-ua", "sec-ch-ua-mobile", "accept", "origin", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-dest", "referer", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err := c.session.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
	err = json.NewDecoder(resp.Body).Decode(&ProblemRes)
	return
}

func (c *WafCaptcha) Verify(res ProblemResponse, solution []int, elapsed int) (verifyRes VerifyRes, err error) {
	url := fmt.Sprintf("https://%s/verify", c.host)
	
	body := VerifyS{
		State:          res.State,
		Key:            res.Key,
		HmacTag:        res.HmacTag,
		ClientSolution: solution,
		Metrics: struct {
			SolveTimeMillis int `json:"solve_time_millis"`
		}{SolveTimeMillis: elapsed},
		GokuProps: c.gokuProps,
		Locale:    "en",
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return
	}
	
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encoded))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = http.Header{
		"sec-ch-ua-platform": {"\"Windows\""},
		"user-agent":         {c.userAgent},
		"sec-ch-ua":          {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"content-type":       {"text/plain;charset=UTF-8"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://huggingface.co"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://huggingface.co/"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey:  {"content-length", "sec-ch-ua-platform", "user-agent", "sec-ch-ua", "content-type", "sec-ch-ua-mobile", "accept", "origin", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-dest", "referer", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err := c.session.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
	err = json.NewDecoder(resp.Body).Decode(&verifyRes)
	
	return
}

func (c *WafCaptcha) VerifyVoucher(voucher string) (verifyRes aws.VerifyRes, err error) {
	url := fmt.Sprintf("https://%s/voucher", c.hostToken)
	
	body := VoucherS{
		CaptchaVoucher: voucher,
		ExistingToken:  c.existingToken,
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return
	}
	
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encoded))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = http.Header{
		"sec-ch-ua-platform": {"\"Windows\""},
		"user-agent":         {c.userAgent},
		"sec-ch-ua":          {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"content-type":       {"text/plain;charset=UTF-8"},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"origin":             {"https://huggingface.co"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://huggingface.co/"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":           {"u=1, i"},
		http.HeaderOrderKey:  {"content-length", "sec-ch-ua-platform", "user-agent", "sec-ch-ua", "content-type", "sec-ch-ua-mobile", "accept", "origin", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-dest", "referer", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err := c.session.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
	err = json.NewDecoder(resp.Body).Decode(&verifyRes)
	
	return
}

func (c *WafCaptcha) Run() (token string, err error) {
	problemRes, err := c.GetCaptcha()
	if err != nil {
		return
	}
	images, err := JSONStringToSlice[string](problemRes.Assets.Images)
	if err != nil {
		return
	}
	start := time.Now()
	solution, err := SolveImage(images, problemRes.LocalizedAssets.Target0)
	if err != nil {
		return
	}
	end := time.Since(start).Milliseconds()
	res, err := c.Verify(problemRes, solution, int(end))
	if err != nil {
		return
	}
	if !res.Success {
		err = errors.New("gemini is retarded and couldn't solve the captcha")
		return
	}
	fmt.Println("[+] Received Voucher", res.CaptchaVoucher[:100])
	
	tokenRes, err := c.VerifyVoucher(res.CaptchaVoucher)
	if err != nil {
		return
	}
	token = tokenRes.Token
	fmt.Println("[+] Received Token", token)
	return
}
