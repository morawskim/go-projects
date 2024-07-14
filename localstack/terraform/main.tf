terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.58.0"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "aws" {
  region  = "us-east-1"
  profile = "localstack"

  s3_use_path_style = true
}

# buckets
resource "aws_s3_bucket" "images" {
  bucket = "images"
}
resource "aws_s3_bucket" "thumbnails" {
  bucket = "thumbnails"
}


# IAM role and policy for lambda
data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
resource "aws_iam_role" "lambda" {
  name               = "AssumeLambdaRole"
  description        = "Role for lambda to assume lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}


# create lambda function
resource "null_resource" "function_binary" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../lambda-test ../main.go"
  }
}
data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = "../lambda-test"
  output_path = "../lambda.zip"
}
resource "aws_lambda_function" "thumbnails" {
  function_name = "thumbnails"
  description   = "Create thumbnails"
  role          = aws_iam_role.lambda.arn
  handler       = "lambda-test"
  memory_size   = 128

  filename         = "../lambda.zip"
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "go1.x"
}


# connect s3 bucket notification with lambda
resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.thumbnails.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.images.arn
}
resource "aws_s3_bucket_notification" "thumbnails" {
  bucket = aws_s3_bucket.images.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.thumbnails.arn
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = ".jpg"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}
