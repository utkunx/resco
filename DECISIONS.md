# Project Decisions and Features

This document tracks important decisions, features, and architectural changes made to the RESCO project.

---

## 2025-10-24

### Enhanced CheckProduct Endpoint Response
**Status**: ✅ Implemented

Added summary statistics and formatted message to `/api/checkproduct/{itemCode}` response:

**New Response Fields**:
- `count-ok`: Total number of products with status "OK" (found in Heihu)
- `count-not`: Total number of products with status "NOT" (not found in Heihu)
- `not-codes`: Formatted string listing all missing products (e.g., "The products that are not app CP20250 + ABC123")

**Example Response**:
```json
{
  "data": [...],
  "count": 3,
  "count-ok": 1,
  "count-not": 2,
  "not-codes": "The products that are not app CP20250 + ABC123",
  "message": "Product check completed successfully"
}
```

**Rationale**:
- Quick summary statistics without parsing entire data array
- Human-readable message for missing products
- Easy integration with reporting and alert systems
- Simplified client-side logic for displaying validation results

**Files**:
- `handlers/bom_handler.go` - Enhanced CheckProduct() with count calculations and message formatting

---

### Git Repository and GitHub Integration
**Status**: ✅ Implemented

Initialized git repository and pushed to GitHub:
- Repository: https://github.com/utkunx/resco
- Created .gitignore to exclude sensitive files (.env, IDE files, binaries)
- Initial commit with 22 files, 3661 lines of code
- Comprehensive commit message documenting all features

**Files Excluded from Git**:
- `.env` - Environment variables with sensitive data
- `.idea/` - IDE configuration files
- Compiled binaries (`resco`)
- OS files (`.DS_Store`)

**Repository Structure**: All source code, documentation, and translation files are version controlled.

---

### Product Check Endpoint with Heihu API Validation
**Status**: ✅ Implemented

Created new `/api/checkproduct/{itemCode}` endpoint that validates all BOM products against Heihu API:

**Functionality**:
1. Takes an item code as input
2. Retrieves all unique product codes from BOM (using `/api/bomtotal/{itemCode}` logic)
3. Queries Heihu API for each product code
4. Returns list with sequence numbers and status ("OK" or "NOT")
5. Implements rate limiting (100ms delay between requests = 10 req/sec, under 20 QPS limit)

**Response Format**:
```json
{
  "data": [
    {"sequence-number": 1, "code": "360004", "status": "OK"},
    {"sequence-number": 2, "code": "CP20250", "status": "NOT"}
  ],
  "count": 2,
  "message": "Product check completed successfully"
}
```

**Rationale**:
- Bulk validation of BOM components against external system
- Identifies which products exist in Heihu system vs. local BOM only
- Helps detect data synchronization issues between systems
- Built-in rate limiting prevents API quota violations

**Rate Limiting Strategy**:
- 100ms delay between consecutive API calls
- Achieves 10 requests/second (50% of allowed 20 QPS)
- Conservative approach to avoid hitting hourly limits on large BOMs

**Files**:
- `services/bom.go` - Added `ProductCheckResult` struct and `CheckProducts()` function
- `handlers/bom_handler.go` - Added `CheckProduct()` handler
- `main.go` - Added route registration

---

### Rate Limiting Configuration for External APIs
**Status**: ✅ Implemented

Created `limits/api_limits.json` to document rate limiting constraints for the Heihu external API:
- **All endpoints**: 10,000 calls/hour, 20 QPS
- **Write endpoints** (create, update, modify): 5,000 calls/hour, 10 QPS
- **Work Order Query / Task List Query**: Once per 10 seconds for identical request parameters

**Rationale**: Need to comply with external API provider's rate limits to avoid service disruption.

**Implementation Note**: Configuration file created. Rate limiting enforcement is now implemented in the CheckProducts function with 100ms delays between requests.

**Files**:
- `limits/api_limits.json`
- `services/bom.go` - Rate limiting implementation in CheckProducts()

---

## Previously Implemented Features

