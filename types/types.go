package types

type Response struct {
	Success bool              `json:"success" description:"Indicates if the request was successful"`
	Context map[string]string `json:"context,omitempty" description:"Context of the response"`
	Message *string           `json:"message,omitempty" description:"Message of the response"`
	JSON    any               `json:"json,omitempty" description:"JSON data of the response"`
}
