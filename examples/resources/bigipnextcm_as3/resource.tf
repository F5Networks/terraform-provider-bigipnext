resource "bigipnext_cm_as3" "test2" {
  as3_json = <<EOT
{
    "class": "AS3",
    "action": "deploy",
    "persist": true,
    "declaration": {
        "class": "ADC",
        "schemaVersion": "3.45.0",
        "id": "example-declaration-01",
        "label": "Sample 1",
        "remark": "Simple HTTP application with round robin pool",
        "target": {
            "address": "xxx.xxx.xxx.xxxx"
        },
        "next-cm-tenant02": {
            "class": "Tenant",
            "next-cm-app02": {
                "class": "Application",
                "template": "http",
                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.0.2.10"
                    ],
                    "pool": "next-cm-pool02"
                },
                "next-cm-pool02": {
                    "class": "Pool",
                    "monitors": [
                        "http"
                    ],
                    "members": [
                        {
                            "servicePort": 80,
                            "serverAddresses": [
                                "192.0.3.100",
                                "192.0.3.110"
                            ]
                        }
                    ]
                }
            }
        }
    }
}

EOT
}
