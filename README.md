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
