package storage

import (
	"database/sql"
	"projects/gobank/types"

	"github.com/lib/pq"
)

type Storage interface {
  CreateAccount(*types.Account) error
  DeleteAccount(int) error
  UpdateAccount(*types.Account) error
  GetAccounts() ([]*types.Account, error)
  GetAccountById(int) (*types.Account, error)
}

type PostgresStore struct {
  db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
  connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
  
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    return nil, err
  }

  if err := db.Ping(); err != nil {
    return nil, err
  }

  return &PostgresStore{
    db: db,
  }, nil
}

func (s *PostgresStore) Init() error {
  return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
  query := `create table if not exists account (
    id serial primary key,
    first_name varchar(50), 
    last_name varchar(50),
    number serial,
    balance serial,
    created_at timestamp
  )`

  _, err := s.db.Exec(query)
  return err
}

func (s *PostgresStore) CreateAccount(account *types.Account) error {
  query := `insert into account
    (first_name, last_name, number, balance, created_at)
    values ($1, $2, $3, $4, $5)`

  _, err := s.db.Query(
    query,
    account.FirstName,
    account.LastName,
    account.Number,
    account.Balance,
    account.CreatedAt,
  )
  
  if err != nil {
    return err
  }

  return nil
}

func (s *PostgresStore) UpdateAccount(*types.Account) error {
  return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
  return nil
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
  rows, err := s.db.Query("select * from account")
  if err != nil {
    return nil, err
  }

  accounts := []*types.Account{}
  for rows.Next() {
    account := new(types.Account)
    err := rows.Scan(
      &account.ID,
      &account.FirstName,
      &account.LastName,
      &account.Number,
      &account.Balance,
      &account.CreatedAt)

    if err != nil {
      return nil, err
    } 

    accounts = append(accounts, account)
  }

  return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*types.Account, error) {
  return nil, nil
}






