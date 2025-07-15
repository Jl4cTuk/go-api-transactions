package postgresql

import (
	"database/sql"
	"fmt"
	"infotex/internal/config"
	"infotex/internal/domain/model"
	walletsmaker "infotex/internal/lib/walletsMaker"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type DBServer struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

func New(storagePath config.DBServer) (*Storage, error) {
	const op = "storage.postgresql.New"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		storagePath.Address, storagePath.Port, storagePath.User, storagePath.Password, storagePath.DBname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddWallet(address string, balance int) (int64, error) {
	const op = "storage.postgresql.AddWallet"

	stmt, err := s.db.Prepare("INSERT INTO wallets (address, balance) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var lastInsertId int64 = 0
	err = stmt.QueryRow(address, balance).Scan(&lastInsertId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return lastInsertId, nil
}

func (s *Storage) GetWalletBalance(address string) (int64, error) {
	const op = "storage.postgresql.GetWalletBalance"

	stmt, err := s.db.Prepare("SELECT balance FROM wallets WHERE address=$1")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var balance int64 = 0
	err = stmt.QueryRow(address).Scan(&balance)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}

func (s *Storage) ProcessTransactions(senderAdress, receiverAdress string, amount int) error {
	const op = "storage.postgresql.GetWalletBalance"

	stmt, err := s.db.Prepare("CALL transfer($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(senderAdress, receiverAdress, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetLastTransactions(count int) ([]model.Transaction, error) {
	const op = "storage.postgresql.GetLastTransactions"

	stmt, err := s.db.Query("SELECT sender, receiver, amount FROM transactions LIMIT $1", count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var transactions []model.Transaction

	for stmt.Next() {
		var tx model.Transaction
		if err := stmt.Scan(&tx.From, &tx.To, &tx.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (s *Storage) GenRandomWallet(amont int) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var err error
	for i := 0; i < amont; i++ {
		_, err = s.AddWallet(walletsmaker.GenAddress(8, r), 100)
		if err != nil {
			return err
		}
	}
	return nil
}
