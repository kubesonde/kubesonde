package types

type NestatInfoRequestBody []NestatInfoRequestBodyItem
type NestatInfoRequestBodyItem struct {
	Fd     int           `json:"fd"`
	Family int           `json:"family"`
	Type   int           `json:"type"`
	Laddr  []interface{} `json:"laddr"`
	Raddr  []interface{} `json:"raddr"`
	Status string        `json:"status"`
	Pid    interface{}   `json:"pid"`
}
