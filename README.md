```markdown
# API de Locadora de Filmes em Go com Docker

Este projeto consiste em uma API RESTful com operações CRUD (Criar, Ler, Atualizar, Excluir) para gerenciar filmes. A aplicação foi desenvolvida em Go e é executada em um ambiente multi-container orquestrado com Docker Compose.

O ambiente é composto por dois serviços:
- `app`: O contêiner da aplicação Go.
- `db`: O contêiner do banco de dados PostgreSQL para persistência dos dados.

## Requisitos da Atividade Atendidos

- **Dockerfile:** Utiliza build de múltiplos estágios e imagens Alpine para otimização.
- **Docker Compose:** Orquestra os serviços da aplicação e do banco de dados.
- **Volumes:** Garante a persistência dos dados do PostgreSQL através de um volume nomeado.
- **Rede Customizada:** Isola a comunicação entre os contêineres em uma rede `bridge` customizada.
- **Variáveis de Ambiente:** Configurações sensíveis (credenciais do banco) são gerenciadas através de um arquivo `.env`, sem valores "hardcoded" no código.
- **Segurança:** A aplicação se conecta ao banco com um usuário específico e não com o superusuário padrão.

## Pré-requisitos

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install) (para testes locais, se desejado)
- [Git](https://git-scm.com/)

## Como Executar

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/Antonio-Elvis-Dev/locadora-api
   cd locadora-api

2. crie o arquivo .env com o nome das environments presentes no .env.exemple

3. Execute o docker
   
   Após garantir que o docker está funcionando execute:
    - docker-compose up --build -d
    e 
    - docker-compose ps
    para verificar se o nosso serviço da aplicação está funcionado. 

4. Após é só fazer as devidas requisições para a nossa API

    - GET -> /movies
    - GET -> /movies/{id}
    - POST -> /movies => precisa de title - director - year 
    - PUT -> /movies/${id} => precisa de title - director - year 
    - DELETE -> /movies/${id}
