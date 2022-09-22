# Vex Terraform Provider

```hcl
terraform {
  required_providers {
    vex = {
      source  = "broswen/vex"
      version = "1.0.0"
    }
  }
}

provider "vex" {
  api_token = "<api token>"
}

resource "vex_account" "test_account" {
  name = "test account"
  description = "test account"
}

resource "vex_project" "test_project" {
  account_id = vex_account.test_account.id
  name = "test project"
  description = "test project"
}

resource "vex_flag" "feature_1" {
  project_id = vex_project.test_project.id
  account_id = vex_account.test_account.id
  key = "feature1"
  type = "STRING"
  value = "feature one"
}

resource "vex_flag" "feature_2" {
  project_id = vex_project.test_project.id
  account_id = vex_account.test_account.id
  key = "feature2"
  type = "NUMBER"
  value = "123.45"
}
```
