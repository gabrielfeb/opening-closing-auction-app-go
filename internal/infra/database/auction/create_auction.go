package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection      *mongo.Collection
	auctionInterval time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection:      database.Collection("auctions"),
		auctionInterval: getAuctionInterval(),
	}

	go repo.closeAuctionsRoutine(context.Background())

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

// CORREÇÃO APLICADA AQUI
func (ar *AuctionRepository) closeAuctionsRoutine(ctx context.Context) {
	ticker := time.NewTicker(ar.auctionInterval)
	defer ticker.Stop() // Garante que o ticker seja parado para evitar vazamento de recursos

	// Utiliza for range, que é a forma idiomática para iterar sobre um ticker.
	for range ticker.C {
		activeAuctions, err := ar.FindAuctions(ctx, auction_entity.Active, "", "")
		if err != nil {
			logger.Error("Error finding active auctions in routine", err)
			continue
		}

		auctionDuration := getAuctionDuration()
		for _, auction := range activeAuctions {
			auctionEndTime := auction.Timestamp.Add(auctionDuration)
			if time.Now().After(auctionEndTime) {
				if err := ar.updateAuctionStatus(ctx, auction.Id, auction_entity.Completed); err != nil {
					logger.Error(fmt.Sprintf("Error closing auction %s", auction.Id), err)
				}
			}
		}
	}
}

func (ar *AuctionRepository) updateAuctionStatus(ctx context.Context, auctionId string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	filter := bson.M{"_id": auctionId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to update auction status for id = %s", auctionId), err)
		return internal_error.NewInternalServerError("Error trying to update auction status")
	}
	return nil
}

func getAuctionDuration() time.Duration {
	durationStr := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Duração padrão caso a variável de ambiente não seja definida ou seja inválida
		return 1 * time.Minute
	}
	return duration
}

func getAuctionInterval() time.Duration {
	intervalStr := os.Getenv("AUCTION_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		// Intervalo padrão de verificação
		return 10 * time.Second
	}
	return interval
}
