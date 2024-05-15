resource "bigipnext_cm_waf_policy_import" "sample" {
  name        = "new_waf_policy"
  description = "new_waf_policy desc"
  file_path   = "/Users/r.chinthalapalli/Downloads/testwaf5_awaf.json"
  file_md5    = md5(file("/Users/r.chinthalapalli/Downloads/testwaf5_awaf.json"))
}