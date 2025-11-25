# 🚀 Plataforma de Marketplace Freelance

Este projeto é uma plataforma completa de conexão entre **Empresas** e **Freelancers**, desenvolvida como requisito para a disciplina de **Project Lab**. O sistema permite a publicação de projetos, envio de propostas, gestão de contratos e comunicação em tempo real.

## 📋 Sobre o Projeto

A aplicação resolve o problema da fragmentação na contratação de serviços, oferecendo um ambiente seguro onde:
1.  **Empresas** publicam demandas e orçamentos.
2.  **Freelancers** enviam propostas técnicas e financeiras.
3.  Um **Contrato** é gerado automaticamente ao aceitar uma proposta.
4.  As partes negociam via **Chat em Tempo Real**.
5.  Ao final, ocorre a **Avaliação Mútua** (Rating) para construir reputação.

## 🛠️ Tecnologias Utilizadas

O projeto utiliza uma arquitetura moderna, separando Backend e Frontend, conteinerizada para fácil execução.

### **Backend (API)**
* **Linguagem:** Go (Golang) 1.21+
* **Framework Web:** Gin Gonic
* **ORM:** GORM
* **Real-time:** Gorilla WebSocket
* **Autenticação:** JWT (JSON Web Tokens)

### **Frontend (Interface)**
* **Framework:** React + Vite
* **Linguagem:** TypeScript
* **Estilização:** Tailwind CSS + Shadcn/UI
* **HTTP Client:** Axios

### **Banco de Dados (Híbrido)**
* **PostgreSQL:** Dados relacionais (Usuários, Projetos, Propostas, Contratos, Avaliações).
* **MongoDB:** Dados não-estruturados e volumosos (Histórico de Chat).

### **Infraestrutura**
* **Docker & Docker Compose:** Orquestração de todos os serviços.

---

## ⚙️ Como Rodar o Projeto

### Pré-requisitos
* Docker e Docker Compose instalados.

### Passo a Passo

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/seu-usuario/mkp-back.git](https://github.com/seu-usuario/mkp-back.git)
    cd mkp-back
    ```

2.  **Configure as Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto (baseado no exemplo abaixo):
    ```env
    # Banco Relacional (Postgres)
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=marketplace
    DB_PORT=5432

    # Banco NoSQL (Mongo)
    MONGO_URI=mongodb://mongodb:27017

    # Aplicação
    PORT=8080
    JWT_SECRET=sua_chave_secreta_aqui
    ```

3.  **Inicie a Aplicação:**
    Execute o comando mágico para subir Backend, Frontend e Bancos de dados:
    ```bash
    docker-compose up --build
    ```

4.  **Acesse:**
    * **Frontend:** [http://localhost:5173](http://localhost:5173)
    * **Backend API:** [http://localhost:8080/api/v1](http://localhost:8080/api/v1)
    * **PgAdmin (Banco de Dados):** [http://localhost:5050](http://localhost:5050)

---

## 🔌 Principais Endpoints da API

A API segue o padrão RESTful (exceto o Chat que usa WebSocket).

### Autenticação
* `POST /api/v1/auth/login` - Login (Retorna Token JWT)
* `POST /api/v1/users/` - Registro de Freelancer
* `POST /api/v1/enterprise/register` - Registro de Empresa

### Fluxo de Trabalho
* `GET /api/v1/projetos` - Listar projetos abertos
* `POST /api/v1/proposals` - Enviar proposta (Freelancer)
* `POST /api/v1/proposals/:id/accept` - Aceitar proposta (Gera Contrato)
* `GET /api/v1/contracts/meus` - Listar contratos ativos
* `PATCH /api/v1/contracts/:id/status` - Finalizar ou Cancelar contrato

### Chat & Social
* `WS /api/v1/chat/ws?token=...` - Conexão WebSocket
* `GET /api/v1/chat/history/:userId` - Histórico de mensagens
* `GET /api/v1/profiles/:id` - Ver perfil público
* `POST /api/v1/avaliacoes` - Avaliar usuário (após contrato concluído)

---

## 📂 Estrutura de Pastas (Backend)
 ```bash ├── chat/ # Lógica do WebSocket (Hub e Client) 
├── controller/ # Handlers das rotas HTTP 
├── database/ # Conexão com Postgres e Mongo 
├── dto/ # Objetos de Transferência de Dados 
├── middleware/ # Autenticação e CORS 
├── models/ # Estruturas do Banco de Dados 
├── routes/ # Definição das rotas da API 
├── service/ # Regras de Negócio 
└── utils/ # Funções auxiliares (Hash, JWT)
```


## 👥 Autores

* Bruno Ricardo Cavalcante Gomes
* João Pedro Carvalho Cândido
* Jorge Raphael Martins Braga Braz

---
*Centro Universitário de Maceió - 2025*
