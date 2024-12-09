package entity

type Route struct {
	Path         string `json:"path"`
	Method       string `json:"method"`
	ResponseFile string `json:"responseFile"`
	HttpStatus   uint32 `json:"httpStatus"`
	Delay        uint32 `json:"delay"`
}
