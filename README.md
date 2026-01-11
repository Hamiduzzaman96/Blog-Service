# Blog SERVICE

A production-ready microservice-based blogging platform built with Golang, following Clean Architecture and Event-Driven Design.  

The platform consists of four main services:

1. User Service – Manages users and authentication (JWT-based).  
2. Author Service – Handles author promotion and author-specific operations.  
3. Blog Service – Manages blog posts and publishes events asynchronously via RabbitMQ.  
4. Notification Service – Sends notifications to users asynchronously.

---

## Table of Contents

A. [Requirements](#requirements)  
B. [Setup & Run](#setup--run)  
C. [Services & Endpoints](#services--endpoints)  
D. [gRPC Endpoints](#grpc-endpoints)  
E. [Architecture Notes](#architecture-notes)  

---

A. Requirements

- Go >= 1.21  
- PostgreSQL >= 14  
- Redis  
- RabbitMQ  
- `protoc` and gRPC Go plugin for gRPC services  

---

B. Setup & Run

1. Clone the repository

```bash
git clone <repository_url>
cd Blog-Service

2.Copy environment file
   cp .env.example .env
3. Edit .env with your credentials for PostgreSQL, Redis, RabbitMQ, and JWT secret.
4. Install Go dependencies
5. Run services (in separate terminals):
 ## User Service
> cd cmd/user
> go run main.go

## Author Service
> cd cmd/author
> go run main.go

## Blog Service
> cd cmd/blog
> go run main.go

## Notification Service
> cd cmd/notification
> go run main.go

C. Services & HTTP Endpoints

>> User Service
Method	  Path	             Description
POST	/register	      Register a new user
POST	/login	             User login

>> Author Service
Method	    Path	            Description
POST	/become-author	  Promote a user to author

D. gRPC Endpoints

Proto files located in proto/ directory:

##       Service	            Method	            Request	                Response
        UserService	         RegisterUser	     RegisterRequest	    RegisterResponse
        UserService	          LoginUser	          LoginRequest	          LoginResponse
        AuthorService	     BecomeAuthor	   BecomeAuthorRequest	   BecomeAuthorResponse
        BlogService	          CreatePost	    CreatePostRequest	      BlogResponse
        NotificationService  SendNotification	NotificationRequest	  NotificationResponse


E. Architecture Notes

>> Clean Architecture: Domain → Repository → Usecase → Handler.

>> JWT: Used for authentication across User, Author, and Blog services.

>> Redis: Used for session and token management.

>> RabbitMQ: Event-driven communication for blog creation and notifications.

>> gRPC + HTTP: All services expose both HTTP and gRPC endpoints.



