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

// --- Função de Setup do Banco de Dados ---
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
		auction, err := auction_entity.CreateAuction("Laptop", "Eletronicos", "Desc", auction_entity.New)
		require.Nil(t, err)
		require.Nil(t, repo.CreateAuction(ctx, auction))

		foundAuction, findErr := repo.FindAuctionById(ctx, auction.Id)
		require.Nil(t, findErr)
		require.NotNil(t, foundAuction)
		require.Equal(t, auction.Id, foundAuction.Id)
	})

	// --- Sub-teste: Falha ao Encontrar por ID ---
	t.Run("FindAuctionById - Not Found", func(t *testing.T) {
		foundAuction, err := repo.FindAuctionById(ctx, "non-existent-id")
		require.NotNil(t, err)
		require.Nil(t, foundAuction)
	})

	// --- Sub-teste: Encontrar Múltiplos Leilões ---
	t.Run("FindAuctions - Success", func(t *testing.T) {
		auction1, _ := auction_entity.CreateAuction("Cadeira Gamer", "Móveis", "Desc", auction_entity.Used)
		auction2, _ := auction_entity.CreateAuction("Monitor Gamer", "Eletronicos", "Desc", auction_entity.New)
		auction2.Status = auction_entity.Completed
		require.Nil(t, repo.CreateAuction(ctx, auction1))
		require.Nil(t, repo.CreateAuction(ctx, auction2))

		auctions, err := repo.FindAuctions(ctx, auction_entity.Active, "Móveis", "")
		require.Nil(t, err)
		require.Len(t, auctions, 1)
		require.Equal(t, auction1.Id, auctions[0].Id)
	})
}

// --- Teste 3: Função de Criação de Leilão ---
func TestCreateAuction(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewAuctionRepository(db)
	ctx := context.Background()

	t.Run("quando cria um leilão com sucesso", func(t *testing.T) {
		// Preparação
		auction, err := auction_entity.CreateAuction(
			"Playstation 5",
			"Eletrônicos",
			"Video game de última geração, novo na caixa.",
			auction_entity.New)
		require.Nil(t, err) // Garante que a entidade foi criada sem erro

		// Execução
		createErr := repo.CreateAuction(ctx, auction)
		require.Nil(t, createErr) // Garante que foi salvo no banco sem erro

		// Verificação
		foundAuction, findErr := repo.FindAuctionById(ctx, auction.Id)
		require.Nil(t, findErr)
		require.NotNil(t, foundAuction)
		require.Equal(t, auction.ProductName, foundAuction.ProductName)
	})

	t.Run("quando a criação falha por contexto inválido", func(t *testing.T) {
		// Preparação
		auction, err := auction_entity.CreateAuction(
			"Xbox Series X",
			"Eletrônicos",
			"Outro video game de última geração.",
			auction_entity.New)
		require.Nil(t, err)

		// Criação de um contexto já cancelado para forçar um erro
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Execução e Verificação
		createErr := repo.CreateAuction(cancelledCtx, auction)
		require.NotNil(t, createErr) // DEVE retornar um erro
	})
}
