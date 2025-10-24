# RESCO - Bill of Materials API

A Go REST API for querying Bill of Materials (BOM) data from SQL Server using recursive queries.

## Features

- RESTful API endpoint for BOM queries
- Recursive BOM traversal (up to 10 levels deep)
- SQL Server integration
- JSON response format

## Project Structure

```
resco/
├── main.go              # Application entry point and HTTP server setup
├── 000.sql              # SQL script for recursive BOM queries
├── db/
│   └── connection.go    # Database connection management
├── handlers/
│   └── bom_handler.go   # HTTP request handlers
└── services/
    └── bom.go           # Business logic for BOM queries
```

## Prerequisites

- Go 1.24 or higher
- SQL Server with RESCO_2019 database
- Access to BOMU01T and STOK00 tables

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Configure database connection by setting environment variables or creating a `.env` file:
```bash
export DB_SERVER=your_server
export DB_PORT=1433
export DB_USER=your_username
export DB_PASSWORD=your_password
export DB_DATABASE=RESCO_2019
export PORT=8080
```

## Running the API

```bash
# Run directly
go run main.go

# Or build and run
go build -o resco
./resco
```

The server will start on `http://localhost:8080` (or your configured PORT).

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "ok",
  "message": "API is running"
}
```

### Get BOM by Item Code
```
GET /api/bom/{itemCode}
```

Parameters:
- `itemCode` (path parameter): The item code to search for (e.g., "360004")

Example:
```bash
curl http://localhost:8080/api/bom/360004
```

Response:
```json
{
  "data": [
    {
      "parent-number": "360004",
      "parent-name": "Item Name",
      "par_pro_spec": "",
      "child-number": "SOURCE123",
      "child-name": "Sub Item Name",
      "sub_pro_spec": "",
      "child-quantity": 1.5,
      "depth": 1
    }
  ],
  "count": 1,
  "message": "BOM data retrieved successfully"
}
```

### Get BOM with Chinese Translations
```
GET /api/bomcn/{itemCode}
```

Returns BOM data with Turkish names translated to Chinese.

### Get BOM Combined (Turkish + Chinese)
```
GET /api/bomcombined/{itemCode}
```

Returns BOM data with both Turkish and Chinese names.

### Get BOM Total (Unique Codes)
```
GET /api/bomtotal/{itemCode}
```

Returns all unique product codes from the BOM with sequential numbers.

Response:
```json
{
  "data": [
    {"sequence-number": 1, "code": "360004"},
    {"sequence-number": 2, "code": "CP20250"}
  ],
  "count": 2,
  "message": "Unique BOM codes retrieved successfully"
}
```

### Query Heihu API
```
GET /api/queryhe/{itemCode}
```

Queries external Heihu API for a specific product code.

### Check Products Against Heihu API
```
GET /api/checkproduct/{itemCode}
```

Checks all BOM products against Heihu API and returns their availability status.

Example:
```bash
curl http://localhost:8080/api/checkproduct/360004
```

Response:
```json
{
  "data": [
    {"sequence-number": 1, "code": "360004", "status": "OK"},
    {"sequence-number": 2, "code": "CP20250", "status": "NOT"},
    {"sequence-number": 3, "code": "ABC123", "status": "NOT"}
  ],
  "count": 3,
  "count-ok": 1,
  "count-not": 2,
  "not-codes": "The products that are not app CP20250 + ABC123",
  "message": "Product check completed successfully"
}
```

**Response Fields**:
- `data`: Array of all products with their check status
- `count`: Total number of products checked
- `count-ok`: Number of products found in Heihu system
- `count-not`: Number of products NOT found in Heihu system
- `not-codes`: Formatted string listing all missing product codes
- `message`: Status message

**Note**: This endpoint implements rate limiting (100ms delay between requests) to comply with Heihu API limits. Response time will scale with the number of products in the BOM.

Error Response:
```json
{
  "error": "Error message here"
}
```

## Configuration

The application reads configuration from environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| DB_SERVER | SQL Server hostname | localhost |
| DB_PORT | SQL Server port | 1433 |
| DB_USER | Database username | sa |
| DB_PASSWORD | Database password | (empty) |
| DB_DATABASE | Database name | RESCO_2019 |
| PORT | HTTP server port | 8080 |

## SQL Query Details

The API executes the recursive SQL query from `000.sql` which:
1. Searches for initial item in BOMREC_CODE
2. Recursively finds related items through BOMREC_KAYNAKCODE
3. Joins with STOK00 table for item names
4. Returns hierarchy with depth levels
5. Limits recursion to 10 levels to prevent infinite loops

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Build
```bash
go build -o resco main.go
```

## License

[Add your license here]