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
   git clone 
   cd 