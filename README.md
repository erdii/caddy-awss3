# caddy-awss3

### Plan:
* proxy GET file requests to s3 bucket

### Maybe:
* proxy PUT requests to s3 bucket
* proxy DELETE requests to s3 bucket
* proxy GET directory listing requests to s3 bucket

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