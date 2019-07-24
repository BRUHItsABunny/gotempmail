package gotempmail

import (
	"github.com/BRUHItsABunny/gOkHttp"
	"regexp"
)

type Mail struct {
	MailId          string  `json:"mail_id"`
	MailAddressId   string  `json:"mail_address_id"`
	MailFrom        string  `json:"mail_from"`
	MailFromAddress string  `json:"mail_from_address"`
	MailSubject     string  `json:"mail_subject"`
	MailText        string  `json:"mail_text"`
	MailTimeStamp   float64 `json:"mail_timestamp"`
}

type Attachment struct {
	Header struct {
		ContentType             string `json:"content-type"`
		ContentDisposition      string `json:"content-disposition"`
		ContentTransferEncoding string `json:"content-transfer-encoding"`
		ContentID               string `json:"x-attachment-id"`
	} `json:"header"`
	Body string `json:"body"`
}

type MailClient struct {
	Client      *gokhttp.HttpClient
	Address     string
	AddressHash string
	Domains     []string
	BaseURL     string
	URLSuffix   string
	Regex       *regexp.Regexp
}
