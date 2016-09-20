
variable "name" {
  type    = "string"
  default = "value4"
}

provider "tower" {
  endpoint    = "http://10.42.0.42/api/v1/"
  username    = "admin"
  password    = "8yYZw3QW647Z"
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

resource "tower_organization" "alpha" {
  name = "alpha"
  description = "alpha"
}

resource "tower_organization" "beta" {
  name = "beta"
  description = "beta"
}

resource "tower_inventory" "default" {
  name = "alpha"
  organization_id = "${tower_organization.alpha.id}"
	variables = <<VARIABLES
---
template_name: ${var.name}
VARIABLES
}

resource "tower_inventory" "alpha" {
  name = "alpha"
  organization_id = "${tower_organization.alpha.id}"
  variables = "${data.template_file.default_json.rendered}"
}

resource "tower_group" "alpha" {
  name = "alpha"
  inventory_id = "${tower_inventory.alpha.id}"
  variables = "${data.template_file.default_yaml.rendered}"
}

resource "tower_host" "alpha" {
  name = "127.0.0.1"
  inventory_id = "${tower_inventory.alpha.id}"
  enabled = true
  variables = "${data.template_file.default_yaml.rendered}"
}

resource "tower_inventory" "beta" {
  name = "beta"
  organization_id = "${tower_organization.beta.id}"
  variables = "${data.template_file.default_json.rendered}"
}

resource "tower_group" "beta" {
  name = "beta"
  inventory_id = "${tower_inventory.beta.id}"
  variables = "${data.template_file.default_yaml.rendered}"
}

resource "tower_host" "beta" {
  name = "127.0.0.1"
  inventory_id = "${tower_inventory.beta.id}"
  enabled = true
  variables = "${data.template_file.default_json.rendered}"
}

resource "tower_job_template" "alpha" {
  name = "alpha"

  inventory_id = "${tower_inventory.alpha.id}"
  project_id = "${tower_project.alpha.id}"
  credential_id = "${tower_credential.alpha.id}"

  job_type = "run"
  playbook = "hello_world.yml"
}

resource "tower_credential" "alpha" {
  name = "alpha"
  kind = "ssh"
  organization_id = "${tower_organization.alpha.id}"
}

resource "tower_project" "alpha" {
  name = "alpha"
  scm_type = "git"
  scm_url = "https://github.com/ansible/ansible-tower-samples"
  scm_update_on_launch = true
  organization_id = "${tower_organization.alpha.id}"
}
