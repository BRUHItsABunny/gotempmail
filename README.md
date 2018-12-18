# gotempmail
Go temp-mail.org wrapper

## Example
```
// Get client
client := gotempmail.GetClient()

// Get domains
domains := client.GetDomains()
fmt.Println(domains)
// Needed, ALWAYS DO THIS
client.Domains = domains

// Set address
address := "bunny" + domains[3]
_, hash := client.SetAddress(address)
if len(hash) > 0 {
    // Needed, ALWAYS DO THIS
    client.AddressHash = hash
    client.Address = address
    fmt.Println("Set email: " + address)
} else {
    fmt.Println("Didn't set email, domain doesnt match up (" + address + ")")
}

// Get mails
for i := 1; i <= 50; i++ {
    time.Sleep(3 * time.Second)
    mails, err := client.CheckMail()
    if err == nil {
        result_, _ := json.Marshal(mails)
        fmt.Println(string(result_))
        break
    }
    fmt.Println(err)
}
```
