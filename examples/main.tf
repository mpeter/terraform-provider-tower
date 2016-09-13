
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


resource "tower_inventory" "default" {
  name = "default8-mp"
  description = "default-changed2-3"
  organization = "1"
	variables_yaml = <<VARIABLES
---
template_name: ${var.name}
VARIABLES
}

resource "tower_organization" "default" {
  name = "default"
  description = "default"
}

resource "tower_inventory" "child_yaml" {
  name = "default12-mp"
  organization = "${tower_organization.default.id}"
  description = "${tower_inventory.default.description}"
  variables_yaml = "${data.template_file.default_yaml.rendered}"
}

resource "tower_inventory" "child_json" {
  name = "default11-mp"
  organization = "${tower_organization.default.id}"
  description = "${tower_inventory.default.description}"
  variables_json = "${data.template_file.default_json.rendered}"
}

resource "tower_host" "myhost" {
  name = "myhost-02"
  inventory = "${tower_inventory.child_yaml.id}"
  enabled = true
  description = "${tower_inventory.default.description}"
  variables_json = "${data.template_file.default_json.rendered}"
}
