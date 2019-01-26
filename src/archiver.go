package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mholt/archiver"
	_"github.com/lib/pq"
	"flag"
	"os"
	"bytes"
	"database/sql"
	"os/exec"
	"fmt"
)

var db *sql.DB


func loadDatabaseDump(path string) error {

	log.Info("[Main] Dropping existing database + create new one")
	//connect as user postgres, in order to drop + re-create database imagemonkey
	localDb, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		return err
	}

	defer localDb.Close()

	//terminate any open database connections (we need to do this, before we can drop the database)
	_, err = localDb.Exec(`SELECT pg_terminate_backend(pid)
					  FROM pg_stat_activity
					  WHERE datname = 'imagemonkey'`)
	if err != nil {
		return err
	}

	_, err = localDb.Exec("DROP DATABASE IF EXISTS imagemonkey")
	if err != nil {
		return err
	}

	_, err = localDb.Exec("CREATE DATABASE imagemonkey OWNER monkey")
	if err != nil {
		return err
	}

	log.Info("[Main] Loading database dump from ", path)
	var out, stderr bytes.Buffer

	//load schema
	cmd := exec.Command("psql", "-f", path, "-d", "imagemonkey", "-U", "postgres")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
	    fmt.Sprintf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	    return err
	}

	return nil
}


func dumpDatabase(toPath string) error {
	var out, stderr bytes.Buffer

	log.Info("[Main] Dumping database to ", toPath)

	cmd := exec.Command("pg_dump", "-U", "postgres", "imagemonkey")
	// open the out file for writing
    outfile, err := os.Create(toPath)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    cmd.Stdout = outfile

    err = cmd.Start() 
    if err != nil {
    	fmt.Sprintf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
        return err
    }
    cmd.Wait()

    return nil
}


func createArchive(donationsDir string, databaseDumpPath string, destination string) error {
	log.Info("[Main] Creating archive: ", destination)
	err := archiver.Archive([]string{databaseDumpPath, donationsDir}, destination)
	return err
}


func main() {
	//dbPassword := flag.String("dbpasswd", "", "database password")
	dbName := flag.String("dbname", "imagemonkey", "database name")
	dbUser := flag.String("dbuser", "monkey", "database user")
	dryRun := flag.Bool("dryrun", true, "Dry run")
	verifyOutput := flag.Bool("verify", true, "Verify that the generated archive is correct")
	donationsDir := flag.String("donationsdir", "", "Path to the donations directory")
	dbDump := flag.String("dbdump", "", "Path to the database dump")
	outputFolder := flag.String("output", "", "Output folder")

	flag.Parse()

	/*if *dbPassword == "" {
		log.Fatal("[Main] Please provide a valid database password!")
	}*/

	if *donationsDir == "" {
		log.Fatal("[Main] Please provide a path to the donations directory!")
	} else {
		if _, err := os.Stat(*donationsDir); os.IsNotExist(err) {
			log.Fatal("[Main] Please provide a VALID path to the donations directory!")
		}
	}

	if *dbDump == "" {
		log.Fatal("[Main] Please provide a path to the database dump!")
	} else {
		if _, err := os.Stat(*donationsDir); os.IsNotExist(err) {
			log.Fatal("[Main] Please provide a VALID path to the database dump!")
		}
	}

	if *outputFolder == "" {
		log.Fatal("[Main] Please provide a valid output directory!")
	}

	err := loadDatabaseDump(*dbDump)
	if err != nil {
		log.Fatal("[Main] Couldn't load database dump: ", err.Error())
	}

	//dbConnectionString := "host=127.0.0.1 user=" + *dbUser + " dbname=" + *dbName + " password=" + *dbPassword + " sslmode=disable"
	dbConnectionString := "host=127.0.0.1 user=" + *dbUser + " dbname=" + *dbName + " sslmode=disable"

	log.Info("[Main] Archiver started")

	//open database and make sure that we can ping it
	db, err = sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatal("[Main] Couldn't open database: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("[Main] Couldn't ping database: ", err.Error())
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("[Main] Couldn't start transaction: ", err.Error())
	}

	
	obfuscate(tx)

	if *dryRun {
		err = tx.Rollback()
		if err != nil {
			log.Fatal("[Main] Couldn't rollback transaction: ", err.Error())
		}
	} else {
		err = tx.Commit()
		if err != nil {
			log.Fatal("[Main] Couldn't commit transaction: ", err.Error())
		}

		//after we've committed transaction, create database dump + archive it
		//together with the donations
		obfuscatedDbDumpPath := *outputFolder + "/" + "imagemonkey.sql"
		err = dumpDatabase(obfuscatedDbDumpPath)
		if err != nil {
			log.Fatal("[Main] Couldn't create database dump: ", err.Error())
		}

		archiveOutputPath := *outputFolder + "/" + "imagemonkey.zip"
		err = createArchive(*donationsDir, obfuscatedDbDumpPath, archiveOutputPath)
		if err != nil {
			log.Fatal("[Main] Couldn't create archive: ", err.Error())
		}

		//after changes are comitted, run verification
		if *verifyOutput {
			verify(archiveOutputPath, *outputFolder)
		}

		log.Info("[Main] All good!")
	}


}