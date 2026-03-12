variable "name" {
  type = string
}

resource "descope_project" "test" {
  name = var.name

  authorization = {
    roles = [
      {
        key         = "app-developer"
        name        = "App Developer"
        description = "Builds apps and uploads new beta builds"
        permissions = ["build-apps", "upload-builds", "install-builds"]
      },
      {
        key         = "app-tester"
        name        = "App Tester"
        description = "Installs and tests beta releases"
        permissions = ["install-builds"]
      },
    ]
    permissions = [
      {
        name        = "build-apps"
        description = "Allowed to build and sign applications"
      },
      {
        name        = "upload-builds"
        description = "Allowed to upload new releases"
      },
      {
        name        = "install-builds"
        description = "Allowed to install beta releases"
      },
    ]
  }
}

output "id" {
  value = descope_project.test.id
}

output "name" {
  value = descope_project.test.name
}
