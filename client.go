package gotempmail

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/BRUHItsABunny/gOkHttp"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetClient() MailClient {
	var domains []string
	options := gokhttp.HttpClientOptions{
		Headers: map[string]string{
			"Accept":        "application/json",
			"Authorization": "Basic OGM4ODA4YTAtYTQ3ZC00MDkxLTllM2QtODhlMDYwM2ViMzY5OmplWTJTVFliMg==",
			"User-Agent":    "okhttp/3.5.0",
		},
	}
	client := gokhttp.GetHTTPClient(&options)
	return MailClient{Client: &client,
		Address:     "",
		AddressHash: "",
		Domains:     domains,
		BaseURL:     "http://api2.temp-mail.org/request/",
		URLSuffix:   "/format/json",
		Regex:       regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)}
}

func (client *MailClient) GetDomains() ([]string, error) {

	var result []string
	var err error
	var resp *gokhttp.HttpResponse
	var req *http.Request

	if client.Domains == nil {
		req, err = client.Client.MakeGETRequest(client.BaseURL+"domains"+client.URLSuffix, url.Values{}, map[string]string{})
		if err == nil {
			resp, err = client.Client.Do(req)
			if err == nil {
				err = resp.Object(&result)
				if err == nil {
					client.Domains = result
					return result, nil
				}
			}
		}
	} else {
		return client.Domains, nil
	}
	return nil, err
}

func (client *MailClient) SetAddress(address string) error {
	validator := strings.Split(address, "@")
	domains, err := client.GetDomains()
	if err == nil {
		if len(validator) == 2 {
			validator[1] = "@" + validator[1]
			err = errors.New("invalid domain")
			for _, element := range domains {
				if element == validator[1] {
					result := md5.Sum([]byte(address))
					hash := hex.EncodeToString(result[:])
					client.Address = address
					client.AddressHash = hash
					err = nil
					break
				}
			}
		} else {
			err = errors.New("invalid email address")
		}
	}
	return err
}

func (client MailClient) CheckMail() ([]Mail, error) {

	var result []Mail
	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse

	if len(client.Address) > 0 {
		req, err = client.Client.MakeGETRequest(client.BaseURL+"mail/id/"+client.AddressHash+client.URLSuffix, url.Values{}, map[string]string{})
		if err == nil {
			resp, err = client.Client.Do(req)
			if err == nil {
				if resp.StatusCode != 404 {
					err = resp.Object(&result)
					if err == nil {
						for index, mail := range result {
							from := strings.Split(mail.MailFrom, " ")
							address := strings.TrimLeft(strings.TrimRight(from[1], ">"), "<")
							result[index].MailFrom = strings.Replace(strings.Replace(result[index].MailFrom, "\u003e", "", 1), "\u003c", "", 1)
							result[index].MailFromAddress = address
							result[index].MailText = client.Regex.ReplaceAllString(mail.MailText, "")
						}
						return result, nil
					}
				} else {
					err = errors.New("no emails yet")
				}
			}
		}
	} else {
		err = errors.New("need to set email address first")
	}
	return nil, err
}

func (client MailClient) DeleteMail(mailId string) error {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse

	if len(client.Address) > 0 {
		req, err = client.Client.MakeGETRequest(client.BaseURL+"delete/id/"+mailId+client.URLSuffix, url.Values{}, map[string]string{})
		if err == nil {
			resp, err = client.Client.Do(req)
			if err == nil {
				_, err = resp.Bytes()
			}
		}
	} else {
		err = errors.New("need to set email address first")
	}
	return err
}

func (client MailClient) GetAttachments(mailId string) ([]Attachment, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result [][]Attachment

	if len(client.Address) > 0 {
		req, err = client.Client.MakeGETRequest(client.BaseURL+"attachments/id/"+mailId+client.URLSuffix, url.Values{}, map[string]string{})
		if err == nil {
			resp, err = client.Client.Do(req)
			if err == nil {
				if resp.StatusCode != 404 {
					err = resp.Object(&result)
					if err == nil {
						for index, attachment := range result[0] {
							result[0][index].Body = client.Regex.ReplaceAllString(attachment.Body, "")
						}
						return result[0], nil
					}
				} else {
					err = errors.New("email doesn't exist")
				}
			}
		}
	} else {
		err = errors.New("need to set email address first")
	}
	return nil, err
}

func (client MailClient) GetRawMail(mailId string) (string, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result string

	if len(client.Address) > 0 {
		req, err = client.Client.MakeGETRequest(client.BaseURL+"source/id/"+mailId+client.URLSuffix, url.Values{}, map[string]string{})
		if err == nil {
			resp, err = client.Client.Do(req)
			if err == nil {
				result, err = resp.Text()
				if err == nil {
					return result, nil
				}
			}
		}
	} else {
		err = errors.New("need to set email address first")
	}
	return "", err
}
