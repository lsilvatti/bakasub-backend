# Bakasub - Tradutor de Legendas por IA 🌸✨

Tradução de Legendas - Alimentado por IA.

> *"N-não vá pensando que eu fiz uma ferramenta pra traduzir suas legendas pra você, baka!"*

Presta atenção! Esse é o **repositório principal**. Você não precisa olhar o código complicado do backend ou do frontend se só quiser usar o aplicativo. O Bakasub é uma ferramenta self-hosted que extrai, traduz usando os melhores modelos de IA (via OpenRouter) e remuxa legendas direto nos seus arquivos de vídeo (MKV, MP4) sem suar a camisa.

Eu orquestro tudo pra você usando Docker, fazendo o build do backend e do frontend direto das branchs principais deles. Então é só seguir as instruções e não fazer besteira!

## 🛠️ Do Que Eu Sou Feita (Arquitetura)
* **Backend:** Go (Golang) + PostgreSQL (Eu evoluí do SQLite, tente acompanhar!)
* **Frontend:** React + TypeScript + Vite (servido via Nginx)
* **Processamento de Vídeo:** FFmpeg + MKVToolNix
* **Orquestração:** Docker Compose

## 🚀 Como Sair Comigo (Instalação)

1. **Clone este repositório:**
   ```bash
   git clone [https://github.com/lsilvatti/bakasub.git](https://github.com/lsilvatti/bakasub.git)
   cd bakasub
   ```

2. **Configure suas variáveis de ambiente:**
   ```bash
   cp .env.example .env
   ```
   Abra o arquivo `.env`. Você *precisa* preencher:
   * A sua `OPENROUTER_API_KEY`
   * A sua `TMDB_API_KEY` (Pra eu poder buscar as capas e metadados dos seus filmes. Não que eu me importe com o que você assiste!)
   * O caminho absoluto para a sua pasta local de vídeos (`VIDEO_DIR`). Se você não me der um diretório de vídeos, eu não vou ter o que traduzir!

3. **Inicie a aplicação:**
   ```bash
   docker compose up -d --build
   ```
   *Nota: Eu vou automaticamente subir um banco de dados PostgreSQL 16, baixar o código mais recente do backend/frontend e conectar tudo. Seja paciente enquanto eu faço o build!*

4. **Acesse o painel:**
   Abra o seu navegador e acesse: **http://localhost:3000**