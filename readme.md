# Task Management API

This project provides a simple RESTful API for managing tasks. It includes endpoints for creating, retrieving, updating, and deleting tasks.

## Endpoints

### GET /tasks

Retrieves a list of all tasks.

#### Response

* **Status Code** : `200 OK`
* **Body** :
* `Array` of tasks, where each task includes:
  * `id`: String, unique identifier of the task
  * `name`: String, name of the task
  * `status`: String, status of the task (`"Incomplete"` or `"Completed"`)

### POST /tasks

Creates a new task.

#### Request Body

* `name`: String, name of the task (required)
* `status`: String, status of the task, either `"Incomplete"` or `"Completed"` (required)

#### Response

* **Status Code** : `201 Created` on success, `400 Bad Request` on failure
* **Body** :
* On success: The created task object, including `id`, `name`, and `status`.
* On failure: An error object with an `error` key containing the error message.

### PUT /tasks/:id

Updates an existing task.

#### Parameters

* `id`: Path parameter, the unique identifier of the task to update.

#### Request Body

* `name`: String, new name of the task (optional)
* `status`: String, new status of the task, either `"Incomplete"` or `"Completed"` (optional)

#### Response

* **Status Code** : `200 OK` on success, `400 Bad Request` or `404 Not Found` on failure
* **Body** :
* On success: The updated task object, including `id`, `name`, and `status`.
* On failure: An error object with an `error` key containing the error message.

### DELETE /tasks/:id

Deletes a task.

#### Parameters

* `id`: Path parameter, the unique identifier of the task to delete.

#### Response

* **Status Code** : `200 OK` on success, `404 Not Found` on failure
* **Body** :
* On success: A message object with a `message` key indicating successful deletion.
* On failure: An error object with an `error` key containing the error message.

## Status Codes

The API uses the following status codes:

* `200 OK`: The request was successful.
* `201 Created`: A new resource was successfully created.
* `400 Bad Request`: The request was not processed due to client error.
* `404 Not Found`: The requested resource was not found.
* `422 Unprocessable Entity`: The request was well-formed but was unable to be followed due to semantic errors.

## Errors

Errors are returned as JSON objects in the following format:

```
{
  "error": "Error message describing the problem"
}
```

## Task Structure

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
