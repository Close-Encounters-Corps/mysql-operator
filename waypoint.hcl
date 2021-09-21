project = "mysql-operator"

variable "registry" {
    type    = string
    default = "cr.yandex/crpit22t8m0nfc4t68p5"
}

variable "tag" {
    type    = string
    default = "latest"
}

variable "dockerpassword" {
    type    = string
    default = ""
}

variable "dockerconfigjson" {
  type    = string
  default = ""
}

variable "mysql_host" {
    type    = string
    default = ""
}

variable "mysql_user" {
    type    = string
    default = "root"
}

variable "mysql_password" {
    type    = string
    default = ""
}

variable "mysql_uri_args" {
    type    = string
    default = ""
}

variable "mysql_default_db" {
    type    = string
    default = ""
}

variable "namespace" {
    type    = string
    default = "operators"
}

app "manager" {
    build {
        use "docker" {}
        registry {
            use "docker" {
                image   = "${var.registry}/mysql-operator"
                tag     = var.tag
                encoded_auth = base64encode(
                    var.dockerpassword != "" ? jsonencode({ "username" : "json_key", "password" : var.dockerpassword }) : ""
                )
            }
        }
    }
    deploy {
        use "kubernetes-apply" {
            path = templatedir("${path.app}/k8s", {
                mysql_host = base64encode(var.mysql_host)
                mysql_user = base64encode(var.mysql_user)
                mysql_password = base64encode(var.mysql_password)
                mysql_uri_args = base64encode(var.mysql_uri_args)
                mysql_default_db = base64encode(var.mysql_default_db)
            })
            prune_label = "app=mysql-operator"
        }
    }
    url {
        auto_hostname = false
    }
}