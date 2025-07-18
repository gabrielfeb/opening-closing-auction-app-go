package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Função de Setup do Banco de Dados (sem alterações) ---
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URL_TEST")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err, "Failed to connect to MongoDB")

	db := client.Database("test_auction_db_final")

	cleanup := func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if err := db.Drop(cleanupCtx); err != nil {
			t.Logf("Could not drop database: %v", err)
		}
		if err := client.Disconnect(cleanupCtx); err != nil {
			t.Logf("Failed to disconnect from MongoDB: %v", err)
		}
	}
	return db, cleanup
}

// --- Teste Principal e Único ---
func TestAuctionRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewAuctionRepository(db)
	ctx := context.Background()

	// --- Sub-teste: Encontrar por ID com Sucesso ---
	t.Run("FindAuctionById - Success", func(t *testing.T) {
		// Preparação
		auction, err := auction_entity.CreateAuction("Laptop", "Eletronicos", "Desc", auction_entity.New)
		require.NoError(t, err)                              // Garante que a entidade foi criada sem erro
		require.NoError(t, repo.CreateAuction(ctx, auction)) // Garante que foi salva sem erro

		// Execução e Validação
		foundAuction, findErr := repo.FindAuctionById(ctx, auction.Id)
		require.NoError(t, findErr) // DEVE encontrar sem erro
		require.NotNil(t, foundAuction)
		require.Equal(t, auction.Id, foundAuction.Id)
	})

	// --- Sub-teste: Falha ao Encontrar por ID ---
	t.Run("FindAuctionById - Not Found", func(t *testing.T) {
		foundAuction, err := repo.FindAuctionById(ctx, "non-existent-id")
		require.Error(t, err) // DEVE retornar um erro
		require.Nil(t, foundAuction)
	})

	// --- Sub-teste: Encontrar Múltiplos Leilões ---
	t.Run("FindAuctions - Success", func(t *testing.T) {
		// Preparação
		auction1, _ := auction_entity.CreateAuction("Cadeira Gamer", "Móveis", "Desc", auction_entity.Used)
		auction2, _ := auction_entity.CreateAuction("Monitor Gamer", "Eletronicos", "Desc", auction_entity.New)
		auction2.Status = auction_entity.Completed
		require.NoError(t, repo.CreateAuction(ctx, auction1))
		require.NoError(t, repo.CreateAuction(ctx, auction2))

		// Execução e Validação
		auctions, err := repo.FindAuctions(ctx, auction_entity.Active, "Móveis", "")
		require.NoError(t, err)
		require.Len(t, auctions, 1)
		require.Equal(t, auction1.Id, auctions[0].Id)
	})
}
