package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// DeliveryCondition structure for delivery_conditions.json
type DeliveryCondition struct {
	ConditionID               int     `json:"condition_id"`
	MinimumFreeDeliveryAmount float64 `json:"minimum_free_delivery_amount"`
	DeliveryFee               float64 `json:"delivery_fee"`
}

// Supplier structure for suppliers.json
type Supplier struct {
	UserID      int `json:"user_id"`
	Name        string
	ConditionID int `json:"condition_id"`
}

// Product structure for products.json
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url"`
	GTIN       int64  `json:"gtin"`
	MinPrice   int    `json:"min_price"`
	SupplierID int    `json:"supplier_id"`
}

// ProductSupplier structure for products_supplier.json
type ProductSupplier struct {
	ProductID     int `json:"product_id"`
	SupplierID    int `json:"supplier_id"`
	Price         int `json:"price"`
	MinSellAmount int `json:"min_sell_amount"`
}

func generatePhoneNumber(userID int) string {
	// For example, we generate a phone number by prefixing with a fixed code
	// You can customize the logic as needed
	return fmt.Sprintf("+1234567%d", userID)
}

// Function to generate a hashed password (same for each user)
func generateHashedPassword() (string, error) {
	// The password we want to hash (same for every user)
	const password = "password1@"

	// Generate the bcrypt hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func generateOrderAmountByIndex(index int) int {
	orderAmounts := []int{30000, 50000, 70000, 100000, 150000}
	return orderAmounts[index%len(orderAmounts)]
}

func insertSupplier(db *sqlx.DB, suppliers []Supplier) {
	for i, supplier := range suppliers {
		phoneNumber := generatePhoneNumber(supplier.UserID)
		hashedPassword, err := generateHashedPassword()
		if err != nil {
			log.Fatalf("Error generating password: %v", err)
		}

		_, err = db.Exec(`INSERT INTO users (id, name, phone_number, hashed_password, role) 
							VALUES ($1, $2, $3, $4, 1)
							ON CONFLICT (id) DO UPDATE
							SET name = EXCLUDED.name,
								phone_number = EXCLUDED.phone_number,
								hashed_password = EXCLUDED.hashed_password,
								role = EXCLUDED.role`, supplier.UserID, supplier.Name, phoneNumber, hashedPassword)
		if err != nil {
			log.Fatalf("Error inserting user: %v", err)
		}

		orderAmount := generateOrderAmountByIndex(i)
		_, err = db.Exec(`INSERT INTO suppliers (user_id, condition_id, name ,order_amount) 
							VALUES ($1, $2, $3, $4)
							ON CONFLICT (user_id) DO UPDATE
							SET condition_id = EXCLUDED.condition_id,
								name = EXCLUDED.name,
								order_amount = EXCLUDED.order_amount`, supplier.UserID, supplier.ConditionID, supplier.Name, orderAmount)
		if err != nil {
			log.Fatalf("Error inserting supplier: %v", err)
		}
	}
}

func updateProductPrices(db *sqlx.DB) error {
	// Query to get the lowest price and the supplier_id for each product
	rows, err := db.Query(`
		SELECT ps.product_id, ps.price AS lowest_price, ps.supplier_id
		FROM products_supplier ps
		WHERE ps.price = (
			SELECT MIN(price)
			FROM products_supplier
			WHERE product_id = ps.product_id
		)
		GROUP BY ps.product_id, ps.supplier_id, ps.price;
	`)
	if err != nil {
		return fmt.Errorf("error querying lowest prices and supplier ids: %v", err)
	}
	defer rows.Close()

	// Prepare the update query for setting the lowest_price and lowest_supplier_id in the 'products' table
	updateQuery := `UPDATE products SET lowest_price = $1, lowest_supplier_id = $2 WHERE id = $3`

	// Iterate through the results and update the 'products' table
	for rows.Next() {
		var productID int
		var lowestPrice int
		var lowestSupplierID int

		// Scan the results
		if err := rows.Scan(&productID, &lowestPrice, &lowestSupplierID); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Update the 'lowest_price' and 'lowest_supplier_id' in the 'products' table
		_, err := db.Exec(updateQuery, lowestPrice, lowestSupplierID, productID)
		if err != nil {
			return fmt.Errorf("error updating product: %v", err)
		}
	}

	// Check for any error that might have occurred during row iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating through rows: %v", err)
	}

	return nil
}

