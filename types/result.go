package types

type GenerationSettings struct {
	FrequencyPenalty float64       `json:"frequency_penalty"`
	Grammar          string        `json:"grammar"`
	IgnoreEOS        bool          `json:"ignore_eos"`
	LogitBias        []interface{} `json:"logit_bias"`
	Mirostat         int           `json:"mirostat"`
	MirostatEta      float64       `json:"mirostat_eta"`
	MirostatTau      float64       `json:"mirostat_tau"`
	Model            string        `json:"model"`
	NCtx             int           `json:"n_ctx"`
	NKeep            int           `json:"n_keep"`
	NPredict         int           `json:"n_predict"`
	NProbs           int           `json:"n_probs"`
	PenalizeNL       bool          `json:"penalize_nl"`
	PresencePenalty  float64       `json:"presence_penalty"`
	RepeatLastN      int           `json:"repeat_last_n"`
	RepeatPenalty    float64       `json:"repeat_penalty"`
	Seed             uint32        `json:"seed"`
	Stop             []interface{} `json:"stop"`
	Stream           bool          `json:"stream"`
	Temp             float64       `json:"temp"`
	TfsZ             float64       `json:"tfs_z"`
	TopK             int           `json:"top_k"`
	TopP             float64       `json:"top_p"`
	TypicalP         float64       `json:"typical_p"`
}

type Timings struct {
	PredictedMs         float64 `json:"predicted_ms"`
	PredictedN          int     `json:"predicted_n"`
	PredictedPerSecond  float64 `json:"predicted_per_second"`
	PredictedPerTokenMs float64 `json:"predicted_per_token_ms"`
	PromptMs            float64 `json:"prompt_ms"`
	PromptN             int     `json:"prompt_n"`
	PromptPerSecond     float64 `json:"prompt_per_second"`
	PromptPerTokenMs    float64 `json:"prompt_per_token_ms"`
}

type Result struct {
	Content            string             `json:"content"`
	GenerationSettings GenerationSettings `json:"generation_settings"`
	Model              string             `json:"model"`
	Prompt             string             `json:"prompt"`
	SlotID             int                `json:"slot_id"`
	Stop               bool               `json:"stop"`
	StoppedEOS         bool               `json:"stopped_eos"`
	StoppedLimit       bool               `json:"stopped_limit"`
	StoppedWord        bool               `json:"stopped_word"`
	StoppingWord       string             `json:"stopping_word"`
	Timings            Timings            `json:"timings"`
	TokensCached       int                `json:"tokens_cached"`
	TokensEvaluated    int                `json:"tokens_evaluated"`
	TokensPredicted    int                `json:"tokens_predicted"`
	Truncated          bool               `json:"truncated"`
}
