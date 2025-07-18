# 🏦 Sistema de Leilões em Go

Um sistema de leilões robusto desenvolvido em Go, utilizando uma arquitetura limpa para garantir escalabilidade e manutenibilidade. A API permite a criação de leilões, o registro de lances e a consulta de resultados, com uma funcionalidade de fechamento automático de leilões baseada em tempo.

---

## ✨ Funcionalidades Principais

- **Criação de Leilões**: Permite registrar novos produtos para leilão.
- **Registro de Lances**: Usuários podem fazer lances em leilões ativos.
- **Consulta de Leilões**: Busca de leilões por status, categoria ou nome do produto.
- **Consulta de Lances**: Visualização de todos os lances feitos para um determinado leilão.
- **Fechamento Automático**: Uma rotina em background (goroutine) monitora e fecha automaticamente os leilões cujo tempo expirou, com duração configurável via variáveis de ambiente.
- **Processamento de Lances em Lote**: Otimização de performance através do processamento de lances em lotes.

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: Go
- **Framework Web**: Gin Gonic
- **Banco de Dados**: MongoDB
- **Containerização**: Docker & Docker Compose
- **Testes**: Pacote `testing` nativo do Go & Testify

---

## 🏗️ Arquitetura

O projeto segue princípios de **Arquitetura Limpa (Clean Architecture)**, separando as responsabilidades em camadas distintas:

- **`internal/entity`**: Camada de entidades e regras de negócio agnósticas.
- **`internal/usecase`**: Implementa os casos de uso e orquestra as operações.
- **`internal/infra`**: Detalhes de implementação (servidor web, banco de dados, etc).

---

## 🚀 Começando

### ✨ Pré-requisitos

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### ⚙️ Configuração

```bash
git clone <URL_DO_SEU_REPOSITORIO>
cd <NOME_DO_REPOSITORIO>
```

Crie um arquivo `.env` dentro de `cmd/auction/` com o seguinte conteúdo:

```env
MONGODB_URL=mongodb://mongodb:27017
MONGODB_DB=auctions
AUCTION_DURATION=1m
AUCTION_INTERVAL=10s
BATCH_INSERT_INTERVAL=3m
MAX_BATCH_SIZE=5
```

---

## 🚀 Executando a Aplicação

```bash
docker-compose up --build
```

A API estará disponível em: [http://localhost:8080](http://localhost:8080)

---

## ✅ Rodando os Testes

### 1. Suba apenas o banco MongoDB:
```bash
docker-compose up -d mongodb
```

### 2. Execute os testes:
```bash
docker build --target builder -t auction-tester . && \
  docker run --rm --network auction-network \
  -e MONGODB_URL_TEST=mongodb://mongodb:27017 \
  auction-tester go test -v ./...
```

### 3. Pare o banco:
```bash
docker-compose down
```

---

## ↔️ Endpoints da API

| Método | Rota                           | Descrição                                 |
|--------|--------------------------------|-------------------------------------------|
| POST   | `/auction`                    | Cria um novo leilão                        |
| GET    | `/auction`                    | Busca leilões (filtros: status, categoria, produto) |
| GET    | `/auction/:auctionId`         | Busca um leilão específico pelo ID         |
| GET    | `/auction/winner/:auctionId`  | Busca o lance vencedor de um leilão       |
| POST   | `/bid`                        | Cria um novo lance                         |
| GET    | `/bid/:auctionId`            | Busca todos os lances de um leilão        |
| GET    | `/user/:userId`              | Busca um usuário pelo ID                   |

---

## 🙏 Contribuição
Pull requests são bem-vindos! Para melhorias maiores, abra uma issue primeiro para discutirmos o que você gostaria de alterar.

---

## ✉️ Licença
[MIT](LICENSE)
