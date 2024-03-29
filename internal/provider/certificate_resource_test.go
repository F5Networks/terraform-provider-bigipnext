package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNextCMCertificateCreateTC1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMCertificateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "duration_in_days", "216"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_type", "ECDSA"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_size", "2048"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "division.*", "test_division"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "organization.*", "test_organization"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "state.*", "test_state"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "email.*", "testemai@gmail.com"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "country.*", "RU"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_curve_name", "secp384r1"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_passphrase", "test_passphrase"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "administrator_email", "admin@gmail.com"),
				),
			},
			{
				Config: testAccNextCMCertificateResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "name", "sample_cert"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name_updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "duration_in_days", "220"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_type", "RSA"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_size", "3072"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "division.*", "test_division_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "organization.*", "test_organization_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "state.*", "test_state_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "country.*", "US"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_curve_name", "secp384r1"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_passphrase", "test_passphrase_updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "administrator_email", "admin_updated@gmail.com"),
				),
			},
		},
	})
}

func TestUnitCMCertificateCreateUnitTC1Resource(t *testing.T) {
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
	mux.HandleFunc("/api/v1/spaces/default/certificates/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/create"}},"path":"/v1/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"id": "497f6eca-6276-4993-bfeb-53cbbbba6f08","message": "string","status": 200 }`)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/4fda14ac-5bc8-46f7-a6ae-f0437cad50a1"}},"cert":{"fingerprint":"a055479187ad62c5eafb776c4c8b952b6f0955065aa8bc293ce0f6cda7e10e97317d71cfeae666ff0647523d1f91aed606e7811e143925445992bc2b34cdff89","checksum":"88eee23299f906cb53abf95f2ff38a0dd53e222ce4346be988472ea3efb8c35158bd927c74a878216da9cf9dd19e126258f9b1017ad2ff236db9d0004cab884d","public_key_type":"ECDSA","public_key_size":384,"public_key_curve_name":"secp384r1","expiration_date_time":"2024-07-02T18:37:28.323908212Z","valid_from":"2023-11-29T18:37:28.323908212Z","issuer":{"Country":["RU"],"Organization":["test_organization"],"OrganizationalUnit":["test_division"],"Locality":["test_locality"],"Province":["test_state"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai@gmail.com"}]},"serial_number":1701283048323,"size":1029,"subject":{"Country":["RU"],"Organization":["test_organization"],"OrganizationalUnit":["test_division"],"Locality":["test_locality"],"Province":["test_state"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai@gmail.com"}]},"version":"0","content":"certificate_hsm_id"},"common_name":"sample_cert_common_name","count":0,"country":["RU"],"creation_date_time":"2023-11-29T18:37:28.098301Z","current_step":"certificate task completed","division":["test_division"],"duration_in_days":216,"email":["testemai@gmail.com"],"hsm":{"storage":"Vault","secret":"bigip-certs","key":"4fda14ac-5bc8-46f7-a6ae-f0437cad50a1","role":"certificate"},"id":"4fda14ac-5bc8-46f7-a6ae-f0437cad50a1","issuer":"Self","key":{"fingerprint":"a055479187ad62c5eafb776c4c8b952b6f0955065aa8bc293ce0f6cda7e10e97317d71cfeae666ff0647523d1f91aed606e7811e143925445992bc2b34cdff89","checksum":"fa09dccf94e2f69797b708840b42df704206ec9747a4033fee15d5e3fd0db6f0eaba86b898c50d409213d9cf71c93aecb17b4b7e6aad3fb32d72545166587297","private_key_type":"ec-private","private_key_size":384,"private_key_curve_name":"secp384r1","passphrase":"key_password_hsm_id","size":481,"content":"key_hsm_id"},"key_curve_name":"secp384r1","key_size":2048,"key_type":"ECDSA","locality":["test_locality"],"modification_date_time":"2023-11-29T18:37:28.723328Z","name":"sample_cert","organization":["test_organization"],"state":["test_state"],"status":"completed","task_id":"45ce37aa-876c-4bfa-bba4-3363cb1ffb1f"}`)
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/renew", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/create"}},"path":"/v1/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMCertificateResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(
				// resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name"),
				),
			},
		},
	})
}

func TestUnitCMCertificateCreateUnitTC2Resource(t *testing.T) {
	var count = 0
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
	mux.HandleFunc("/api/v1/spaces/default/certificates/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/create"}},"path":"/v1/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"id": "497f6eca-6276-4993-bfeb-53cbbbba6f08","message": "string","status": 200 }`)
		} else {
			if count == 0 || count == 1 {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/4fda14ac-5bc8-46f7-a6ae-f0437cad50a1"}},"cert":{"fingerprint":"a055479187ad62c5eafb776c4c8b952b6f0955065aa8bc293ce0f6cda7e10e97317d71cfeae666ff0647523d1f91aed606e7811e143925445992bc2b34cdff89","checksum":"88eee23299f906cb53abf95f2ff38a0dd53e222ce4346be988472ea3efb8c35158bd927c74a878216da9cf9dd19e126258f9b1017ad2ff236db9d0004cab884d","public_key_type":"ECDSA","public_key_size":384,"public_key_curve_name":"secp384r1","expiration_date_time":"2024-07-02T18:37:28.323908212Z","valid_from":"2023-11-29T18:37:28.323908212Z","issuer":{"Country":["RU"],"Organization":["test_organization"],"OrganizationalUnit":["test_division"],"Locality":["test_locality"],"Province":["test_state"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai@gmail.com"}]},"serial_number":1701283048323,"size":1029,"subject":{"Country":["RU"],"Organization":["test_organization"],"OrganizationalUnit":["test_division"],"Locality":["test_locality"],"Province":["test_state"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai@gmail.com"}]},"version":"0","content":"certificate_hsm_id"},"common_name":"sample_cert_common_name","count":0,"country":["RU"],"creation_date_time":"2023-11-29T18:37:28.098301Z","current_step":"certificate task completed","division":["test_division"],"duration_in_days":216,"email":["testemai@gmail.com"],"hsm":{"storage":"Vault","secret":"bigip-certs","key":"4fda14ac-5bc8-46f7-a6ae-f0437cad50a1","role":"certificate"},"id":"4fda14ac-5bc8-46f7-a6ae-f0437cad50a1","issuer":"Self","key":{"fingerprint":"a055479187ad62c5eafb776c4c8b952b6f0955065aa8bc293ce0f6cda7e10e97317d71cfeae666ff0647523d1f91aed606e7811e143925445992bc2b34cdff89","checksum":"fa09dccf94e2f69797b708840b42df704206ec9747a4033fee15d5e3fd0db6f0eaba86b898c50d409213d9cf71c93aecb17b4b7e6aad3fb32d72545166587297","private_key_type":"ec-private","private_key_size":384,"private_key_curve_name":"secp384r1","passphrase":"key_password_hsm_id","size":481,"content":"key_hsm_id"},"key_curve_name":"secp384r1","key_size":2048,"key_type":"ECDSA","locality":["test_locality"],"modification_date_time":"2023-11-29T18:37:28.723328Z","name":"sample_cert","organization":["test_organization"],"state":["test_state"],"status":"completed","task_id":"45ce37aa-876c-4bfa-bba4-3363cb1ffb1f"}`)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}},"administrator_email":"admin_updated@gmail.com","cert":{"fingerprint":"16913f31fc4954bf79feb2d55f64a8aff47915ea1c28c22d581b7752e5094e6a7d1b9e5972a9092737aba558146dffa0d8e05433b2b65b3e26364b79d7732d62","checksum":"132fac2729ac27d9ecac31e5d509662ed0d9a3268f45eafa227b344c1c5f6c5898e945d61aa370b47b084297dc926dcdf98a7c0f0791acc27a9cea3c18d2d821","public_key_type":"RSA","public_key_size":3072,"public_key_curve_name":"secp384r1","expiration_date_time":"2024-07-20T11:07:49.859634592Z","valid_from":"2023-12-13T11:07:49.859634592Z","issuer":{"Country":["US"],"Organization":["test_organization_updated"],"OrganizationalUnit":["test_division_updated"],"Locality":["test_locality_updated"],"Province":["test_state_updated"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name_updated","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai_updated@gmail.com"}]},"serial_number":1702465669859,"size":1968,"subject":{"Country":["US"],"Organization":["test_organization_updated"],"OrganizationalUnit":["test_division_updated"],"Locality":["test_locality_updated"],"Province":["test_state_updated"],"StreetAddress":null,"PostalCode":null,"SerialNumber":"","CommonName":"sample_cert_common_name_updated","Names":null,"ExtraNames":[{"Type":[1,2,840,113549,1,9,1],"Value":"testemai_updated@gmail.com"}]},"version":"0","content":"certificate_hsm_id"},"common_name":"sample_cert_common_name_updated","count":0,"country":["US"],"creation_date_time":"2023-12-13T11:04:18.889429Z","current_step":"certificate task completed","division":["test_division_updated"],"duration_in_days":220,"email":["testemai_updated@gmail.com"],"hsm":{"storage":"Vault","secret":"bigip-certs","key":"43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2","role":"certificate"},"id":"43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2","issuer":"Self","key":{"fingerprint":"16913f31fc4954bf79feb2d55f64a8aff47915ea1c28c22d581b7752e5094e6a7d1b9e5972a9092737aba558146dffa0d8e05433b2b65b3e26364b79d7732d62","checksum":"7eae0daa45ee626a64e21b523ec294cdc653ddb9c6d4b243c8d6972de6c678851c01a6e42a797c84f612c46e7c5f56784cd128163efb74ae8cd33211132f7c1e","private_key_type":"rsa-private","private_key_size":3072,"private_key_curve_name":"secp384r1","passphrase":"key_password_hsm_id","size":2670,"content":"key_hsm_id"},"key_curve_name":"secp384r1","key_size":3072,"key_type":"RSA","locality":["test_locality_updated"],"modification_date_time":"2023-12-13T11:07:51.266171Z","name":"sample_cert","organization":["test_organization_updated"],"state":["test_state_updated"],"status":"completed","task_id":"41e7c59a-88b2-4600-8ace-7fab6947bcbd"}`)
			}
			count++
		}
	})
	mux.HandleFunc("/api/v1/spaces/default/certificates/renew", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"_links":{"self":{"href":"/api/v1/spaces/default/certificates/renew"}},"path":"/v1/certificates/43b7bd5b-5b61-4a64-8fe4-68ef8ed910f2"}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { testAccPreUnitCheck(t) },
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNextCMCertificateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "duration_in_days", "216"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_type", "ECDSA"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_size", "2048"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "division.*", "test_division"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "organization.*", "test_organization"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "state.*", "test_state"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "email.*", "testemai@gmail.com"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "country.*", "RU"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_curve_name", "secp384r1"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_passphrase", "test_passphrase"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "administrator_email", "admin@gmail.com"),
				),
			},
			{
				Config: testAccNextCMCertificateResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "common_name", "sample_cert_common_name_updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "duration_in_days", "220"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_type", "RSA"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_size", "3072"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "division.*", "test_division_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "organization.*", "test_organization_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "state.*", "test_state_updated"),
					resource.TestCheckTypeSetElemAttr("bigipnext_cm_certificate.sample_cert", "country.*", "US"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_curve_name", "secp384r1"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "key_passphrase", "test_passphrase_updated"),
					resource.TestCheckResourceAttr("bigipnext_cm_certificate.sample_cert", "administrator_email", "admin_updated@gmail.com"),
				),
			},
		},
	})
}

const testAccNextCMCertificateResourceConfig = `
resource "bigipnext_cm_certificate" "sample_cert" {
	issuer              = "Self"
	name                = "sample_cert"
	common_name         = "sample_cert_common_name"
	duration_in_days    = 216
	key_type            = "ECDSA"
	key_size            = 2048
	division            = ["test_division"]
	organization        = ["test_organization"]
	locality            = ["test_locality"]
	state               = ["test_state"]
	email               = ["testemai@gmail.com"]
	country             = ["RU"]
	key_curve_name      = "secp384r1"
	key_passphrase      = "test_passphrase"
	administrator_email = "admin@gmail.com"
}`

const testAccNextCMCertificateResourceUpdateConfig = `
resource "bigipnext_cm_certificate" "sample_cert" {
	issuer              = "Self" // Certificate Authority , Self
	name                = "sample_cert"
	common_name         = "sample_cert_common_name_updated"
	duration_in_days    = 220
	key_type            = "RSA"
	key_size            = 3072
	division            = ["test_division_updated"]
	organization        = ["test_organization_updated"]
	locality            = ["test_locality_updated"]
	state               = ["test_state_updated"]
	email               = ["testemai_updated@gmail.com"]
	country             = ["US"]
	key_curve_name      = "secp384r1"
	key_passphrase      = "test_passphrase_updated"
	administrator_email = "admin_updated@gmail.com"
}`
