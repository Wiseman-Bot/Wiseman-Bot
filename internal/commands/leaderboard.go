package commands

import (
	"context"
	"fmt"
	"time"
	"wiseman/internal/db"
	"wiseman/internal/discord"
	"wiseman/internal/shared"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	Helpers = append(Helpers, Helper{
		Name:        "leaderboard",
		Category:    "This is a category",
		Description: "This is a descriptio",
		Usage:       "This is a usage",
	})

	discord.Commands["leaderboard"] = Leaderboard
}

func Leaderboard(s *discordgo.Session, m *discordgo.MessageCreate, client *mongo.Client, args []string) error {

	ctx := context.TODO()
	collection := client.Database(shared.DB_NAME).Collection(shared.USERS_INFIX)

	findOptions := options.Find()
	// Sort by `rank` field descending
	findOptions.SetSort(bson.D{primitive.E{Key: "rank", Value: -1}})
	// Limit by 10 documents only
	findOptions.SetLimit(10)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil
	}

	var leaderboard db.UserType
	var fields []*discordgo.MessageEmbedField

	for cursor.Next(ctx) {
		err := cursor.Decode(&leaderboard)
		if err != nil {
			return nil
		} else {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  string(leaderboard.Username),
				Value: fmt.Sprint(leaderboard.Rank),
			})

		}
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       9004799, // Green
		Description: "top 10 active users.",
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:       "Leaderboard",
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return nil
}
