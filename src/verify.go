package main

import (
	"database/sql"
	"github.com/mholt/archiver"
	log "github.com/Sirupsen/logrus"
	"errors"
	"strings"
	"os"
)

func verifyObfuscatedUsernames(tx *sql.Tx) error {
	log.Info("[Verification] Verify obfuscated usernames")
	rows, err := tx.Query(`SELECT name FROM account`)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(username, "imagemonkey-user-") {
			return errors.New("[Verification] Username not valid: ", )
		}
	}

	return nil
}

func verifyRemovedEmailAddresses(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed email addresses")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM account WHERE email is not null`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] Email addresses not valid!")
	}

	return nil
}

func verifyRemovedHashedPasswords(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed hashed passwords")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM account WHERE hashed_password is not null`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] Hashed passwords not valid!")
	}

	return nil
}

func verifyRemovedApiTokens(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed API tokens")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM api_token`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] API tokens not valid!")
	}

	return nil
}

func verifyRemovedAccessTokens(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed access tokens")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM access_token`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] access tokens not valid!")
	}

	return nil
}

func verifyRemovedUnverifiedDonations(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed unverified donations")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM image WHERE unlocked = false`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] unverified donations not valid!")
	}

	return nil
}

func verifyRemovedDonationsInQuarantine(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed donations in quarantine")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM image_quarantine`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] image quarantine not valid!")
	}

	return nil
}

func verifyRemovedLabelSuggestions(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed label suggestions")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM label_suggestion`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] label suggestions not valid!")
	}

	return nil
}

func verifyRemovedTrendingLabelSuggestions(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed trending label suggestions")
	var num int
	err := tx.QueryRow(`SELECT COUNT(*) FROM trending_label_suggestion`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] trending label suggestions not valid!")
	}

	return nil
}

func verifyRemovedBlogSubscriptions(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed blog subscriptions")
	var num int
	err := tx.QueryRow(`SELECT count(*) FROM blog.subscription`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] blog subscriptions not empty!")
	}

	return nil
}

func verifyRemovedPendingImageDescriptions(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed pending image descriptions")

	var num int
	err := tx.QueryRow(`SELECT count(*) FROM image_description WHERE state != 'unlocked'`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] image descriptions not empty!")
	}

	return nil
}

func verifyRemovedImageReports(tx *sql.Tx) error {
	log.Info("[Verification] Verify removed image reports")

	var num int
	err := tx.QueryRow(`SELECT count(*) FROM image_report`).Scan(&num)
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New("[Verification] image reports not empty!")
	}

	return nil
}

/*func verifyChangedMonkeyUserPassword(tx *sql.Tx) error {
	var currentPasswordHash string
	err := tx.QueryRow(`SELECT rolpassword FROM pg_authid WHERE rolname = 'monkey'`).Scan(&currentPasswordHash)
	if err != nil {
		return err
	}

	var expectedPasswordHash string
	err = tx.QueryRow(`SELECT md5('imagemonkey' || 'monkey')`).Scan(&expectedPasswordHash)

	if currentPasswordHash != expectedPasswordHash {
		return errors.New("[Verification] Passwords do not match!")
	}

	return nil
}*/

func removeArchive(path string) error {
	err := os.Remove(path)
	if err != nil {
		log.Error("[Verification] Couldn't remove archive: ", err.Error())
	}
	return err
}

func removeTempFiles(outputFolder string) error {
	p := outputFolder + "/" + "imagemonkey.sql"
	err := os.Remove(p)
	if err != nil {
		log.Error("[Verification] Couldn't remove temp file: ", err.Error())
		return err
	}

	/*p = outputFolder + "/" + "donations"
	err = os.Remove(p)
	if err != nil {
		log.Error("[Verification] Couldn't remove temp file: ", err.Error())
		return err
	}*/

	return nil
}

func extractArchive(archivePath, outputFolder string) error {
	return archiver.Zip.Open(archivePath, outputFolder)
}

func handleVerificationError(tx *sql.Tx, path string, err error, tempFilesFolder string) {
	log.Error("[Verification] Couldn't verify dataset: ", err.Error())
	err = tx.Rollback()
	if err != nil {
		//do not use Fatal() here, as we want to remove the archive later
		log.Error("[Verification] Couldn't rollback transaction: ", err.Error())
	}
	removeArchive(path)
	removeTempFiles(tempFilesFolder)
	os.Exit(1)
}


func verify(path string, outputFolder string) {
	extractedOutputFolder := outputFolder + "/verify"
	err := extractArchive(path, extractedOutputFolder)
	if err != nil {
		removeArchive(path)
		log.Fatal("[Verification] Couldn't extract archive: ", err.Error())
	}

	dbDumpPath := extractedOutputFolder + "/" + "imagemonkey.sql"
	err = loadDatabaseDump(dbDumpPath)
	if err != nil {
		log.Error("[Verification] Couldn't load database dump: ", err.Error())
		removeArchive(path)
		removeTempFiles(extractedOutputFolder)
	}


	tx, err := db.Begin()
	if err != nil {
		log.Error("[Verification] Couldn't start transaction: ", err.Error())
		removeArchive(path)
		removeTempFiles(extractedOutputFolder)
	}

	log.Info("[Verification] Starting verification")


	if err := verifyObfuscatedUsernames(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedEmailAddresses(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedHashedPasswords(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedApiTokens(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedAccessTokens(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedUnverifiedDonations(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedDonationsInQuarantine(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedLabelSuggestions(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedTrendingLabelSuggestions(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedBlogSubscriptions(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedPendingImageDescriptions(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	if err := verifyRemovedImageReports(tx); err != nil {
		handleVerificationError(tx, path, err, extractedOutputFolder)
	}
	

	/*if err := verifyChangedMonkeyUserPassword(tx); err != nil {
		handleVerificationError(tx, path, err)
	}*/

	err = tx.Commit()
	if err != nil {
		log.Error("[Verification] Couldn't commit transaction: ", err.Error())
		removeArchive(path)
		removeTempFiles(extractedOutputFolder)
	}

	log.Info("[Verification] Cleaning up temp files")
	removeTempFiles(extractedOutputFolder)
}