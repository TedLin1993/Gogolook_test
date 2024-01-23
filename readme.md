# Task Management API

This project provides a simple RESTful API for managing tasks. It includes endpoints for creating, retrieving, updating, and deleting tasks.

## Endpoints

* **GET /tasks:** Retrieve a list of all tasks.
* **POST /tasks:** Create a new task.
* **PUT /tasks/{id}:** Update an existing task.
* **DELETE /tasks/{id}:** Delete a task.

### Task Structure

A task consists of the following fields:

* **ID:** Unique identifier for the task.
* **Name:** Task name.
* **Status:** Task status (0 for incomplete, 1 for completed).

## Running the Application

### Using Docker

You can run the application using Docker and Docker Compose.

```bash
docker-compose up -d
```

This command will build and start the application in detached mode, making it accessible at `http://localhost:8080`.

### Running Tests

To run tests, you can use the following command:

```bash
go test ./handler
```

This will execute tests in the `handler` package and provide feedback on the test results.
