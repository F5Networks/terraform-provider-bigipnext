package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMImportCertificateCreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMImportCertificateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "cert_text", "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYuUWsPoMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTEuY29t\nMB4XDTIzMTEwMzA4NDQxOFoXDTI0MTEwMjA4NDQxOFowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCnXg+hc28I3A6QDQUSgL3+1/6u\ndpriakWd8g4fQOT+2S5xuAQNq07SwJZiYMW2VAxjO3IvG6Vjg0zp5cqm0LtuxhFD\norXFDuSYGiItdy5yBTnNcIRZ8CBVY7HXtvo/TURv1u/PB63ArpI+TPgjEGl9LWsH\nSZamlTLI7YIeQPpUbUUroEBHoVxieoa6keFwrtfr2WqwMuO8sAMYux68Qcyi7X2M\nA+Imp8TxkqI7lfvRkix0aeAlE4Me1BugICvXVrR3rDd62jVxvEXITbaH4jCgFevF\nH0B79UPGLyJ4XTGov3LHQL86B2sv74H95A6fdPfDrkdCBcCg5tD0m8RaPzMBAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAg+efqrxQ\nu6PyBCK/SMV98Xo9AtHmozmfSBDAIWJ0ZcYDPLrN8bj7cRC32F7tJg/w+f+hbtHm\nw+MrNlX1PMLJ4IdYNkgHRN0AYMn6gVfNLpL5+B27JeQaYkX63bZeAsvvy02/bbco\nCi+ntqBy7qg8RyavDTG6XyC4qUF0AZZk6O+HHQV/kGqJJMTLLplDyb6rrZecWN4/\nPSBxeGRT8KiLC9war9cP+esAxomo8+ckg5IJDPjBPCVYuPir99gNwylTj2lNCRuT\nu0uq+u9Xhtf5Grd4vz6DiT8uwFXvBDvY2ot1JjxHpP90nLBohr8H6q417gfxoZhy\naoxawGQyUTXDhw==\n-----END CERTIFICATE-----\n"),
				),
			},
			{
				Config: testAccNextCMImportCertificateResourceUpdateConfig,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "cert_text", "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYxYcT/WMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTIuY29t\nMB4XDTIzMTIxMTEwMzQyNloXDTI0MTIxMDEwMzQyNlowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMi5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDUB8g3B0UifUeo+09LTEPIP5zg\nQXFsjNwY5ixbG/Ch0VZCLaqbJtxtoYhKQhtxZE5xTa3Lx20/OG3uE2askoBdKd3A\nmgIoGQMntdxTAgqvNEhEZs2yM+fgOQOl2z0ipK8wjE0IWjYGbF++ZOFAsKFoG3pK\n80793F4Eb9C/vLL5IOeQedkJFr80vyCdT34llTriJrXl0vcyFbKwLtt2zf3KKto6\nGGSErGFyE8/MuRRsKnxVrUtZkoO+uNskqL+hHfA5eepPBmr3SHWAOi1s+B6uHqO+\n/JAb1ZagtUXjbaQ+KA/SCxaBDr528IuTIyFO0nfbBW30Yz42AfSbOse7q1wpAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAyEeAPa/K\nk0JY37wwnrmPaP54B8Z2hIPcDljKa/IJwByRiqYiAwnLwQGYKoMjhh7DwTNgH+94\nrHLKdZxPSZiiAM6nvT4r6Zqq1ptkdNA1k449lLoi6Q4XE5/k/6EMaGir3aN+Bg19\nIVjIZArk/lxD9xPrkTLar8z+275XTC7weDnWrWzyDCCG7DLWcUoUhAImVSyqUKJb\nVwj3X569wCqAogntLe84opVnr1vO5jHR6w8CKPHs+FgEBtXWRoXU/tK3TtCbEHX2\nYNLVL/ydON3ACprzVU9RSNZhs9Er+qKj2ApTtLPA/VGY7G4/1CnMgDklxj5WNBhB\nce82y3bVBt7Orw==\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

func TestUnitCMImportCertificateCreateResource(t *testing.T) {
	testAccPreUnitCheck(t)
	var count = 0
	var count1 = 0
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
			"tokenType": "Bearer",
			"expiresIn": 3600,
			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
			"refreshExpiresIn": 1209600
		}`)
	})

	mux.HandleFunc("/api/v1/spaces/default/certificates/import", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/spaces/default/certificates/import" {
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/certificates/import"
					}
				},
				"path": "/v1/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e"
			}`)
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"id": "6bfa74a3-4e67-4ef5-ad9d-68dcdb3cf603","message": "deleted","status": 200 }`)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e"}},"cert":{"fingerprint":"f2c966cdf0e7067bbd0e0e26167bf4aeb650be8618ca8e3a08852f8418b537b7792dd374e73fa7cbf73e9210300f94ac6830d4d9eb6341ffb9d41e9548a115ce","checksum":"f7f861824087e78a041767e7f7b434223cc8932e14b8f78e7714fab75405027ddb1e085de99308304a89557a14b9b142d3a251b0e302bb5a680c0d1668cb8545","public_key_type":"RSA","public_key_size":2048,"expiration_date_time":"2024-12-10T10:34:26Z","valid_from":"2023-12-11T10:34:26Z","issuer":{"Country":["IN"],"Organization":["F"],"OrganizationalUnit":["INF"],"Locality":["Telangana"],"Province":["Hyderabad"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"www.example2.com","Names":[{"Type":[2,5,4,6],"Value":"IN"},{"Type":[2,5,4,8],"Value":"Hyderabad"},{"Type":[2,5,4,7],"Value":"Telangana"},{"Type":[2,5,4,10],"Value":"F"},{"Type":[2,5,4,11],"Value":"INF"},{"Type":[2,5,4,3],"Value":"www.example2.com"}],"ExtraNames":null},"serial_number":1702290866134,"size":1240,"subject":{"Country":["IN"],"Organization":["F"],"OrganizationalUnit":["INF"],"Locality":["Telangana"],"Province":["Hyderabad"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"www.example2.com","Names":[{"Type":[2,5,4,6],"Value":"IN"},{"Type":[2,5,4,8],"Value":"Hyderabad"},{"Type":[2,5,4,7],"Value":"Telangana"},{"Type":[2,5,4,10],"Value":"F"},{"Type":[2,5,4,11],"Value":"INF"},{"Type":[2,5,4,3],"Value":"www.example2.com"}],"ExtraNames":null},"version":"3","content":"certificate_hsm_id"},"common_name":"www.example2.com","count":0,"country":["IN"],"creation_date_time":"2023-12-11T11:21:40.245163Z","current_step":"certificate task completed","division":["INF"],"duration_in_days":365,"hsm":{"storage":"Vault","secret":"bigip-certs","key":"a519317c-fe11-4699-b7ef-e0c6996d632e","role":"certificate"},"id":"a519317c-fe11-4699-b7ef-e0c6996d632e","issuer":"Self","key":{"fingerprint":"f2c966cdf0e7067bbd0e0e26167bf4aeb650be8618ca8e3a08852f8418b537b7792dd374e73fa7cbf73e9210300f94ac6830d4d9eb6341ffb9d41e9548a115ce","checksum":"ded8793ab9e2bdbe4bde39d72fea784b9181c986695468070c02cc967950c445b6bbb7a2a26cd15d6e72a592c31fdc5998d614ee64296c507c2601f6bf001df8","private_key_type":"rsa-private","private_key_size":2048,"size":1675,"content":"key_hsm_id"},"key_size":2048,"key_type":"RSA","locality":["Telangana"],"modification_date_time":"2023-12-11T11:21:55.600512Z","name":"sample_cert","organization":["F"],"state":["Hyderabad"],"status":"completed","task_id":"e6efcb56-205f-43f3-b8b5-e4592f308b8e"}`)
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e/crt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if count == 0 || count == 1 {
			_, _ = fmt.Fprintf(w, "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYuUWsPoMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTEuY29t\nMB4XDTIzMTEwMzA4NDQxOFoXDTI0MTEwMjA4NDQxOFowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCnXg+hc28I3A6QDQUSgL3+1/6u\ndpriakWd8g4fQOT+2S5xuAQNq07SwJZiYMW2VAxjO3IvG6Vjg0zp5cqm0LtuxhFD\norXFDuSYGiItdy5yBTnNcIRZ8CBVY7HXtvo/TURv1u/PB63ArpI+TPgjEGl9LWsH\nSZamlTLI7YIeQPpUbUUroEBHoVxieoa6keFwrtfr2WqwMuO8sAMYux68Qcyi7X2M\nA+Imp8TxkqI7lfvRkix0aeAlE4Me1BugICvXVrR3rDd62jVxvEXITbaH4jCgFevF\nH0B79UPGLyJ4XTGov3LHQL86B2sv74H95A6fdPfDrkdCBcCg5tD0m8RaPzMBAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAg+efqrxQ\nu6PyBCK/SMV98Xo9AtHmozmfSBDAIWJ0ZcYDPLrN8bj7cRC32F7tJg/w+f+hbtHm\nw+MrNlX1PMLJ4IdYNkgHRN0AYMn6gVfNLpL5+B27JeQaYkX63bZeAsvvy02/bbco\nCi+ntqBy7qg8RyavDTG6XyC4qUF0AZZk6O+HHQV/kGqJJMTLLplDyb6rrZecWN4/\nPSBxeGRT8KiLC9war9cP+esAxomo8+ckg5IJDPjBPCVYuPir99gNwylTj2lNCRuT\nu0uq+u9Xhtf5Grd4vz6DiT8uwFXvBDvY2ot1JjxHpP90nLBohr8H6q417gfxoZhy\naoxawGQyUTXDhw==\n-----END CERTIFICATE-----\n")
		} else {
			_, _ = fmt.Fprintf(w, "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYxYcT/WMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTIuY29t\nMB4XDTIzMTIxMTEwMzQyNloXDTI0MTIxMDEwMzQyNlowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMi5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDUB8g3B0UifUeo+09LTEPIP5zg\nQXFsjNwY5ixbG/Ch0VZCLaqbJtxtoYhKQhtxZE5xTa3Lx20/OG3uE2askoBdKd3A\nmgIoGQMntdxTAgqvNEhEZs2yM+fgOQOl2z0ipK8wjE0IWjYGbF++ZOFAsKFoG3pK\n80793F4Eb9C/vLL5IOeQedkJFr80vyCdT34llTriJrXl0vcyFbKwLtt2zf3KKto6\nGGSErGFyE8/MuRRsKnxVrUtZkoO+uNskqL+hHfA5eepPBmr3SHWAOi1s+B6uHqO+\n/JAb1ZagtUXjbaQ+KA/SCxaBDr528IuTIyFO0nfbBW30Yz42AfSbOse7q1wpAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAyEeAPa/K\nk0JY37wwnrmPaP54B8Z2hIPcDljKa/IJwByRiqYiAwnLwQGYKoMjhh7DwTNgH+94\nrHLKdZxPSZiiAM6nvT4r6Zqq1ptkdNA1k449lLoi6Q4XE5/k/6EMaGir3aN+Bg19\nIVjIZArk/lxD9xPrkTLar8z+275XTC7weDnWrWzyDCCG7DLWcUoUhAImVSyqUKJb\nVwj3X569wCqAogntLe84opVnr1vO5jHR6w8CKPHs+FgEBtXWRoXU/tK3TtCbEHX2\nYNLVL/ydON3ACprzVU9RSNZhs9Er+qKj2ApTtLPA/VGY7G4/1CnMgDklxj5WNBhB\nce82y3bVBt7Orw==\n-----END CERTIFICATE-----")
		}
		count++
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e/key", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if count1 == 0 || count1 == 1 {
			_, _ = fmt.Fprintf(w, "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAp14PoXNvCNwOkA0FEoC9/tf+rnaa4mpFnfIOH0Dk/tkucbgE\nDatO0sCWYmDFtlQMYztyLxulY4NM6eXKptC7bsYRQ6K1xQ7kmBoiLXcucgU5zXCE\nWfAgVWOx17b6P01Eb9bvzwetwK6SPkz4IxBpfS1rB0mWppUyyO2CHkD6VG1FK6BA\nR6FcYnqGupHhcK7X69lqsDLjvLADGLsevEHMou19jAPiJqfE8ZKiO5X70ZIsdGng\nJRODHtQboCAr11a0d6w3eto1cbxFyE22h+IwoBXrxR9Ae/VDxi8ieF0xqL9yx0C/\nOgdrL++B/eQOn3T3w65HQgXAoObQ9JvEWj8zAQIDAQABAoIBAG7Rp0Pd7Q1OuA3F\nsTAZMrSFTEs0mUWF3GbCmXs0OuxL3abKk1SBs4V0o56sOP2LFVC3UdnDUnVbwRe2\nYcKmvxSP7Wp9WCIMxGu6EhtMqOIyd52B/QCYMLCZfG465+P2Q3RSOyM4EGJetNKv\noDWHbnHGGvIOVcQjUiccrGVf3OD9DtqDCFd25t6x0V3WkQR44Fr+OKC+Es9yhXFE\nb9UvnsRgaVqiGoW/7dmx0wo9zk7gP51+feBJgT8D0wrPlhVdD4gnYuVhIAjJjMk1\nFdWheUYWHYlo01Z1eU9caAzEajG4DOC2be56tYGbNZGMHZN+lySuyZymui4HNHQj\n43dwR4ECgYEAx1xWQeYkv9ZAzAr0RFxm6v6Dp+qZbY7eMuOc7p3x+ZlMJ6uNNRiv\nf69vME+Kj/npW4HsTewByt9ygt0TvNNHfQ1q7BTSfCItQbMAQSjlWE6ukvSWJHV6\n60LVIm9dPvE4IGo4szP0kuj4V7C5lYdhq/fmNIYCMFn+foQjyUf56pkCgYEA1urU\nCQ5MI8nq76Eplx1a25hnQsKiQSJU2KUaBzzQqxMrJdZj+tJe/muQ8indnWTFRazI\nPPWwhf4+xCF7pP11Bgqbp+WrkP6b7SpUM4d2PPoJZRCK0+QpeknPZv83PmTgSBsO\nZxsJUr9haDeoj5XBg64DbnlTIVbvGoE4lndmdKkCgYAuzbTKf9d82jYYMTIrom3f\nGaWbFG600+fClsFPG/GlIaJJZfMe1g7NsUgvVV04c/mfLB9oI9I/6Lmfk3uAxzFv\nYGkLx8+qqPNrCzUyFwHQ+5fslFNzd8lF1kjnbrG7hzIgGg/5smbm3p7/J1RKkKAT\nmX2IMzXsWBRxa2BjbuxzcQKBgQCvxMh7S40r6/TP3K/mHiTz2fYB3KrUuF5J/OWH\nq85BS+ELBgco2KrGS3T1CRZtpj/M1x3A9XNUcvYkc/nqmzv9H+nj6+tgH0upMOhC\naHRkNF5AoMHZwA3ILNuKMgqdZeUkM7SY0LzURx9EG9ko7WKh7kxyKpm5d57/v1Vn\ngelyWQKBgDCNWzvHoRMzNQDYrod0W6Tb6npCwjHp2THTllj6lidaF8lb7SAwzOqr\nUByXENFnYnYax1RVWajsV0bJFSysyJVaoLTQ0AcESW1bCK8JPDlOv7cxI6RnBNVu\nK9KcbiJ9p1uUJQm62HbwYVEUNzWJXJgUL9bi+quY5YSoN9pdm4s5\n-----END RSA PRIVATE KEY-----\n")
		} else {
			_, _ = fmt.Fprintf(w, "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA1AfINwdFIn1HqPtPS0xDyD+c4EFxbIzcGOYsWxvwodFWQi2q\nmybcbaGISkIbcWROcU2ty8dtPzht7hNmrJKAXSndwJoCKBkDJ7XcUwIKrzRIRGbN\nsjPn4DkDpds9IqSvMIxNCFo2BmxfvmThQLChaBt6SvNO/dxeBG/Qv7yy+SDnkHnZ\nCRa/NL8gnU9+JZU64ia15dL3MhWysC7bds39yiraOhhkhKxhchPPzLkUbCp8Va1L\nWZKDvrjbJKi/oR3wOXnqTwZq90h1gDotbPgerh6jvvyQG9WWoLVF422kPigP0gsW\ngQ6+dvCLkyMhTtJ32wVt9GM+NgH0mzrHu6tcKQIDAQABAoIBAHT81juGh17AQQm7\nn8SsD7otXyFc+ngqQEZ8uXyLrfmaxz08dSWmC3lx5wER+JJYBe/+LCaDooN/Xyg9\nDCmrq6e7sd7TGt2E73i5nxctyTdiYX1cO6JXgVj2HE0m6lRzCTaAMwCoxaZFpY4n\nmyFWU7hdcNxp4uuu6zEHgOZJ93X88JhsoT90hDrmGKz5Qi1S25n8urXwPOmKUp4F\nOd2F8Ork8EwLg7TmKug0bVZc3I7IEodO86DIZtZ4sAQK8Ccrp0Av6aTl289mRK/S\n73HnxTlqfawFdshmD/uqIqvKbsGDT0KFZBgiScbqjIpV6y739qbFpFBtklmmFNjt\n1fjqZQECgYEA4wIoKoCdEDymIRcuixBTo852HOBB/I8YUhEU7XXO44ZNht5wMygq\nGiMMQ/AGzXx+9FPLripxN4ydsyUoNpDKwJwpf9auudxSQEkfzUWkHIu1ncZRDExh\nPExra0A677PpZVr3K7eS2D8Y1EhKANGhaUB4LyuMBNWkpSNiPeYU2NMCgYEA7xvu\nf41OKPhAR25w0F/eTHOS//gfWD7N1FQr4m4xv0zFknHjb+LISJD1uRhu27I+GVe7\nlvS1aBG2hc/2bJ3ykvxM8rR+pzaWAk/bLHdcrhOSxXjf0+kTq6XjPUHqeY2D3vuw\naln6Esh2nYtlqUoyJmZ4xwK15tyRFmKIAhV42ZMCgYBT/L1NlE4H7thsD76Zls3L\nIhzS7Cmdvnd6DXXXsSl9RngyeOO8GZUSHHtyO0DZD8GMtd/6rRs8ORszZ4DsRz+s\naVp1QMFeZGROAn/wm15vbUUhfXkI+s1S2Nc5VAc6Hi8w36npE78RoK6YA7LVgLme\nTkro8MyaEU0cB+5WBmUaHwKBgFf3z7vPkdTS2FsvT4Pp8U/xKUDQ2T9PA8y9FtQc\nNGMr7HgfPEyag5Lm+fAaBBcBsYUDWPmFmAPsmFkMlJ2LUoYvGmQkcYA1PeUl2f23\nADru6o2KFdbRpjH9Ouf7izcjEEQGFvZZmf41ECaP9Vvd9ytXgdG9toxz01EH+P/D\nRI3PAoGBAIqxC6uop1QIJV2VrMgMYgAPK1iW1+//FVRJZyQbXvt3o+n81yIrc+23\nceCPOSu+Oc3NoJBZEodmtKldSOW1rpzDlOC5gdQJtly4iTQCW8UjhUAdwv/eiBYX\nQrIZGa9VPfwiUFsvxJ0Edz25UuyPKvEZ/36B6tbm4JSipCaMlXnD\n-----END RSA PRIVATE KEY-----\n")
		}
		count1++
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMImportCertificateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "cert_text", "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYuUWsPoMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTEuY29t\nMB4XDTIzMTEwMzA4NDQxOFoXDTI0MTEwMjA4NDQxOFowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCnXg+hc28I3A6QDQUSgL3+1/6u\ndpriakWd8g4fQOT+2S5xuAQNq07SwJZiYMW2VAxjO3IvG6Vjg0zp5cqm0LtuxhFD\norXFDuSYGiItdy5yBTnNcIRZ8CBVY7HXtvo/TURv1u/PB63ArpI+TPgjEGl9LWsH\nSZamlTLI7YIeQPpUbUUroEBHoVxieoa6keFwrtfr2WqwMuO8sAMYux68Qcyi7X2M\nA+Imp8TxkqI7lfvRkix0aeAlE4Me1BugICvXVrR3rDd62jVxvEXITbaH4jCgFevF\nH0B79UPGLyJ4XTGov3LHQL86B2sv74H95A6fdPfDrkdCBcCg5tD0m8RaPzMBAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAg+efqrxQ\nu6PyBCK/SMV98Xo9AtHmozmfSBDAIWJ0ZcYDPLrN8bj7cRC32F7tJg/w+f+hbtHm\nw+MrNlX1PMLJ4IdYNkgHRN0AYMn6gVfNLpL5+B27JeQaYkX63bZeAsvvy02/bbco\nCi+ntqBy7qg8RyavDTG6XyC4qUF0AZZk6O+HHQV/kGqJJMTLLplDyb6rrZecWN4/\nPSBxeGRT8KiLC9war9cP+esAxomo8+ckg5IJDPjBPCVYuPir99gNwylTj2lNCRuT\nu0uq+u9Xhtf5Grd4vz6DiT8uwFXvBDvY2ot1JjxHpP90nLBohr8H6q417gfxoZhy\naoxawGQyUTXDhw==\n-----END CERTIFICATE-----\n"),
				),
			},
			{
				Config: testAccNextCMImportCertificateResourceUpdateConfig,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_import_certitficate.sample_cert", "cert_text", "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYxYcT/WMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTIuY29t\nMB4XDTIzMTIxMTEwMzQyNloXDTI0MTIxMDEwMzQyNlowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMi5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDUB8g3B0UifUeo+09LTEPIP5zg\nQXFsjNwY5ixbG/Ch0VZCLaqbJtxtoYhKQhtxZE5xTa3Lx20/OG3uE2askoBdKd3A\nmgIoGQMntdxTAgqvNEhEZs2yM+fgOQOl2z0ipK8wjE0IWjYGbF++ZOFAsKFoG3pK\n80793F4Eb9C/vLL5IOeQedkJFr80vyCdT34llTriJrXl0vcyFbKwLtt2zf3KKto6\nGGSErGFyE8/MuRRsKnxVrUtZkoO+uNskqL+hHfA5eepPBmr3SHWAOi1s+B6uHqO+\n/JAb1ZagtUXjbaQ+KA/SCxaBDr528IuTIyFO0nfbBW30Yz42AfSbOse7q1wpAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAyEeAPa/K\nk0JY37wwnrmPaP54B8Z2hIPcDljKa/IJwByRiqYiAwnLwQGYKoMjhh7DwTNgH+94\nrHLKdZxPSZiiAM6nvT4r6Zqq1ptkdNA1k449lLoi6Q4XE5/k/6EMaGir3aN+Bg19\nIVjIZArk/lxD9xPrkTLar8z+275XTC7weDnWrWzyDCCG7DLWcUoUhAImVSyqUKJb\nVwj3X569wCqAogntLe84opVnr1vO5jHR6w8CKPHs+FgEBtXWRoXU/tK3TtCbEHX2\nYNLVL/ydON3ACprzVU9RSNZhs9Er+qKj2ApTtLPA/VGY7G4/1CnMgDklxj5WNBhB\nce82y3bVBt7Orw==\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

func TestUnitCMImportCertificateCreateFailResource(t *testing.T) {
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
			"tokenType": "Bearer",
			"expiresIn": 3600,
			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
			"refreshExpiresIn": 1209600
		}`)
	})

	mux.HandleFunc("/api/v1/spaces/default/certificates/import", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/spaces/default/certificates/import" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, `{
				"code": 400,
				"message": "The certificate could not be imported. The certificate is not valid."
			}`)
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccNextCMImportCertificateResourceConfig,
				ExpectError: regexp.MustCompile(`Failed to Import Certificate, got error`),
			},
		},
	})
}

func TestUnitCMImportCertificateReadFailResource(t *testing.T) {
	testAccPreUnitCheck(t)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"token": "eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzM4NCIsImtpZCI6IjJiMGE4MjEwLWJhYmQtNDRhZi04MmMyLTI2YWE4Yjk3OWYwMCIsInR5cCI6IkpXVCJ9.AbY1hUw8wHO8Vt1qxRd5xQj_21EQ1iaH6q9Z2XgRwQl98M7aCpyjiF2J16S4HrZ-",
			"tokenType": "Bearer",
			"expiresIn": 3600,
			"refreshToken": "ODA0MmQzZTctZTk1Mi00OTk1LWJmMjUtZWZmMjc1NDE3YzliOt4bKlRr6g7RdTtnBKhm2vzkgJeWqfvow68gyxTipleCq4AjR4nxZDBYKQaWyCWGeA",
			"refreshExpiresIn": 1209600
		}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/import", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/v1/spaces/default/certificates/import" {
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintf(w, `{
				"_links": {
					"self": {
						"href": "/v1/certificates/import"
					}
				},
				"path": "/v1/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e"
			}`)
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"id": "6bfa74a3-4e67-4ef5-ad9d-68dcdb3cf603","message": "deleted","status": 200 }`)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e"}},"cert":{"fingerprint":"f2c966cdf0e7067bbd0e0e26167bf4aeb650be8618ca8e3a08852f8418b537b7792dd374e73fa7cbf73e9210300f94ac6830d4d9eb6341ffb9d41e9548a115ce","checksum":"f7f861824087e78a041767e7f7b434223cc8932e14b8f78e7714fab75405027ddb1e085de99308304a89557a14b9b142d3a251b0e302bb5a680c0d1668cb8545","public_key_type":"RSA","public_key_size":2048,"expiration_date_time":"2024-12-10T10:34:26Z","valid_from":"2023-12-11T10:34:26Z","issuer":{"Country":["IN"],"Organization":["F"],"OrganizationalUnit":["INF"],"Locality":["Telangana"],"Province":["Hyderabad"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"www.example2.com","Names":[{"Type":[2,5,4,6],"Value":"IN"},{"Type":[2,5,4,8],"Value":"Hyderabad"},{"Type":[2,5,4,7],"Value":"Telangana"},{"Type":[2,5,4,10],"Value":"F"},{"Type":[2,5,4,11],"Value":"INF"},{"Type":[2,5,4,3],"Value":"www.example2.com"}],"ExtraNames":null},"serial_number":1702290866134,"size":1240,"subject":{"Country":["IN"],"Organization":["F"],"OrganizationalUnit":["INF"],"Locality":["Telangana"],"Province":["Hyderabad"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"www.example2.com","Names":[{"Type":[2,5,4,6],"Value":"IN"},{"Type":[2,5,4,8],"Value":"Hyderabad"},{"Type":[2,5,4,7],"Value":"Telangana"},{"Type":[2,5,4,10],"Value":"F"},{"Type":[2,5,4,11],"Value":"INF"},{"Type":[2,5,4,3],"Value":"www.example2.com"}],"ExtraNames":null},"version":"3","content":"certificate_hsm_id"},"common_name":"www.example2.com","count":0,"country":["IN"],"creation_date_time":"2023-12-11T11:21:40.245163Z","current_step":"certificate task completed","division":["INF"],"duration_in_days":365,"hsm":{"storage":"Vault","secret":"bigip-certs","key":"a519317c-fe11-4699-b7ef-e0c6996d632e","role":"certificate"},"id":"a519317c-fe11-4699-b7ef-e0c6996d632e","issuer":"Self","key":{"fingerprint":"f2c966cdf0e7067bbd0e0e26167bf4aeb650be8618ca8e3a08852f8418b537b7792dd374e73fa7cbf73e9210300f94ac6830d4d9eb6341ffb9d41e9548a115ce","checksum":"ded8793ab9e2bdbe4bde39d72fea784b9181c986695468070c02cc967950c445b6bbb7a2a26cd15d6e72a592c31fdc5998d614ee64296c507c2601f6bf001df8","private_key_type":"rsa-private","private_key_size":2048,"size":1675,"content":"key_hsm_id"},"key_size":2048,"key_type":"RSA","locality":["Telangana"],"modification_date_time":"2023-12-11T11:21:55.600512Z","name":"sample_cert","organization":["F"],"state":["Hyderabad"],"status":"completed","task_id":"e6efcb56-205f-43f3-b8b5-e4592f308b8e"}`)
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/a519317c-fe11-4699-b7ef-e0c6996d632e/crt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestTimeout)
		_, _ = fmt.Fprintf(w, `{
			"code": 408,
			"message": "The certificate could not be imported. The certificate is not valid."
		}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccNextCMImportCertificateResourceConfig,
				ExpectError: regexp.MustCompile(`Failed to Read Certificate, Key Data`),
			},
		},
	})
}

const testAccNextCMImportCertificateResourceConfig = `
resource "bigipnext_cm_import_certitficate" "sample_cert" {
	name      = "sample_cert"
	cert_text = "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYuUWsPoMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTEuY29t\nMB4XDTIzMTEwMzA4NDQxOFoXDTI0MTEwMjA4NDQxOFowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCnXg+hc28I3A6QDQUSgL3+1/6u\ndpriakWd8g4fQOT+2S5xuAQNq07SwJZiYMW2VAxjO3IvG6Vjg0zp5cqm0LtuxhFD\norXFDuSYGiItdy5yBTnNcIRZ8CBVY7HXtvo/TURv1u/PB63ArpI+TPgjEGl9LWsH\nSZamlTLI7YIeQPpUbUUroEBHoVxieoa6keFwrtfr2WqwMuO8sAMYux68Qcyi7X2M\nA+Imp8TxkqI7lfvRkix0aeAlE4Me1BugICvXVrR3rDd62jVxvEXITbaH4jCgFevF\nH0B79UPGLyJ4XTGov3LHQL86B2sv74H95A6fdPfDrkdCBcCg5tD0m8RaPzMBAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAg+efqrxQ\nu6PyBCK/SMV98Xo9AtHmozmfSBDAIWJ0ZcYDPLrN8bj7cRC32F7tJg/w+f+hbtHm\nw+MrNlX1PMLJ4IdYNkgHRN0AYMn6gVfNLpL5+B27JeQaYkX63bZeAsvvy02/bbco\nCi+ntqBy7qg8RyavDTG6XyC4qUF0AZZk6O+HHQV/kGqJJMTLLplDyb6rrZecWN4/\nPSBxeGRT8KiLC9war9cP+esAxomo8+ckg5IJDPjBPCVYuPir99gNwylTj2lNCRuT\nu0uq+u9Xhtf5Grd4vz6DiT8uwFXvBDvY2ot1JjxHpP90nLBohr8H6q417gfxoZhy\naoxawGQyUTXDhw==\n-----END CERTIFICATE-----\n"
	key_text  = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAp14PoXNvCNwOkA0FEoC9/tf+rnaa4mpFnfIOH0Dk/tkucbgE\nDatO0sCWYmDFtlQMYztyLxulY4NM6eXKptC7bsYRQ6K1xQ7kmBoiLXcucgU5zXCE\nWfAgVWOx17b6P01Eb9bvzwetwK6SPkz4IxBpfS1rB0mWppUyyO2CHkD6VG1FK6BA\nR6FcYnqGupHhcK7X69lqsDLjvLADGLsevEHMou19jAPiJqfE8ZKiO5X70ZIsdGng\nJRODHtQboCAr11a0d6w3eto1cbxFyE22h+IwoBXrxR9Ae/VDxi8ieF0xqL9yx0C/\nOgdrL++B/eQOn3T3w65HQgXAoObQ9JvEWj8zAQIDAQABAoIBAG7Rp0Pd7Q1OuA3F\nsTAZMrSFTEs0mUWF3GbCmXs0OuxL3abKk1SBs4V0o56sOP2LFVC3UdnDUnVbwRe2\nYcKmvxSP7Wp9WCIMxGu6EhtMqOIyd52B/QCYMLCZfG465+P2Q3RSOyM4EGJetNKv\noDWHbnHGGvIOVcQjUiccrGVf3OD9DtqDCFd25t6x0V3WkQR44Fr+OKC+Es9yhXFE\nb9UvnsRgaVqiGoW/7dmx0wo9zk7gP51+feBJgT8D0wrPlhVdD4gnYuVhIAjJjMk1\nFdWheUYWHYlo01Z1eU9caAzEajG4DOC2be56tYGbNZGMHZN+lySuyZymui4HNHQj\n43dwR4ECgYEAx1xWQeYkv9ZAzAr0RFxm6v6Dp+qZbY7eMuOc7p3x+ZlMJ6uNNRiv\nf69vME+Kj/npW4HsTewByt9ygt0TvNNHfQ1q7BTSfCItQbMAQSjlWE6ukvSWJHV6\n60LVIm9dPvE4IGo4szP0kuj4V7C5lYdhq/fmNIYCMFn+foQjyUf56pkCgYEA1urU\nCQ5MI8nq76Eplx1a25hnQsKiQSJU2KUaBzzQqxMrJdZj+tJe/muQ8indnWTFRazI\nPPWwhf4+xCF7pP11Bgqbp+WrkP6b7SpUM4d2PPoJZRCK0+QpeknPZv83PmTgSBsO\nZxsJUr9haDeoj5XBg64DbnlTIVbvGoE4lndmdKkCgYAuzbTKf9d82jYYMTIrom3f\nGaWbFG600+fClsFPG/GlIaJJZfMe1g7NsUgvVV04c/mfLB9oI9I/6Lmfk3uAxzFv\nYGkLx8+qqPNrCzUyFwHQ+5fslFNzd8lF1kjnbrG7hzIgGg/5smbm3p7/J1RKkKAT\nmX2IMzXsWBRxa2BjbuxzcQKBgQCvxMh7S40r6/TP3K/mHiTz2fYB3KrUuF5J/OWH\nq85BS+ELBgco2KrGS3T1CRZtpj/M1x3A9XNUcvYkc/nqmzv9H+nj6+tgH0upMOhC\naHRkNF5AoMHZwA3ILNuKMgqdZeUkM7SY0LzURx9EG9ko7WKh7kxyKpm5d57/v1Vn\ngelyWQKBgDCNWzvHoRMzNQDYrod0W6Tb6npCwjHp2THTllj6lidaF8lb7SAwzOqr\nUByXENFnYnYax1RVWajsV0bJFSysyJVaoLTQ0AcESW1bCK8JPDlOv7cxI6RnBNVu\nK9KcbiJ9p1uUJQm62HbwYVEUNzWJXJgUL9bi+quY5YSoN9pdm4s5\n-----END RSA PRIVATE KEY-----\n"
}
`
const testAccNextCMImportCertificateResourceUpdateConfig = `
resource "bigipnext_cm_import_certitficate" "sample_cert" {
    name      = "sample_cert"
    cert_text = "-----BEGIN CERTIFICATE-----\nMIIDZjCCAk6gAwIBAgIGAYxYcT/WMA0GCSqGSIb3DQEBDQUAMGoxCzAJBgNVBAYT\nAklOMRIwEAYDVQQIEwlIeWRlcmFiYWQxEjAQBgNVBAcTCVRlbGFuZ2FuYTEKMAgG\nA1UEChMBRjEMMAoGA1UECxMDSU5GMRkwFwYDVQQDExB3d3cuZXhhbXBsZTIuY29t\nMB4XDTIzMTIxMTEwMzQyNloXDTI0MTIxMDEwMzQyNlowajELMAkGA1UEBhMCSU4x\nEjAQBgNVBAgTCUh5ZGVyYWJhZDESMBAGA1UEBxMJVGVsYW5nYW5hMQowCAYDVQQK\nEwFGMQwwCgYDVQQLEwNJTkYxGTAXBgNVBAMTEHd3dy5leGFtcGxlMi5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDUB8g3B0UifUeo+09LTEPIP5zg\nQXFsjNwY5ixbG/Ch0VZCLaqbJtxtoYhKQhtxZE5xTa3Lx20/OG3uE2askoBdKd3A\nmgIoGQMntdxTAgqvNEhEZs2yM+fgOQOl2z0ipK8wjE0IWjYGbF++ZOFAsKFoG3pK\n80793F4Eb9C/vLL5IOeQedkJFr80vyCdT34llTriJrXl0vcyFbKwLtt2zf3KKto6\nGGSErGFyE8/MuRRsKnxVrUtZkoO+uNskqL+hHfA5eepPBmr3SHWAOi1s+B6uHqO+\n/JAb1ZagtUXjbaQ+KA/SCxaBDr528IuTIyFO0nfbBW30Yz42AfSbOse7q1wpAgMB\nAAGjEjAQMA4GA1UdDwEB/wQEAwIFoDANBgkqhkiG9w0BAQ0FAAOCAQEAyEeAPa/K\nk0JY37wwnrmPaP54B8Z2hIPcDljKa/IJwByRiqYiAwnLwQGYKoMjhh7DwTNgH+94\nrHLKdZxPSZiiAM6nvT4r6Zqq1ptkdNA1k449lLoi6Q4XE5/k/6EMaGir3aN+Bg19\nIVjIZArk/lxD9xPrkTLar8z+275XTC7weDnWrWzyDCCG7DLWcUoUhAImVSyqUKJb\nVwj3X569wCqAogntLe84opVnr1vO5jHR6w8CKPHs+FgEBtXWRoXU/tK3TtCbEHX2\nYNLVL/ydON3ACprzVU9RSNZhs9Er+qKj2ApTtLPA/VGY7G4/1CnMgDklxj5WNBhB\nce82y3bVBt7Orw==\n-----END CERTIFICATE-----"
    key_text  = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA1AfINwdFIn1HqPtPS0xDyD+c4EFxbIzcGOYsWxvwodFWQi2q\nmybcbaGISkIbcWROcU2ty8dtPzht7hNmrJKAXSndwJoCKBkDJ7XcUwIKrzRIRGbN\nsjPn4DkDpds9IqSvMIxNCFo2BmxfvmThQLChaBt6SvNO/dxeBG/Qv7yy+SDnkHnZ\nCRa/NL8gnU9+JZU64ia15dL3MhWysC7bds39yiraOhhkhKxhchPPzLkUbCp8Va1L\nWZKDvrjbJKi/oR3wOXnqTwZq90h1gDotbPgerh6jvvyQG9WWoLVF422kPigP0gsW\ngQ6+dvCLkyMhTtJ32wVt9GM+NgH0mzrHu6tcKQIDAQABAoIBAHT81juGh17AQQm7\nn8SsD7otXyFc+ngqQEZ8uXyLrfmaxz08dSWmC3lx5wER+JJYBe/+LCaDooN/Xyg9\nDCmrq6e7sd7TGt2E73i5nxctyTdiYX1cO6JXgVj2HE0m6lRzCTaAMwCoxaZFpY4n\nmyFWU7hdcNxp4uuu6zEHgOZJ93X88JhsoT90hDrmGKz5Qi1S25n8urXwPOmKUp4F\nOd2F8Ork8EwLg7TmKug0bVZc3I7IEodO86DIZtZ4sAQK8Ccrp0Av6aTl289mRK/S\n73HnxTlqfawFdshmD/uqIqvKbsGDT0KFZBgiScbqjIpV6y739qbFpFBtklmmFNjt\n1fjqZQECgYEA4wIoKoCdEDymIRcuixBTo852HOBB/I8YUhEU7XXO44ZNht5wMygq\nGiMMQ/AGzXx+9FPLripxN4ydsyUoNpDKwJwpf9auudxSQEkfzUWkHIu1ncZRDExh\nPExra0A677PpZVr3K7eS2D8Y1EhKANGhaUB4LyuMBNWkpSNiPeYU2NMCgYEA7xvu\nf41OKPhAR25w0F/eTHOS//gfWD7N1FQr4m4xv0zFknHjb+LISJD1uRhu27I+GVe7\nlvS1aBG2hc/2bJ3ykvxM8rR+pzaWAk/bLHdcrhOSxXjf0+kTq6XjPUHqeY2D3vuw\naln6Esh2nYtlqUoyJmZ4xwK15tyRFmKIAhV42ZMCgYBT/L1NlE4H7thsD76Zls3L\nIhzS7Cmdvnd6DXXXsSl9RngyeOO8GZUSHHtyO0DZD8GMtd/6rRs8ORszZ4DsRz+s\naVp1QMFeZGROAn/wm15vbUUhfXkI+s1S2Nc5VAc6Hi8w36npE78RoK6YA7LVgLme\nTkro8MyaEU0cB+5WBmUaHwKBgFf3z7vPkdTS2FsvT4Pp8U/xKUDQ2T9PA8y9FtQc\nNGMr7HgfPEyag5Lm+fAaBBcBsYUDWPmFmAPsmFkMlJ2LUoYvGmQkcYA1PeUl2f23\nADru6o2KFdbRpjH9Ouf7izcjEEQGFvZZmf41ECaP9Vvd9ytXgdG9toxz01EH+P/D\nRI3PAoGBAIqxC6uop1QIJV2VrMgMYgAPK1iW1+//FVRJZyQbXvt3o+n81yIrc+23\nceCPOSu+Oc3NoJBZEodmtKldSOW1rpzDlOC5gdQJtly4iTQCW8UjhUAdwv/eiBYX\nQrIZGa9VPfwiUFsvxJ0Edz25UuyPKvEZ/36B6tbm4JSipCaMlXnD\n-----END RSA PRIVATE KEY-----\n"
}
`
