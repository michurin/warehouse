Building:

```
git clone https://github.com/ollama/ollama.git
cd ollama
go generate ./...
go build .
```

Running server:

```
./ollama serve
```

Command line interface:

```
./ollama run gemma3:1b
./ollama run gemma3:270m
```

OpenAI API calls:

```
curl localhost:11434/api/generate -d '{"model": "gemma3:1b", "prompt":"Why is the sky blue?"}'
curl localhost:11434/v1/chat/completions -d '{"model": "gemma3:270m", "messages":[{"role": "user", "content": "Γιατί ο ουρανός είναι μπλε?"}]}'
curl -qs localhost:11434/v1/models | jq -r .data[].id
curl -qs localhost:11434/api/ps | jq
```

Where is models:

```
du -h ~/.ollama
```

References:

- <https://platform.openai.com/docs/api-reference/introduction>
- <https://ollama.com/library/gemma3>
- <https://github.com/ollama/ollama/blob/main/docs/api.md>
