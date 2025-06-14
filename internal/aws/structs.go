package aws

type WebGL struct {
	WebGLUnmaskedVendor string `json:"webgl_unmasked_vendor"`
	WebGLExtensions     string `json:"webgl_extensions"`
}
type GPUInfo struct {
	WebGLRenderer string  `json:"web_gl_renderer"`
	WebGL         []WebGL `json:"webgl"`
}
type Metrics struct {
	Fp2          int `json:"fp2"`
	Browser      int `json:"browser"`
	Capabilities int `json:"capabilities"`
	GPU          int `json:"gpu"`
	DNT          int `json:"dnt"`
	Math         int `json:"math"`
	Screen       int `json:"screen"`
	Navigator    int `json:"navigator"`
	Auto         int `json:"auto"`
	Stealth      int `json:"stealth"`
	Subtle       int `json:"subtle"`
	Canvas       int `json:"canvas"`
	FormDetector int `json:"formdetector"`
	BE           int `json:"be"`
}
type Fingerprint struct {
	Metrics      Metrics       `json:"metrics"`
	Start        int64         `json:"start"`
	FlashVersion interface{}   `json:"flashVersion"`
	Plugins      []Plugin      `json:"plugins"`
	DupedPlugins string        `json:"dupedPlugins"`
	ScreenInfo   string        `json:"screenInfo"`
	Referrer     string        `json:"referrer"`
	UserAgent    string        `json:"userAgent"`
	Location     string        `json:"location"`
	WebDriver    bool          `json:"webDriver"`
	Capabilities Capabilities  `json:"capabilities"`
	GPU          GPUBlock      `json:"gpu"`
	DNT          interface{}   `json:"dnt"`
	Math         MathBlock     `json:"math"`
	Automation   Automation    `json:"automation"`
	Stealth      Stealth       `json:"stealth"`
	Crypto       CryptoBlock   `json:"crypto"`
	Canvas       CanvasBlock   `json:"canvas"`
	FormDetected bool          `json:"formDetected"`
	NumForms     int           `json:"numForms"`
	NumFormElems int           `json:"numFormElements"`
	BE           BEBlock       `json:"be"`
	End          int64         `json:"end"`
	Errors       []interface{} `json:"errors"`
	Version      string        `json:"version"`
	ID           string        `json:"id"`
}

type Plugin struct {
	Name string `json:"name"`
	Str  string `json:"str"`
}
type CSSCapabilities struct {
	TextShadow       int `json:"textShadow"`
	WebkitTextStroke int `json:"WebkitTextStroke"`
	BoxShadow        int `json:"boxShadow"`
	BorderRadius     int `json:"borderRadius"`
	BorderImage      int `json:"borderImage"`
	Opacity          int `json:"opacity"`
	Transform        int `json:"transform"`
	Transition       int `json:"transition"`
}

type JSCapabilities struct {
	Audio        bool   `json:"audio"`
	Geolocation  bool   `json:"geolocation"`
	LocalStorage string `json:"localStorage"`
	Touch        bool   `json:"touch"`
	Video        bool   `json:"video"`
	WebWorker    bool   `json:"webWorker"`
}

type Capabilities struct {
	CSS     CSSCapabilities `json:"css"`
	JS      JSCapabilities  `json:"js"`
	Elapsed int             `json:"elapsed"`
}

type GPUBlock struct {
	Vendor     string   `json:"vendor"`
	Model      string   `json:"model"`
	Extensions []string `json:"extensions"`
}

type MathBlock struct {
	Tan string `json:"tan"`
	Sin string `json:"sin"`
	Cos string `json:"cos"`
}

type AutomationProperties struct {
	Document  []string `json:"document"`
	Window    []string `json:"window"`
	Navigator []string `json:"navigator"`
}
type AutomationWD struct {
	Properties AutomationProperties `json:"properties"`
}

type AutomationPhantom struct {
	Properties PhantomProperties `json:"properties"`
}

type PhantomProperties struct {
	Window []string `json:"window"`
}

type Automation struct {
	WD      AutomationWD      `json:"wd"`
	Phantom AutomationPhantom `json:"phantom"`
}

type Stealth struct {
	T1  int  `json:"t1"`
	T2  int  `json:"t2"`
	I   int  `json:"i"`
	MTE int  `json:"mte"`
	MTD bool `json:"mtd"`
}

type CryptoBlock struct {
	Crypto        int  `json:"crypto"`
	Subtle        int  `json:"subtle"`
	Encrypt       bool `json:"encrypt"`
	Decrypt       bool `json:"decrypt"`
	WrapKey       bool `json:"wrapKey"`
	UnwrapKey     bool `json:"unwrapKey"`
	Sign          bool `json:"sign"`
	Verify        bool `json:"verify"`
	Digest        bool `json:"digest"`
	DeriveBits    bool `json:"deriveBits"`
	DeriveKey     bool `json:"deriveKey"`
	GetRandomVals bool `json:"getRandomValues"`
	RandomUUID    bool `json:"randomUUID"`
}

type CanvasBlock struct {
	Hash          int         `json:"hash"`
	EmailHash     interface{} `json:"emailHash"`
	HistogramBins []int       `json:"histogramBins"`
}

type BEBlock struct {
	SI bool `json:"si"`
}
type Challenge struct {
	Input  string `json:"input"`
	Hmac   string `json:"hmac"`
	Region string `json:"region"`
}
type Inputs struct {
	Challenge     Challenge `json:"challenge"`
	ChallengeType string    `json:"challenge_type"`
	Difficulty    int       `json:"difficulty"`
}
type VerifyMetrics struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type ValueVerifySignals struct {
	Present string `json:"Present"`
}
type VerifySignals struct {
	Name  string             `json:"name"`
	Value ValueVerifySignals `json:"value"`
}

type GokuProps struct {
	Key     string `json:"key"`
	Iv      string `json:"iv"`
	Context string `json:"context"`
}
type Verify struct {
	Challenge     Challenge       `json:"challenge"`
	Solution      string          `json:"solution"`
	Signals       []VerifySignals `json:"signals"`
	Checksum      string          `json:"checksum"`
	ExistingToken []string        `json:"existing_token"`
	Client        string          `json:"client"`
	Domain        string          `json:"domain"`
	Metrics       []VerifyMetrics `json:"metrics"`
	GokuProps     GokuProps       `json:"goku_props"`
}

type VerifyRes struct {
	Token  string      `json:"token"`
	Inputs interface{} `json:"inputs"`
}
