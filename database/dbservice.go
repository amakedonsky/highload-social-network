package database

import (
	"amakedonsky/highload-social-network/helpers"
	"amakedonsky/highload-social-network/models"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	SocialNetworkDB *sql.DB
)

func Init() {
	err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = createPersonalPageTable(SocialNetworkDB)
	if err != nil {
		log.Fatal(err)
	}
	err = createFriendshipTable(SocialNetworkDB)
	if err != nil {
		log.Fatal(err)
	}
}

func dbConnection() error {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "db:3306",
		DBName: "highload_social_network",
	}

	var err error
	SocialNetworkDB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := SocialNetworkDB.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+cfg.DBName)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return err
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return err
	}
	log.Printf("rows affected %d\n", no)

	SocialNetworkDB.SetMaxOpenConns(20)
	SocialNetworkDB.SetMaxIdleConns(20)
	SocialNetworkDB.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err = SocialNetworkDB.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return err
	}

	log.Printf("Connected to DB %s successfully\n", cfg.DBName)
	return nil
}

func createUserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS user(user_id int primary key auto_increment, username text, password text, 
        created_at datetime default CURRENT_TIMESTAMP)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating user table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when creating table user: %d", rows)
	return nil
}

func createPersonalPageTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS personal_page(
    	page_id int primary key auto_increment, 
    	password text not null,
    	email text not null,
    	first_name text, 
    	last_name text, 
        age int, 
        sex text, 
        address text, 
        created_at datetime default CURRENT_TIMESTAMP, 
        updated_at datetime default CURRENT_TIMESTAMP)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating personal_page table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when creating table personal_page: %d", rows)
	return nil
}

func createFriendshipTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS person_friendship (
		person_id INT NOT NULL,
		friend_id INT NOT NULL,
		PRIMARY KEY (person_id, friend_id),
		CONSTRAINT constr_friendship_person_fk
			FOREIGN KEY person_fk (person_id) REFERENCES personal_page(page_id)
			ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT constr_friendship_friend_fk
			FOREIGN KEY friend_fk (friend_id) REFERENCES personal_page(page_id)
			ON DELETE CASCADE ON UPDATE CASCADE
	) `

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating person_friendship table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}

	log.Printf("Rows affected when creating table person_friendship: %d", rows)
	return nil
}

// GetUserByEmail handles fetching user details by email
func GetUserByEmail(c context.Context, email string) (page models.PersonalPage, err error) {
	row := SocialNetworkDB.QueryRowContext(c, "SELECT page_id, email, first_name, last_name FROM personal_page WHERE email = ?", email)

	err = row.Scan(&page.Id, &page.Email, &page.FirstName, &page.LastName)

	if err != nil {
		log.Fatal(err)
		return page, err
	}

	return page, nil
}

// GetPasswordByEmail handles fetching hashed password by email
func GetPasswordByEmail(c context.Context, email string) (page models.PersonalPage, err error) {
	row := SocialNetworkDB.QueryRowContext(c, "SELECT page_id, password, email FROM personal_page WHERE email = ?", email)

	err = row.Scan(&page.Id, &page.Password, &page.Email)

	if err != nil {
		log.Fatal(err)
		return page, err
	}

	return page, nil
}

func CreatePersonalPage(c context.Context, newPage models.PersonalPage) (int64, error) {
	result, err := SocialNetworkDB.ExecContext(c,
		"INSERT INTO personal_page (email, password, first_name, last_name, age, sex, address) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newPage.Email,
		helpers.GeneratePasswordHash(newPage.Password),
		newPage.FirstName,
		newPage.LastName,
		newPage.Age,
		newPage.Sex,
		newPage.Address)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func UpdatePersonalPage(c context.Context, page models.PersonalPage) (int64, error) {
	result, err := SocialNetworkDB.ExecContext(c,
		"update personal_page set first_name = ?, last_name = ?, age = ?, sex = ?, address = ? where page_id = ?",
		page.FirstName,
		page.LastName,
		page.Age,
		page.Sex,
		page.Address,
		page.Id)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func FetchPersonalPage(c context.Context, id string) (page models.PersonalPage, err error) {
	row := SocialNetworkDB.QueryRowContext(c, "SELECT email, first_name, last_name, age, sex, address FROM personal_page WHERE page_id = ?", id)

	err = row.Scan(&page.Email, &page.FirstName, &page.LastName, &page.Age, &page.Sex, &page.Address)

	if err != nil {
		log.Fatal(err)
		return page, err
	}

	return page, nil
}

func AddToFriends(c context.Context, personId string, friendId string) error {
	_, err := SocialNetworkDB.ExecContext(c,
		"INSERT INTO person_friendship (person_id, friend_id) VALUES (?, ?)", personId, friendId)

	return err
}

func DelFromFriends(c context.Context, personId string, friendId string) error {
	_, err := SocialNetworkDB.ExecContext(c,
		"DELETE FROM person_friendship WHERE person_id = ? AND friend_id = ?", personId, friendId)

	return err
}

func GetAllFriends(c context.Context, personId string) ([]models.PersonalPage, error) {
	rows, err := SocialNetworkDB.QueryContext(c,
		"select email, first_name, last_name, age, sex, address from personal_page "+
			"join person_friendship pf on personal_page.page_id = pf.friend_id where pf.person_id = ?", personId)

	if err != nil {
		return []models.PersonalPage{}, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var pages []models.PersonalPage
	for rows.Next() {
		var page models.PersonalPage
		if err := rows.Scan(&page.Email, &page.FirstName, &page.LastName, &page.Age, &page.Sex, &page.Address); err != nil {
			return []models.PersonalPage{}, err
		}
		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return []models.PersonalPage{}, err
	}

	return pages, nil
}
