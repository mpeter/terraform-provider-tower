
variable "name" {
  type    = "string"
  default = "value4"
}

provider "tower" {
  endpoint    = "http://10.42.0.42/api/v1/"
  username    = "admin"
  password    = "pDxKehFK2HAC"
}

data "template_file" "default_yaml" {
    template = "${file("${path.module}/example.yaml")}"
    vars {
       name = "${var.name}"
    }
}

data "template_file" "default_json" {
    template = "${file("${path.module}/example.json")}"
    vars {
       name = "${var.name}"
    }
}

resource "tower_organization" "default" {
  name = "default"
  description = "default"
}

resource "tower_inventory" "default" {
  name = "default"
  organization_id = "${tower_organization.default.id}"
	variables_yaml = <<VARIABLES
---
template_name: ${var.name}
VARIABLES
}

resource "tower_inventory" "alpha" {
  name = "alpha"
  organization_id = "${tower_organization.default.id}"
  variables_json = "${data.template_file.default_json.rendered}"
}

resource "tower_group" "alpha" {
  name = "alpha"
  inventory_id = "${tower_inventory.alpha.id}"
  variables_yaml = "${data.template_file.default_yaml.rendered}"
}

resource "tower_host" "alpha" {
  name = "127.0.0.1"
  inventory_id = "${tower_inventory.alpha.id}"
  enabled = true
  variables_yaml = "${data.template_file.default_yaml.rendered}"
}

resource "tower_inventory" "beta" {
  name = "beta"
  organization_id = "${tower_organization.default.id}"
  variables_json = "${data.template_file.default_json.rendered}"
}

resource "tower_group" "beta" {
  name = "beta"
  inventory_id = "${tower_inventory.beta.id}"
  variables_yaml = "${data.template_file.default_yaml.rendered}"
}

resource "tower_host" "beta" {
  name = "127.0.0.1"
  inventory_id = "${tower_inventory.beta.id}"
  enabled = true
  variables_json = "${data.template_file.default_json.rendered}"
}
