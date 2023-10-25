package types

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
)

type Model_Request struct {
	Model       string `json:"model" llama:"model" default:""`
	ModelName   string `json:"modelName,omitempty" llama:"alias" default:""`
	ContextSize int    `json:"contextSize,omitempty" llama:"ctx-size" default:"4096"`

	NBatch        int     `json:"nBatch,omitempty" llama:"batch-size" default:"512"`
	F32Memory     bool    `json:"f16Memory,omitempty" llama:"memory-f32" default:"false"` //should be false
	MLock         bool    `json:"mLock,omitempty" llama:"mlock" default:"false"`
	NoMMap        bool    `json:"noMMap,omitempty" llama:"no-mmap" default:"false"`
	Embeddings    bool    `json:"embeddings,omitempty" llama:"embedding" default:"false"`
	NUMA          bool    `json:"numa,omitempty" llama:"numa" default:"false"`
	NGPULayers    int     `json:"nGpuLayers,omitempty" llama:"n-gpu-layers" default:"0"`
	MainGPU       string  `json:"mainGpu,omitempty" llama:"main-gpu" default:""`
	TensorSplit   string  `json:"tensorSplit,omitempty" llama:"tensor-split" default:""`
	FreqRopeBase  float32 `json:"freqRopeBase,omitempty" llama:"rope-freq-base" default:"1000.0"`
	FreqRopeScale float32 `json:"freqRopeScale,omitempty" llama:"rope-freq-scale" default:"1.0"`

	LoraBase    string `json:"loraBase,omitempty" llama:"lora-base" default:""`
	LoraAdapter string `json:"loraAdapter,omitempty" llama:"lora" default:""`

	ParallelSlots   int    `json:"parallelSlots,omitempty" llama:"parallel" default:"1"` //defaults to 1
	Port            int    `json:"port,omitempty" llama:"port" default:"8080"`           //defaults to 8080
	Host            string `json:"host,omitempty" llama:"host" default:"localhost"`      //defaults to localhost
	SystemPromtFile string `json:"systemPromptFile,omitempty" llama:"system-prompt-file" default:""`
}

func NewModelRequestWithDefaults() *Model_Request {
	mr := &Model_Request{}
	t := reflect.TypeOf(mr)
	v := reflect.ValueOf(&mr).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		def := field.Tag.Get("default")
		fv := v.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			fv.SetString(def)
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
	return mr
}
func (mr *Model_Request) ToMap() map[string]string {
	m := make(map[string]string)
	t := reflect.TypeOf(*mr)
	v := reflect.ValueOf(*mr)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		def := field.Tag.Get("default")
		llama := field.Tag.Get("llama")
		key := "--" + llama
		fv := v.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			val := fv.String()
			if val != def {
				m[key] = val
			}
		case reflect.Int, reflect.Int64:
			val := fv.Int()
			if defVal, err := strconv.ParseInt(def, 10, 64); err == nil && val != defVal {
				m[key] = strconv.FormatInt(val, 10)
			}
		case reflect.Float32:
			val := fv.Float()
			if defVal, err := strconv.ParseFloat(def, 32); err == nil && float32(val) != float32(defVal) {
				m[key] = strconv.FormatFloat(val, 'f', 2, 32)
			}
		case reflect.Bool:
			val := fv.Bool()
			if defVal, err := strconv.ParseBool(def); err == nil && val != defVal && val {
				m[key] = ""
			}
		}
	}
	return m
}
func (m *Model_Request) Save(ModelPath string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(ModelPath+"/"+m.ModelName+".json", data, 0644)
}

func (m *Model_Request) Load(ModelPath string, name string) error {
	data, err := os.ReadFile(ModelPath + "/" + name + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, m)
}

func (m *Model_Request) Delete(ModelPath string) error {
	err := os.Remove(ModelPath + "/" + m.ModelName + ".json")
	return err
}
