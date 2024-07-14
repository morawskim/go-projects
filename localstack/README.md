# localstack

Localstack is a tool that allows you to run and test applications relying on AWS services locally, 
without the need to connect to a real AWS environment.

This demo provision S3 buckets (images and thumbnails) and lambda function (thumbnails)
via terraform. 
When you upload a new object into images bucket, the same object will be uploaded to thumbnails bucket.
Thanks to S3 bucket notifications.

## Usage

**You need configured aws-cli with localstack.**
**The profile should be named localstack**

Call command - `docker compose up -d`
Go to directory terraform and apply plan - `cd terraform && tofu apply`
Approve changes by type "yes"
Return to main directory - "cd .."

Run `aws s3 --profile localstack cp ./216-800x900.jpg s3://images/photo1.jpg`
The image should exist `aws s3 ls --profile localstack s3://images/`
The thumbnail also should exist `aws s3 ls --profile localstack s3://thumbnails/`

To check log group of lambda function - `aws lambda get-function --function-name thumbnails --profile localstack`
To get all log streams `aws logs describe-log-streams --profile localstack --log-group-name <value of field Configuration.LoggingConfig.LogGroup from one of previous command eg. /aws/lambda/thumbnails>`
To get events from log stream `aws logs get-log-events --profile localstack --log-group-name <value of field Configuration.LoggingConfig.LogGroup from one of previous command eg. /aws/lambda/thumbnails> --log-stream-name <value of field logStreams.0.logStreamName from previous command>`
