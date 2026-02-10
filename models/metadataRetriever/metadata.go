package metadataRetriever

import "time"

import "encoding/json"

func UnmarshalMetadata(data []byte) (Metadata, error) {
	var r Metadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Metadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Metadata struct {
	Status   int64         `json:"status"`
	Result   []Result      `json:"result"`
	Warnings []interface{} `json:"warnings"`
}

type Result struct {
	CreatedByID        string           `json:"createdById"`
	CreatedByName      string         	`json:"createdByName"`
	CreatedDate        time.Time        `json:"createdDate"`
	FileName           string           `json:"fileName"`
	FullName           string           `json:"fullName"`
	ID                 string           `json:"id"`
	LastModifiedByID   string           `json:"lastModifiedById"`
	LastModifiedByName string         	`json:"lastModifiedByName"`
	LastModifiedDate   time.Time        `json:"lastModifiedDate"`
	ManageableState    ManageableState  `json:"manageableState"`
	Type               string             `json:"type"`
	NamespacePrefix    *string `json:"namespacePrefix,omitempty"`
}

type ManageableState string

const (
	Installed ManageableState = "installed"
	Unmanaged ManageableState = "unmanaged"
)

// type MetadataType string
// const (
// 	ApexClass MetadataType = "ApexClass"
// )
