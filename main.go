package main

import (
	"awswaf/internal/aws"
	"awswaf/internal/aws/captcha"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"log"
	"net/url"
	"time"
)

func solveHuggingFace() {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
	}
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		panic(err)
	}
	
	req, err := http.NewRequest(http.MethodGet, "https://huggingface.com/join", nil)
	if err != nil {
		panic(err)
	}
	
	req.Header = http.Header{
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-ch-ua":                 {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey:         {"upgrade-insecure-requests", "user-agent", "accept", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	html := string(body)
	
	goku, host, err := aws.Extract(html)
	if err != nil {
		panic(err)
	}
	waf, err := aws.NewAwsWaf(
		host,
		"huggingface.co",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		goku, "",
	)
	if err != nil {
		panic(err)
	}
	
	token, err := waf.Run()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("[+] Solved Token %s\n", token[len(token)-100:])
	
	parsed, _ := url.Parse("https://huggingface.co/")
	cookie := &http.Cookie{
		Name:     "aws-waf-token",
		Value:    token,
		Domain:   "huggingface.co",
		Path:     "/",
		HttpOnly: true,
	}
	client.SetCookies(parsed, []*http.Cookie{cookie})
	
	req, err = http.NewRequest(http.MethodGet, "https://huggingface.co/join", nil)
	if err != nil {
		panic(err)
	}
	
	req.Header = http.Header{
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-ch-ua":                 {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey:         {"upgrade-insecure-requests", "user-agent", "accept", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	goku, host, err = aws.ExtractCaptcha(string(body))
	if err != nil {
		panic(err)
	}
	
	wafCaptcha := captcha.NewAwsWafCaptcha(waf, goku, host, token)
	token2, err := wafCaptcha.Run()
	cookie = &http.Cookie{
		Name:     "aws-waf-token",
		Value:    token2,
		Domain:   "huggingface.co",
		Path:     "/",
		HttpOnly: true,
	}
	client.SetCookies(parsed, []*http.Cookie{cookie})
	
	req, err = http.NewRequest(http.MethodGet, "https://huggingface.co/join", nil)
	if err != nil {
		panic(err)
	}
	
	req.Header = http.Header{
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-ch-ua":                 {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey:         {"upgrade-insecure-requests", "user-agent", "accept", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language", "priority"},
	}
	
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func solveBinance(proxy string) {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
		tlsclient.WithProxyUrl(proxy),
		tlsclient.WithInsecureSkipVerify(),
	}
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodGet, "https://www.binance.com/", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cache-control":             {"no-cache"},
		"pragma":                    {"no-cache"},
		"priority":                  {"u=0, i"},
		"sec-ch-ua":                 {`"Google Chrome";v="138", "Chromium";v="138", "Not/A)Brand";v="24"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"same-origin"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"accept-encoding",
			"cache-control",
			"pragma",
			"priority",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	
	gokuProps, host, err := aws.Extract(string(body))
	if err != nil {
		log.Println(err)
		return
	}
	
	waf, err := aws.NewAwsWaf(
		host,
		"www.binance.com",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		gokuProps, proxy,
	)
	if err != nil {
		log.Println(err)
		return
	}
	
	start := time.Now()
	
	token, err := waf.Run()
	if err != nil {
		log.Println(err)
		return
	}
	
	end := time.Now()
	
	parsed, _ := url.Parse("https://www.binance.com/")
	cookie := &http.Cookie{
		Name:     "aws-waf-token",
		Value:    token,
		Domain:   "www.binance.com",
		Path:     "/",
		HttpOnly: true,
	}
	client.SetCookies(parsed, []*http.Cookie{cookie})
	
	req, err = http.NewRequest(http.MethodGet, "https://www.binance.com/", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cache-control":             {"no-cache"},
		"pragma":                    {"no-cache"},
		"priority":                  {"u=0, i"},
		"sec-ch-ua":                 {`"Google Chrome";v="138", "Chromium";v="138", "Not/A)Brand";v="24"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"same-origin"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
	}
	
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	
	if len(string(body)) > 5000 {
		fmt.Printf("[+] Solved! %s in %s\n", token[len(token)-100:], end.Sub(start).String())
	} else {
		fmt.Println("[-] Failed to solve!")
	}
}

func main() {
	solveHuggingFace()
}
