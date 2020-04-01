

locals {
  foo = "${var.bar}"
}

output "foo_bar" {
  value = "${local.foo}"
}
