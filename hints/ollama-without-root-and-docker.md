# install

```
go install github.com/ollama/ollama@latest
```

# or build

```
git clone https://github.com/ollama/ollama.git
cd ollama
go generate ./...
go build .
```

It seems master is little bit broken. Do not follow instructions
in logs. Do something like this:

```
cmake -S llama/server --preset cpu
cmake --build build/llama-server-cpu
```

```
sudo pacman -S vulkan-headers
sudo pacman -S spirv-headers
cmake -S llama/server --preset vulkan
cmake --build build/llama-server-vulkan
```

You may need to add swap:

```
free -h
sudo fallocate -l 8G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
free -h
```

To avoid OOM, you may want to use `-j1`

# run server

```
./ollama serve
```

# command line interface

```
./ollama run gemma3:1b
./ollama run gemma3:270m
```

# OpenAI API calls

```
curl localhost:11434/api/generate -d '{"model": "gemma3:1b", "prompt":"Why is the sky blue?"}'
curl localhost:11434/v1/chat/completions -d '{"model": "gemma3:270m", "messages":[{"role": "user", "content": "Γιατί ο ουρανός είναι μπλε?"}]}'
curl -qs localhost:11434/v1/models | jq -r .data[].id
curl -qs localhost:11434/api/ps | jq
```

# where is models

```
du -h ~/.ollama
```

# references

- <https://ollama.com/search> (<https://ollama.com/library/gemma3>)
- <https://platform.openai.com/docs/api-reference/introduction>
- <https://github.com/ollama/ollama/blob/main/docs/api.md>
