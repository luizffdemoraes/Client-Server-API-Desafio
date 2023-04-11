package config

import (
	"context"
	"fmt"
	"time"

	"github.com/luizffdemoraes/desafio/client-server-api/schemas"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() error {
	var err error

	// Inicialize SQLite
	db, err = InitializeSQLite()

	if err != nil {
		return fmt.Errorf("error initialize sqlite: %v", err)
	}

	return nil
}

func GetSQLite() *gorm.DB {
	return db
}

func PersistDataBase(exchange *schemas.UsdBrls) error {

	// Defining duration
	// of Nanoseconds method
	nano, _ := time.ParseDuration("10s")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(nano.Nanoseconds()))
	defer cancel()

	db.AutoMigrate(&schemas.UsdBrl{})

	// CREATE
	return db.WithContext(ctx).Create(&schemas.UsdBrl{
		Code:       exchange.USDBRL.Code,
		Codein:     exchange.USDBRL.Codein,
		Name:       exchange.USDBRL.Name,
		High:       exchange.USDBRL.High,
		Low:        exchange.USDBRL.Low,
		VarBid:     exchange.USDBRL.VarBid,
		PctChange:  exchange.USDBRL.PctChange,
		Bid:        exchange.USDBRL.Bid,
		Ask:        exchange.USDBRL.Ask,
		Timestamp:  exchange.USDBRL.Timestamp,
		CreateDate: exchange.USDBRL.CreateDate,
	}).Error
}
