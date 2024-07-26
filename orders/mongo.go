package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap/zapcore"
)


type MongoCore struct {
    collection *mongo.Collection
    level      zapcore.Level
}

func NewMongoCore(uri, dbName, collectionName string, level zapcore.Level) (*MongoCore, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    collection := client.Database(dbName).Collection(collectionName)
    return &MongoCore{collection: collection, level: level}, nil
}

func (m *MongoCore) Enabled(level zapcore.Level) bool {
    return level >= m.level
}

func (m *MongoCore) With(fields []zapcore.Field) zapcore.Core {
    return &MongoCore{
        collection: m.collection,
        level:      m.level,
    }
}

func (m *MongoCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
    if m.Enabled(entry.Level) {
        return checkedEntry.AddCore(entry, m)
    }
    return checkedEntry
}

func (m *MongoCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
    logEntry := bson.M{
        "level":     entry.Level.String(),
        "message":   entry.Message,
        "timestamp": entry.Time,
        "caller":    entry.Caller.TrimmedPath(),
    }

    for _, field := range fields {
        logEntry[field.Key] = field.Interface
    }

    _, err := m.collection.InsertOne(context.Background(), logEntry)
    return err
}

func (m *MongoCore) Sync() error {
    return nil
}


