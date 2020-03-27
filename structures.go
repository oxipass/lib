package bslib

// BSItem - item structure
type BSItem struct {
	ID      int64     `json:"item_id"`
	Name    string    `json:"item_name"`
	Icon    string    `json:"item_icon"`
	Created string    `json:"created"`
	Updated string    `json:"updated"`
	Deleted bool      `json:"deleted"`
	Fields  []BSField `json:"fields"`
}

// BSField - fields definitions
type BSField struct {
	ID        int64  `json:"field_id"`
	Name      string `json:"field_name"`
	Icon      string `json:"field_icon"`
	ValueType string `json:"value_type"`
	Value     string `json:"field_value"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Deleted   bool   `json:"deleted"`
}

// CommonResponse - common response header structure
type CommonResponse struct {
	Status string `json:"status"`
	MsgNum string `json:"msg_num"`
	MsgTxt string `json:"msg_text"`
}

// ItemResponse - response returning one item
type ItemResponse struct {
	CommonResponse
	BSItem
}

// ItemsResponse - response returning many items
type ItemsResponse struct {
	CommonResponse
	Items []BSItem `json:"items"`
}

// UpdateFieldForm input structure to add or update the field
type UpdateFieldForm struct {
	ItemID int64 `json:"item_id"`
	BSField
}

// FieldAddedResponse - response structure for adding field
type FieldAddedResponse struct {
	CommonResponse
	FieldID int64 `json:"field_id"`
}

// UpdateItemForm - input structure to add the item
type UpdateItemForm struct {
	Action string
	BSItem
}

// ItemAddedResponse - response structure for adding item
type ItemAddedResponse struct {
	CommonResponse
	ItemID int64 `json:"item_id"`
}

// InitStorageForm - initializing the database
type InitStorageForm struct {
	FileName   string `json:"filename"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"`
}

// ReadAllForm parameters for reading from database
type ReadAllForm struct {
	ReadDeleted bool `json:"read_deleted"`
}
