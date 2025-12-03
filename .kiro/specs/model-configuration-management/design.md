# Design Document: Model Configuration Management

## Overview

This design document outlines the implementation of a Model Configuration Management feature for Veritas. The feature enables users to configure, manage, and switch between multiple LLM providers and models through a UI panel. The backend will provide RESTful APIs for CRUD operations on model configurations, securely store credentials, and use these configurations when making LLM requests.

The implementation will extend the existing Veritas architecture by:
- Adding a new `ModelConfig` database model
- Creating new API endpoints for model configuration management
- Modifying the LLM client to use stored configurations instead of environment variables
- Enhancing the frontend UI with a model configuration panel
- Implementing encryption for sensitive credentials

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                    │
│  ┌──────────────────┐  ┌─────────────────────────────────┐ │
│  │  Chat Interface  │  │  Model Config Management Panel  │ │
│  │  - Model Selector│  │  - Create/Edit/Delete Configs   │ │
│  │  - Chat Messages │  │  - Test Connection              │ │
│  └──────────────────┘  └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST API
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go/Gin)                        │
│  ┌──────────────────┐  ┌─────────────────────────────────┐ │
│  │  Chat Handler    │  │  Model Config Handler           │ │
│  │  - Uses Config   │  │  - CRUD Operations              │ │
│  └──────────────────┘  │  - Encryption/Decryption        │ │
│                        │  - Connection Testing            │ │
│  ┌──────────────────┐  └─────────────────────────────────┘ │
│  │  LLM Client      │                                       │
│  │  - Dynamic Config│                                       │
│  └──────────────────┘                                       │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                       │
│  - model_configs (encrypted API keys)                       │
│  - conversations                                             │
│  - messages (with model_config_id reference)                │
└─────────────────────────────────────────────────────────────┘
```

### Component Interaction Flow

1. **Configuration Management Flow**:
   - User creates/edits model config via UI panel
   - Frontend sends request to `/api/model-configs` endpoint
   - Backend validates input, encrypts API key, stores in database
   - Backend returns masked configuration to frontend

2. **Chat Flow with Model Selection**:
   - User selects model from dropdown (now shows user-configured models)
   - User sends message
   - Backend retrieves model config by ID
   - Backend decrypts API key in memory
   - Backend creates LLM client with config parameters
   - Backend makes request to LLM provider
   - Backend stores message with model_config_id reference

## Components and Interfaces

### Backend Components

#### 1. Model Configuration Model (`server/models/model_config.go`)

```go
type ModelConfig struct {
    ID          string    `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"not null" json:"name"`
    Provider    string    `gorm:"not null" json:"provider"` // e.g., "openai", "anthropic", "custom"
    BaseURL     string    `json:"baseUrl"`
    ModelID     string    `gorm:"not null" json:"modelId"`
    APIKey      string    `gorm:"not null" json:"-"` // Encrypted, never sent to client
    IsDefault   bool      `gorm:"default:false" json:"isDefault"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

// ModelConfigResponse is the sanitized version sent to clients
type ModelConfigResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Provider  string    `json:"provider"`
    BaseURL   string    `json:"baseUrl"`
    ModelID   string    `json:"modelId"`
    IsDefault bool      `json:"isDefault"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

#### 2. Message Model Update (`server/models/chat.go`)

```go
type Message struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    ConversationID string    `json:"conversationId"`
    Role           string    `json:"role"`
    Content        string    `json:"content"`
    ModelConfigID  string    `json:"modelConfigId"` // NEW: Track which model was used
    CreatedAt      time.Time `json:"createdAt"`
}
```

#### 3. Model Config Handler (`server/api/model_config_handler.go`)

```go
// CreateModelConfig creates a new model configuration
func CreateModelConfig(c *gin.Context)

// GetModelConfigs returns all model configurations (with masked API keys)
func GetModelConfigs(c *gin.Context)

// GetModelConfig returns a specific model configuration by ID
func GetModelConfig(c *gin.Context)

// UpdateModelConfig updates an existing model configuration
func UpdateModelConfig(c *gin.Context)

// DeleteModelConfig deletes a model configuration
func DeleteModelConfig(c *gin.Context)

// TestModelConfig tests a model configuration without saving
func TestModelConfig(c *gin.Context)
```

#### 4. Encryption Service (`server/services/encryption.go`)

```go
// EncryptAPIKey encrypts an API key using AES-256-GCM
func EncryptAPIKey(plaintext string) (string, error)

// DecryptAPIKey decrypts an encrypted API key
func DecryptAPIKey(ciphertext string) (string, error)
```

#### 5. Updated LLM Client (`server/api/llm_client.go`)

```go
// createLLMClientFromConfig creates an OpenAI client from a ModelConfig
func createLLMClientFromConfig(config *models.ModelConfig) (*openai.Client, error)
```

### Frontend Components

#### 1. Model Config Management Panel (`web/src/components/model-config-panel.tsx`)

A new component that provides:
- List view of all configured models
- Create/Edit form with fields: Name, Provider, Base URL, Model ID, API Key
- Delete confirmation dialog
- Test connection button
- Set as default option

#### 2. Updated Chat Component (`web/src/components/chat.tsx`)

Modifications:
- Model selector now fetches from `/api/model-configs` instead of `/api/models`
- Display model config name in the selector
- Store selected model config ID with messages
- Show which model was used for each message in history

### API Endpoints

#### Model Configuration Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/model-configs` | Create a new model configuration |
| GET | `/api/model-configs` | Get all model configurations |
| GET | `/api/model-configs/:id` | Get a specific model configuration |
| PUT | `/api/model-configs/:id` | Update a model configuration |
| DELETE | `/api/model-configs/:id` | Delete a model configuration |
| POST | `/api/model-configs/test` | Test a model configuration |

#### Request/Response Schemas

**Create/Update Model Config Request**:
```json
{
  "name": "GPT-4 Turbo",
  "provider": "openai",
  "baseUrl": "https://api.openai.com/v1",
  "modelId": "gpt-4-turbo-preview",
  "apiKey": "sk-...",
  "isDefault": false
}
```

**Model Config Response**:
```json
{
  "id": "uuid",
  "name": "GPT-4 Turbo",
  "provider": "openai",
  "baseUrl": "https://api.openai.com/v1",
  "modelId": "gpt-4-turbo-preview",
  "isDefault": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Test Model Config Request**:
```json
{
  "baseUrl": "https://api.openai.com/v1",
  "modelId": "gpt-4-turbo-preview",
  "apiKey": "sk-..."
}
```

**Test Model Config Response**:
```json
{
  "success": true,
  "message": "Connection successful",
  "details": {
    "responseTime": "1.2s",
    "modelAvailable": true
  }
}
```

## Data Models

### Database Schema

#### model_configs Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | VARCHAR(36) | PRIMARY KEY | UUID |
| name | VARCHAR(255) | NOT NULL | User-friendly name |
| provider | VARCHAR(50) | NOT NULL | Provider identifier |
| base_url | VARCHAR(512) | | API base URL |
| model_id | VARCHAR(255) | NOT NULL | Model identifier |
| api_key | TEXT | NOT NULL | Encrypted API key |
| is_default | BOOLEAN | DEFAULT false | Default model flag |
| created_at | TIMESTAMP | | Creation timestamp |
| updated_at | TIMESTAMP | | Last update timestamp |

#### messages Table (Updated)

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Auto-increment ID |
| conversation_id | VARCHAR(36) | FOREIGN KEY | Reference to conversation |
| role | VARCHAR(20) | NOT NULL | "user" or "assistant" |
| content | TEXT | NOT NULL | Message content |
| model_config_id | VARCHAR(36) | FOREIGN KEY | Reference to model config |
| created_at | TIMESTAMP | | Creation timestamp |

### Encryption Strategy

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Source**: Environment variable `ENCRYPTION_KEY` (32 bytes, base64 encoded)
- **Key Rotation**: Support for key rotation by storing key version with encrypted data
- **Format**: `version:nonce:ciphertext:tag` (all base64 encoded)

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Configuration persistence
*For any* valid model configuration (with name, base URL, model ID, and API key), when submitted to the system, the configuration should be persisted to the database and retrievable with all fields intact except the API key which should be encrypted.
**Validates: Requirements 1.1**

### Property 2: API key masking in responses
*For any* request to retrieve model configurations (list or individual), all returned configurations should have their API keys masked or omitted, never exposing the plaintext or encrypted API key to the client.
**Validates: Requirements 1.2, 5.2, 5.4, 6.2**

### Property 3: Configuration update validation
*For any* existing model configuration and valid update parameters, when the update is submitted, the system should validate the new parameters, persist the changes, and return the updated configuration with masked credentials.
**Validates: Requirements 1.3, 6.3**

### Property 4: Configuration deletion completeness
*For any* model configuration, when deleted, the configuration should be removed from the database and any subsequent attempts to use or retrieve it should fail appropriately.
**Validates: Requirements 1.4, 6.4**

### Property 5: Validation error on missing fields
*For any* model configuration submission with one or more required fields missing (name, model ID, or API key), the system should reject the request and return a validation error indicating which fields are missing.
**Validates: Requirements 1.5**

### Property 6: Model selection persistence
*For any* model configuration selected by the user, that configuration should become the active model and be used for all subsequent messages in new conversations until changed.
**Validates: Requirements 2.2, 2.4**

### Property 7: Model switch with history preservation
*For any* active conversation, when the user switches to a different model configuration, the conversation history should remain unchanged and the new model should be used for subsequent messages.
**Validates: Requirements 3.1, 3.2, 3.3**

### Property 8: Message model tracking
*For any* message generated by the system, the message record should include a reference to the model configuration ID that was used to generate it.
**Validates: Requirements 3.4**

### Property 9: Active model display update
*For any* model configuration switch, the UI should immediately update to display the newly selected model's name.
**Validates: Requirements 4.2**

### Property 10: Message history model indication
*For any* message in the conversation history, the UI should indicate which model configuration was used to generate that message.
**Validates: Requirements 4.3**

### Property 11: API key encryption at rest
*For any* model configuration stored in the database, the API key field should be encrypted (not plaintext) and should be decryptable back to the original value.
**Validates: Requirements 5.1**

### Property 12: Configuration creation with unique ID
*For any* valid model configuration submitted via POST request, the system should return the created configuration with a unique identifier that can be used to reference it in future requests.
**Validates: Requirements 6.1**

### Property 13: Configuration retrieval by ID
*For any* existing model configuration ID, when requested via GET, the system should return that specific configuration with masked credentials.
**Validates: Requirements 6.5**

### Property 14: Connection test attempt
*For any* set of model configuration parameters submitted for testing, the system should attempt to establish a connection to the specified LLM endpoint.
**Validates: Requirements 7.1**

### Property 15: Test timeout enforcement
*For any* model configuration test, the system should complete (either successfully or with error) within a reasonable timeout period and not hang indefinitely.
**Validates: Requirements 7.4**

## Error Handling

### Validation Errors

- **Missing Required Fields**: Return 400 Bad Request with details about which fields are missing
- **Invalid URL Format**: Return 400 Bad Request with message about invalid base URL
- **Empty String Values**: Return 400 Bad Request for empty name, model ID, or API key

### Database Errors

- **Connection Failure**: Return 500 Internal Server Error with generic message (log details server-side)
- **Duplicate Configuration**: Return 409 Conflict if a configuration with the same name already exists
- **Foreign Key Violation**: Return 400 Bad Request if trying to delete a config that's referenced by messages

### Encryption Errors

- **Missing Encryption Key**: Fail fast on startup if ENCRYPTION_KEY environment variable is not set
- **Decryption Failure**: Return 500 Internal Server Error and log the error (may indicate key rotation needed)
- **Invalid Key Format**: Fail fast on startup if encryption key is not valid base64 or wrong length

### LLM Provider Errors

- **Connection Timeout**: Return 408 Request Timeout when testing configuration
- **Authentication Failure**: Return 401 Unauthorized when API key is invalid during test
- **Invalid Model ID**: Return 400 Bad Request when model ID doesn't exist at provider
- **Rate Limiting**: Return 429 Too Many Requests and include retry-after information

### Not Found Errors

- **Configuration Not Found**: Return 404 Not Found when requesting non-existent configuration ID
- **No Default Configuration**: Return 400 Bad Request if user tries to chat without selecting a model and no default exists

## Testing Strategy

### Unit Testing

The implementation will include unit tests for:

**Backend Unit Tests**:
- Encryption/decryption functions with various input sizes
- API key masking logic
- Validation functions for model configuration fields
- Database CRUD operations with mock database
- Error handling for various failure scenarios
- LLM client creation from configuration

**Frontend Unit Tests**:
- Model configuration form validation
- Model selector component rendering
- API response handling and error display
- State management for selected model

### Property-Based Testing

Property-based testing will be implemented using **Rapid** (Go's property-based testing library) for the backend. Each property-based test will:
- Run a minimum of 100 iterations with randomly generated inputs
- Be tagged with a comment referencing the specific correctness property from this design document
- Use the format: `**Feature: model-configuration-management, Property {number}: {property_text}**`

**Property Test Coverage**:
- Property 1: Generate random valid configurations, persist them, verify retrieval
- Property 2: Generate random configurations, retrieve them, verify API keys are masked
- Property 3: Generate random configurations and updates, verify persistence
- Property 4: Generate random configurations, delete them, verify removal
- Property 5: Generate configurations with missing fields, verify rejection
- Property 6: Generate random model selections, verify they become active
- Property 7: Generate conversations with model switches, verify history preservation
- Property 8: Generate random messages, verify model config ID is recorded
- Property 11: Generate random API keys, encrypt them, verify they're not plaintext and can be decrypted
- Property 12: Generate random configurations, create them, verify unique IDs
- Property 13: Generate random configurations, retrieve by ID, verify correctness
- Property 14: Generate random configuration parameters, test connections
- Property 15: Generate random configurations, verify tests complete within timeout

**Frontend Property Testing**:
Property-based testing for the frontend will use **fast-check** (TypeScript/JavaScript property-based testing library):
- Property 9: Generate random model switches, verify UI updates
- Property 10: Generate random message histories, verify model indicators

### Integration Testing

Integration tests will verify:
- End-to-end flow: Create config → Select model → Send message → Verify correct model used
- Model switching during conversation maintains history
- Encryption/decryption round-trip with actual database
- API endpoint integration with database operations
- Frontend-backend integration for all CRUD operations

### Security Testing

Security-focused tests will verify:
- API keys are never logged in plaintext
- API keys are never transmitted to client unencrypted
- Encrypted API keys in database cannot be decrypted without proper key
- SQL injection attempts are properly handled
- XSS attempts in configuration names are sanitized

## Implementation Notes

### Migration Strategy

1. **Database Migration**: Add `model_configs` table and `model_config_id` column to `messages` table
2. **Backward Compatibility**: Keep existing environment variable support as fallback
3. **Default Configuration**: On first run, create a default configuration from existing `OPENAI_API_KEY` and `OPENAI_BASE_URL` environment variables if they exist
4. **Gradual Rollout**: Frontend can detect if backend supports new endpoints and fall back to old behavior if not

### Performance Considerations

- **Caching**: Cache decrypted API keys in memory for a short period (5 minutes) to avoid repeated decryption
- **Connection Pooling**: Reuse LLM client instances when possible for the same configuration
- **Lazy Loading**: Only decrypt API keys when actually needed for LLM requests
- **Indexing**: Add database index on `model_configs.name` for faster lookups

### Security Considerations

- **Key Rotation**: Support key rotation by storing key version with encrypted data
- **Audit Logging**: Log all configuration changes (creation, updates, deletion) with timestamps and user info
- **Rate Limiting**: Implement rate limiting on test endpoint to prevent abuse
- **Input Sanitization**: Sanitize all user inputs to prevent XSS and SQL injection
- **HTTPS Only**: Enforce HTTPS in production to protect API keys in transit

