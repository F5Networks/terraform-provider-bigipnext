resource "bigipnext_cm_waf_policy" "sample" {
  name                 = "new_waf_policy"
  description          = "new_waf_policy desc"
  tags                 = ["test3", "test4"]
  enforcement_mode     = "blocking"
  application_language = "utf-8"
  template_name        = "Rating-Based-Template"
}