func main() {
	// Open connection to your PostgreSQL database
	connStr := connectionString()
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Read JSON files
	// Delivery Conditions
	deliveryConditions := []DeliveryCondition{}
	loadJSON("delivery_conditions.json", &deliveryConditions)

	// Insert into delivery_conditions
	for _, condition := range deliveryConditions {
		_, err := db.Exec(`INSERT INTO delivery_conditions (condition_id, minimum_free_delivery_amount, delivery_fee)
							VALUES ($1, $2, $3)
							ON CONFLICT (condition_id) DO UPDATE
							SET minimum_free_delivery_amount = EXCLUDED.minimum_free_delivery_amount,
								delivery_fee = EXCLUDED.delivery_fee`, condition.ConditionID, condition.MinimumFreeDeliveryAmount, condition.DeliveryFee)
		if err != nil {
			log.Fatalf("Error inserting delivery condition: %v", err)
		}
	}

	// // Suppliers
	suppliers := []Supplier{}
	loadJSON("suppliers.json", &suppliers)
	fmt.Println("Suppliers")
	insertSupplier(db, suppliers)

	// Products
	products := []Product{}
	loadJSON("products.json", &products)

	fmt.Println(products)
	// Insert into products
	for _, product := range products {
		_, err := db.Exec(`INSERT INTO products (id, name, image_url, gtin)
							VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO UPDATE
							SET name = EXCLUDED.name,
								image_url = EXCLUDED.image_url,
								gtin = EXCLUDED.gtin`, product.ID, product.Name, product.ImageURL, product.GTIN)
		if err != nil {
			fmt.Println(product)
			log.Fatalf("Error inserting product: %v", err)
		}
	}

	// Products Supplier
	productSuppliers := []ProductSupplier{}
	loadJSON("products_supplier.json", &productSuppliers)
	fmt.Print("Products Supplier")
	// Insert into products_supplier
	for _, ps := range productSuppliers {
		_, err := db.Exec(`INSERT INTO products_supplier (product_id, supplier_id, price, sell_amount)
							VALUES ($1, $2, $3, $4)
							ON CONFLICT (product_id, supplier_id) DO UPDATE
							SET price = EXCLUDED.price,
								sell_amount = EXCLUDED.sell_amount`, ps.ProductID, ps.SupplierID, ps.Price, ps.MinSellAmount)
		if err != nil {
			fmt.Println(ps)
			log.Fatalf("Error inserting product supplier: %v", err)
		}
	}

	fmt.Println("updating products_supplier")
	err = updateProductPrices(db)
	if err != nil {
		log.Fatalf("Error updating product prices: %v", err)
	}

	syncSequence(db, "users", "id")
	syncSequence(db, "products", "id")
	syncSequence(db, "delivery_conditions", "condition_id")

	fmt.Println("Data migration completed successfully!")
}

// LoadJSON reads the JSON file and unmarshals it into the provided interface{}
func loadJSON(fileName string, data interface{}) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
}

func connectionString() string {
	host := getenv("PGHOST", "localhost")
	port := getenvInt("PGPORT", 5432)
	dbName := getenv("PGDATABASE", "catalog")
	user := getenv("PGUSER", "note-user")
	password := getenv("PGPASSWORD", "note-password")

	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host,
		port,
		dbName,
		user,
		password,
	)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid %s value: %v", key, err)
	}

	return intValue
}

func syncSequence(db *sqlx.DB, tableName, columnName string) {
	query := fmt.Sprintf(
		`SELECT setval(pg_get_serial_sequence('%s', '%s'), COALESCE(MAX(%s), 1), true) FROM %s`,
		tableName,
		columnName,
		columnName,
		tableName,
	)

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error syncing sequence for %s.%s: %v", tableName, columnName, err)
	}
}
