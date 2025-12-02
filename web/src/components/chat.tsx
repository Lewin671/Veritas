"use client"

import { useState, useEffect, useRef } from "react"
import { Send, Bot, User, Plus, MessageSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

interface Message {
  role: "user" | "assistant"
  content: string
}

interface Conversation {
  id: string
  title: string
  createdAt: string
  messages?: Message[]
}

interface Model {
  id: string
  name: string
  description: string
}

export function Chat() {
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const [models, setModels] = useState<Model[]>([])
  const [selectedModel, setSelectedModel] = useState<string>("")
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Fetch models
    fetch("http://localhost:8080/api/models")
      .then((res) => res.json())
      .then((data) => {
        setModels(data)
        if (data.length > 0) {
          setSelectedModel(data[0].id)
        }
      })
      .catch((err) => console.error("Failed to fetch models:", err))

    // Fetch conversations
    fetchConversations()
  }, [])

  const fetchConversations = () => {
    fetch("http://localhost:8080/api/conversations")
      .then((res) => res.json())
      .then((data) => setConversations(data || []))
      .catch((err) => console.error("Failed to fetch conversations:", err))
  }

  const loadConversation = (id: string) => {
    setLoading(true)
    fetch(`http://localhost:8080/api/conversations/${id}`)
      .then((res) => res.json())
      .then((data) => {
        setCurrentConversationId(data.id)
        setMessages(data.messages || [])
      })
      .catch((err) => console.error("Failed to load conversation:", err))
      .finally(() => setLoading(false))
  }

  const startNewChat = () => {
    setCurrentConversationId(null)
    setMessages([])
  }

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || !selectedModel) return

    const userMessage: Message = { role: "user", content: input }
    setMessages((prev) => [...prev, userMessage])
    setInput("")
    setLoading(true)

    try {
      const res = await fetch("http://localhost:8080/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          modelId: selectedModel,
          message: userMessage.content,
          conversationId: currentConversationId,
        }),
      })
      const data = await res.json()
      
      // If it was a new chat, the backend created a conversation ID
      // We should refresh the list and set the current ID if we can get it from response
      // But currently the chat API doesn't return the conversation ID explicitly in the top level
      // Ideally we should update the backend to return it.
      // For now, let's just refresh the list.
      fetchConversations()

      const botMessage: Message = { role: "assistant", content: data.response }
      setMessages((prev) => [...prev, botMessage])
    } catch (err) {
      console.error("Failed to send message:", err)
      setMessages((prev) => [
        ...prev,
        { role: "assistant", content: "Error: Failed to get response from server." },
      ])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex h-screen bg-background text-foreground">
      {/* Sidebar */}
      <div className="w-64 border-r bg-muted/20 p-4 hidden md:flex flex-col">
        <Button onClick={startNewChat} variant="outline" className="w-full justify-start gap-2 mb-4">
          <Plus className="h-4 w-4" />
          New Chat
        </Button>
        <div className="flex-1 overflow-auto space-y-2">
          <div className="text-sm font-medium text-muted-foreground mb-2">History</div>
          {conversations.map((conv) => (
            <Button
              key={conv.id}
              variant={currentConversationId === conv.id ? "secondary" : "ghost"}
              className="w-full justify-start gap-2 text-sm font-normal truncate"
              onClick={() => loadConversation(conv.id)}
            >
              <MessageSquare className="h-4 w-4 shrink-0" />
              <span className="truncate">{conv.title}</span>
            </Button>
          ))}
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <div className="h-14 border-b flex items-center px-4 justify-between">
          <div className="font-semibold">Veritas</div>
          <select
            className="bg-transparent border rounded-md px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            value={selectedModel}
            onChange={(e) => setSelectedModel(e.target.value)}
          >
            {models.map((model) => (
              <option key={model.id} value={model.id}>
                {model.name}
              </option>
            ))}
          </select>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-auto p-4 space-y-4">
          {messages.length === 0 && (
            <div className="h-full flex flex-col items-center justify-center text-muted-foreground">
              <Bot className="h-12 w-12 mb-4" />
              <p className="text-lg font-medium">How can I help you today?</p>
            </div>
          )}
          {messages.map((msg, i) => (
            <div
              key={i}
              className={cn(
                "flex gap-3 max-w-3xl mx-auto",
                msg.role === "user" ? "justify-end" : "justify-start"
              )}
            >
              {msg.role === "assistant" && (
                <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                  <Bot className="h-5 w-5 text-primary" />
                </div>
              )}
              <div
                className={cn(
                  "rounded-lg px-4 py-2 max-w-[80%] whitespace-pre-wrap",
                  msg.role === "user"
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted"
                )}
              >
                {msg.content}
              </div>
              {msg.role === "user" && (
                <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center shrink-0">
                  <User className="h-5 w-5" />
                </div>
              )}
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        {/* Input Area */}
        <div className="p-4 border-t">
          <form onSubmit={handleSubmit} className="max-w-3xl mx-auto flex gap-2">
            <Input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Message Veritas..."
              disabled={loading}
              className="flex-1"
            />
            <Button type="submit" disabled={loading || !input.trim()}>
              <Send className="h-4 w-4" />
            </Button>
          </form>
        </div>
      </div>
    </div>
  )
}
