package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"auctions-service/pb"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAuction(ctx context.Context, auction *pb.Auction) error {
	query := `
		INSERT INTO tb_auctions
			(id, title, description, organizer_id, max_participants, participants_num, starts_at, ends_at, status)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query, auction.Id, auction.Title, auction.Description, auction.OrganizerID, auction.MaxParticipants,
		auction.ParticipantsNumber, auction.StartsAt, auction.EndsAt, "WAITING")
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) UpdateAuction(ctx context.Context, auction *pb.Auction) error {
	query := `
		UPDATE tb_auctions 
		SET 
		title=$1, description=$2, organizer=$3, max_participants=$4, participants_num=$5, starts_at=$6, ends_at=$7, status=$8
		WHERE id=$9
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) CancelAuction(ctx context.Context, id string) error {
	query := `UPDATE tb_auctions SET status=$1 WHERE id=$2`

	_, err := r.db.Exec(query, "", id)
	if err != nil {

		return err
	}

	return nil
}

func (r *Repository) StartAuction(ctx context.Context, id string) error {
	query := `UPDATE tb_auctions SET status=$1 WHERE id=$2`

	_, err := r.db.Exec(query, "", id)
	if err != nil {

		return err
	}

	return nil
}

func (r *Repository) AddParticipant(ctx context.Context, auction_id, participant_id string) error {
	query := `INSERT INTO tb_participants (auction_id, participant_id) VALUES ($1, $2)`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, auction_id, participant_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) RemoveParticipant(ctx context.Context, auction_id, participant_id string) error {
	query := `DELETE FROM tb_participants WHERE auction_id=$1 AND participant_id=$2`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, auction_id, participant_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) AddProduct(ctx context.Context, product *pb.Product) error {
	query := `
		INSERT INTO tb_products 
			(id, auction_id, title, description, start_price, params)
		VALUES
			($1, $2, $3, $4, $5, $6)
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	params := product.Params

	_, err = tx.Exec(query, product.Id, product.AuctionId, product.Title, product.Description, product.StartPrice, params)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product *pb.Product) error {
	query := `
		UPDATE tb_products SET title=$1, description=$2, start_price=$3, params=$4 WHERE id=$5
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, product.Title, product.Description, product.StartPrice, product.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM tb_products WHERE id=$1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
