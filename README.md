# Project Description

API in Golang for controlling executed tasks, with user control.
Access validation using JWT, integration with RabbitMq for sending notifications.

# Technologies Used

**Backend**
*Go
*MySql
*RabbitMq
*Docker
*Mockery
*JWT
*Go frameworks(Gin, Gorm, etc)

# How to run the project

**Run Docker environment**

1. Configure environment variables in the .env file. Like **.env.example** file.

2. Access the root folder of the project where the file (docker-compose.yml).

3. Execute the command: **docker-compose up** (NOTE: this command will create the images and run the applications' docker containers.)
 
The API will be available at: http://localhost:8080

# Consume the endpoint's

With the containers running, the endpoints are ready to be consumed.

The api has two services(User and Task):

**User:**

1. **GET** http://localhost:8080/user  List of Users
2. **POST** http://localhost:8080/user   Create new User
3. **POST** http://localhost:8080/user/login  Login returns token JWT


**Task:**

1. **GET** http://localhost:8080/task  List of Tasks
2. **GET** http://localhost:8080/task/:id  List Task By ID
3. **POST** http://localhost:8080/task  Create a Task
4. **PATCH** http://localhost:8080/task/:id  Update Task Info
5. **PATCH** http://localhost:8080/task/execute/:id  Complete a task
6. **DELETE** http://localhost:8080/task/:id  Delete a Task