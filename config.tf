terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.56.1"
    }
    opensearch = {
      source  = "opensearch-project/opensearch"
      version = "2.3.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-2"
}

provider "aws" {
  alias  = "us"
  region = "us-east-1"
}

provider "opensearch" {
  url         = aws_opensearchserverless_collection.kb.collection_endpoint
  aws_region  = "us-east-1"
  healthcheck = false
}

data "aws_caller_identity" "main" {}
data "aws_partition" "main" {}

variable "pg_username" {
  type = string
}

variable "pg_password" {
  type      = string
  sensitive = true
}

locals {
  account_id = data.aws_caller_identity.main.account_id
  partition  = data.aws_partition.main.partition
  content_type_map = {
    "js"   = "text/javascript"
    "html" = "text/html"
    "css"  = "text/css"
    "jpeg" = "image/jpeg"
    "png"  = "image/png"
  }
  app = {
    "chart" : "chart",
    "dashboard" : "dashboard",
    "search-basic" : "search/basic",
    "search-advanced" : "search/advanced",
    "search-lyrics" : "search/lyrics"
  }
  dbcred = {
    host     = aws_rds_cluster.main.reader_endpoint
    port     = aws_rds_cluster.main.port
    user     = aws_rds_cluster.main.master_username
    password = var.pg_password
    dbname   = aws_rds_cluster.main.database_name
  }
  kbinfo = {
    kb_id = aws_bedrockagent_knowledge_base.kb.id
  }
}