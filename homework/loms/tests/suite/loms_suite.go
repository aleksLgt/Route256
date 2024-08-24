package lomssuite

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/db/orders"
	"route256/loms/internal/repository/db/stocks"
)

type ItemS struct {
	suite.Suite
	stocksStorage *stocks.Storage
	ordersStorage *orders.Storage
	ctx           context.Context
	conn          *pgx.Conn
}

func (s *ItemS) SetupSuite() {
	initEnv()

	ctx := context.Background()

	const dbConnEnv = "DB_CONN_TEST"
	dbConnStr := os.Getenv(dbConnEnv)
	conn, err := pgx.Connect(ctx, dbConnStr)

	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = ctx
	s.ordersStorage = orders.NewStorage(conn, conn)
	s.stocksStorage = stocks.NewStorage(conn, conn)
	s.conn = conn
}

func (s *ItemS) SetupTest() {
	// To fill in the database from the stock-data.json file
	err := stocks.FillStocks(stocks.New(s.conn))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ItemS) TearDownTest() {
	const query = `
	TRUNCATE TABLE orders, order_items, stocks;`

	_, err := s.conn.Exec(s.ctx, query)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ItemS) TestReserveStocksDB() {
	items := []domain.Item{{
		SKU:   872821,
		Count: 8,
	}}

	err := s.stocksStorage.Reserve(s.ctx, items)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestReserveRemoveStocksDB() {
	items := []domain.Item{{
		SKU:   872821,
		Count: 8,
	}}

	err := s.stocksStorage.ReserveRemove(s.ctx, items)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestReserveCancelStocksDB() {
	items := []domain.Item{{
		SKU:   872821,
		Count: 8,
	}}

	err := s.stocksStorage.ReserveCancel(s.ctx, items)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestGetStockBySKUDB() {
	var sku uint32 = 1076963

	_, err := s.stocksStorage.GetBySKU(s.ctx, sku)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestCreateOrderDB() {
	var userID int64 = 727

	items := []domain.Item{{
		SKU:   872821,
		Count: 8,
	}}

	_, err := s.ordersStorage.Create(s.ctx, userID, items)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestSetOrderStatusDB() {
	var orderID int64 = 231

	err := s.ordersStorage.SetStatus(s.ctx, orderID, orderStatus.AwaitingPayment)
	require.NoError(s.T(), err)
}

func (s *ItemS) TestGetOrderByIDDB() {
	var userID int64 = 727

	items := []domain.Item{{
		SKU:   872821,
		Count: 8,
	}}

	orderID, err := s.ordersStorage.Create(s.ctx, userID, items)
	require.NoError(s.T(), err)

	_, err = s.ordersStorage.GetByID(s.ctx, orderID)
	require.NoError(s.T(), err)
}

func initEnv() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
