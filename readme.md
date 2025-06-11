# Short URL Service

A simple and efficient URL shortening service built with Go. This service allows you to convert long URLs into shorter, more manageable URLs.

## Features

- URL shortening with random 8-character codes
- URL validation
- In-memory storage of URL mappings
- Simple and RESTful API endpoints

## API Endpoints

### 1. Create Short URL
- **Endpoint**: `POST /shorten`
- **Request Body**:
  ```json
  {
    "url": "https://example.com/very/long/url"
  }
  ```
- **Response**:
  ```json
  {
    "short_url": "http://localhost:8080/AbCdEfGh",
    "original_url": "https://example.com/very/long/url"
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid URL format or request body
  - 405: Method not allowed

### 2. Redirect to Original URL
- **Endpoint**: `GET /{shortCode}`
- **Description**: Redirects to the original URL associated with the short code
- **Status Codes**:
  - 307: Temporary redirect to original URL
  - 400: Missing short URL
  - 404: URL not found

### 3. Check Storage
- **Endpoint**: `GET /check`
- **Description**: Returns all stored URL mappings
- **Response**: JSON object containing all short URL to original URL mappings

## Running the Application

1. Make sure you have Go installed (version 1.24.0 or later)
2. Clone the repository
3. Run the application:
   ```bash
   go run cmd/main.go
   ```
4. The server will start on `http://localhost:8080`

## Example Usage

1. Create a short URL:
   ```bash
   curl -X POST http://localhost:8080/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com/very/long/url"}'
   ```

2. Access the shortened URL:
   - Open `http://localhost:8080/{shortCode}` in your browser
   - Or use curl: `curl -L http://localhost:8080/{shortCode}`

3. Check all stored URLs:
   ```bash
   curl http://localhost:8080/check
   ```

## Notes

- The service uses in-memory storage, so all URL mappings will be lost when the server restarts
- Short URLs are generated using a random 8-character base64 encoded string
- The service validates URLs to ensure they have a valid scheme and host
