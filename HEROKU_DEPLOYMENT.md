# Heroku Deployment Guide for CLIProxyAPI

This guide covers deploying CLIProxyAPI to Heroku, including all available endpoints, authentication (Google Antigravity OAuth), and usage examples.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Deployment Methods](#deployment-methods)
   - [Method 1: Docker Container (Recommended)](#method-1-docker-container-recommended)
   - [Method 2: Go Buildpack](#method-2-go-buildpack)
3. [Configuration](#configuration)
4. [Available API Endpoints](#available-api-endpoints)
5. [Authentication Guide](#authentication-guide)
   - [Google Antigravity OAuth](#google-antigravity-oauth)
   - [Other OAuth Providers](#other-oauth-providers)
6. [Usage Examples](#usage-examples)
7. [Persistent Storage](#persistent-storage)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- A [Heroku](https://heroku.com) account
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) installed
- Git installed
- A Google account with Antigravity access (for Antigravity OAuth)

---

## Deployment Methods

### Method 1: Docker Container (Recommended)

This method uses the included `Dockerfile` and `heroku.yml`.

```bash
# 1. Clone the repository
git clone https://github.com/Anggahrm/CLIProxyAPI.git
cd CLIProxyAPI

# 2. Login to Heroku
heroku login

# 3. Create a new Heroku app
heroku create your-app-name

# 4. Set the stack to container
heroku stack:set container -a your-app-name

# 5. Set required environment variables
heroku config:set DEPLOY=cloud -a your-app-name

# 6. Push to Heroku
git push heroku main
```

### Method 2: Go Buildpack

This method uses the Go buildpack with the included `Procfile`.

```bash
# 1. Clone the repository
git clone https://github.com/Anggahrm/CLIProxyAPI.git
cd CLIProxyAPI

# 2. Login to Heroku
heroku login

# 3. Create a new Heroku app
heroku create your-app-name

# 4. Set the Go buildpack
heroku buildpacks:set heroku/go -a your-app-name

# 5. Set required environment variables
heroku config:set DEPLOY=cloud -a your-app-name

# 6. Push to Heroku
git push heroku main
```

> **Note:** Heroku automatically assigns a port via the `PORT` environment variable. CLIProxyAPI reads this variable and binds to it automatically — no manual port configuration is needed.

---

## Configuration

### Environment Variables

Set these via `heroku config:set`:

| Variable | Required | Description |
|----------|----------|-------------|
| `DEPLOY` | Yes | Set to `cloud` for cloud deployment mode |
| `PORT` | Auto | Automatically set by Heroku (do not set manually) |
| `MANAGEMENT_PASSWORD` | No | Password for the management web UI |
| `PGSTORE_DSN` | No | PostgreSQL connection string for persistent token storage |
| `PGSTORE_SCHEMA` | No | PostgreSQL schema (default: `public`) |
| `GITSTORE_GIT_URL` | No | Git repository URL for config/token storage |
| `GITSTORE_GIT_USERNAME` | No | Git username for git-backed store |
| `GITSTORE_GIT_TOKEN` | No | Git token for git-backed store |
| `OBJECTSTORE_ENDPOINT` | No | S3-compatible object store endpoint |
| `OBJECTSTORE_BUCKET` | No | Object store bucket name |
| `OBJECTSTORE_ACCESS_KEY` | No | Object store access key |
| `OBJECTSTORE_SECRET_KEY` | No | Object store secret key |

### Config File Setup

For cloud deployment, upload your config via the Management API or use a persistent storage backend (PostgreSQL, Git, or Object Store).

**Minimal `config.yaml` example:**

```yaml
port: 8317  # Overridden by PORT env var on Heroku

api-keys:
  - "your-secret-api-key"

remote-management:
  allow-remote: true
  secret-key: "your-management-secret"

auth-dir: "/tmp/cli-proxy-api"
```

To upload config via the Management API:

```bash
curl -X PUT "https://your-app-name.herokuapp.com/v0/management/config.yaml" \
  -H "Authorization: Bearer your-management-secret" \
  -H "Content-Type: text/yaml" \
  --data-binary @config.yaml
```

---

## Available API Endpoints

### Core API Endpoints (require API key auth)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Server info and endpoint list |
| `GET` | `/v1/models` | List available models |
| `POST` | `/v1/chat/completions` | OpenAI-compatible chat completions |
| `POST` | `/v1/completions` | OpenAI-compatible completions |
| `POST` | `/v1/messages` | Claude-compatible messages |
| `POST` | `/v1/messages/count_tokens` | Claude token counting |
| `POST` | `/v1/responses` | OpenAI Responses API |
| `POST` | `/v1/responses/compact` | OpenAI Responses API (compact) |

### Gemini API Endpoints (require API key auth)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1beta/models` | List Gemini models |
| `POST` | `/v1beta/models/*action` | Gemini generate/stream actions |
| `GET` | `/v1beta/models/*action` | Gemini model info |

### OAuth Callback Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/google/callback` | Google/Gemini OAuth callback |
| `GET` | `/antigravity/callback` | Antigravity OAuth callback |
| `GET` | `/anthropic/callback` | Claude/Anthropic OAuth callback |
| `GET` | `/codex/callback` | OpenAI Codex OAuth callback |
| `GET` | `/iflow/callback` | iFlow OAuth callback |

### Management API Endpoints (require management secret)

All management endpoints are under `/v0/management/` and require the `Authorization: Bearer <secret-key>` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v0/management/config` | Get current config (JSON) |
| `GET` | `/v0/management/config.yaml` | Get current config (YAML) |
| `PUT` | `/v0/management/config.yaml` | Update config (YAML) |
| `GET` | `/v0/management/usage` | Get usage statistics |
| `GET` | `/v0/management/api-keys` | List API keys |
| `PUT` | `/v0/management/api-keys` | Update API keys |
| `GET` | `/v0/management/auth-files` | List auth token files |
| `POST` | `/v0/management/auth-files` | Upload auth token files |
| `DELETE` | `/v0/management/auth-files` | Delete auth token files |
| `GET` | `/v0/management/auth-files/models` | List models from auth files |
| `GET` | `/v0/management/antigravity-auth-url` | Get Antigravity OAuth URL |
| `GET` | `/v0/management/gemini-cli-auth-url` | Get Gemini CLI OAuth URL |
| `GET` | `/v0/management/codex-auth-url` | Get Codex OAuth URL |
| `GET` | `/v0/management/anthropic-auth-url` | Get Claude OAuth URL |
| `GET` | `/v0/management/qwen-auth-url` | Get Qwen OAuth URL |
| `GET` | `/v0/management/kimi-auth-url` | Get Kimi OAuth URL |
| `GET` | `/v0/management/iflow-auth-url` | Get iFlow OAuth URL |
| `POST` | `/v0/management/oauth-callback` | Post OAuth callback data |
| `GET` | `/v0/management/get-auth-status` | Get OAuth auth status |
| `GET` | `/v0/management/logs` | Get application logs |

### Management Web UI

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/management.html` | Web-based management control panel |

---

## Authentication Guide

### Google Antigravity OAuth

Google Antigravity uses OAuth2 to authenticate your Google account. This gives you access to Gemini models through the Antigravity channel.

#### Option A: Via Management API (Remote/Heroku)

This is the recommended method for Heroku deployments since you cannot run CLI commands directly on the server.

**Step 1: Get the OAuth URL**

```bash
curl "https://your-app-name.herokuapp.com/v0/management/antigravity-auth-url" \
  -H "Authorization: Bearer your-management-secret"
```

This returns a JSON response with an OAuth URL and a session state identifier:

```json
{
  "url": "https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=...&redirect_uri=...&response_type=code&scope=...&state=...",
  "state": "random-state-string"
}
```

**Step 2: Authenticate in your browser**

Open the URL from the response in your browser. Sign in with your Google account and grant the requested permissions:
- Cloud Platform access
- User info (email/profile)
- Code Assist logging
- Experiments and configs

**Step 3: Complete the callback**

After authenticating, Google redirects to the callback URL. If the redirect URL points to your Heroku app (which it should for remote deployments), the server handles it automatically.

If the callback goes to `localhost`, you need to manually post the callback data:

```bash
curl -X POST "https://your-app-name.herokuapp.com/v0/management/oauth-callback" \
  -H "Authorization: Bearer your-management-secret" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "antigravity",
    "state": "the-state-from-step-1",
    "code": "the-authorization-code-from-callback-url"
  }'
```

**Step 4: Verify authentication**

```bash
curl "https://your-app-name.herokuapp.com/v0/management/get-auth-status" \
  -H "Authorization: Bearer your-management-secret"
```

```bash
curl "https://your-app-name.herokuapp.com/v0/management/auth-files" \
  -H "Authorization: Bearer your-management-secret"
```

#### Option B: Via CLI (Local, then upload tokens)

If you run CLIProxyAPI locally first, you can authenticate and then upload the token files to your Heroku deployment.

**Step 1: Run locally and authenticate**

```bash
# Build and run the login command
go build -o CLIProxyAPI ./cmd/server/
./CLIProxyAPI -antigravity-login
```

This opens your browser to Google's OAuth consent page. Sign in and authorize the application.

**Step 2: Upload token files to Heroku**

After local authentication, token files are saved in your auth directory (default: `~/.cli-proxy-api/`). Upload them to your Heroku deployment:

```bash
# List local auth files
ls ~/.cli-proxy-api/

# Upload each auth file via Management API
curl -X POST "https://your-app-name.herokuapp.com/v0/management/auth-files" \
  -H "Authorization: Bearer your-management-secret" \
  -F "file=@$HOME/.cli-proxy-api/antigravity_tokens_yourname@gmail.com.json"
```

#### What Antigravity provides

Once authenticated, Antigravity gives access to these models (among others):

- `gemini-2.5-pro` / `gemini-2.5-flash`
- `gemini-3-pro-preview` / `gemini-3-flash-preview`
- `claude-sonnet-4-5` / `claude-opus-4-5` (via Antigravity routing)
- And other models available through Google's Antigravity/Code Assist platform

### Other OAuth Providers

The same pattern applies for other providers:

| Provider | CLI Flag | Management Endpoint | Callback Path |
|----------|----------|---------------------|---------------|
| **Gemini CLI** | `-login` | `/v0/management/gemini-cli-auth-url` | `/google/callback` |
| **Antigravity** | `-antigravity-login` | `/v0/management/antigravity-auth-url` | `/antigravity/callback` |
| **OpenAI Codex** | `-codex-login` | `/v0/management/codex-auth-url` | `/codex/callback` |
| **Claude** | `-claude-login` | `/v0/management/anthropic-auth-url` | `/anthropic/callback` |
| **Qwen** | `-qwen-login` | `/v0/management/qwen-auth-url` | N/A |
| **iFlow** | `-iflow-login` | `/v0/management/iflow-auth-url` | `/iflow/callback` |
| **Kimi** | `-kimi-login` | `/v0/management/kimi-auth-url` | N/A |

---

## Usage Examples

### List Available Models

```bash
curl "https://your-app-name.herokuapp.com/v1/models" \
  -H "Authorization: Bearer your-api-key"
```

### Chat Completions (OpenAI-compatible)

```bash
curl -X POST "https://your-app-name.herokuapp.com/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-2.5-pro",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "stream": false
  }'
```

### Streaming Chat Completions

```bash
curl -X POST "https://your-app-name.herokuapp.com/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-2.5-flash",
    "messages": [
      {"role": "user", "content": "Write a short poem about coding"}
    ],
    "stream": true
  }'
```

### Claude-compatible Messages

```bash
curl -X POST "https://your-app-name.herokuapp.com/v1/messages" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-sonnet-4-5",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Explain quantum computing in simple terms"}
    ]
  }'
```

### Gemini Native API

```bash
curl -X POST "https://your-app-name.herokuapp.com/v1beta/models/gemini-2.5-pro:generateContent" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [
      {"parts": [{"text": "What is the meaning of life?"}]}
    ]
  }'
```

### Using with OpenAI Python SDK

```python
from openai import OpenAI

client = OpenAI(
    api_key="your-api-key",
    base_url="https://your-app-name.herokuapp.com/v1"
)

response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)
print(response.choices[0].message.content)
```

### Using with Claude Code

Set the following environment variables to point Claude Code at your proxy:

```bash
export ANTHROPIC_BASE_URL="https://your-app-name.herokuapp.com"
export ANTHROPIC_API_KEY="your-api-key"
```

---

## Persistent Storage

Heroku has an ephemeral filesystem — files are lost on dyno restart. To persist OAuth tokens and configuration, use one of these backends:

### Option 1: Heroku Postgres (Recommended)

```bash
# Add Heroku Postgres addon
heroku addons:create heroku-postgresql:essential-0 -a your-app-name

# Get the DATABASE_URL
heroku config:get DATABASE_URL -a your-app-name

# Set the PGSTORE_DSN
heroku config:set PGSTORE_DSN="your-database-url" -a your-app-name
```

### Option 2: Git-backed Store

```bash
heroku config:set \
  GITSTORE_GIT_URL="https://github.com/your-org/cli-proxy-config.git" \
  GITSTORE_GIT_USERNAME="your-username" \
  GITSTORE_GIT_TOKEN="ghp_your_token" \
  -a your-app-name
```

### Option 3: S3/Object Store

```bash
heroku config:set \
  OBJECTSTORE_ENDPOINT="https://s3.amazonaws.com" \
  OBJECTSTORE_BUCKET="your-bucket" \
  OBJECTSTORE_ACCESS_KEY="your-access-key" \
  OBJECTSTORE_SECRET_KEY="your-secret-key" \
  -a your-app-name
```

---

## Troubleshooting

### App crashes on startup

- Make sure you have set `DEPLOY=cloud`
- Check logs: `heroku logs --tail -a your-app-name`
- Ensure a config file is available (via persistent storage or Management API upload)

### Port binding errors

- **Do not** set the `PORT` environment variable manually — Heroku assigns it automatically
- CLIProxyAPI reads `$PORT` at startup and overrides the config file port

### OAuth callbacks not working

- For Heroku deployments, OAuth callbacks redirect to `localhost` by default
- Use the Management API (`/v0/management/oauth-callback`) to manually post callback data
- Alternatively, authenticate locally first and upload token files

### Tokens lost after dyno restart

- Heroku's filesystem is ephemeral; configure a persistent storage backend (Postgres, Git, or Object Store)
- See the [Persistent Storage](#persistent-storage) section

### Management API returns 404

- Ensure `remote-management.secret-key` is set in your config
- Ensure `remote-management.allow-remote` is `true` for Heroku (non-localhost access)
