package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	generate()
}

func generate() {
	var agency int
	var clusters int
	var siorc bool

	fmt.Println("Do you want a cluster (true or false)?")
	fmt.Scanln(&siorc)

	if siorc {

		fmt.Println("number of agency : ")
		_, err := fmt.Scanln(&agency)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("number of DB & coordinator : ")
		fmt.Scanln(&clusters)

		fmt.Println("You have created agencies : ", agency)
		fmt.Println("You have created DB & coordinator : ", clusters)
	}

	f, err := os.Create("docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if !siorc {
		f.WriteString("version: '3.7'\n")
		f.WriteString("services:\n")
		f.WriteString("  arangodb:\n")
		f.WriteString("    image: arangodb:latest\n")
		f.WriteString("    environment:\n")
		f.WriteString("      - ARANGO_RANDOM_ROOT_PASSWORD=1\n")
		f.WriteString("    ports: \n")
		f.WriteString("      - 8529:8529\n")
		f.WriteString("    volumes: \n")
		f.WriteString("      - arangodb_data_container:/var/lib/arangodb3\n")
		f.WriteString("      - arangodb_apps_data_container:/var/lib/arangodb3-apps\n")
		f.WriteString("volumes:\n")
		f.WriteString("  arangodb_data_container:\n")
		f.WriteString("  arangodb_apps_data_container:\n")
		return
	}

	f.WriteString("version: \"3.7\"\n")
	f.WriteString("services:\n")

	for i := 0; i < agency; i++ {
		f.WriteString("  arangodb-agency" + strconv.Itoa(i+1) + ":\n")
		f.WriteString("    image: arangodb/arangodb\n")
		f.WriteString("    environment:\n")
		f.WriteString("      - ARANGO_RANDOM_ROOT_PASSWORD=1\n")
		f.WriteString("    command: >\n")
		f.WriteString("      --server.endpoint tcp://0.0.0.0:8530\n")
		f.WriteString("      --server.jwt-secret secret\n")
		f.WriteString("      --server.authentication true\n")
		f.WriteString("      --server.authentication-system-only false\n")
		f.WriteString("      --server.statistics false\n")
		f.WriteString("      --foxx.queues false\n")
		f.WriteString("      --agency.size " + strconv.Itoa(agency) + "\n")
		f.WriteString("      --agency.supervision true\n")
		f.WriteString("      --agency.activate true\n")
		f.WriteString("      --agency.my-address tcp://arangodb-agency" + strconv.Itoa(i+1) + ":8530\n")
		f.WriteString("      --log.level info\n")
		for j := 0; j < agency; j++ {
			f.WriteString("      --agency.endpoint tcp://arangodb-agency" + strconv.Itoa(j+1) + ":8530\n")
		}
	}

	for i := 0; i < clusters; i++ {
		f.WriteString("  arangodb-coordinator" + strconv.Itoa(i+1) + ":\n")
		f.WriteString("    image: arangodb/arangodb\n")
		f.WriteString("    environment:\n")
		f.WriteString("      - ARANGO_RANDOM_ROOT_PASSWORD=1\n")
		f.WriteString("    command: >\n")
		f.WriteString("      --server.endpoint tcp://0.0.0.0:8529\n")
		f.WriteString("      --server.jwt-secret secret\n")
		f.WriteString("      --server.authentication true\n")
		f.WriteString("      --server.authentication-system-only false\n")
		f.WriteString("      --server.statistics true\n")
		f.WriteString("      --foxx.queues true\n")
		for j := 0; j < agency; j++ {
			f.WriteString("      --cluster.agency-endpoint tcp://arangodb-agency" + strconv.Itoa(j+1) + ":8530\n")
		}
		f.WriteString("      --cluster.my-address tcp://arangodb-coordinator" + strconv.Itoa(i+1) + ":8529\n")
		f.WriteString("      --cluster.my-local-info COORD" + strconv.Itoa(i+1) + "\n")
		f.WriteString("      --cluster.my-role COORDINATOR\n")
		f.WriteString("      --log.level info\n")
		if i != clusters-1 {
			f.WriteString("    ports: ['800" + strconv.Itoa(i) + ":8529']\n")
		}
	}

	for i := 0; i < clusters; i++ {
		f.WriteString("  arangodb-dbserver" + strconv.Itoa(i+1) + ":\n")
		f.WriteString("    image: arangodb/arangodb\n")
		f.WriteString("    environment:\n")
		f.WriteString("      - ARANGO_RANDOM_ROOT_PASSWORD=1\n")
		f.WriteString("    command: >\n")
		for j := 0; j < agency; j++ {
			f.WriteString("      --cluster.agency-endpoint tcp://arangodb-agency" + strconv.Itoa(j+1) + ":8530\n")
		}
		f.WriteString("      --cluster.my-address tcp://arangodb-dbserver" + strconv.Itoa(i+1) + ":8529\n")
		f.WriteString("      --server.endpoint tcp://0.0.0.0:8529\n")
		f.WriteString("      --server.jwt-secret secret\n")
		f.WriteString("      --server.authentication true\n")
		f.WriteString("      --server.authentication-system-only false\n")
		f.WriteString("      --server.statistics true\n")
		f.WriteString("      --foxx.queues false\n")
		f.WriteString("      --cluster.my-local-info ")
		if i == 0 {
			f.WriteString("DB1")
		} else {
			f.WriteString("HAVDB" + strconv.Itoa(i))
		}
		f.WriteString("\n")
		f.WriteString("      --cluster.my-role PRIMARY\n")
		f.WriteString("      --log.level info\n")

	}

	fmt.Println("done")

}
