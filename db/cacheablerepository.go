package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/domain"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxOpenDbConns = 10
	maxDbLifetime  = 5 * time.Minute
	dbTimeout      = 1 * time.Minute

	ignorePriceValue = 250 * 1000000
)

type CacheableRepository struct {
	DbPool *pgxpool.Pool

	dataCenters *map[int]readertype.DataCenter

	worlds *map[int]readertype.World
}

func InitRepository(
	dsn string, worlds *map[int]readertype.World, dataCenters *map[int]readertype.DataCenter,
) (*CacheableRepository, error) {
	repo := &CacheableRepository{}

	// Connect then return repository
	err := repo.Connect(dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		return nil, err
	}

	repo.dataCenters = dataCenters
	repo.worlds = worlds

	return repo, nil
}

func (c *CacheableRepository) CreateListing(listing domain.Listing) (*domain.Listing, error) {
	if listing.PricePer > ignorePriceValue {
		return &listing, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO listings 
			(universalis_listing_id, 
			 item_id, region_id, 
			 data_center_id, world_id, 
			 price_per_unit, quantity, 
			 total_price, is_high_quality, 
			 retainer_name, retainer_city, 
			 last_review_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (universalis_listing_id) DO UPDATE SET
		    item_id = EXCLUDED.item_id,
			price_per_unit = EXCLUDED.price_per_unit,
			quantity = EXCLUDED.quantity,
			total_price = EXCLUDED.total_price,
			last_review_time = EXCLUDED.last_review_time
		RETURNING listing_id`

	serverInfo := (*c.worlds)[listing.WorldId]

	_, err := c.DbPool.Exec(
		ctx, query,
		listing.UniversalisId, listing.ItemId,
		serverInfo.RegionId, serverInfo.DataCenterId,
		listing.WorldId, listing.PricePer,
		listing.Quantity, listing.Total,
		listing.IsHighQuality, listing.RetainerName,
		listing.RetainerCity, listing.LastReview,
	)

	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (c *CacheableRepository) CreateListings(listings *[]domain.Listing) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := c.DbPool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := new(pgx.Batch)
	for _, listing := range *listings {
		// Don't bother uploading listing information if the price is ridiculous
		if listing.PricePer > ignorePriceValue {
			continue
		}

		query := `INSERT INTO listings 
			(universalis_listing_id, 
			 item_id, region_id, 
			 data_center_id, world_id, 
			 price_per_unit, quantity, 
			 total_price, is_high_quality, 
			 retainer_name, retainer_city, 
			 last_review_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (universalis_listing_id, data_center_id) DO UPDATE SET
		    item_id = EXCLUDED.item_id,
			price_per_unit = EXCLUDED.price_per_unit,
			quantity = EXCLUDED.quantity,
			total_price = EXCLUDED.total_price,
			last_review_time = EXCLUDED.last_review_time
		RETURNING listing_id`

		listingWorldRelation := (*c.worlds)[listing.WorldId]

		_ = batch.Queue(
			query,
			listing.UniversalisId, listing.ItemId,
			listingWorldRelation.RegionId, listingWorldRelation.DataCenterId,
			listing.WorldId, listing.PricePer,
			listing.Quantity, listing.Total,
			listing.IsHighQuality, listing.RetainerName,
			listing.RetainerCity, listing.LastReview,
		)
	}

	// Execute the query with all arguments
	result := tx.SendBatch(ctx, batch)
	defer func() {
		if e := result.Close(); e != nil {
			log.Printf("Failed to close batch: %s", e)
			err = e
		}

		if err != nil {
			err = tx.Rollback(ctx)

			if err != nil {
				log.Printf("Failed to rollback transaction: %s", err)
			}
		} else {
			if e := tx.Commit(ctx); e != nil {
				log.Printf("Failed to commit transaction: %s", e)
				err = e
			}
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to execute batch: %s", err)
	}

	return nil
}

func (c *CacheableRepository) DeleteListingByUniversalisId(listingUniversalisId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM listings WHERE universalis_listing_id = $1`

	_, err := c.DbPool.Exec(
		ctx, query, listingUniversalisId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *CacheableRepository) DeleteListings(universalisListingIds []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := c.DbPool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := new(pgx.Batch)
	for _, listingId := range universalisListingIds {
		query := `DELETE FROM listings WHERE universalis_listing_id = $1`

		_ = batch.Queue(
			query,
			listingId,
		)
	}

	// Execute the query with all arguments
	result := tx.SendBatch(ctx, batch)
	defer func() {
		if e := result.Close(); e != nil {
			log.Printf("Failed to close batch: %s", e)
			err = e
		}

		if err != nil {
			err = tx.Rollback(ctx)

			if err != nil {
				log.Printf("Failed to rollback transaction: %s", err)
			}
		} else {
			if e := tx.Commit(ctx); e != nil {
				log.Printf("Failed to commit transaction: %s", e)
				err = e
			}
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to execute batch: %s", err)
	}

	return nil
}

func (c *CacheableRepository) CreateSale(sale domain.Sale) (*domain.Sale, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO sales
			(item_id, world_id, 
			 price_per_unit, quantity, 
			 total_price, is_high_quality, 
			 buyer_name, sale_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (sales_id, sale_time) DO NOTHING
		RETURNING sales_id
	`

	returnedRow := c.DbPool.QueryRow(
		ctx, query,
		sale.ItemId, sale.WorldId,
		sale.PricePer, sale.Quantity,
		sale.TotalPrice, sale.IsHighQuality,
		sale.BuyerName, sale.Timestamp,
	)

	err := returnedRow.Scan(&sale.Id)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (c *CacheableRepository) CreateSales(sales *[]domain.Sale) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := c.DbPool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := new(pgx.Batch)
	for _, listing := range *sales {
		query := `INSERT INTO sales 
			(item_id, world_id, 
			 price_per_unit, quantity, 
			 total_price, is_high_quality, 
			 buyer_name, sale_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (sales_id, sale_time) DO NOTHING
		RETURNING sales_id`

		_ = batch.Queue(
			query,
			listing.ItemId, listing.WorldId,
			listing.PricePer, listing.Quantity,
			listing.TotalPrice, listing.IsHighQuality,
			listing.BuyerName, listing.Timestamp,
		)
	}

	// Execute the query with all arguments
	result := tx.SendBatch(ctx, batch)
	defer func() {
		if e := result.Close(); e != nil {
			log.Printf("Failed to close batch: %s", e)
			err = e
		}

		if err != nil {
			err = tx.Rollback(ctx)

			if err != nil {
				log.Printf("Failed to rollback transaction: %s", err)
			}
		} else {
			if e := tx.Commit(ctx); e != nil {
				log.Printf("Failed to commit transaction: %s", e)
				err = e
			}
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to execute batch: %s", err)
	}

	return nil
}

func (c *CacheableRepository) DeleteSaleById(saleId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM sales WHERE sales_id = $1`

	_, err := c.DbPool.Exec(
		ctx, query, saleId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *CacheableRepository) Connect(connectionInfo string) error {
	pgxConfig, err := pgxpool.ParseConfig(connectionInfo)

	if err != nil {
		return err
	}

	pgxConfig.MinConns = 0
	pgxConfig.MaxConns = maxOpenDbConns
	pgxConfig.MaxConnLifetime = maxDbLifetime

	d, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)

	if err != nil {
		return err
	}

	err = testConnection(d)
	if err != nil {
		return err
	}

	c.DbPool = d
	return nil
}

func (c *CacheableRepository) createConstraint(
	batch *pgx.Batch, tableName string, columnName string, constraintInfo string,
) {
	// We can only drop constraints if they already exist, but can't skip creating them if they exist
	// Dropping then re-creating in case we accidentally run the setup process twice.
	dropExistingConstraint := fmt.Sprintf(
		`ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s_%s_index`,
		tableName,
		tableName,
		columnName,
	)
	_ = batch.Queue(dropExistingConstraint)

	uniqueIndexQuery := fmt.Sprintf(
		`ALTER TABLE %s ADD CONSTRAINT %s_%s_index %s`,
		tableName,
		tableName,
		columnName,
		constraintInfo,
	)
	_ = batch.Queue(uniqueIndexQuery)
}

func (c *CacheableRepository) CreatePartitions() error {
	// Create a batch
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := c.DbPool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := new(pgx.Batch)
	year := time.Now().Year()

	// Create partitions for last and current year
	c.createSalePartionsForYear(year, batch)

	// Create partitions for listings by data center
	c.createListingPartitionsForDc(batch)

	// Execute the query with all arguments
	result := tx.SendBatch(ctx, batch)
	defer func() {
		if e := result.Close(); e != nil {
			log.Printf("Failed to close batch: %s", e)
			err = e
		}

		if err != nil {
			err = tx.Rollback(ctx)

			if err != nil {
				log.Printf("Failed to rollback transaction: %s", err)
			}
		} else {
			if e := tx.Commit(ctx); e != nil {
				log.Printf("Failed to commit transaction: %s", e)
				err = e
			}
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to execute batch: %s", err)
	}

	return nil
}

func (c *CacheableRepository) createListingPartitionsForDc(
	batch *pgx.Batch,
) {
	for dcId, dataCenter := range *c.dataCenters {
		partitionName := fmt.Sprintf("listings_%s", dataCenter.Name)

		partitionQuery := fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s PARTITION OF listings FOR VALUES IN (%d)`,
			partitionName,
			dcId,
		)
		_ = batch.Queue(partitionQuery)

		// Create check to enforce partition
		dataCenterCheck := fmt.Sprintf(
			`CHECK (data_center_id = %d)`,
			dcId,
		)

		c.createConstraint(batch, partitionName, "data_center_id", dataCenterCheck)

		// Create unique universalis_listing_id index
		// We can only drop constraints if they already exist, but can't skip creating them if they exist
		// Dropping then re-creating in case we accidentally run the setup process twice.
		dropExistingConstraint := fmt.Sprintf(
			`ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s_universalis_listing_id_index`,
			partitionName,
			partitionName,
		)
		_ = batch.Queue(dropExistingConstraint)

		uniqueIndexQuery := fmt.Sprintf(
			`ALTER TABLE %s ADD CONSTRAINT %s_universalis_listing_id_index UNIQUE (universalis_listing_id)`,
			partitionName,
			partitionName,
		)
		_ = batch.Queue(uniqueIndexQuery)

		c.createConstraint(batch, partitionName, "sale_time", dataCenterCheck)

		// Create world, datacenter and region indexes on the partition
		worldIdIndex := fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS %s_world_index
    on %s (item_id desc, world_id asc) INCLUDE (total_price, quantity, price_per_unit)`, partitionName, partitionName,
		)
		dataCenterIndex := fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS %s_data_center_index
    on %s (item_id desc, data_center_id asc) INCLUDE (total_price, quantity, price_per_unit)`,
			partitionName,
			partitionName,
		)
		regionIndex := fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS %s_region_index
    on %s (item_id desc, region_id asc) INCLUDE (total_price, quantity, price_per_unit)`, partitionName, partitionName,
		)

		_ = batch.Queue(worldIdIndex)
		_ = batch.Queue(dataCenterIndex)
		_ = batch.Queue(regionIndex)
	}
}

func (c *CacheableRepository) createSalePartionsForYear(year int, batch *pgx.Batch) {
	for month := 1; month <= 12; month++ {
		partitionName := fmt.Sprintf("sales_%d_%d", year, month)

		// Determine the start and end date for the month.
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		partitionQuery :=
			fmt.Sprintf(
				`CREATE TABLE IF NOT EXISTS %s PARTITION OF sales FOR VALUES FROM ('%s') TO ('%s')`,
				partitionName,
				startDate.Format(time.RFC3339),
				endDate.Format(time.RFC3339),
			)
		_ = batch.Queue(partitionQuery)

		// Create a check that enforces the values of the partition
		timeConstraintCheck := fmt.Sprintf(
			`CHECK (sale_time >= '%s' AND sale_time < '%s')`,
			startDate.Format(time.RFC3339),
			endDate.Format(time.RFC3339),
		)

		c.createConstraint(batch, partitionName, "sale_time", timeConstraintCheck)

		// Also queue the creation of indexes for this table
		indexQuery := fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS %s_index ON %s 
(item_id desc, world_id asc) INCLUDE (total_price, price_per_unit, quantity, sale_time)`, partitionName, partitionName,
		)
		_ = batch.Queue(indexQuery)

	}
}

func testConnection(pool *pgxpool.Pool) error {
	err := pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Pinged successfully.")
	return nil
}

func (c *CacheableRepository) GetListingsByItemAndWorldId(itemId, worldId int) (*[]*domain.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM listings WHERE item_id = $1 AND world_id = $2 ORDER BY total_price LIMIT 100`

	rows, err := c.DbPool.Query(ctx, query, itemId, worldId)
	if err != nil {
		return nil, err
	}

	listings, err := extractListings(rows, err)
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func (c *CacheableRepository) GetListingsForItemOnDataCenter(itemId, dataCenterId int) (*[]*domain.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM listings WHERE item_id = $1 AND data_center_id = $2 ORDER BY total_price LIMIT 100`

	rows, err := c.DbPool.Query(ctx, query, itemId, dataCenterId)
	if err != nil {
		return nil, err
	}

	listings, err := extractListings(rows, err)
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func (c *CacheableRepository) GetSalesByItemAndWorldId(itemId, worldId int) (*[]*domain.Sale, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM sales WHERE item_id = $1 AND world_id = $2 ORDER BY sale_time DESC LIMIT 100`

	rows, err := c.DbPool.Query(ctx, query, itemId, worldId)
	if err != nil {
		return nil, err
	}

	sales, err := extractSales(rows, err)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (c *CacheableRepository) GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]*domain.Sale, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	worldsOnDc := getWorldsOnDc(c, dataCenterId)

	query := `SELECT * FROM sales WHERE item_id = $1 AND world_id = ANY($2) ORDER BY sale_time DESC LIMIT 100`

	rows, err := c.DbPool.Query(ctx, query, itemId, worldsOnDc)
	if err != nil {
		return nil, err
	}

	sales, err := extractSales(rows, err)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func getWorldsOnDc(c *CacheableRepository, dataCenterId int) *pgtype.Array[int] {
	// The smallest data center currently has 4 worlds on it, and the largest has 8
	worldsOnDc := make([]int, 4, 8)
	for _, world := range *c.worlds {
		if len(worldsOnDc) == cap(worldsOnDc) {
			break
		}

		if world.DataCenterId == dataCenterId {
			worldsOnDc = append(worldsOnDc, world.Id)
		}
	}

	return &pgtype.Array[int]{
		Elements: worldsOnDc,
		Dims: []pgtype.ArrayDimension{
			{Length: int32(len(worldsOnDc)), LowerBound: 1},
		},
		Valid: true,
	}
}

func extractListings(rows pgx.Rows, err error) (*[]*domain.Listing, error) {
	var listings []*domain.Listing
	for rows.Next() {
		var listing domain.Listing
		err = rows.Scan(
			&listing.Id, &listing.UniversalisId,
			&listing.ItemId, &listing.RegionId,
			&listing.DataCenterId, &listing.WorldId,
			&listing.PricePer, &listing.Quantity,
			&listing.Total, &listing.IsHighQuality,
			&listing.RetainerName, &listing.RetainerCity,
			&listing.LastReview,
		)
		if err != nil {
			return nil, err
		}

		listings = append(listings, &listing)
	}

	return &listings, nil
}

func extractSales(rows pgx.Rows, err error) (*[]*domain.Sale, error) {
	var sales []*domain.Sale
	for rows.Next() {
		var sale domain.Sale
		err = rows.Scan(
			&sale.Id, &sale.ItemId,
			&sale.WorldId, &sale.PricePer,
			&sale.Quantity, &sale.TotalPrice,
			&sale.IsHighQuality, &sale.BuyerName,
			&sale.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		sales = append(sales, &sale)
	}

	return &sales, nil
}
