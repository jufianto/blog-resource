package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

var sqlConnStr string

func init() {
	viper.SetConfigFile("../env/env.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper config not found: %v", err)
	}

	sqlConnStr = fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
	)
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("Usage: go run migrate.go [up|down|force] [version]")
	}

	conn, err := sql.Open("pgx", sqlConnStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}
	defer conn.Close()

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("failed to create db instance: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get migration version: %v", err)
	}

	fmt.Printf("Current migration version: %d, Dirty: %v\n", version, dirty)

	switch args[1] {
	case "up":
		fmt.Println("Running migrations up")
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		} else if err == migrate.ErrNoChange {
			fmt.Println("No schema changes to apply")
		} else {
			fmt.Println("Migration completed successfully")
		}

	case "down":
		fmt.Printf("Downgrade database, are you sure? (Y/n): ")
		var down string
		fmt.Scan(&down)
		if down == "Y" {
			fmt.Println("Running migrations down")
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Downgrade failed: %v", err)
			} else if err == migrate.ErrNoChange {
				fmt.Println("No migrations to roll back")
			} else {
				fmt.Println("Downgrade completed successfully")
			}
		} else {
			fmt.Println("Abort downgrade")
		}

	case "force":
		if len(args) <= 2 {
			log.Fatal("Usage: go run migrate.go force VERSION")
		}
		version, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		fmt.Printf("Forcing migration version to %d\n", version)
		if err := m.Force(version); err != nil {
			log.Fatalf("Force migration failed: %v", err)
		}
		fmt.Println("Force migration completed successfully")

	default:
		log.Fatal("Usage: go run migrate.go [up|down|force] [version]")
	}
}
