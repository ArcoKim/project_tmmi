data "aws_bedrock_foundation_model" "claude" {
  model_id = "anthropic.claude-3-sonnet-20240229-v1:0"
}

data "aws_bedrock_foundation_model" "titan" {
  model_id = "amazon.titan-embed-text-v2:0"
}

resource "aws_bedrockagent_knowledge_base" "kb" {
  name     = "tmmi-kb"
  role_arn = aws_iam_role.bedrock_kb.arn
  knowledge_base_configuration {
    vector_knowledge_base_configuration {
      embedding_model_arn = data.aws_bedrock_foundation_model.titan.model_arn
    }
    type = "VECTOR"
  }
  storage_configuration {
    type = "OPENSEARCH_SERVERLESS"
    opensearch_serverless_configuration {
      collection_arn    = aws_opensearchserverless_collection.kb.arn
      vector_index_name = opensearch_index.kb.name
      field_mapping {
        vector_field   = "bedrock-knowledge-base-default-vector"
        text_field     = "AMAZON_BEDROCK_TEXT_CHUNK"
        metadata_field = "AMAZON_BEDROCK_METADATA"
      }
    }
  }
  depends_on = [
    time_sleep.aws_iam_role_policy_kb_oss
  ]
}

resource "aws_bedrockagent_data_source" "kb" {
  knowledge_base_id = aws_bedrockagent_knowledge_base.kb.id
  name              = "tmmi-songs"
  data_source_configuration {
    type = "S3"
    s3_configuration {
      bucket_arn = aws_s3_bucket.kb.arn
    }
  }
}

resource "aws_iam_role" "bedrock_kb" {
  name = "AmazonBedrockExecutionRoleForKnowledgeBase_tmmi"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "bedrock.amazonaws.com"
        }
        Condition = {
          StringEquals = {
            "aws:SourceAccount" = local.account_id
          }
          ArnLike = {
            "aws:SourceArn" = "arn:${local.partition}:bedrock:ap-northeast-2:${local.account_id}:knowledge-base/*"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "kb_model" {
  name = "AmazonBedrockFoundationModelPolicyForKnowledgeBase_tmmi"
  role = aws_iam_role.bedrock_kb.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "bedrock:InvokeModel"
        Effect   = "Allow"
        Resource = data.aws_bedrock_foundation_model.titan.model_arn
      }
    ]
  })
}

resource "aws_iam_role_policy" "kb_oss" {
  name = "AmazonBedrockOSSPolicyForKnowledgeBase_tmmi"
  role = aws_iam_role.bedrock_kb.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "aoss:APIAccessAll"
        Effect   = "Allow"
        Resource = aws_opensearchserverless_collection.kb.arn
      }
    ]
  })
}

resource "time_sleep" "aws_iam_role_policy_kb_oss" {
  create_duration = "30s"
  depends_on      = [aws_iam_role_policy.kb_oss]
}

resource "aws_iam_role_policy" "kb_s3" {
  name = "AmazonBedrockS3PolicyForKnowledgeBase_tmmi"
  role = aws_iam_role.bedrock_kb.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "S3ListBucketStatement"
        Action   = "s3:ListBucket"
        Effect   = "Allow"
        Resource = aws_s3_bucket.kb.arn
        Condition = {
          StringEquals = {
            "aws:PrincipalAccount" = local.account_id
          }
      } },
      {
        Sid      = "S3GetObjectStatement"
        Action   = "s3:GetObject"
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.kb.arn}/*"
        Condition = {
          StringEquals = {
            "aws:PrincipalAccount" = local.account_id
          }
        }
      }
    ]
  })
}

resource "aws_opensearchserverless_collection" "kb" {
  name     = "bedrock-knowledge-base-tmmi"
  type     = "VECTORSEARCH"
  depends_on = [
    aws_opensearchserverless_access_policy.kb,
    aws_opensearchserverless_security_policy.kb_encryption,
    aws_opensearchserverless_security_policy.kb_network
  ]
}

resource "opensearch_index" "kb" {
  name                           = "bedrock-knowledge-base-default-index"
  number_of_shards               = "2"
  number_of_replicas             = "0"
  index_knn                      = true
  index_knn_algo_param_ef_search = "512"
  mappings                       = <<-EOF
    {
      "properties": {
        "bedrock-knowledge-base-default-vector": {
          "type": "knn_vector",
          "dimension": 1024,
          "method": {
            "name": "hnsw",
            "engine": "faiss",
            "parameters": {
              "m": 16,
              "ef_construction": 512
            },
            "space_type": "l2"
          }
        },
        "AMAZON_BEDROCK_METADATA": {
          "type": "text",
          "index": "false"
        },
        "AMAZON_BEDROCK_TEXT_CHUNK": {
          "type": "text",
          "index": "true"
        }
      }
    }
  EOF
  force_destroy                  = true
  depends_on                     = [aws_opensearchserverless_collection.kb]
}

resource "aws_opensearchserverless_access_policy" "kb" {
  name     = "bedrock-knowledge-base-tmmi"
  type     = "data"
  policy = jsonencode([
    {
      Rules = [
        {
          ResourceType = "index"
          Resource = [
            "index/bedrock-knowledge-base-tmmi/*"
          ]
          Permission = [
            "aoss:UpdateIndex",
            "aoss:DescribeIndex",
            "aoss:ReadDocument",
            "aoss:WriteDocument",
            "aoss:CreateIndex"
          ]
        },
        {
          ResourceType = "collection"
          Resource = [
            "collection/bedrock-knowledge-base-tmmi"
          ]
          Permission = [
            "aoss:CreateCollectionItems",
            "aoss:DescribeCollectionItems",
            "aoss:UpdateCollectionItems"
          ]
        }
      ],
      Principal = [
        aws_iam_role.bedrock_kb.arn,
        data.aws_caller_identity.main.arn
      ]
    }
  ])
}

resource "aws_opensearchserverless_security_policy" "kb_encryption" {
  name     = "bedrock-knowledge-base-tmmi"
  type     = "encryption"
  policy = jsonencode({
    Rules = [
      {
        Resource = [
          "collection/bedrock-knowledge-base-tmmi"
        ]
        ResourceType = "collection"
      }
    ],
    AWSOwnedKey = true
  })
}

resource "aws_opensearchserverless_security_policy" "kb_network" {
  name     = "bedrock-knowledge-base-tmmi"
  type     = "network"
  policy = jsonencode([
    {
      Rules = [
        {
          ResourceType = "collection"
          Resource = [
            "collection/bedrock-knowledge-base-tmmi"
          ]
        },
        {
          ResourceType = "dashboard"
          Resource = [
            "collection/bedrock-knowledge-base-tmmi"
          ]
        }
      ]
      AllowFromPublic = true
    }
  ])
}