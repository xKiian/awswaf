package main

import (
	"awswaf/internal/aws"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"log"
	"net/url"
	"strings"
)

func solveHuggingFace(client tlsclient.HttpClient) {
	
	resp, err := client.Get("https://huggingface.co/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	html := string(body)
	idx := strings.Index(html, `src="https://`)
	if idx == -1 {
		panic("couldn't find idx")
	}
	idx += 13
	tail := html[idx:]
	host := tail[:strings.Index(tail, "/challenge.js")]
	
	waf, err := aws.NewAwsWaf(
		host,
		"huggingface.co",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		aws.GokuProps{},
	)
	if err != nil {
		panic(err)
	}
	
	token, err := waf.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	
}

func solveBinance(client tlsclient.HttpClient) {
	req, err := http.NewRequest(http.MethodGet, "https://www.binance.com/", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("accept", "* /*")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	
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
	fmt.Println(string(body))
	
	gokuProps, host, err := aws.Extract(string(body))
	if err != nil {
		panic(err)
	}
	
	waf, err := aws.NewAwsWaf(
		host,
		"huggingface.co",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		gokuProps,
	)
	if err != nil {
		panic(err)
	}
	
	token, err := waf.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	
	parsed, _ := url.Parse("https://huggingface.co")
	cookie := &http.Cookie{
		Name:     "aws-waf-token",
		Value:    token,
		Domain:   "huggingface.co",
		Path:     "/",
		HttpOnly: true,
	}
	client.SetCookies(parsed, []*http.Cookie{cookie})
	
	resp, err = client.Get("https://www.binance.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(string(body)) > 5000)
}

func main() {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
	}
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		panic(err)
	}
	for range 2 {
		go func() {
			for {
				solveBinance(client)
			}
		}()
	}
	select {}
}
