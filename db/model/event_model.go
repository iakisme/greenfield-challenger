package model

import (
	"gorm.io/gorm"
)

type EventStartChallenge struct {
	SpOperatorAddress string `gorm:"NOT NULL"`
	BucketName        string `gorm:"NOT NULL"`
	ObjectName        string `gorm:"NOT NULL"`
	ObjectId          uint64 `gorm:"NOT NULL"`
	Segment           int    `gorm:"NOT NULL"`
	Piece             int    `gorm:"NOT NULL"`

	ChallengeId uint64 `gorm:"NOT NULL"`

	//ChannelId uint8  `gorm:"NOT NULL;index:idx_inscription_relay_transaction_channel_seq_status"`
	//Height    uint64 `gorm:"NOT NULL;index:idx_inscription_relay_transaction_height"`
	//PayLoad    string `gorm:"type:text"`
	//Status      db.TxStatus `gorm:"NOT NULL;index:idx_inscription_relay_transaction_status"`
}

func (*EventStartChallenge) TableName() string {
	return "event_start_challenge"
}

func InitEventTables(db *gorm.DB) {
	if !db.Migrator().HasTable(&EventStartChallenge{}) {
		err := db.Migrator().CreateTable(&EventStartChallenge{})
		if err != nil {
			panic(err)
		}
	}
}
