package db

// target, subject, organization, preDate, postDate, certificateLife, protocol
import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/eogns47/SSL_Info_Crawler/src/models"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func Connect() (*sql.DB, error) {
	os.Clearenv()

	dbPath := "./config/dbConfig.env"
	err := godotenv.Load(dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading DB config file: "+err.Error())
	}
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  os.Getenv("DB_NETWORK"),
		Addr:                 os.Getenv("DB_ADDRESS"),
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               os.Getenv("DB_NAME"),
	}
	os.Clearenv()
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "mysql.NewConnector failed")
	}
	db := sql.OpenDB(connector)
	err = CreateCertInfoTableIfNotExists(db)
	if err != nil {
		return nil, errors.Wrap(err, "check and create table failed : ")
	}

	return db, err
}

func CreateCertInfoTableIfNotExists(db *sql.DB) error {
	isExist, err := tableExists(db, "cert_info")
	if isExist {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "check table exist failed : ")
	}

	// 테이블 생성 쿼리
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			cert_info_id INT AUTO_INCREMENT PRIMARY KEY,
			url_id VARCHAR(100),
			subject VARCHAR(50),
			issuer VARCHAR(50),
            valid_since VARCHAR(50),
            valid_until VARCHAR(50),
            certificate_life VARCHAR(50),
            protocol VARCHAR(50),
            insert_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`, "cert_info")

	// 테이블 생성

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return errors.Wrap(err, "Create cert_info Table Query failed")
	}

	fmt.Printf("\nTable cert_info created successfully.\n")
	return nil
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	rows, err := db.Query(query)
	if err != nil {
		errors.Wrap(err, "Error checking if table exists: ")
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func InsertDBInfo(certInfo models.Certificate) error {
	// db := db.DB{}
	// db.SaveCertificate(target, subject, organization, preDate, postDate, certificateLife, protocol)
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO cert_info (url_id, subject, issuer, valid_since, valid_until, certificate_life, protocol) VALUES (?, ?, ?, ?, ?, ?, ?)", certInfo.Target, certInfo.Subject, certInfo.Organization, certInfo.PreDate, certInfo.PostDate, certInfo.CertificateLife, certInfo.Protocol)

	if err != nil {
		return err
	}
	fmt.Println("\nDB Insert Success!")
	return nil
}
