resource "aws_lambda_function" "main" {
  for_each = local.app

  filename      = data.archive_file.lambda[each.key].output_path
  function_name = "tmmi-${each.key}"
  role          = each.key == "search-advanced" ? aws_iam_role.lambda_im.arn : (each.key == "search-lyrics" ? aws_iam_role.lambda_kb.arn : aws_iam_role.lambda.arn)
  handler       = "bootstrap"

  source_code_hash = data.archive_file.lambda[each.key].output_base64sha256

  runtime       = "provided.al2023"
  architectures = ["arm64"]
  timeout       = 30

  environment {
    variables = each.key == "search-lyrics" ? local.kbinfo : local.dbcred
  }

  logging_config {
    log_format = "Text"
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda,
  ]
}

resource "null_resource" "backend_binary" {
  for_each = local.app

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C backend/${each.value} -o bootstrap -tags lambda.norpc ."
  }
}

data "archive_file" "lambda" {
  for_each    = local.app
  type        = "zip"
  source_file = "backend/${each.value}/bootstrap"
  output_path = "backend/${each.value}/boostrap.zip"

  depends_on = [null_resource.backend_binary]
}

resource "aws_cloudwatch_log_group" "lambda" {
  for_each          = local.app
  name              = "/aws/lambda/tmmi-${each.key}"
  retention_in_days = 30
}

data "aws_iam_policy_document" "lambda" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy" "lambda_basic_execution" {
  name = "AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "invoke" {
  name        = "bedrock_claude_access"
  path        = "/"
  description = "Bedrock Claude Model Access Policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["bedrock:InvokeModel"]
        Effect   = "Allow"
        Resource = data.aws_bedrock_foundation_model.claude.model_arn
      },
    ]
  })
}

resource "aws_iam_policy" "kb" {
  name        = "bedrock_kb_access"
  path        = "/"
  description = "Bedrock Knowledge Base Access Policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["bedrock:Retrieve", "bedrock:RetrieveAndGenerate"]
        Effect   = "Allow"
        Resource = aws_bedrockagent_knowledge_base.kb.arn
      },
    ]
  })
}

resource "aws_iam_role" "lambda" {
  name                = "FunctionExecutionRoleForLambda"
  assume_role_policy  = data.aws_iam_policy_document.lambda.json
  managed_policy_arns = [data.aws_iam_policy.lambda_basic_execution.arn]
}

resource "aws_iam_role" "lambda_im" {
  name                = "ModelInvokeRoleForLambda"
  assume_role_policy  = data.aws_iam_policy_document.lambda.json
  managed_policy_arns = [data.aws_iam_policy.lambda_basic_execution.arn, aws_iam_policy.invoke.arn]
}

resource "aws_iam_role" "lambda_kb" {
  name                = "RagRoleForLambda"
  assume_role_policy  = data.aws_iam_policy_document.lambda.json
  managed_policy_arns = [data.aws_iam_policy.lambda_basic_execution.arn, aws_iam_policy.invoke.arn, aws_iam_policy.kb.arn]
}

resource "aws_lambda_permission" "api" {
  for_each = local.app

  statement_id  = "allowInvokeFromAPIGatewayRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main[each.key].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:ap-northeast-2:${local.account_id}:${aws_apigatewayv2_api.main.id}/*"
}