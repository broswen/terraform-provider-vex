# Vex Terraform Provider

```hcl
provider "vex" {
  api_token = "api token"
  account_id = "account id"
}

resource "vex_account" "main" {
  name = "account name"
  description = "account description"
}

resource "vex_project" "app1" {
  account_id = vex_account.main.id
  name = "project name"
  description = "project description"
}

resource "vex_flag" "feature1" {
  project_id = vex_project.app1.id
  key = "flag key"
  type = "flag type"
  value = "flag raw value"
}

```
