package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"infotex/internal/config"
	"infotex/internal/domain/model"
	wallet "infotex/internal/lib/random"
	"infotex/internal/storage"
	"math/rand"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New initializes a new PostgreSQL storage connection using the provided configuration.
// It returns a pointer to the Storage struct or an error if the connection fails.
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

// AddWallet inserts a new wallet record into the database with the specified address and balance.
// It returns the ID of the newly inserted wallet and an error, if any occurs during the operation.
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

// GetWalletBalance retrieves the balance of a wallet identified by the given address.
// Returns the wallet balance as int64 if found, or -1 and an error if the query fails or the address does not exist.
func (s *Storage) GetWalletBalance(address string) (int64, error) {
	const op = "storage.postgresql.GetWalletBalance"

	isExists, err := s.walletExists(address)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	if !isExists {
		return -1, fmt.Errorf("%s: %w", op, storage.ErrWalletNotFound)
	}

	stmt, err := s.db.Prepare("SELECT balance FROM wallets WHERE address=$1")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var balance int64
	err = stmt.QueryRow(address).Scan(&balance)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}

// walletExists checks does the wallet exists in database
// Returns true or false and an error if something is wrong
func (s *Storage) walletExists(address string) (bool, error) {
	const op = "storage.postgresql.WalletExists"

	stmt, err := s.db.Prepare("SELECT EXISTS(SELECT 1 FROM wallets WHERE address=$1)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	err = stmt.QueryRow(address).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

// ProcessTransactions executes a stored procedure to transfer an amount from one wallet to another.
func (s *Storage) ProcessTransactions(senderAdress, receiverAdress string, amount int) error {
	const op = "storage.postgresql.GetWalletBalance"

	stmt, err := s.db.Prepare("CALL transfer($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(senderAdress, receiverAdress, amount)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "P0010":
				return fmt.Errorf("%s: %w", op, storage.ErrInvalidWallet)
			case "P0011":
				return fmt.Errorf("%s: %w", op, storage.ErrInsufficientFunds)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetLastTransactions retrieves the last N transactions from the database.
// It returns a slice of Transaction models or an error if the query fails.
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

// GenRandomWallet generates a specified number of random wallets with an initial balance of 100.
func (s *Storage) GenRandomWallet(amont int) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var err error
	for i := 0; i < amont; i++ {
		_, err = s.AddWallet(wallet.GenAddress(8, r), 100)
		if err != nil {
			return err
		}
	}
	return nil
}
