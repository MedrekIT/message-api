package automated

import (
	"context"
	"fmt"
	"log"

	"github.com/MedrekIT/message-api/internal/database"
)

func DbCleanup(ctx context.Context, db *database.Queries) error {
	log.Println("Cleaning up the database!")
	err := db.ClearRefreshTokens(ctx)
	if err != nil {
		return fmt.Errorf("error while removing expired refresh tokens - %w\n", err)
	}
	err = db.ClearInvitationLinks(ctx)
	if err != nil {
		return fmt.Errorf("error while removing expired invitation links - %w\n", err)
	}
	return nil
}