### 1. REST API with Multiple BOM Endpoints
**Status**: ✅ Implemented

Created comprehensive REST API with Gorilla Mux router providing multiple endpoints:

**Endpoints**:
- `GET /health` - Health check endpoint
- `GET /api/bom/{itemCode}` - Basic BOM data with Turkish names
- `GET /api/bomcn/{itemCode}` - BOM data with Chinese translations
- `GET /api/bomcombined/{itemCode}` - BOM data with both Turkish and Chinese names
- `GET /api/bomtotal/{itemCode}` - Unique BOM codes with sequential numbering
- `GET /api/queryhe/{itemCode}` - Query external Heihu API

**Rationale**: Different consumers need different formats - some need Turkish only, some need Chinese translations, some need both for comparison purposes.

**Files**:
- `main.go` - Route definitions
- `handlers/bom_handler.go` - HTTP handlers
- `services/bom.go` - Business logic

---

### 2. Turkish to Chinese Translation Service with Fallback Logic
**Status**: ✅ Implemented

Implemented a smart translation system with two-tier lookup:

**Features**:
- Direct translation lookup from `tr-to-cn.json`
- Fallback translation using first 4 digits of item code from `fallback-tr-to-cn.json`
- In-memory caching of translations after first load
- Thread-safe access with RWMutex

**Translation Flow**:
1. First tries direct Turkish text → Chinese translation
2. If not found, extracts first 4 digits from item code
3. Looks up in fallback dictionary by prefix
4. Returns original text if no translation found

**Rationale**: Not all product names have direct translations. Products with similar codes (same prefix) often have similar naming patterns, so fallback by code prefix provides better coverage.

**Files**:
- `services/translation.go`
- `translate/tr-to-cn.json` - Direct translations
- `translate/fallback-tr-to-cn.json` - Prefix-based fallback translations

---

### 3. Heihu External API Integration
**Status**: ✅ Implemented

Integrated with external Heihu system for product queries:

**Features**:
- POST request to Heihu API with product code
- X-Auth header authentication
- Exact match filtering (only returns products with exact code match)
- Error handling for non-OK responses
- Configuration via environment variables

**Rationale**: Heihu API returns multiple similar products. We filter to exact matches only to avoid confusion and ensure data accuracy.

**Configuration**:
```
HEIHU_LINK=<base_url>
HEIHU_SUB_LINK=<endpoint_path>
X_AUTH=<auth_token>
```

**Files**:
- `services/heihu.go`
- `handlers/bom_handler.go` (QueryHeihu handler)

---

### 4. Recursive BOM Query with SQL Server
**Status**: ✅ Implemented

Implemented recursive Bill of Materials querying using SQL Server CTEs:

**Features**:
- Recursive traversal up to 10 levels deep
- Parameterized queries to prevent SQL injection
- Temporary tables for result processing
- String trimming and data cleaning
- Joins with STOK00 table for item names

**Database Tables**:
- `BOMU01T` - BOM records
- `STOK00` - Stock/item master data

**Rationale**: BOM structures are inherently hierarchical. CTEs provide efficient recursive traversal while depth limit prevents infinite loops in circular references.

**Files**:
- `services/bom.go` (GetBOMByCodeParameterized function)
- `000.sql` - Original SQL script
- `db/connection.go` - Database connection

---

### 5. Multiple BOM Response Formats
**Status**: ✅ Implemented

Created different response formats for different use cases:

**BOMResult** (Turkish only):
```json
{
  "parent-number": "360004",
  "parent-name": "Turkish Name",
  "child-number": "123456",
  "child-name": "Turkish Child Name",
  "child-quantity": 2.5,
  "depth": 1
}
```

**BOMResultCombined** (Turkish + Chinese):
```json
{
  "parent-number": "360004",
  "parent-name": "Turkish Name",
  "parent-name-cn": "中文名称",
  "child-number": "123456",
  "child-name": "Turkish Child Name",
  "child-name-cn": "中文子项名称",
  "child-quantity": 2.5,
  "depth": 1
}
```

