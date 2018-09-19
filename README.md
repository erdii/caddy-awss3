# caddy-awss3

### Plan:
* proxy GET file requests to s3 bucket

### Todo:
* use GetObject with a limited range instead of HeadObject
	* does this work?
	* idea is:
		1. fetch first chunk of file
		1. if chunk is "full", take headers for file, stream rest of file via s3manager.downloader
		1. if chunk is "not full", send response
	* https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.GetObject
* tests

### Maybe:
* PUT requests to upload objects
* DELETE requests to delete objects
* GET directory listings

### IAM policy:
```JSON
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::your-bucket-name"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject"
            ],
            "Resource": [
                "arn:aws:s3:::your-bucket-name/*"
            ]
        }
    ]
}
```

### Attribution

Some code has been inspired by or taken directly from [coopernurse/caddy-awslambda](https://github.com/coopernurse/caddy-awslambda). Check out this nice project if you want an alternative API-Gateway for AWS Lambda!

* config.go
	* .ToAwsConfig()
	* .ParseConfigs()
	* .StripPathPrefix()
* config_test.go
	* .TestToAwsConfigStaticCreds()
	* .TestToAwsConfigDefaults()
	* .TestParseConfigs()
* middleware.go
	* .match()

The copied code is licensed under the [MIT License](https://github.com/coopernurse/caddy-awslambda/blob/761c41e19aed5db2d4bffcc624f353bf70cb2407/LICENSE), too.