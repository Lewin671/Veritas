## Introduction

This project implements **Veritas**, a simple information retrieval and news analysis agent.


## Tech Stack

### Web

- **Framework**: Next.js 16 (App Router, Server Actions)
- **Language**: TypeScript
- **UI & Styling**: shadcn/ui, Tailwind CSS, Lucide React and Vercel AI SDK for prebuilt AI elements
- **State Management**: Zustand (Client-side), React Server Components (Server-side)
- **Validation**: Zod
- **Infrastructure**: Vercel
- **Lint/Format**: Biome
- **Package Management**: pnpm


### Server

- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (access + refresh token)
- **Cache / Token Store**: Redis
- **HTTP Framework**: Gin
- **Reverse Proxy / TLS Termination**: Nginx
- **Streaming**: Server-Sent Events (SSE) or HTTP Chunked Responses
- **Observability**: Prometheus for metrics + OpenTelemetry for tracing
- **Containerization**: Docker + Docker Compose
- **Hosting**: Kubernetes (optional for scaling) or VM + containers
- **AI / LLM SDK**: [OpenAI Go SDK](https://github.com/openai/openai-go)
- **Config Management**: `.env`
- **Lint/Format**: golangci-lint


## Setup commands

### Style & Lint check

#### Web
Always run the following command to check for style and lint issues either before committing or after modifying code in the web directory:

```bash
cd web && npm run check
```

#### Server
Always run the following command to check for style and lint issues in the server directory:

```bash
cd server && golangci-lint run
```
