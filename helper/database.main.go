package helper

import (
	"log"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/models"

	"gorm.io/gorm"
)

func InitializeDatabase() {
	database.ConnectDB()
	database.InitMongoDB()

	schemaName := "marketplace"
	createSchema(database.DB, schemaName)

	migrateModels(database.DB)

	SeedData(database.DB)
	log.Println("Inicialização do banco de dados concluída com sucesso.")
}

func createSchema(db *gorm.DB, schemaName string) {
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName).Error; err != nil {
		log.Fatalf("Erro ao criar o schema %s: %v", schemaName, err)
	}
	log.Printf("Schema %s verificado/criado.", schemaName)
}

func migrateModels(db *gorm.DB) {
	log.Println("Iniciando migração dos modelos...")
	err := db.AutoMigrate(
		&models.Skill{},
		&models.Empresa{},
		&models.User{},
		&models.Perfil{},
		&models.Digest{},
		&models.Projeto{},
		&models.Proposta{},
		&models.Contrato{},
		&models.Avaliacao{},

		&models.PortfolioItem{},
		&models.SocialLink{},
	)

	if err != nil {
		log.Fatalf("Erro durante a migração: %v", err)
	}
	log.Println("Migração dos modelos concluída.")
}

func SeedData(db *gorm.DB) {
	log.Println("Passo de popular dados (seeding) pulado (não implementado).")
}
