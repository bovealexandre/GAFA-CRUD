package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var DB driver.Database

// waitUntilServerAvailable keeps waiting until the server/cluster that the client is addressing is available.
func waitUntilServerAvailable(ctx context.Context, c driver.Client) bool {
	instanceUp := make(chan bool)
	go func() {
		for {
			verCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			if _, err := c.Version(verCtx); err == nil {
				cancel()
				instanceUp <- true
				return
			} else {
				cancel()
				time.Sleep(time.Second)
			}
		}
	}()
	select {
	case up := <-instanceUp:
		return up
	case <-ctx.Done():
		return false
	}
}

func ConnectDB() {
	// Create an HTTP connection to the database
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8000", "http://localhost:8001", "http://localhost:8002"},
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", ""),
	})
	if err != nil {
		log.Fatalf("Failed to get client info: %v", err)
	}

	waitUntilServerAvailable(context.Background(), c)

	// Ask the version of the server
	versionInfo, err := c.Version(nil)
	if err != nil {
		log.Fatalf("Failed to get version info: %v", err)
	}
	fmt.Printf("Database has version '%s' and license '%s'\n", versionInfo.Version, versionInfo.License)

	ctx := context.Background()
	db, err := c.Database(ctx, "housecms")
	if err != nil {
		// handle error
		log.Fatalf("Failed to get database info: %v", err)
	}
	fmt.Printf("housecms database exists: %v\n", db)

	DB = db
}
