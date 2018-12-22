# gotempmail
Go temp-mail.org wrapper

## Example
```
// Get client
client := gotempmail.GetClient()

// Get domains
domains := client.GetDomains()
fmt.Println(domains)

// Set address
address := "bunny" + domains[0]
res := client.SetAddress(address)
if res {
    fmt.Println(client.Address)
    fmt.Println(client.AddressHash)
} else {
    fmt.Println("Didn't set email, domain doesnt match up (" + address + ")")
}

// Get mails
for i := 1; i <= 50; i++ {
    time.Sleep(3 * time.Second)
    mails, err := client.CheckMail()
    if err == nil {
        for _, mail := range mails{
            result_, _ := json.Marshal(mail)
            fmt.Println(string(result_))

            // Get attachments
            attachments, err2 := client.GetAttachments(mail.MailId)
            if err2 == nil {
                for _, attachment := range attachments{
                    /*
                    File is in attachment.Body as base64 encoded byte array
                    */
                    header, _ := json.Marshal(attachment.Header)
                    fmt.Println(string(header))
                }
            }

            // Delete email
            client.DeleteMail(mail.MailId)
        }
    } else {
        fmt.Println(err)
    }
}
```
