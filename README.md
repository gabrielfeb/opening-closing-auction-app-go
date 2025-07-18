# Sistema de Leilões em Go

Este é um sistema de leilões robusto desenvolvido em Go, utilizando uma arquitetura limpa para garantir escalabilidade e manutenibilidade. A API permite a criação de leilões, o registro de lances e a consulta de resultados, incluindo uma funcionalidade de fechamento automático de leilões baseada em tempo.

## Origem do Projeto

Este projeto é baseado no repositório original da [Full Cycle Tecnologia](https://faculdadefullcycle.edu.br/), clonado a partir de:

[https://github.com/devfullcycle/labs-auction-goexpert](https://github.com/devfullcycle/labs-auction-goexpert)

A lógica principal do sistema foi mantida, com melhorias desenvolvidas conforme o enunciado da atividade da pós-graduação **Go Expert** da Full Cycle. As principais alterações incluem:

- Adição de funcionalidade de fechamento automático de leilões após tempo definido via variáveis de ambiente.
- Implementação de goroutines para monitoramento de leilões ativos.
- Criação de testes para validar o encerramento automático.

## Funcionalidades

- Criação de leilões com tempo configurável.
- Registro de lances por usuários.
- Consulta de leilões por status, categoria ou nome do produto.
- Consulta de lances por leilão.
- Fechamento automático de leilões por tempo.
- Processamento de lances em lote para maior eficiência.

## Tecnologias Utilizadas

- Linguagem: Go
- Framework Web: Gin Gonic
- Banco de Dados: MongoDB
- Containerização: Docker & Docker Compose
- Testes: testing & Testify

## Arquitetura

O projeto segue os princípios da Clean Architecture:

- `internal/entity`: Entidades e regras de negócio.
- `internal/usecase`: Casos de uso.
- `internal/infra`: Implementação de API, banco de dados e adaptadores.

## Como Executar

### Requisitos

- Docker
- Docker Compose

### Clonagem e Configuração

```bash
git clone <URL_DO_SEU_REPOSITORIO>
cd <NOME_DO_REPOSITORIO>
```

Crie um arquivo `.env` dentro da pasta `cmd/auction/` com o seguinte conteúdo:

```env
MONGODB_URL=mongodb://mongodb:27017
MONGODB_DB=auctions
AUCTION_DURATION=1m
AUCTION_INTERVAL=10s
BATCH_INSERT_INTERVAL=3m
MAX_BATCH_SIZE=5
```

### Execução

```bash
docker-compose up --build
```

A API estará acessível em: http://localhost:8080

## Executando os Testes

### Etapa 1: Subir apenas o MongoDB
```bash
docker-compose up -d mongodb
```

### Etapa 2: Rodar os testes
```bash
docker build --target builder -t auction-tester . && \
  docker run --rm --network auction-network \
  -e MONGODB_URL_TEST=mongodb://mongodb:27017 \
  auction-tester go test -v ./...
```

### Etapa 3: Encerrar serviços
```bash
docker-compose down
```

## Endpoints da API

| Método | Rota                           | Descrição                                 |
|--------|--------------------------------|-------------------------------------------|
| POST   | `/auction`                    | Cria um novo leilão                        |
| GET    | `/auction`                    | Busca leilões (filtros: status, categoria, produto) |
| GET    | `/auction/:auctionId`         | Busca um leilão específico pelo ID         |
| GET    | `/auction/winner/:auctionId`  | Busca o lance vencedor de um leilão       |
| POST   | `/bid`                        | Cria um novo lance                         |
| GET    | `/bid/:auctionId`            | Busca todos os lances de um leilão        |
| GET    | `/user/:userId`              | Busca um usuário pelo ID                   |

## Licença

Este projeto é distribuído sob a Licença MIT.

---

Para mais informações sobre o curso e materiais complementares, acesse:
[https://faculdadefullcycle.edu.br](https://faculdadefullcycle.edu.br)
