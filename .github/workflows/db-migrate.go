const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "expenses"
)

type db struct {
	pool *pgxpool.Pool
}

func newDb() error {
	url := "postgresql://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		slog.Error("Could not connect to database -", "Error", err)
		return err
	}
	config.MaxConns = 4
	config.MinConns = 0
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Minute * 1

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Could create database pool-", "Error", err)
		return err
	}
	DB.pool = pool
	return nil
}

// Run migrations on test database

func main() {
	newDb()

	## Create expenses table
	_, err := DB.pool.Exec(context.Background(), "CREATE DATABASE expenses WITH  OWNER = postgres")
	if err != nil {
		fmt.Printf("Error creating expenses database", err)
		os.exit(1)
	}
	
	_, err := DB.pool.Exec(context.Background(), "CREATE TABLE expenses  (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		price NUMERIC,
		sku VARCHAR(255),
		dateadded timestamptz,
		lastupdate timestamptz)")
	if err != nil {
		fmt.Printf("Error creating expenses table", err)
		os.exit(1)
	}

	_, err := DB.pool.Exec(context.Background(), "CREATE TABLE budget  (
		id SERIAL PRIMARY KEY,
		budget NUMERIC,
		dateadded timestamptz,
		lastupdate timestamptz)")
	if err != nil {
		fmt.Printf("Error creating budget table", err)
		os.exit(1)
	}

}
