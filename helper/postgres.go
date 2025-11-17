package helper

import "github.com/jackc/pgx/v5/pgtype"

func ParseUUID(str string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(str)
	return uuid, err
}
