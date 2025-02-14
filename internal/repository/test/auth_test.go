package repository

/*import "testing"

func TestInsertAndRetrieveUser(t *testing.T) {
	_, err := testDB.Exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = testDB.Exec(`INSERT INTO users (name) VALUES ($1)`, "Alice")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	var name string
	err = testDB.QueryRow(`SELECT name FROM users WHERE name = $1`, "Alice").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", name)
	}
}
*/
