package gotempmail

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
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
	Client      *http.Client
	Address     string
	AddressHash string
	Domains     []string
	BaseURL     string
	URLSuffix   string
	Regex       *regexp.Regexp
}

func (client MailClient) makeRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic OGM4ODA4YTAtYTQ3ZC00MDkxLTllM2QtODhlMDYwM2ViMzY5OmplWTJTVFliMg==")
	req.Header.Add("User-Agent", "okhttp/3.5.0")
	return req
}

func GetClient() MailClient {
	var domains []string
	return MailClient{Client: &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		}},
		Address:     "",
		AddressHash: "",
		Domains:     domains,
		BaseURL:     "http://api2.temp-mail.org/request/",
		URLSuffix:   "/format/json",
		Regex:       regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)}
}

func (client MailClient) GetDomains() []string {
	if client.Domains == nil {
		var result []string
		resp, err := client.Client.Do(client.makeRequest(client.BaseURL + "domains" + client.URLSuffix))
		if err != nil {
			log.Fatalln(err)
			return nil
		}
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Fatalln(err2)
			return nil
		}
		_ = json.Unmarshal(bodyBytes, &result)
		client.Domains = result
	}
	return client.Domains
}

func (client MailClient) SetAddress(address string) (string, string) {
	validator := strings.Split(address, "@")
	domains := client.GetDomains()
	if len(validator) == 2 {
		validator[1] = "@" + validator[1]
		for _, element := range domains {
			if element == validator[1] {
				result := md5.Sum([]byte(address))
				hash := hex.EncodeToString(result[:])
				return address, hash
			}
		}
	}
	return "", ""
}

func (client MailClient) CheckMail() ([]Mail, error) {
	var result []Mail
	if len(client.Address) > 0 {
		resp, err := client.Client.Do(client.makeRequest(client.BaseURL + "mail/id/" + client.AddressHash + client.URLSuffix))
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 404 {
			if err2 != nil {
				log.Fatalln(err2)
				return nil, err2
			}
			_ = json.Unmarshal(bodyBytes, &result)
			for index, mail := range result {
				from := strings.Split(mail.MailFrom, " ")
				address := strings.TrimLeft(strings.TrimRight(from[1], ">"), "<")
				result[index].MailFromAddress = address
				result[index].MailText = client.Regex.ReplaceAllString(mail.MailText, "")
			}
			return result, nil
		}
		return nil, errors.New("no emails yet")
	}
	return nil, errors.New("need to set email address first")
}

func (client MailClient) DeleteMail(mailId string) {
	if len(client.Address) > 0 {
		resp, err := client.Client.Do(client.makeRequest(client.BaseURL + "delete/id/" + mailId + client.URLSuffix))
		if err != nil {
			log.Fatalln(err)
		}
		_, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Fatalln(err2)
		}
	}
}

func (client MailClient) GetAttachments(mailId string) ([]Attachment, error) {
	if len(client.Address) > 0 {
		url := client.BaseURL + "attachments/id/" + mailId + client.URLSuffix
		resp, err := client.Client.Do(client.makeRequest(url))
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 404 {
			if err2 != nil {
				log.Fatalln(err2)
				return nil, err2
			}
			var result [][]Attachment
			_ = json.Unmarshal(bodyBytes, &result)
			for index, attachment := range result[0] {
				result[0][index].Body = client.Regex.ReplaceAllString(attachment.Body, "")
			}
			return result[0], nil
		}
		return nil, errors.New("no attachments in this email")
	}
	return nil, errors.New("need to set email address first")
}

func (client MailClient) GetRawMail(mailId string) (string, error) {
	if len(client.Address) > 0 {
		url := client.BaseURL + "source/id/" + mailId + client.URLSuffix
		resp, err := client.Client.Do(client.makeRequest(url))
		if err != nil {
			log.Fatalln(err)
			return "", err
		}
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Fatalln(err2)
			return "", err2
		}
		return string(bodyBytes), nil
	}
}
