## Veritas – Real-Time News & Search Agent

This project implements **Veritas**, a simple information retrieval and news analysis agent.

- **Goal**: Answer user questions with accurate, real-time information, with a strong focus on **fresh news** and **reliable sources**.
- **Core idea**: Use a ReAct-style loop (Reasoning + Acting) where the agent:
  - thinks about the user’s question,
  - calls a `search_tool` to look up current web pages and news,
  - inspects the results (dates, sources, consistency),
  - and then synthesizes a final answer.

### What this agent does

- **Real-time news retrieval**: Looks up the latest articles and discards outdated results when the user asks for “latest” or “current” info.
- **Knowledge synthesis**: Combines multiple trustworthy sources into a clear, sourced answer.
- **Hallucination prevention**: If the search returns nothing reliable, the agent explicitly says it couldn’t find good information instead of guessing.
- **Neutral tone**: Presents facts objectively and can mention multiple viewpoints on controversial topics.


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
