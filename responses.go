package smsman

type PhoneDetail struct {
	Country       string
	Code          string
	CountryId     interface{} `json:"country_id"`
	ApplicationId interface{} `json:"application_id"`
	Phone         string      `json:"number"`
	Taskid        int         `json:"request_id"`
}

type SmsDetail struct {
	CountryId     interface{} `json:"country_id"`
	ApplicationId interface{} `json:"application_id"`
	Code          string      `json:"sms_code"`
	Phone         string      `json:"number"`
	Taskid        int         `json:"request_id"`
}

type ReleaseDetail struct {
	Message string `json:"code"`
}
