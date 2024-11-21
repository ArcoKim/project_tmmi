resource "aws_rds_cluster" "main" {
  cluster_identifier  = "tmmi-postgres"
  engine              = "aurora-postgresql"
  engine_mode         = "provisioned"
  engine_version      = "15.4"
  database_name       = "tmmi"
  master_username     = var.pg_username
  master_password     = var.pg_password
  storage_encrypted   = true
  skip_final_snapshot = true

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  serverlessv2_scaling_configuration {
    max_capacity = 4.0
    min_capacity = 0.5
  }
}

resource "aws_rds_cluster_role_association" "main" {
  db_cluster_identifier = aws_rds_cluster.main.id
  feature_name          = "s3Export"
  role_arn              = aws_iam_role.export.arn
}

resource "aws_iam_role" "export" {
  name = "rds-s3-export-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "rds.amazonaws.com"
        }
        Condition = {
          StringEquals = {
            "aws:SourceAccount" = local.account_id
            "aws:SourceArn"     = aws_rds_cluster.main.arn
          }
        }
      }
    ]
  })
  managed_policy_arns = [aws_iam_policy.export.arn]
}

resource "aws_iam_policy" "export" {
  name   = "rds-s3-export-policy"
  policy = data.aws_iam_policy_document.export.json
}

data "aws_iam_policy_document" "export" {
  statement {
    effect = "Allow"

    actions = [
      "s3:PutObject",
      "s3:AbortMultipartUpload"
    ]

    resources = ["${aws_s3_bucket.kb.arn}/*"]
  }
}

resource "aws_rds_cluster_instance" "main" {
  cluster_identifier  = aws_rds_cluster.main.id
  instance_class      = "db.serverless"
  engine              = aws_rds_cluster.main.engine
  engine_version      = aws_rds_cluster.main.engine_version
  publicly_accessible = true
}

resource "aws_db_subnet_group" "main" {
  name       = "tmmi-subnets"
  subnet_ids = [aws_subnet.public-a.id, aws_subnet.public-b.id]
}

resource "aws_security_group" "rds" {
  name        = "tmmi-rds-sg"
  description = "Allow PostgreSQL inbound traffic and all outbound traffic"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_postgre" {
  security_group_id = aws_security_group.rds.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 5432
  ip_protocol       = "tcp"
  to_port           = 5432
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4" {
  security_group_id = aws_security_group.rds.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}