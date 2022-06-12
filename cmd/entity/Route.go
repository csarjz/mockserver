package entity

type Route struct {
	Path         string `json:"path"`
	Method       string `json:"method"`
	ResponseFile string `json:"responseFile"`
	Delay        uint32 `json:"delay"`
}
