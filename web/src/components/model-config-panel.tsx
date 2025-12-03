'use client';

import { ArrowLeft, Loader2, Plus, Settings, Trash2, X } from 'lucide-react';
import Link from 'next/link';
import { useCallback, useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

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

interface ModelConfigFormData {
  name: string;
  provider: string;
  baseUrl: string;
  modelId: string;
  apiKey: string;
  isDefault: boolean;
}

export function ModelConfigPanel() {
  const [configs, setConfigs] = useState<ModelConfig[]>([]);
  const [loading, setLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);
  const [formData, setFormData] = useState<ModelConfigFormData>({
    name: '',
    provider: 'openai',
    baseUrl: 'https://api.openai.com/v1',
    modelId: '',
    apiKey: '',
    isDefault: false,
  });
  const [error, setError] = useState<string | null>(null);

  const fetchConfigs = useCallback(async () => {
    try {
      const res = await fetch('http://localhost:8080/api/model-configs');
      const data = await res.json();
      setConfigs(data || []);
    } catch (err) {
      console.error('Failed to fetch configs:', err);
      setError('Failed to load configurations');
    }
  }, []);

  useEffect(() => {
    fetchConfigs();
  }, [fetchConfigs]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const url = editingId
        ? `http://localhost:8080/api/model-configs/${editingId}`
        : 'http://localhost:8080/api/model-configs';
      const method = editingId ? 'PUT' : 'POST';

      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || 'Failed to save configuration');
      }

      await fetchConfigs();
      resetForm();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save configuration');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this configuration?')) {
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/model-configs/${id}`, {
        method: 'DELETE',
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || 'Failed to delete configuration');
      }

      await fetchConfigs();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete configuration');
    }
  };

  const handleEdit = (config: ModelConfig) => {
    setEditingId(config.id);
    setFormData({
      name: config.name,
      provider: config.provider,
      baseUrl: config.baseUrl,
      modelId: config.modelId,
      apiKey: '', // Don't populate API key for security
      isDefault: config.isDefault,
    });
    setShowForm(true);
  };

  const handleTest = async () => {
    setTesting(true);
    setTestResult(null);
    setError(null);

    try {
      const res = await fetch('http://localhost:8080/api/model-configs/test', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          baseUrl: formData.baseUrl,
          modelId: formData.modelId,
          apiKey: formData.apiKey,
        }),
      });

      const data = await res.json();
      setTestResult(data);
    } catch (_err) {
      setTestResult({
        success: false,
        message: 'Failed to test connection',
      });
    } finally {
      setTesting(false);
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      provider: 'openai',
      baseUrl: 'https://api.openai.com/v1',
      modelId: '',
      apiKey: '',
      isDefault: false,
    });
    setEditingId(null);
    setShowForm(false);
    setTestResult(null);
    setError(null);
  };

  return (
    <div className="flex h-screen flex-col bg-background">
      {/* Header */}
      <div className="flex h-14 items-center justify-between border-b px-4">
        <div className="flex items-center gap-2">
          <Link href="/">
            <Button variant="ghost" size="sm">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <Settings className="h-5 w-5" />
          <span className="font-semibold">Model Configurations</span>
        </div>
        <Button onClick={() => setShowForm(!showForm)} variant="outline" size="sm">
          {showForm ? (
            <>
              <X className="mr-2 h-4 w-4" />
              Cancel
            </>
          ) : (
            <>
              <Plus className="mr-2 h-4 w-4" />
              Add Model
            </>
          )}
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-4">
        {error && (
          <div className="mb-4 rounded-lg border border-red-200 bg-red-50 p-3 text-red-800 text-sm">
            {error}
          </div>
        )}

        {/* Form */}
        {showForm && (
          <div className="mb-6 rounded-lg border bg-card p-4">
            <h3 className="mb-4 font-medium text-lg">
              {editingId ? 'Edit Configuration' : 'New Configuration'}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="name" className="mb-1 block font-medium text-sm">
                  Name *
                </label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="My GPT-4 Config"
                  required
                />
              </div>

              <div>
                <label htmlFor="provider" className="mb-1 block font-medium text-sm">
                  Provider *
                </label>
                <select
                  id="provider"
                  className="w-full rounded-md border bg-background px-3 py-2 text-sm"
                  value={formData.provider}
                  onChange={(e) => setFormData({ ...formData, provider: e.target.value })}
                  required
                >
                  <option value="openai">OpenAI</option>
                  <option value="anthropic">Anthropic</option>
                  <option value="custom">Custom</option>
                </select>
              </div>

              <div>
                <label htmlFor="baseUrl" className="mb-1 block font-medium text-sm">
                  Base URL
                </label>
                <Input
                  id="baseUrl"
                  value={formData.baseUrl}
                  onChange={(e) => setFormData({ ...formData, baseUrl: e.target.value })}
                  placeholder="https://api.openai.com/v1"
                />
              </div>

              <div>
                <label htmlFor="modelId" className="mb-1 block font-medium text-sm">
                  Model ID *
                </label>
                <Input
                  id="modelId"
                  value={formData.modelId}
                  onChange={(e) => setFormData({ ...formData, modelId: e.target.value })}
                  placeholder="gpt-4o"
                  required
                />
              </div>

              <div>
                <label htmlFor="apiKey" className="mb-1 block font-medium text-sm">
                  API Key *
                </label>
                <Input
                  id="apiKey"
                  type="password"
                  value={formData.apiKey}
                  onChange={(e) => setFormData({ ...formData, apiKey: e.target.value })}
                  placeholder="sk-..."
                  required
                />
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="isDefault"
                  checked={formData.isDefault}
                  onChange={(e) => setFormData({ ...formData, isDefault: e.target.checked })}
                  className="h-4 w-4"
                />
                <label htmlFor="isDefault" className="text-sm">
                  Set as default model
                </label>
              </div>

              {testResult && (
                <div
                  className={cn(
                    'rounded-lg border p-3 text-sm',
                    testResult.success
                      ? 'border-green-200 bg-green-50 text-green-800'
                      : 'border-red-200 bg-red-50 text-red-800'
                  )}
                >
                  {testResult.message}
                </div>
              )}

              <div className="flex gap-2">
                <Button
                  type="button"
                  onClick={handleTest}
                  disabled={testing || !formData.apiKey}
                  variant="outline"
                >
                  {testing ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Testing...
                    </>
                  ) : (
                    'Test Connection'
                  )}
                </Button>
                <Button type="submit" disabled={loading}>
                  {loading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    'Save'
                  )}
                </Button>
                <Button type="button" onClick={resetForm} variant="ghost">
                  Cancel
                </Button>
              </div>
            </form>
          </div>
        )}

        {/* List */}
        <div className="space-y-3">
          {configs.length === 0 && !showForm && (
            <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
              <Settings className="mx-auto mb-2 h-12 w-12 opacity-50" />
              <p className="mb-2 font-medium">No model configurations yet</p>
              <p className="text-sm">Add your first model configuration to get started</p>
            </div>
          )}

          {configs.map((config) => (
            <div key={config.id} className="rounded-lg border bg-card p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <h4 className="font-medium">{config.name}</h4>
                    {config.isDefault && (
                      <span className="rounded-full bg-primary/10 px-2 py-0.5 text-primary text-xs">
                        Default
                      </span>
                    )}
                  </div>
                  <div className="mt-1 space-y-1 text-muted-foreground text-sm">
                    <div>Provider: {config.provider}</div>
                    <div>Model: {config.modelId}</div>
                    {config.baseUrl && <div>Base URL: {config.baseUrl}</div>}
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button onClick={() => handleEdit(config)} variant="ghost" size="sm">
                    Edit
                  </Button>
                  <Button
                    onClick={() => handleDelete(config.id)}
                    variant="ghost"
                    size="sm"
                    className="text-red-600 hover:text-red-700"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
