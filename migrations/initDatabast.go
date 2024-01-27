package main

func main() {
	Connect()
	defer Disconnect()

	DB.Exec(
		`CREATE TABLE Deliveries (
			id serial PRIMARY KEY, 
			name VARCHAR(191) NOT NULL,
			phone VARCHAR(191) NOT NULL,
			zip VARCHAR(191) NOT NULL,
			city VARCHAR(191) NOT NULL,
			address VARCHAR(191) NOT NULL,
			region VARCHAR(191) NOT NULL,
			email VARCHAR(191) NOT NULL
			);`)

	DB.Exec(
		`CREATE TABLE Products (
			id serial PRIMARY KEY, 
			chrt_id INTEGER NOT NULL,
			track_number VARCHAR(191) NOT NULL,
			price FLOAT NOT NULL,
			rid VARCHAR(191) NOT NULL,
			name VARCHAR(191) NOT NULL,
			sale INTEGER NOT NULL,
			size VARCHAR(191) NOT NULL,
			total_price FLOAT NOT NULL,
			nm_id INTEGER NOT NULL,
			brand VARCHAR(191) NOT NULL,
			status INTEGER NOT NULL
			);`)

	DB.Exec(
		`CREATE TABLE Payments
			(id serial PRIMARY KEY, 
			transaction VARCHAR(191) NOT NULL, 
			request_id VARCHAR(191) NOT NULL, 
			currency VARCHAR(191) NOT NULL, 
			provider VARCHAR(191) NOT NULL, 
			amount FLOAT NOT NULL, 
			payment_dt INTEGER NOT NULL, 
			bank VARCHAR(191) NOT NULL, 
			delivery_cost INTEGER NOT NULL, 
			goods_total INTEGER NOT NULL, 
			custom_fee INTEGER NOT NULL
			);`)

	DB.Exec(
		`CREATE TABLE Orders (
			id serial PRIMARY KEY, 
			order_uid VARCHAR(191) NOT NULL,
			track_number VARCHAR(191) NOT NULL,
			entry VARCHAR(191) NOT NULL,
			locale VARCHAR(191) NOT NULL,
			internal_signature VARCHAR(191) NOT NULL,
			customer_id VARCHAR(191) NOT NULL,
			delivery_service VARCHAR(191) NOT NULL,
			shardkey VARCHAR(191) NOT NULL,
			sm_id INTEGER NOT NULL,
			date_created VARCHAR(191) NOT NULL,
			delivery_id INTEGER NOT NULL,
			payment_id INTEGER NOT NULL,
			cof_shard VARCHAR(191) NOT NULL);`)

	DB.Exec(`
		CREATE TABLE OrderProducts (
			order_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL
		);`)

	DB.Exec(`
		CREATE TABLE cache_state (
			id serial PRIMARY KEY,
			order_id INTEGER NOT NULL
		);`)
}
