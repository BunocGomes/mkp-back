package database

import (
	"fmt"
	//"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema" // <-- ADICIONE ESTE NOVO IMPORT
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// AJUSTE AQUI: Adicionamos uma configuração ao GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Esta configuração diz ao GORM para adicionar "marketplace."
		// antes do nome de todas as tabelas.
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "marketplace.",
		},
	})

	if err != nil {
		panic("Erro ao conectar ao banco de dados: " + err.Error())
	}

	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso.")
	DB = db
}
