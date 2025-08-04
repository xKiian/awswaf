package captcha

import "awswaf/internal/aws"

type State struct {
	Iv      string `json:"iv"`
	Payload string `json:"payload"`
}

type ProblemResponse struct {
	ProblemType string `json:"problem_type"`
	State       State  `json:"state"`
	Key         string `json:"key"`
	HmacTag     string `json:"hmac_tag"`
	Assets      struct {
		Target string `json:"target"`
		Images string `json:"images"`
		Res    string `json:"res"`
	} `json:"assets"`
	LocalizedAssets struct {
		Target0 string `json:"target0"`
	} `json:"localized_assets"`
	ValidityDuration int `json:"validity_duration"`
}

type VerifyS struct {
	State          State  `json:"state"`
	Key            string `json:"key"`
	HmacTag        string `json:"hmac_tag"`
	ClientSolution []int  `json:"client_solution"`
	Metrics        struct {
		SolveTimeMillis int `json:"solve_time_millis"`
	} `json:"metrics"`
	GokuProps aws.GokuProps `json:"goku_props"`
	Locale    string        `json:"locale"`
}

type VerifyRes struct {
	Success              bool        `json:"success"`
	CaptchaVoucher       string      `json:"captcha_voucher"`
	NumSolutionsRequired int         `json:"num_solutions_required"`
	NumSolutionsProvided int         `json:"num_solutions_provided"`
	Reason               interface{} `json:"reason"`
	Problem              interface{} `json:"problem"`
}

type VoucherS struct {
	CaptchaVoucher string `json:"captcha_voucher"`
	ExistingToken  string `json:"existing_token"`
}
