package aws

import (
	"encoding/json"
	"fmt"
	"github.com/bogdanfinn/tls-client/profiles"
	"log"
	"math/rand"
	"strings"
	
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
)

type Waf struct {
	Session   tlsclient.HttpClient
	gokuProps GokuProps
	Host      string
	Domain    string
	UserAgent string
}

func NewAwsWaf(
	host, domain, userAgent string,
	gokuProps GokuProps, proxy string,
) (*Waf, error) {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(5),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
	}
	if proxy != "" {
		options = append(options, tlsclient.WithProxyUrl(proxy))
	}
	//options = append(options, tlsclient.WithCharlesProxy("", "53961"))
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}
	
	return &Waf{
		Session:   client,
		gokuProps: gokuProps,
		Host:      host,
		Domain:    domain,
		UserAgent: userAgent,
	}, nil
}

func Extract(html string) (GokuProps, string, error) {
	const marker = "window.gokuProps = "
	start := strings.Index(html, marker)
	if start == -1 {
		return GokuProps{}, "", fmt.Errorf("gokuProps not found")
	}
	
	start += len(marker)
	end := strings.Index(html[start:], ";")
	if end == -1 {
		return GokuProps{}, "", fmt.Errorf("end of gokuProps not found")
	}
	
	var gokuProps GokuProps
	if err := json.Unmarshal([]byte(html[start:start+end]), &gokuProps); err != nil {
		return GokuProps{}, "", err
	}
	
	idx := strings.Index(html, `src="https://`)
	if idx == -1 {
		return gokuProps, "", nil
	}
	idx += 13
	tail := html[idx:]
	host := tail[:strings.Index(tail, "/challenge.js")]
	
	return gokuProps, host, nil
}

func ExtractCaptcha(html string) (GokuProps, string, error) {
	const marker = "window.gokuProps = "
	start := strings.Index(html, marker)
	if start == -1 {
		return GokuProps{}, "", fmt.Errorf("gokuProps not found")
	}
	
	start += len(marker)
	end := strings.Index(html[start:], ";")
	if end == -1 {
		return GokuProps{}, "", fmt.Errorf("end of gokuProps not found")
	}
	
	var gokuProps GokuProps
	if err := json.Unmarshal([]byte(html[start:start+end]), &gokuProps); err != nil {
		return GokuProps{}, "", err
	}
	
	captchaMarker := `src="https://`
	captchaIndex := strings.Index(html, `/captcha.js`)
	if captchaIndex == -1 {
		return gokuProps, "", fmt.Errorf("captcha.js not found")
	}
	
	startQuote := strings.LastIndex(html[:captchaIndex], captchaMarker)
	if startQuote == -1 {
		return gokuProps, "", fmt.Errorf("captcha src not found")
	}
	
	startURL := startQuote + len(captchaMarker)
	host := html[startURL:captchaIndex]
	
	return gokuProps, host, nil
}

func (a *Waf) GetInputs() (Inputs, error) {
	url := fmt.Sprintf("https://%s/inputs?client=browser", a.Host)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return Inputs{}, err
	}
	req.Header = http.Header{
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"accept-language": {"en-US,en;q=0.9"},
		"pragma":          {"no-cache"},
		"priority":        {"u=0, i"},
		//"sec-ch-ua":          {`"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"cross-site"},
		"user-agent":         {a.UserAgent},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"accept-encoding",
			"pragma",
			"priority",
			//"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
		},
	}
	
	resp, err := a.Session.Do(req)
	if err != nil {
		log.Println(err)
		return Inputs{}, err
	}
	defer resp.Body.Close()
	
	var out Inputs
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Inputs{}, err
	}
	return out, nil
}

func (a *Waf) BuildPayload(inputs Inputs) (*Verify, error) {
	checksum, fpPayload, err := GetFP(a.UserAgent)
	if err != nil {
		return nil, err
	}
	
	sol, err := SolveChallenge(
		inputs.ChallengeType, inputs.Challenge.Input,
		checksum, inputs.Difficulty,
	)
	
	if err != nil {
		return nil, err
	}
	
	signals := []VerifySignals{{
		Name: "Zoey",
		Value: ValueVerifySignals{
			Present: fpPayload,
		},
	}}
	
	var metrics []VerifyMetrics
	for _, m := range []VerifyMetrics{
		{"2", rand.Float64(), "2"},
		{"100", 0, "2"}, {"101", 0, "2"},
		{"102", 0, "2"}, {"103", 8, "2"},
		{"104", 0, "2"}, {"105", 0, "2"},
		{"106", 0, "2"}, {"107", 0, "2"},
		{"108", 1, "2"}, {"undefined", 0, "2"},
		{"110", 0, "2"}, {"111", 2, "2"},
		{"112", 0, "2"}, {"undefined", 0, "2"},
		{"3", 4, "2"}, {"7", 0, "4"},
		{"1", rand.Float64()*(20-10) + 10, "2"},
		{"4", 36.5, "2"},
		{"5", rand.Float64(), "2"},
		{"6", rand.Float64()*(60-50) + 50, "2"},
		{"0", rand.Float64()*(140-130) + 130, "2"},
		{"8", 1, "4"},
	} {
		metrics = append(metrics, VerifyMetrics{
			Name:  m.Name,
			Value: m.Value,
			Unit:  m.Unit,
		})
	}
	
	return &Verify{
		Challenge:     inputs.Challenge,
		Solution:      sol,
		Signals:       signals,
		Checksum:      checksum,
		ExistingToken: nil,
		Client:        "Browser",
		Domain:        a.Domain,
		Metrics:       metrics,
		GokuProps:     a.gokuProps,
	}, nil
}

func (a *Waf) Verify(payload *Verify) (string, error) {
	url := fmt.Sprintf("https://%s/verify", a.Host)
	
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header = http.Header{
		"accept":          {"*/*"},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"connection":      {"keep-alive"},
		"accept-language": {"en-US,en;q=0.9"},
		"content-type":    {"text/plain;charset=UTF-8"},
		"priority":        {"u=1, i"},
		//"sec-ch-ua":          {`"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"cross-site"},
		"user-agent":         {a.UserAgent},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"connection",
			"content-length",
			"content-type",
			"priority",
			//"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
		},
	}
	
	resp, err := a.Session.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	
	defer resp.Body.Close()
	
	var out VerifyRes
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	
	return out.Token, nil
}

func (a *Waf) Run() (string, error) {
	inputs, err := a.GetInputs()
	if err != nil {
		return "", err
	}
	
	payload, err := a.BuildPayload(inputs)
	if err != nil {
		return "", err
	}
	return a.Verify(payload)
}
