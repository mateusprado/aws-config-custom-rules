build:
	GOOS=linux go build -o ./bin/ec2_require_tags && zip -r main.zip bin/
