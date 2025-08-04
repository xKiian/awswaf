package aws

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"hash/crc32"
	"math/rand"
	"os"
	"strings"
	"time"
)

var gpus []GPUInfo

func init() {
	file, err := os.Open("webgl.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gpus)
	if err != nil {
		panic(err)
	}
}

func encodeWithCRC(obj any) ([]byte, []byte, error) {
	payload, err := json.Marshal(obj)
	if err != nil {
		return nil, nil, err
	}
	crc := crc32.ChecksumIEEE(payload)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, crc)
	
	checksum := []byte(strings.ToUpper(hex.EncodeToString(buf)))
	result := append(checksum, '#')
	result = append(result, payload...)
	
	return checksum, result, nil
}

func GetFP(userAgent string, ) (string, string, error) {
	ts := time.Now().UnixMilli()
	bins := make([]int, 256)
	for i := range bins {
		bins[i] = rand.Intn(40)
	}
	bins[0] = rand.Intn(2100) + 14473
	bins[255] = rand.Intn(2100) + 14473
	
	gpu := gpus[rand.Intn(len(gpus))]
	
	fp := Fingerprint{
		Metrics: Metrics{
			Fp2: 1, Browser: 0, Capabilities: 1, GPU: 7, DNT: 0, Math: 0, Screen: 0,
			Navigator: 0, Auto: 1, Stealth: 0, Subtle: 0, Canvas: 113, FormDetector: 1, BE: 0,
		},
		Start:        ts,
		FlashVersion: nil,
		Plugins: []Plugin{
			{"PDF Viewer", "PDF Viewer "},
			{"Chrome PDF Viewer", "Chrome PDF Viewer "},
			{"Chromium PDF Viewer", "Chromium PDF Viewer "},
			{"Microsoft Edge PDF Viewer", "Microsoft Edge PDF Viewer "},
			{"WebKit built-in PDF", "WebKit built-in PDF "},
		},
		DupedPlugins: "PDF Viewer Chrome PDF Viewer Chromium PDF Viewer Microsoft Edge PDF Viewer WebKit built-in PDF ||1920-1080-1032-24-*-*-*",
		ScreenInfo:   "1920-1080-1032-24-*-*-*",
		Referrer:     "",
		UserAgent:    userAgent,
		Location:     "", //they probably don't check this
		WebDriver:    false,
		Capabilities: Capabilities{
			CSS: CSSCapabilities{
				TextShadow: 1, WebkitTextStroke: 1, BoxShadow: 1, BorderRadius: 1,
				BorderImage: 1, Opacity: 1, Transform: 1, Transition: 1,
			},
			JS: JSCapabilities{
				Audio: true, Geolocation: rand.Intn(2) == 1,
				LocalStorage: "supported", Touch: false, Video: true,
				WebWorker: rand.Intn(2) == 1,
			},
			Elapsed: 1,
		},
		GPU: GPUBlock{
			Vendor:     gpu.WebGL[0].WebGLUnmaskedVendor,
			Model:      gpu.WebGLRenderer,
			Extensions: strings.Split(gpu.WebGL[0].WebGLExtensions, ";"),
		},
		DNT: nil,
		Math: MathBlock{
			Tan: "-1.4214488238747245", Sin: "0.8178819121159085", Cos: "-0.5753861119575491",
		},
		Automation: Automation{
			WD: AutomationWD{
				Properties: AutomationProperties{
					Document:  []string{},
					Window:    []string{},
					Navigator: []string{},
				},
			},
			Phantom: AutomationPhantom{
				Properties: PhantomProperties{
					Window: []string{},
				},
			},
		},
		Stealth: Stealth{T1: 0, T2: 0, I: 1, MTE: 0, MTD: false},
		Crypto: CryptoBlock{
			Crypto: 1, Subtle: 1, Encrypt: true, Decrypt: true, WrapKey: true,
			UnwrapKey: true, Sign: true, Verify: true, Digest: true,
			DeriveBits: true, DeriveKey: true, GetRandomVals: true, RandomUUID: true,
		},
		Canvas: CanvasBlock{
			Hash:          rand.Intn(90020000) + 645172295,
			EmailHash:     nil,
			HistogramBins: bins,
		},
		FormDetected: false,
		NumForms:     0,
		NumFormElems: 0,
		BE:           BEBlock{SI: false},
		End:          ts + 1,
		Errors:       []any{},
		Version:      "2.4.0",
		ID:           uuid.New().String(),
	}
	
	checksum, payload, err := encodeWithCRC(fp)
	if err != nil {
		return "", "", err
	}
	encrypted, err := Encrypt(payload)
	return string(checksum), encrypted, err
}
