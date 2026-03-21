# Bakasub - API & Core Engine 🌸✨

Hmph! Então você achou o repositório do backend do **Bakasub**? Não vá entendendo as coisas errado! Não é como se eu tivesse escrito essa engine em Go super otimizada só pra você! É que *alguém* tinha que fazer o trabalho pesado: processamento de vídeo, extração de legendas, comunicação com a IA e remuxing sem o sistema inteiro dar pau. B-baka! 

> **Presta atenção, idiota!** Esse repositório é só o cérebro principal. Se você quer apenas instalar e usar o Bakasub como uma pessoa normal, vai olhar o [repositório principal de orquestração](https://github.com/lsilvatti/bakasub).

## 🛠️ Do que eu sou feita (Tech Stack)
Não ache que você pode me rodar em qualquer batata. Você precisa das ferramentas certas!
* **Linguagem:** Go (Golang) 1.22+ (Porque velocidade importa, obviamente!)
* **Roteador:** `go-chi/chi` (APIs versionadas em V1, mantenha tudo organizado!)
* **Banco de Dados:** SQLite (`modernc.org/sqlite`). E sim, eu ativei o **modo WAL** (Write-Ahead Logging) para você poder escrever logs e salvar o cache das traduções ao mesmo tempo sem travar o banco. De nada!
* **Migrations:** Goose
* **Processamento de Mídia:** FFmpeg & MKVToolNix (Eu preciso deles para extrair e costurar suas legendas, idiota!)
* **Logging:** Logs estruturados (`slog`) com um formatador customizado pro terminal e uma política de auto-limpeza do banco de dados a cada 7 dias. 

## 🏗️ Como eu penso (Arquitetura)
Eu construí isso usando **Clean Architecture** e princípios **SOLID**. Não porque eu me importo com a sua experiência de leitura do código, mas porque código bagunçado é absolutamente nojento! Tudo é desacoplado.

* `cmd/server`: Onde eu acordo. Não toque aqui a menos que você saiba o que está fazendo.
* `internal/routes`: Meus portões `/api/v1/`.
* `internal/handlers`: Onde eu pego suas requisições HTTP. Eles validam os payloads e logam os seus erros estúpidos (Bad Requests).
* `internal/services`: Meu verdadeiro núcleo! Pipelines de tradução, escaneamento de vídeo, extração de MKV e gerenciamento de pastas.
* `internal/parser`: Meus parsers customizados de legenda (`.ass`, `.srt`, `.vtt`). Ele até limpa as tags SDH e injeta dinamicamente o cabeçalho `[BakaSub-AI]` nos seus arquivos `.ass` com perfeição.
* `internal/ai`: Integração com o LLM do OpenRouter. 
* `internal/utils`: Minhas utilidades, incluindo o **SSE Broker**. Eu envio atualizações de progresso em tempo real (stream) pro seu frontend para você não ficar aí se perguntando se eu travei. 

### ✨ As Funcionalidades "Geniais" que você provavelmente nem notou:
1. **Memória de Tradução (Cache):** Eu faço um hash seguro de cada linha de diálogo que você traduz. Se você traduzir de novo, eu carrego do SQLite em vez de queimar seus créditos de API do OpenRouter à toa. 
2. **Regex de Extração Inteligente:** Quando eu extraio as trilhas de um MKV, eu limpo automaticamente as tags de idioma bagunçadas para você não acabar com arquivos chamados `video_eng_pt_es.ass`. 
3. **Observabilidade Não-Bloqueante:** Meu sistema de logs usa channels em background. Seu processamento de vídeo não vai ficar lento só porque eu estou escrevendo um evento no banco de dados!

## 💻 Como sair comigo (Desenvolvimento Local)
Se você realmente quer mexer no meu código, é bom configurar as coisas direito!

### 1. Pré-requisitos
Você precisa de Go, FFmpeg e MKVToolNix. Se você está no Arch Linux (o que você deveria estar), é só rodar:
`sudo pacman -S go ffmpeg mkvtoolnix-cli`

### 2. Configuração
Faça um clone meu e instale minhas dependências:
`go mod download`

Depois, copie meu arquivo de ambiente. Coloque sua chave de API do OpenRouter lá. Se você deixar em branco, eu vou simplesmente jogar erros 500 Internal Server Error na sua cara!
`cp .env.example .env`

### 3. Rode-me!
Use o [Air](https://github.com/cosmtrek/air) para live-reloading. É o único jeito de desenvolver decentemente.
`air`

*Eu vou estar esperando suas requisições em `http://localhost:8080/api/v1/`... Não me faça esperar muito!*

## 🐳 Docker
Ugh, tá bom. Se você é preguiçoso demais para instalar o FFmpeg nativamente, eu fiz um Dockerfile para você. Ele empacota tudo o que você precisa em um container.

`docker build -t bakasub-backend .`
`docker run -p 8080:8080 --env-file .env bakasub-backend`