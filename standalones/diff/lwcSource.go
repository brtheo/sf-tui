package diff

import "encoding/json"

func UnmarshalLwcSource(data []byte) (LwcSource, error) {
	var r LwcSource
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LwcSource) Marshal() ([]byte, error) {
	return json.Marshal(r)

}

type LwcSource struct {
	Status   int64         `json:"status"`
	Result   LwcResult        `json:"result"`
	Warnings []any         `json:"warnings"`
}

type LwcResult struct {
	Records   []LwcRecord       `json:"records"`
	TotalSize int64          `json:"totalSize"`
	Done      bool           `json:"done"`
}

type LwcRecord struct {
	Source   string         `json:"source"`
	Attributes LwcAttributes   `json:"attributes"`
}

type LwcAttributes struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}
