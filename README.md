
# Go SQLite API Example

This is a simple example of a Go API using SQLite for storing and retrieving data. The API has endpoints for posting data (`/post-data`) and getting data (`/get-data`).

## Setup

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/Gulshan256/golang-API-save-data-sql.git
   cd golang-API-save-data-sql
   ```

2. **Run the Application:**
   ```bash
   go run main.go
   ```

   The server will start listening on `http://localhost:8080`.

3. **Test Endpoints:**
   - POST data: Send a POST request to `http://localhost:8080/post-data` with JSON data in the request body.
   - GET data: Send a GET request to `http://localhost:8080/get-data` to retrieve stored data.

## API Endpoints

### 1. POST Data (`/post-data`)

#### Request:

```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "123-456-7890",
  "from": "City A",
  "componentName": "Component X"
}
```

#### Response:

```json
{
  "message": "Data received and stored successfully"
}
```

### 2. GET Data (`/get-data`)

#### Response:

```json
[
  {
    "date": "2022-01-10",
    "time": "08:30:00",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "123-456-7890",
    "from": "City A",
    "componentName": "Component X"
  },
  // Additional data entries...
]
```

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux): A powerful URL router and dispatcher for Go.

- [go-sqlite3](https://github.com/mattn/go-sqlite3): SQLite3 driver for Go using database/sql.


