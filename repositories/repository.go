package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/valikhan03/command-service/pb"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAuction(ctx context.Context, auction *pb.Auction) (string, error) {
	query := `
		INSERT INTO tb_auctions
			(id, title, description, organizer_id, max_participants, participants_number, starts_at, ends_at)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8) 
	`
	newid := uuid.New()
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(query, newid, auction.Title, auction.Description, auction.OrganizerID, auction.MaxParticipants,
		auction.ParticipantsNumber, auction.StartsAt, auction.EndsAt)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	query = `INSERT INTO tb_attempt_requirements 
				(auction_id, approve_required, enter_fee_required, enter_fee_payment)
			VALUES
				($1, $2, $3, $4)`
	_, err = tx.Exec(query, newid, auction.Attempt_Requirements.ApproveRequired, auction.Attempt_Requirements.EnterFeeRequired, auction.Attempt_Requirements.EnterFeeAmount)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return newid.String(), nil
}

func (r *Repository) UpdateAuction(ctx context.Context, auction *pb.Auction) error {
	query := `
		UPDATE tb_auctions 
		SET 
		title=$1, description=$2, organizer_id=$3, max_participants=$4, 
		participants_number=$5, starts_at=$6, ends_at=$7
		WHERE id=$8
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query, auction.Title, auction.Description, auction.OrganizerID, auction.MaxParticipants,
		auction.ParticipantsNumber, auction.StartsAt, auction.EndsAt, auction.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) DeleteAuction(ctx context.Context, id string) error {
	query := `DELETE from tb_auctions WHERE id=$1`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) AddParticipant(ctx context.Context, auction_id string, participant_id int32) (bool, error) {
	query := `SELECT approve_required, enter_fee_required, enter_fee_amount 
				FROM tb_attempt_requirements WHERE auction_id=$1`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return false, err
	}
	var requirements map[string]int 

	rows, err := tx.Query(query, auction_id)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&requirements)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if requirements["approve_required"] == 1 || requirements["enter_fee_amount"] == 1{
		query = `INSERT INTO tb_attempt_requests (auction_id, user_id) VALUES ($1, $2)`
		_, err = tx.Exec(query, auction_id, participant_id)
		if err != nil{
			tx.Rollback()
			return false, err
		}
		return true, nil
	}

	
	query = `INSERT INTO tb_participants (auction_id, user_id) VALUES ($1, $2)`

	_, err = tx.Exec(query, auction_id, participant_id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return false, nil
}

func (r *Repository) RemoveParticipant(ctx context.Context, auction_id string, participant_id int32) error {
	query := `DELETE FROM tb_participants WHERE auction_id=$1 AND user_id=$2`

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

func (r *Repository) AddLot(ctx context.Context, lot *pb.Lot) (int, error) {
	query := `
		INSERT INTO tb_lots 
			(auction_id, title, description, start_price, params)
		VALUES
			($1, $2, $3, $4, $5) 
		RETURNING ID
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return -1, err
	}

	var lot_params json.RawMessage
	lot_params, err = json.Marshal(lot.Params)

	rows, err := tx.Query(query, lot.AuctionId, lot.Title, lot.Description, lot.StartPrice, lot_params)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer rows.Close()

	var id int = 0
	if rows.Next() {
		rows.Scan(&id)
	} else {
		tx.Rollback()
		return -1, errors.New("No lot id returned from DB ")
	}
	tx.Commit()
	return id, nil
}

func (r *Repository) UpdateLot(ctx context.Context, lot *pb.Lot) error {
	query := `
		UPDATE tb_lots SET title=$1, description=$2, start_price=$3, params=$4 WHERE id=$5
	`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	var lot_params json.RawMessage
	lot_params, err = json.Marshal(lot.Params)

	_, err = tx.Exec(query, lot.Title, lot.Description, lot.StartPrice, lot_params, lot.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) DeleteLot(ctx context.Context, id string) error {
	query := `DELETE FROM tb_lots WHERE id=$1`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) AddMediaInfo(ctx context.Context, mediaInfo *pb.MediaInfo) error {
	return nil
}
