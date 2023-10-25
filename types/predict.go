package types

import (
	"reflect"
	"strconv"
	"strings"
)

type SystemPrompt struct {
	Prompt        string `json:"prompt"`
	AntiPrompt    string `json:"anti_prompt"`
	AssistantName string `json:"assistant_name"`
}

type Prediction_Request struct {
	Prompt      string  `json:"prompt" default:""`
	InputPrefix string  `json:"input_prefix,omitempty" default:""`
	InputSuffix string  `json:"input_suffix,omitempty" default:""`
	Temperature float32 `json:"temperature,omitempty" default:"0.8"`
	TopK        int     `json:"top_k,omitempty" default:"40"`
	TopP        float32 `json:"top_p,omitempty" default:"0.95"`

	NPredict int `json:"n_predict,omitempty" default:"-1"` //-1 = infinity
	NKeep    int `json:"n_keep,omitempty" default:"0"`     // 0 none are kept, -1 all are kept

	Stream bool `json:"stream,omitempty" default:"false"`

	Stop []string `json:"stop,omitempty" default:""` //default=stringlist, commaseperated list of stop words

	Tfs_z            float32 `json:"tfs_z,omitempty" default:"1.0"`     //default=1.0 => disabled
	TypicalP         float32 `json:"typical_p,omitempty" default:"1.0"` //default=1.0 => disabled
	RepeatPenalty    float32 `json:"repeat_penalty,omitempty" default:"1.1"`
	RepeatLastN      int     `json:"repeat_last_n,omitempty" default:"64"`      // 0=> disabled, -1 = ctx-size
	PenalizeNl       bool    `json:"penalize_nl,omitempty" default:"true"`      //penalize new line
	PresencePenalty  float32 `json:"presence_penalty,omitempty" default:"0.0"`  //default=0.0 => disabled
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty" default:"0.0"` //default=0.0 => disabled

	Mirostat    int     `json:"mirostat,omitempty" default:"0"`       //default=0 => disabled, 1=enabled, 2=mirostat 2.0
	MirostatTAU float32 `json:"mirostat_tau,omitempty" default:"5.0"` //default=5.0
	MirostatETA float32 `json:"mirostat_eta,omitempty" default:"0.1"` //default=0.1

	Grammar     string        `json:"grammar,omitempty" default:""`
	Seed        int           `json:"seed,omitempty" default:"-1"`          //default=-1 => rng
	IgnoreEOS   bool          `json:"ignore_eos,omitempty" default:"false"` //default=false
	NProbs      int           `json:"n_probs,omitempty" default:"0"`
	SlotID      int           `json:"slot_id" default:"-1"`         //default=-1 => idle slot
	CachePrompt bool          `json:"cache_prompt" default:"false"` //default=false
	SystemPromt *SystemPrompt `json:"system_prompt"`

	Model string `json:"model,omitempty" default:""`  //not set for request
	NCtx  int    `json:"n_ctx,omitempty" default:"0"` //not set for request
}

func NewPredictionRequestWithDefaults() Prediction_Request {
	pr := Prediction_Request{}
	t := reflect.TypeOf(pr)
	v := reflect.ValueOf(&pr).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		def := field.Tag.Get("default")
		fv := v.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			fv.SetString(def)
		case reflect.Slice:
			// Assuming a slice of strings
			if def != "" {
				fv.Set(reflect.ValueOf(strings.Split(def, ",")))
			}
		case reflect.Int, reflect.Int64:
			if iVal, err := strconv.Atoi(def); err == nil {
				fv.SetInt(int64(iVal))
			}
		case reflect.Float32:
			if fVal, err := strconv.ParseFloat(def, 32); err == nil {
				fv.SetFloat(fVal)
			}
		case reflect.Bool:
			if bVal, err := strconv.ParseBool(def); err == nil {
				fv.SetBool(bVal)
			}
		}
	}
	return pr
}
