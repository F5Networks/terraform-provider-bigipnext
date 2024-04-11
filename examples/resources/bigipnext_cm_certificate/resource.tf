resource "bigipnext_cm_certitficate" "test" {
  name            = "test"
  cert_text       = "-----BEGIN CERTIFICATE-----MIIFYjCCA0oCCQCEqd7Lp8dxVDANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJSVTELMAkGA1UECAwCTFUxDjAMBgNVBAcMBVBhdGNoMQ8wDQYDVQQKDAZVcGRhdGUxDzANBgNVBAsMBlVwZGF0ZTELMAkGA1UEAwwCdXAxGDAWBgkqhkiG9w0BCQEWCWlwQHVwLmNvbTAeFw0yMzExMDMwNTUwNDNaFw0yNDExMDIwNTUwNDNaMHMxCzAJBgNVBAYTAlJVMQswCQYDVQQIDAJMVTEOMAwGA1UEBwwFUGF0Y2gxDzANBgNVBAoMBlVwZGF0ZTEPMA0GA1UECwwGVXBkYXRlMQswCQYDVQQDDAJ1cDEYMBYGCSqGSIb3DQEJARYJaXBAdXAuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA5M0EI59m0so215qeWFfalwikf0jlcKqiIdCsnRVml1RVgxUC3/CQrmkfot6UJKq62UITQlrIXPCt3eRERa5JuAxayducq0ygPxnE+qMGPb+C4vSuRnhQj9U8OmgPKp4OIidjkb2UJ1hLcqdMQzh4gyBTYUPoKYGHHIgQ/bTUq13zqaT/HtXdnO+OQmeyCuPaSPBU-----END CERTIFICATE-----"
  key_text        = "-----BEGIN ENCRYPTED PRIVATE KEY-----MIIJnzBJBgkqhkiG9w0BBQ0wPDAbBgkqhkiG9w0BBQwwDgQIqx+JNssqck0CAggAMB0GCWCGSAFlAwQBKgQQqw1PviHJvexXVB8ewne6GQSCCVA7Htfh73O1bmwh8Dkt2sjQXrDhSpLadDp35qMgUvcMvJFicxEfiX+/ON98//lIa9bbF0CN12jM1R0cmY24/6lYgMMgFpRSN2tIxOTcF2Q8pucxN8TZDZWC04kuMyJDTWUxmVOm2/yJGL4a4LH11ExHPkspeWviyUFv44vCIu/R4kfxjb+/dn6MzfiDNK4rS75YnpcvvsS4NHWxQLEFzr+bEMZq2YRWrZ3gHoiF53JHac2OKur/bIH0ylH6AzIkg2sPo+JFtWlu9SK1t9QQ6XNh/4O2ovnL3FiOp9JzKgGRgwRxwJm8jXHXec4AbkzAMYIcL6NmpBoXEIQBSP1lGoe3r6tj+E0MCM+51MilAfbY4wlyf7doTbZRjIWW7zxwA5sT275DDPl7DPXisdIijq6XjYXMPD5u1oxqRNbFRynNIly1JOwD5gKtd/Lcp3c8zXT3XRUzkfCudNjcqdIVu4zRdRmaPCCQmlkB1WBYNuzVtNCCGbSQSQfvfLyZs7LjLSkt/ZYVVUwySI+drF//DqF+NjL2RKjs1xibbrWv2sgW1Q==-----END ENCRYPTED PRIVATE KEY-----"
  key_passphrase  = "P@ss"
  cert_passphrase = "P@ssss"
  import_type     = "PKCS12"
}