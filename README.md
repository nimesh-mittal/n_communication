# n_communication

[![Go Report Card](https://goreportcard.com/badge/github.com/nimesh-mittal/n_communication)](https://goreportcard.com/report/github.com/nimesh-mittal/n_communication)

## Background

Aim of this service is to provide ability to send messages to channels like email, sms, push or messaging queue. This service avoids the need for every other micro services to integrate with these channels.

## Requirments

Service should provide following abilities

- Ability to send email message
- Ability to send sms message
- Ability to send message to messaging bus like Kafka or Kinesis
- Ability to send messages in async mode with an ability to retry on failure

## Service SLA

- Availability
  - Communication service should target 99.99% uptime
- Latency
  - Communication service should aim for less than 10 ms P95 latency
- Throughput
  - Communication service should provide 5000 QPS per node
- Freshness
  - Communication service should be able to send messages in less than 2 mins. Please refere SLA of each channel below

## Architecture

![image](https://github.com/nimesh-mittal/n_communication/blob/main/.github/images/arch.png)

## Implementation

### API Interface

```go

type Interface{
  // Send message to channel with provided payload
  Send(channel string, to string, from string, payload string, title string) (bool, error)
}

```

## Data Model

| Table Name | Description | Columns |
| ------- | ---- | ---- |
| messages | Represents messages sent | (tenant_name, msg_id, channel, *attributes...*, *who...*)

where attributes are:

- To              string
- From            string
- Payload         string
- Status          string
- RetryCount      int
- LastSendAtemptAt    time.Time

Who columns are"

- Active    bool
- CreatedBy string
- CreatedAt int64
- UpdatedBy string
- UpdatedAt int64
- DeletedBy string
- DeletedAt int64

## Database choice

A close look at the API request reveals large amount to write than read requests. Data consistency is important.

Postgres can be used to store all the messages.

## Channels
### Email
For email use Amazon AWS SES.

### SMS
For sms use Twilio service. Signup for trial account to test.

## Scalability and Fault tolerance

Inorder to survive host failure, multiple instances of the service needs to be deployed behind load balancer. Load balance should detect host failure and transfer any incoming request to only healthy node. One choice is to use ngnix server to perform load balancing.

Given one instance of service can serve not more than 5000 request per second, one must deploy more than one instance to achive required throughput from the system

Load balancer should also rate limiting incoming requests to avoid single user/client penalising all other user/client due to heavy load.

Given the service is going to perform more write request than read, we can perform async writes to increase throughput.

## Functional and Load testing

Service should implement good code coverage and write functional and load test cases to maintain high engineering quality standards.

## Logging, Alerting and Monitoring

Service should also expose health end-point to accuretly monitor the health of the service.

Service should also integrate with alerting framework like new relic to accuretly alert developers on unexpected failured and downtime

Service should also integrate with logging library like zap and distributed tracing library like Jager for easy debugging

## Security

Service should validated the request using Oauth token to ensure request is coming from authentic source

## Documentation

This README file provides complete documentation. Link to any other documentation will be provided in the Reference section of this document.

## Local Development Setup

- Setup environment variables
  - export AWS_REGION="ap-south-1"
  - export AWS_SECRET_ID="<key here>"
  - export AWS_SECRET="<secret here>"
  - export TWILIO_ACCOUNT_SID="<AC....>"
  - export TWILIO_AUTH_TOKEN="<token here>"
  - export TWILIO_URL="< https://api.twilio.com/2010-04-01/Accounts/<account sid>/Messages.json >"

- Start service

```shell
  go run main.go
```

- Run testcases with coverage

```shell
go test ./... -cover
```

- How to generate mocks in local

```shell
  go get github.com/golang/mock/gomock
```

```shell
  go get github.com/golang/mock/mockgen
```

```shell
  <path to bin>/mockgen -destination=mocks/mock_profilerepo.go -package=mocks n_users/repo ProfileRepo
```
