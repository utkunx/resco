# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**RESCO** is a production-ready Go REST API for querying Bill of Materials (BOM) data from SQL Server with Turkish-to-Chinese translation support and external API integration (Heihu system).

### Key Features
- RESTful API with multiple endpoints for BOM queries
- Recursive BOM traversal (up to 10 levels deep)
- Turkish to Chinese translation with fallback logic
- External Heihu API integration for product queries
- Rate limiting configuration for external API calls
- SQL Server database integration

## Development Commands

### Build and Run
```bash
# Run the main application
go run main.go

# Build the application
go build -o resco main.go

# Run the compiled binary
./resco
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -run TestName ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint installation)
golangci-lint run

# Vet code for potential issues
go vet ./...
```

## Project Structure

```
resco/
├── main.go                    # Application entry point, HTTP server, routes
├── 000.sql                    # SQL script for recursive BOM queries
├── go.mod                     # Go module definition (Go 1.24)
├── go.sum                     # Dependency checksums
├── .env                       # Environment configuration (not in git)
├── README.md                  # Project documentation
├── CLAUDE.md                  # This file - AI assistant context
│
├── db/
│   └── connection.go          # SQL Server database connection management
│
├── handlers/
│   └── bom_handler.go         # HTTP request handlers for all endpoints
│
├── services/
│   ├── bom.go                 # Business logic for BOM queries
│   ├── translation.go         # Turkish to Chinese translation service
│   └── heihu.go               # External Heihu API integration
│
├── translate/
│   ├── tr-to-cn.json          # Direct Turkish to Chinese translations
│   ├── fallback-tr-to-cn.json # Fallback translations by item code prefix (first 4 digits)
│   └── .backupfallback-tr-to-cn.json
│
├── limits/
│   └── api_limits.json        # External API rate limiting configuration
│
├── outputs/
│   ├── first.json             # Sample API responses
│   ├── third.json
│   ├── fourth.json
│   └── final.json
│
└── other/
    └── .env.example           # Environment configuration template
```

## Database Context

### Database: RESCO_2019 (SQL Server)

**Main Tables:**
- **BOMU01T**: BOM records (Bill of Materials entries)
- **STOK00**: Stock/item master data (product information)

**Key Fields:**
- `BOMREC_CODE`: Item/product code
- `BOMREC_KAYNAKCODE`: Source/parent item code (recursive relationship)
- `BOMREC_KAYNAK0`: Quantity
- `BOMREC_INPUTTYPE`: Input type filter (uses 'H')
- `AD`: Turkish item name

**Query Logic:**
The SQL script uses CTEs (Common Table Expressions) to traverse parent-child relationships in the BOM hierarchy:
1. Starts from a specific item code
2. Recursively finds all related components
3. Joins with STOK00 for item names
4. Returns hierarchy with depth levels
5. Limits recursion to 10 levels to prevent infinite loops

## API Endpoints

The API runs on port 8080 (configurable via PORT env var) and provides:

### Health Check
```
GET /health
```
Returns server status.

### BOM Endpoints
```
GET /api/bom/{itemCode}
```
Returns BOM data with Turkish names only.

```
GET /api/bomcn/{itemCode}
```
Returns BOM data with Chinese translations applied.

```
GET /api/bomcombined/{itemCode}
```
Returns BOM data with both Turkish and Chinese names (combined).

```
GET /api/bomtotal/{itemCode}
```
Returns unique BOM codes with sequential numbers (flattened list).

### External API Integration
```
GET /api/queryhe/{itemCode}
```
Queries the external Heihu API for product information. Returns exact matches only.

```
GET /api/checkproduct/{itemCode}
```
Checks all BOM products against Heihu API. Returns a list with sequence numbers and status ("OK" if product exists in Heihu, "NOT" if not found). Implements rate limiting with 100ms delays between requests.

## Services Architecture

### BOM Service (`services/bom.go`)
- **GetBOMByCodeParameterized**: Executes parameterized recursive BOM query
- **GetBOMByCodeWithTranslation**: Returns BOM with Chinese translations
- **GetBOMByCodeCombined**: Returns BOM with both Turkish and Chinese
- **GetBOMTotal**: Returns unique codes with sequential numbering

### Translation Service (`services/translation.go`)
- **LoadTranslations**: Loads direct Turkish to Chinese mappings
- **LoadFallbackTranslations**: Loads fallback translations by item code prefix
- **Translate**: Simple direct translation
- **TranslateWithFallback**: Smart translation with fallback logic based on first 4 digits of item code
- **ApplyTranslationsToBOM**: Applies translations to BOM result sets

Translation logic:
1. First tries direct translation from `tr-to-cn.json`
2. If not found, uses first 4 digits of item code to find fallback in `fallback-tr-to-cn.json`
3. Returns original text if no translation found

### Heihu Service (`services/heihu.go`)
- **QueryHeihu**: Queries external Heihu API
- Filters results to return only exact product code matches
- Uses environment variables for API configuration

## Environment Configuration

Required environment variables (set in `.env` file):

### Database Configuration
```
DB_SERVER=localhost          # SQL Server hostname
DB_PORT=1433                 # SQL Server port
DB_USER=sa                   # Database username
DB_PASSWORD=your_password    # Database password
DB_DATABASE=RESCO_2019       # Database name
```

### Server Configuration
```
PORT=8080                    # HTTP server port
```

### Heihu API Configuration
```
HEIHU_LINK=https://...       # Heihu API base URL
HEIHU_SUB_LINK=/api/...      # Heihu API endpoint path
X_AUTH=your_token            # Heihu API authentication token
```

## External API Rate Limits

Configuration stored in `limits/api_limits.json`:

### All Endpoints
- Hourly limit: 10,000 calls
- QPS: 20 queries per second

### Write Endpoints (create, update, modify)
- Hourly limit: 5,000 calls
- QPS: 10 queries per second

### Specific Endpoints
- Work Order Query Interface: Once per 10 seconds for identical parameters
- Task List Query Interface: Once per 10 seconds for identical parameters

**Note**: Rate limiting is currently documented but not enforced in code. Implementation needed if required.

## Dependencies

```go
github.com/gorilla/mux              // HTTP router
github.com/joho/godotenv            // Environment variable loader
github.com/microsoft/go-mssqldb     // SQL Server driver
```

## Data Flow

1. **HTTP Request** → Handler (`handlers/bom_handler.go`)
2. **Handler** → Service (`services/bom.go`)
3. **Service** → Database (`db/connection.go`)
4. **Database** → Returns raw BOM data
5. **Service** → Applies translations if needed (`services/translation.go`)
6. **Service** → Returns processed data
7. **Handler** → Formats JSON response
8. **HTTP Response** → Client

## Common Tasks

### Adding New Translations
1. Edit `translate/tr-to-cn.json` for direct translations
2. Edit `translate/fallback-tr-to-cn.json` for code-prefix-based fallbacks
3. Translations are loaded once at first use and cached in memory

### Modifying BOM Query Logic
1. Update SQL in `services/bom.go` (GetBOMByCodeParameterized function)
2. Or modify `000.sql` if using file-based approach

### Adding New Endpoints
1. Add route in `main.go`
2. Create handler in `handlers/bom_handler.go`
3. Add business logic in appropriate service file

## Important Notes

- The application uses parameterized queries to prevent SQL injection
- Translations are cached in memory after first load
- Heihu API integration filters results to exact matches only
- BOM recursion depth is hardcoded to 10 levels maximum
- Database connection uses standard `database/sql` with SQL Server driver
