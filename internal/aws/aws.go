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
	session   tlsclient.HttpClient
	gokuProps GokuProps
	host      string
	domain    string
	userAgent string
}

func NewAwsWaf(
	host, domain, userAgent string,
	gokuProps GokuProps,
) (*Waf, error) {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithClientProfile(profiles.Chrome_133),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
	}
	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}
	
	return &Waf{
		session:   client,
		gokuProps: gokuProps,
		host:      host,
		domain:    domain,
		userAgent: userAgent,
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

func (a *Waf) GetInputs() (Inputs, error) {
	url := fmt.Sprintf("https://%s/inputs?client=browser", a.host)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return Inputs{}, err
	}
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("user-agent", a.userAgent)
	req.Header.Set("sec-ch-ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("accept", "* /*")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	
	resp, err := a.session.Do(req)
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

func (a *Waf) BuildPayload(inputs Inputs) (map[string]interface{}, error) {
	checksum, fpPayload, err := GetFP(a.userAgent)
	if err != nil {
		return nil, err
	}
	fmt.Println("type", inputs.ChallengeType)
	sol, err := SolveChallenge(
		inputs.ChallengeType, inputs.Challenge.Input,
		checksum, inputs.Difficulty,
	)
	
	if err != nil {
		return nil, err
	}
	
	signals := []VerifySignals{{
		Name: "KramerAndRio",
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
	_ = Verify{
		Challenge:     inputs.Challenge,
		Solution:      sol,
		Signals:       signals,
		Checksum:      checksum,
		ExistingToken: "",
		Client:        "Browser",
		Domain:        a.domain,
		Metrics:       metrics,
		GokuProps:     a.gokuProps,
	}
	return map[string]interface{}{
		"challenge":      inputs.Challenge,
		"checksum":       checksum,
		"solution":       sol,
		"signals":        signals,
		"existing_token": nil,
		"client":         "Browser",
		"domain":         a.domain,
		"metrics":        metrics,
	}, nil
}

func (a *Waf) Verify(payload map[string]interface{}) (string, error) {
	url := fmt.Sprintf("https://%s/verify", a.host)
	
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("user-agent", a.userAgent)
	req.Header.Set("content-type", "text/plain;charset=UTF-8")
	req.Header.Set("sec-ch-ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("accept", "* /*")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	
	resp, err := a.session.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	
	defer resp.Body.Close()
	
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	
	token, _ := out["token"].(string)
	return token, nil
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
