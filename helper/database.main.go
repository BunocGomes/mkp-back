package helper

import (
	"log"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/models"

	"gorm.io/gorm"
)

func InitializeDatabase() {
	database.ConnectDB()
	database.ConnectMongo()

	schemaName := "marketplace"
	createSchema(database.DB, schemaName)

	// --- REMOVA OU COMENTE ESTA LINHA ---
	// A configuração TablePrefix no gorm.Config já cuida disso de forma mais eficaz.
	/*
		if err := database.DB.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
			log.Fatalf("Falha ao definir o search_path: %v", err)
		}
	*/

	migrateModels(database.DB)
	SeedData(database.DB)
	log.Println("Inicialização do banco de dados concluída com sucesso.")
}

// O resto do arquivo (createSchema, migrateModels, etc.) continua igual.

func createSchema(db *gorm.DB, schemaName string) {
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName).Error; err != nil {
		log.Fatalf("Erro ao criar o schema %s: %v", schemaName, err)
	}
	log.Printf("Schema %s verificado/criado.", schemaName)
}

func migrateModels(db *gorm.DB) {
	log.Println("Iniciando migração dos modelos...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Perfil{},
		&models.Digest{},
	)
	if err != nil {
		log.Fatalf("Erro durante a migração: %v", err)
	}
	log.Println("Migração dos modelos concluída.")
}

func SeedData(db *gorm.DB) {
	log.Println("Passo de popular dados (seeding) pulado (não implementado).")
}
