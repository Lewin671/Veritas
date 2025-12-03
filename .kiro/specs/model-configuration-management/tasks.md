# Implementation Plan

- [x] 1. Set up encryption service and database models
  - Create encryption service with AES-256-GCM for API key encryption/decryption
  - Add ENCRYPTION_KEY environment variable to .env.example
  - Create ModelConfig database model with all required fields
  - Update Message model to include model_config_id field
  - Update database initialization to auto-migrate new models
  - _Requirements: 5.1, 5.3_

- [ ]* 1.1 Write property test for encryption round-trip
  - **Property 11: API key encryption at rest**
  - **Validates: Requirements 5.1**

- [x] 2. Implement model configuration CRUD API endpoints
  - Create model_config_handler.go with handler functions
  - Implement CreateModelConfig endpoint (POST /api/model-configs)
  - Implement GetModelConfigs endpoint (GET /api/model-configs)
  - Implement GetModelConfig endpoint (GET /api/model-configs/:id)
  - Implement UpdateModelConfig endpoint (PUT /api/model-configs/:id)
  - Implement DeleteModelConfig endpoint (DELETE /api/model-configs/:id)
  - Add input validation for all endpoints
  - Ensure API keys are masked in all responses
  - Register routes in main.go
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ]* 2.1 Write property test for configuration persistence
  - **Property 1: Configuration persistence**
  - **Validates: Requirements 1.1**

- [ ]* 2.2 Write property test for API key masking
  - **Property 2: API key masking in responses**
  - **Validates: Requirements 1.2, 5.2, 5.4, 6.2**

- [ ]* 2.3 Write property test for configuration updates
  - **Property 3: Configuration update validation**
  - **Validates: Requirements 1.3, 6.3**

- [ ]* 2.4 Write property test for configuration deletion
  - **Property 4: Configuration deletion completeness**
  - **Validates: Requirements 1.4, 6.4**

- [ ]* 2.5 Write property test for validation errors
  - **Property 5: Validation error on missing fields**
  - **Validates: Requirements 1.5**

- [ ]* 2.6 Write property test for configuration creation with unique ID
  - **Property 12: Configuration creation with unique ID**
  - **Validates: Requirements 6.1**

- [ ]* 2.7 Write property test for configuration retrieval by ID
  - **Property 13: Configuration retrieval by ID**
  - **Validates: Requirements 6.5**

- [x] 3. Implement model configuration test endpoint
  - Create TestModelConfig endpoint (POST /api/model-configs/test)
  - Implement connection testing logic with timeout
  - Return success/failure status with diagnostic information
  - Handle various error scenarios (timeout, auth failure, invalid model)
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [ ]* 3.1 Write property test for connection test attempts
  - **Property 14: Connection test attempt**
  - **Validates: Requirements 7.1**

- [ ]* 3.2 Write property test for test timeout enforcement
  - **Property 15: Test timeout enforcement**
  - **Validates: Requirements 7.4**

- [x] 4. Update LLM client to use model configurations
  - Modify createLLMClient to accept ModelConfig parameter
  - Implement decryption of API key in memory
  - Update Chat handler to retrieve model config by ID
  - Update Chat handler to pass model config to LLM client
  - Update message saving to include model_config_id
  - _Requirements: 2.4, 3.4, 5.3_

- [ ]* 4.1 Write property test for message model tracking
  - **Property 8: Message model tracking**
  - **Validates: Requirements 3.4**

- [x] 5. Create default configuration migration
  - Implement migration logic to create default config from environment variables
  - Run migration on first startup if no configs exist
  - Set the created config as default
  - _Requirements: 2.3_

- [x] 6. Build frontend model configuration management panel
  - Create ModelConfigPanel component with list view
  - Implement create/edit form with all required fields
  - Add delete confirmation dialog
  - Implement test connection button
  - Add set as default toggle
  - Integrate API calls for all CRUD operations
  - Add error handling and user feedback
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 7.1, 7.2, 7.3_

- [ ]* 6.1 Write unit tests for ModelConfigPanel component
  - Test form validation
  - Test API call handling
  - Test error display

- [x] 7. Update frontend chat component for model selection
  - Update model selector to fetch from /api/model-configs
  - Display model config name in selector
  - Update chat API calls to include selected model config ID
  - Add visual indicator of currently active model
  - Display which model was used for each message in history
  - Implement model switching during conversation
  - _Requirements: 2.1, 2.2, 3.1, 3.2, 3.3, 4.1, 4.2, 4.3_

- [ ]* 7.1 Write property test for model selection persistence
  - **Property 6: Model selection persistence**
  - **Validates: Requirements 2.2, 2.4**

- [ ]* 7.2 Write property test for model switch with history preservation
  - **Property 7: Model switch with history preservation**
  - **Validates: Requirements 3.1, 3.2, 3.3**

- [ ]* 7.3 Write property test for active model display update
  - **Property 9: Active model display update**
  - **Validates: Requirements 4.2**

- [ ]* 7.4 Write property test for message history model indication
  - **Property 10: Message history model indication**
  - **Validates: Requirements 4.3**

- [x] 8. Add navigation and routing for model config panel
  - Add settings/config route to Next.js app
  - Create navigation link to model config panel
  - Ensure proper routing between chat and config views
  - _Requirements: 1.1_

- [x] 9. Implement error handling and user feedback
  - Add error boundaries for frontend components
  - Implement toast notifications for success/error messages
  - Add loading states for all async operations
  - Implement proper error messages for all backend error scenarios
  - _Requirements: 1.5, 7.3_

- [x] 10. Update CORS configuration
  - Add PUT and DELETE methods to allowed methods
  - Update allowed headers if needed for new endpoints
  - _Requirements: 1.3, 1.4_

- [ ] 11. Add security enhancements
  - Implement input sanitization for all user inputs
  - Add rate limiting to test endpoint
  - Ensure HTTPS enforcement in production configuration
  - Add audit logging for configuration changes
  - _Requirements: 5.1, 5.2, 5.4_

- [ ]* 11.1 Write security tests
  - Test SQL injection prevention
  - Test XSS prevention in configuration names
  - Test API key never appears in logs

- [ ] 12. Documentation and configuration
  - Update README with model configuration feature documentation
  - Add example configurations to .env.example
  - Document API endpoints in API documentation
  - Add user guide for model configuration panel
  - _Requirements: All_

- [x] 13. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
