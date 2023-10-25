package types

type Result struct {
	ContentType        string              `json:"content_type"`
	Stop               bool                `json:"stop"`
	GenerationSettings *Prediction_Request `json:"generation_settings"`
	Model              string              `json:"model"`
	Prompt             string              `json:"prompt"`
	StoppedEOS         bool                `json:"stopped_eos"`
	StoppedLimit       bool                `json:"stopped_limit"`
	StoppedWord        bool                `json:"stopped_word"`
	StoppingWord       string              `json:"stopping_word"`
	//timings ==> unkown type
	Timings         map[string]interface{} `json:"timings"`
	TokensCached    int                    `json:"tokens_cached"`
	TokensEvaluated int                    `json:"tokens_evaluated"`
	Truncated       bool                   `json:"truncated"`
}
