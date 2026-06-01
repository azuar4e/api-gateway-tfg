# API-Gateway

![Go](https://img.shields.io/badge/Go-1.24-blue)
![AWS](https://img.shields.io/badge/AWS-EKS%20%7C%20RDS%20%7C%20DynamoDB%20%7C%20SQS%20%7C%20SNS-orange)

The source code of this repository corresponds to an API gateway microservice for my TFG (bachelor's thesis). The API is the entry point to the system, it uses the Gin library — a Go framework orientated to the development of APIs REST — and the AWS SDK for Go, necessary for the integration with the managed services used in the architecture.

The complete explanation of the project can be found at my [TFG](https://github.com/azuar4e/tfg) repository.

## Overview

The API has the following features:

- Is exposed in the port `9090`.
- Stores the user data in a PostgreSQL RDS.
- Stores the jobs information in a DynamoDB table.
- Subscribes the users in a Simple Notification Service (SNS) topic.
- Queues the jobs in a Simple Queue Service (SQS) queue, which are process by the [scraper service](https://github.com/azuar4e/scraper-tfg).

## How it works?

The API exposes the following *endpoints*:

```Go
v1 := r.Group("/api/v1")

v1.POST("/signin", controllers.SigninHandler)
v1.POST("/signup", controllers.RegisterHandler)

v1.Use(middleware.AuthMiddleware())

v1.POST("/jobs", handlers.CreateJobHandler)
v1.GET("/jobs", handlers.GetJobsHandler)
v1.DELETE("/jobs", handlers.DeleteJobsHandler)
v1.GET("/jobs/:id", handlers.GetJobByIdHandler)
v1.DELETE("/jobs/:id", handlers.DeleteJobHandler)
v1.GET("/validate", controllers.Validate)
```

First of all, the user signup and signin, where it will require a username and a password, then a cookie will be generated through **JWT**, which is needed for the authentication to the protected endpoints. The middleware function will check the reliablity of the cookie. If everything is correct then the user will be able to create, list and delete jobs.

## Requirements

- Go 1.24+
- AWS credentials configured
