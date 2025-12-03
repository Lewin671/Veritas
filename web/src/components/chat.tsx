'use client';

import { Bot, MessageSquare, Plus, Send, Settings, User } from 'lucide-react';
import Link from 'next/link';
import { useCallback, useEffect, useRef, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

interface Message {
  role: 'user' | 'assistant';
  content: string;
  modelConfigId?: string;
}

interface Conversation {
  id: string;
  title: string;
  createdAt: string;
  messages?: Message[];
}

interface ModelConfig {
  id: string;
  name: string;
  provider: string;
  baseUrl: string;
  modelId: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

interface ChatApiResponse {
  response: string;
  conversationId?: string;
}

export function Chat() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [currentConversationId, setCurrentConversationId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [modelConfigs, setModelConfigs] = useState<ModelConfig[]>([]);
  const [selectedModelConfigId, setSelectedModelConfigId] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const fetchConversations = useCallback(() => {
    fetch('http://localhost:8080/api/conversations')
      .then((res) => res.json())
      .then((data) => setConversations(data || []))
      .catch((err) => console.error('Failed to fetch conversations:', err));
  }, []);

  useEffect(() => {
    // Fetch model configs
    fetch('http://localhost:8080/api/model-configs')
      .then((res) => res.json())
      .then((data) => {
        setModelConfigs(data);
        if (data.length > 0) {
          // Select default model or first model
          const defaultModel = data.find((m: ModelConfig) => m.isDefault);
          setSelectedModelConfigId(defaultModel?.id || data[0].id);
        }
      })
      .catch((err) => console.error('Failed to fetch model configs:', err));

    // Fetch conversations
    fetchConversations();
  }, [fetchConversations]);

  const loadConversation = (id: string) => {
    setLoading(true);
    fetch(`http://localhost:8080/api/conversations/${id}`)
      .then((res) => res.json())
      .then((data) => {
        setCurrentConversationId(data.id);
        setMessages(data.messages || []);
      })
      .catch((err) => console.error('Failed to load conversation:', err))
      .finally(() => setLoading(false));
  };

  const startNewChat = () => {
    setCurrentConversationId(null);
    setMessages([]);
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || !selectedModelConfigId) return;

    const userMessage: Message = { role: 'user', content: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      const res = await fetch('http://localhost:8080/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          modelConfigId: selectedModelConfigId,
          message: userMessage.content,
          conversationId: currentConversationId,
        }),
      });
      const data: ChatApiResponse = await res.json();

      // Update current conversation ID based on backend response
      if (data.conversationId) {
        setCurrentConversationId(data.conversationId);
      }

      fetchConversations();

      const botMessage: Message = {
        role: 'assistant',
        content: data.response,
        modelConfigId: selectedModelConfigId,
      };
      setMessages((prev) => [...prev, botMessage]);
    } catch (err) {
      console.error('Failed to send message:', err);
      setMessages((prev) => [
        ...prev,
        {
          role: 'assistant',
          content: 'Error: Failed to get response from server.',
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const getModelName = (modelConfigId?: string) => {
    if (!modelConfigId) return null;
    const config = modelConfigs.find((m) => m.id === modelConfigId);
    return config?.name;
  };

  return (
    <div className="flex h-screen bg-background text-foreground">
      {/* Sidebar */}
      <div className="hidden w-64 flex-col border-r bg-muted/20 p-4 md:flex">
        <Button
          onClick={startNewChat}
          variant="outline"
          className="mb-4 w-full justify-start gap-2"
        >
          <Plus className="h-4 w-4" />
          New Chat
        </Button>
        <div className="flex-1 space-y-2 overflow-auto">
          <div className="mb-2 font-medium text-muted-foreground text-sm">History</div>
          {conversations.map((conv) => (
            <Button
              key={conv.id}
              variant={currentConversationId === conv.id ? 'secondary' : 'ghost'}
              className="w-full justify-start gap-2 truncate font-normal text-sm"
              onClick={() => loadConversation(conv.id)}
            >
              <MessageSquare className="h-4 w-4 shrink-0" />
              <span className="truncate">{conv.title}</span>
            </Button>
          ))}
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex flex-1 flex-col">
        {/* Header */}
        <div className="flex h-14 items-center justify-between border-b px-4">
          <div className="font-semibold">Veritas</div>
          <div className="flex items-center gap-2">
            <select
              className="rounded-md border bg-transparent px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              value={selectedModelConfigId}
              onChange={(e) => setSelectedModelConfigId(e.target.value)}
            >
              {modelConfigs.map((config) => (
                <option key={config.id} value={config.id}>
                  {config.name} {config.isDefault && '(Default)'}
                </option>
              ))}
            </select>
            <Link href="/settings">
              <Button variant="ghost" size="sm">
                <Settings className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>

        {/* Messages */}
        <div className="flex-1 space-y-4 overflow-auto p-4">
          {messages.length === 0 && (
            <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
              <Bot className="mb-4 h-12 w-12" />
              <p className="font-medium text-lg">How can I help you today?</p>
            </div>
          )}
          {messages.map((msg, i) => (
            <div
              key={`${msg.role}-${i}-${msg.content.slice(0, 20)}`}
              className={cn(
                'mx-auto flex max-w-3xl gap-3',
                msg.role === 'user' ? 'justify-end' : 'justify-start'
              )}
            >
              {msg.role === 'assistant' && (
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                  <Bot className="h-5 w-5 text-primary" />
                </div>
              )}
              <div className="flex flex-col gap-1">
                <div
                  className={cn(
                    'max-w-[80%] whitespace-pre-wrap rounded-lg px-4 py-2',
                    msg.role === 'user' ? 'bg-primary text-primary-foreground' : 'bg-muted'
                  )}
                >
                  {msg.content}
                </div>
                {msg.role === 'assistant' && msg.modelConfigId && (
                  <div className="text-muted-foreground text-xs">
                    {getModelName(msg.modelConfigId)}
                  </div>
                )}
              </div>
              {msg.role === 'user' && (
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-muted">
                  <User className="h-5 w-5" />
                </div>
              )}
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        {/* Input Area */}
        <div className="border-t p-4">
          <form onSubmit={handleSubmit} className="mx-auto flex max-w-3xl gap-2">
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
  );
}