**BOMTotalResult** (Unique codes):
```json
{
  "sequence-number": 1,
  "code": "360004"
}
```

**Rationale**: Different consumers have different needs - manufacturing needs Turkish, Chinese factory needs translations, reporting needs unique code lists.

**Files**:
- `services/bom.go` - Struct definitions and functions

---

### 6. Environment-Based Configuration
**Status**: ✅ Implemented

Implemented environment variable configuration using godotenv:

**Configuration Groups**:
- Database connection (DB_SERVER, DB_PORT, DB_USER, DB_PASSWORD, DB_DATABASE)
- Server settings (PORT)
- External API credentials (HEIHU_LINK, HEIHU_SUB_LINK, X_AUTH)

**Files**:
- `.env` - Active configuration (not in git)
- `other/.env.example` - Template for setup
- `main.go` - Configuration loading

**Rationale**: Separates configuration from code, allows different settings per environment (dev/staging/prod), keeps secrets out of source control.

---

## Pending / Future Features

### Rate Limiting Enforcement
**Status**: 📋 Planned

Implement actual rate limiting middleware to enforce limits defined in `limits/api_limits.json`.

**Approach**:
- Use token bucket or sliding window algorithm
- Track requests per endpoint per time window
- Add middleware to check limits before calling Heihu API
- Return 429 Too Many Requests when limit exceeded

**Files to modify**:
- Create `middleware/rate_limiter.go`
- Update `services/heihu.go` to use rate limiter

---

### Unit Tests
**Status**: 📋 Planned

Add comprehensive unit tests for all services and handlers.

**Priority Areas**:
- Translation service with various fallback scenarios
- BOM query result parsing
- Heihu API response filtering
- Error handling

---

### Logging and Monitoring
**Status**: 📋 Planned

Add structured logging for better observability.

**Requirements**:
- Request/response logging
- Database query logging
- External API call logging
- Performance metrics

---

### Authentication/Authorization
**Status**: 📋 To Be Decided

Currently API has no authentication. Decide if needed based on deployment context.

**Options**:
- API key authentication
- JWT tokens
- IP whitelist
- No auth (if behind secure gateway)

---

## Architectural Decisions

### Why Gorilla Mux?
Standard library's `http.ServeMux` is sufficient for basic routing, but Gorilla Mux provides:
- Path variables (`{itemCode}`)
- Method-based routing
- Better middleware support
- Mature and stable

### Why Separate Translation Files?
Direct translations and fallback translations serve different purposes:
- Direct translations are high-confidence, verified mappings
- Fallback translations are probabilistic based on code patterns
- Separation allows easier maintenance and auditing

### Why In-Memory Translation Cache?
Translation dictionaries are:
- Read-heavy (never modified during runtime)
- Relatively small (fits in memory)
- Critical path (used in every translated request)

Caching provides significant performance improvement with minimal memory cost.

### Why Parameterized Queries?
Original implementation used string replacement in SQL file. Switched to parameterized queries to:
- Prevent SQL injection attacks
- Better performance (query plan caching)
- Cleaner code

---

## Known Issues and Limitations

### 1. No Rate Limiting Enforcement
Rate limits are documented but not enforced in code. Application could violate external API limits.

### 2. No Request Caching
Repeated requests for same item code hit database every time. Could benefit from response caching.

### 3. No Pagination
BOM queries can return large result sets. No pagination implemented.

### 4. Single Database Connection
Uses single database connection pool. May need tuning for high concurrency.

### 5. No Graceful Shutdown
Server doesn't handle shutdown signals gracefully. Database connections may not close cleanly.

---

## Change Log Format

When adding new entries, use this format:

```markdown
## YYYY-MM-DD

### Feature/Decision Name
**Status**: ✅ Implemented | 📋 Planned | ⚠️ Deprecated | ❌ Rejected

Brief description of what was done.

**Rationale**: Why this decision was made.

**Files**:
- List of affected files

**Breaking Changes** (if any):
- Description of breaking changes
```