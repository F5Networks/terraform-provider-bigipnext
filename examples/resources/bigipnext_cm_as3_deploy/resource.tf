resource "bigipnext_cm_as3_deploy" "test" {
  target_address = "10.xxx.xxx.xxx"
  as3_json       = <<EOT
{
    "class": "ADC",
    "schemaVersion": "3.45.0",
    "id": "example-declaration-01",
    "label": "Sample 1",
    "remark": "Simple HTTP application with round robin pool",
    "next-cm-tenant01": {
        "class": "Tenant",
        "next-cm-app01": {
            "class": "Application",
            "template": "http",
            "serviceMain": {
                "class": "Service_HTTP",
                "virtualAddresses": [
                    "10.0.12.10"
                ],
                "pool": "next-cm-pool01"
            },
            "next-cm-pool01": {
                "class": "Pool",
                "monitors": [
                    "http"
                ],
                "members": [
                    {
                        "servicePort": 80,
                        "serverAddresses": [
                            "192.0.2.100",
                            "192.0.2.110"
                        ]
                    }
                ]
            }
        }
    }
}
EOT
}