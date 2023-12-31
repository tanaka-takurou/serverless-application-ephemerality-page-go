# serverless-application-ephemerality kit
Simple kit for serverless application page using AWS Lambda.


## Dependence
- aws-lambda-go
- aws-sdk-go


## Requirements
- AWS (Lambda, API Gateway)
- aws-sam-cli
- golang environment


## Usage

### Deploy
```bash
make clean build
AWS_PROFILE={profile} AWS_DEFAULT_REGION={region} make bucket={bucket} stack={stack name} deploy
```
