## Veritas – ChatGPT-like ReAct Agent

**Veritas** is a conversational AI interface powered by a ReAct (Reasoning + Acting) agent. It combines real-time information retrieval, web search, and research capabilities to provide accurate, sourced answers through a familiar chat interface.

### Key capabilities

- **Real-time information retrieval**: Fetches the latest articles and discards outdated results when users ask for current information
- **Knowledge synthesis**: Combines multiple trustworthy sources into clear, well-sourced answers
- **Hallucination prevention**: Explicitly states when reliable information cannot be found instead of guessing
- **Objective analysis**: Presents facts neutrally and can discuss multiple viewpoints on controversial topics

## Model Configuration Management

Veritas now supports configuring multiple LLM models through a user-friendly interface:

### Features

- **Multiple Model Configurations**: Add and manage multiple LLM providers (OpenAI, Anthropic, custom endpoints)
- **Secure Credential Storage**: API keys are encrypted using AES-256-GCM before storage
- **Model Switching**: Switch between different models during conversations
- **Connection Testing**: Test model configurations before saving
- **Default Model**: Set a default model for new conversations

### Setup

1. Generate an encryption key:
```bash
openssl rand -base64 32
```

2. Add the key to your `.env` file:
```
ENCRYPTION_KEY=your-generated-key-here
```

3. Start the server and navigate to Settings (gear icon) to configure your models

### API Endpoints

- `POST /api/model-configs` - Create a new model configuration
- `GET /api/model-configs` - List all configurations
- `GET /api/model-configs/:id` - Get a specific configuration
- `PUT /api/model-configs/:id` - Update a configuration
- `DELETE /api/model-configs/:id` - Delete a configuration
- `POST /api/model-configs/test` - Test a configuration

### Migration

On first startup, if `OPENAI_API_KEY` is set in environment variables, a default configuration will be automatically created.
