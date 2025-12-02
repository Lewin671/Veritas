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
- **AI / LLM SDK**: [Eino](https://github.com/cloudwego/eino)
- **Config Management**: `.env`