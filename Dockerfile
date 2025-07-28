# --- Estágio 1: O Ambiente de Construção (Builder) ---
# Usamos uma imagem oficial do Go com Alpine Linux, que é leve.
FROM golang:1.21-alpine AS builder

# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências primeiro.
# Isso aproveita o cache do Docker: as dependências só são baixadas novamente
# se o go.mod ou go.sum mudar.
COPY go.mod go.sum ./

# Baixa as dependências do projeto.
RUN go mod download

# Copia todo o resto do código-fonte do nosso projeto.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 cria um binário estático, sem depender de bibliotecas do sistema.
# Isso é crucial para que ele funcione na nossa imagem final mínima.
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# --- Estágio 2: A Imagem Final (Produção) ---
# Começamos com uma imagem Alpine zerada, que é super pequena.
FROM alpine:latest

# Copia APENAS o binário compilado que foi gerado no estágio anterior.
# Nenhum código-fonte ou ferramenta do Go vem junto.
COPY --from=builder /main /main

# Expõe a porta em que nossa aplicação vai rodar dentro do contêiner.
EXPOSE 8080

# Define o comando que será executado quando o contêiner iniciar.
CMD ["/main"]