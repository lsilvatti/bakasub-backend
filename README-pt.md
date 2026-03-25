# Bakasub - API & Core Engine 🌸✨

Hmph! Então você achou o repositório backend do **Bakasub**? Não vá pensando besteira! Não é como se eu tivesse escrito essa engine super otimizada em Go só pra você! É que *alguém* tinha que fazer o trabalho pesado: processamento de vídeo, extração de legendas, comunicação com IA e muxing sem travar tudo. B-baka! 

> **Presta atenção!** Esse repositório é só o cérebro principal. Se você quer apenas instalar e usar o Bakasub como uma pessoa normal, vai olhar o [repositório principal de orquestração](https://github.com/lsilvatti/bakasub).

## 🛠️ Do Que Eu Sou Feita (Tech Stack)
Não ache que pode me rodar numa batata. Você precisa das ferramentas certas!
* **Linguagem:** Go (Golang) 1.22+ (Porque velocidade importa, obviamente!)
* **Router:** `go-chi/chi` (APIs versionadas em V1, mantenha isso organizado!)
* **Banco de Dados:** **PostgreSQL**. Eu evoluí do SQLite porque agora eu lido com traduções em lote concorrentes e uma memória de tradução pesada. Tente acompanhar!
* **Processamento de Mídia:** FFmpeg & MKVToolNix (Eu preciso disso para extrair e costurar suas legendas de volta, idiota!)
* **Provedor de IA:** OpenRouter API (Claude 3.5 Sonnet, GPT-4o, Gemini 1.5 Pro).
* **Logging:** Logs estruturados (`slog`) com um formatador customizado pro terminal.

## 🏗️ Como Eu Penso (Arquitetura)
Eu construí isso usando **Clean Architecture** e princípios **SOLID**. Não porque eu me importo com a sua experiência de leitura, mas porque código bagunçado é absolutamente nojento! Tudo é desacoplado.

* `cmd/server`: Onde eu acordo. Não toque nisso a menos que saiba o que está fazendo.
* `internal/routes`: Meus portões `/api/v1/`.
* `internal/handlers`: Onde eu pego suas requisições HTTP. Eles validam os dados e logam os seus erros (Bad Requests).
* `internal/services`: Meu verdadeiro núcleo! Pipelines de tradução, escaneamento de vídeo, extração de MKV e gerenciamento de pastas.
* `internal/parser`: Meus parsers customizados de legenda (`.ass`, `.srt`, `.vtt`). Ele até limpa tags SDH e injeta cabeçalhos `[BakaSub-AI]` nos seus arquivos `.ass` com perfeição.
* `internal/ai`: Integração com os LLMs do OpenRouter com cálculo dinâmico de preços.
* `internal/utils`: Minhas utilidades, incluindo o **SSE Broker**. Eu envio atualizações de progresso em tempo real para o seu frontend para você não ficar se perguntando se eu travei. 

### ✨ As Funcionalidades "Geniais" Que Você Provavelmente Nem Notou:
1. **Inferência Inteligente de Idiomas:** Eu detecto automaticamente os idiomas de origem a partir de sufixos bagunçados (como `_pt-BR` ou `-spa`) usando um sistema de mapeamento no PostgreSQL antes de mandar pra IA.
2. **Presets Contextuais:** Eu não traduzo às cegas. Eu uso prompts de sistema altamente ajustados para Anime, Filmes, Documentários e Comédia, lidando automaticamente com termos de gênero neutro e ajustando a criatividade (`temperature`) na hora.
3. **Memória de Tradução (Cache):** Eu faço um hash seguro de cada linha de diálogo que você traduz. Se você traduzir de novo, eu carrego do PostgreSQL em vez de torrar seus créditos da API do OpenRouter. 
4. **Regex de Extração Inteligente:** Quando eu extraio faixas de um MKV, eu limpo automaticamente as tags de idioma zoadas.

## 💻 Como Sair Comigo (Desenvolvimento Local)
Se você realmente quer mexer no meu código, é bom configurar as coisas direito!

### 1. Pré-requisitos
Você precisa de Go, PostgreSQL, FFmpeg e MKVToolNix. Se você está no Arch Linux (o que você obviamente deveria estar), é só rodar:
`sudo pacman -S go postgresql ffmpeg mkvtoolnix-cli`

### 2. Configuração
Faça um clone meu e instale minhas dependências:
`go mod download`

Depois, copie meu arquivo de ambiente. Coloque sua chave da API do OpenRouter e as credenciais do PostgreSQL lá. Se você deixar em branco, eu vou jogar erros 500 na sua cara!
`cp .env.example .env`

*Nota: Você não precisa rodar as migrations manualmente. Eu executo os arquivos SQL puros de `internal/db/migrations/` automaticamente quando eu inicio, porque eu sou inteligente assim.*

### 3. Me Rode!
Use o [Air](https://github.com/cosmtrek/air) para live-reloading. É o único jeito de desenvolver direito.
`air`

*Eu estarei esperando por suas requisições em `http://localhost:8080/api/v1/`... Não me faça esperar muito!*

## 🐳 Docker
Ugh, tá bom. Se você é preguiçoso demais para instalar o PostgreSQL e o FFmpeg nativamente, eu fiz um Dockerfile pra você. Ele empacota tudo o que você precisa num container só.

`docker build -t bakasub-backend .`
`docker run -p 8080:8080 --env-file .env bakasub-backend`