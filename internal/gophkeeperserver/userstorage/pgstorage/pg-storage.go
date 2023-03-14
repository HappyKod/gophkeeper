package pgstorage

import (
	"context"
	"database/sql"
	"errors"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/models"

	_ "github.com/lib/pq"
)

type PgStorage struct {
	connect *sql.DB
}

func New(uri string) (*PgStorage, error) {
	connect, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &PgStorage{connect: connect}, nil
}

func (PS *PgStorage) Ping() error {
	if err := PS.connect.Ping(); err != nil {
		return err
	}
	err := createTables(PS.connect)
	if err != nil {
		return err
	}
	return nil
}

func (PS *PgStorage) Close() error {
	if err := PS.connect.Close(); err != nil {
		return err
	}
	return nil
}

func createTables(connect *sql.DB) error {
	_, err := connect.Exec(`
	create table if not exists public.users(
		login_user text primary key,
		password_user text,
		create_user timestamp default now()
	);
	
	create table if not exists public.orders(
		 number_order text primary key,
		 login_user text,
		 status_order varchar(50),
		 accrual_order double precision,
		 uploaded_order timestamp default now(),
		 created_order timestamp default now(),
		 foreign key (login_user) references public.users (login_user)
	);
	
	create table if not exists public.withdraws(
		 login_user text,
		 number_order text,
		 sum double precision,
		 uploaded_order timestamp default now()
	);
	`)
	if err != nil {
		return err
	}
	return nil
}

func (PS *PgStorage) AddUser(ctx context.Context, user models.User) error {
	result, err := PS.connect.ExecContext(ctx,
		`insert into public.users (login_user, password_user) values ($1, $2) on conflict do nothing`,
		user.Login, user.Password)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return constans.ErrorNoUNIQUE
	}
	return nil
}

func (PS *PgStorage) AuthenticationUser(ctx context.Context, user models.User) (bool, error) {
	var done int
	err := PS.connect.QueryRowContext(ctx, `select count(1) from public.users where login_user=$1 and password_user=$2`,
		user.Login, user.Password).Scan(&done)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if done == 0 {
		return false, nil
	}
	return true, nil
}
