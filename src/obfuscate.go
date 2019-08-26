package main

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"os"
)

func obfuscateUsernames(tx *sql.Tx) error {
	log.Info("[Obfuscation] Obfuscating usernames")

	_, err := tx.Exec(`UPDATE account
						SET name = 'imagemonkey-user-' || uuid_generate_v4()`)
	if err != nil {
		return err
	}

	return nil
}

func removeEmailAddresses(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing email addresses")

	_, err := tx.Exec(`UPDATE account
						SET email = null`)

	if err != nil {
		return err
	}

	return nil
}

func removeHashedPasswords(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing hashed passwords")

	_, err := tx.Exec(`UPDATE account
						SET hashed_password = null`)

	if err != nil {
		return err
	}

	return nil
}

func removeApiTokens(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing API tokens")

	_, err := tx.Exec(`DELETE FROM api_token`)

	if err != nil {
		return err
	}

	return nil
}

func removeAccessTokens(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing access tokens")

	_, err := tx.Exec(`DELETE FROM access_token`)

	if err != nil {
		return err
	}

	return nil
}

func removeUnverifiedDonations(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing unverified donations")
	
	_, err := tx.Exec(`DELETE  
							FROM imagehunt_task`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM image_validation v
							USING image i
							WHERE i.id = v.image_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM image_label_suggestion l
							USING image i
							WHERE i.id = l.image_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM user_image u
							USING image i
							WHERE i.id = u.image_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE
							FROM annotation_data d
							WHERE d.image_annotation_revision_id IN (
								SELECT image_annotation_revision_id 
                                 FROM image_annotation_revision r
                                 JOIN image_annotation a ON a.id = r.image_annotation_id
                                 JOIN image i ON i.id = a.image_id
                                 JOIN annotation_data d ON d.image_annotation_revision_id = r.id
                                 WHERE i.unlocked = false
							)`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM annotation_data d
							USING image_annotation a
							JOIN image i ON i.id = a.image_id
							WHERE a.id = d.image_annotation_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM user_image_annotation u
							USING image_annotation a
							JOIN image i ON i.id = a.image_id
							WHERE a.id = u.image_annotation_id AND i.unlocked = false`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM image_annotation_revision r
							WHERE r.id IN (
								SELECT r.id 
								FROM image_annotation_revision r
								JOIN image_annotation a ON r.image_annotation_id = a.id
								JOIN image i ON i.id = a.image_id
								WHERE i.unlocked = false
							)`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM image_annotation a
							USING image i
							WHERE i.id = a.image_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE  
							FROM image_description d
							USING image i
							WHERE i.id = d.image_id AND i.unlocked = false`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE
							FROM image_collection_image c
							USING image i
							WHERE i.id = c.image_id AND i.unlocked = false`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM image WHERE unlocked = false`)

	if err != nil {
		return err
	}

	return nil
}

func removeDonationsInQuarantine(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing donations in quarantine")

	_, err := tx.Exec(`DELETE FROM image_quarantine`)

	if err != nil {
		return err
	}

	return nil
}

func removeBlogSubscriptions(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing blog subscriptions")
	_, err := tx.Exec(`DELETE FROM blog.subscription`)

	return err
}

func removeLabelSuggestions(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing label suggestions")

	_, err := tx.Exec(`DELETE FROM image_label_suggestion`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM label_suggestion`)

	if err != nil {
		return err
	}

	return nil
}

func removeTrendingLabelSuggestions(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing trending label suggestions")

	_, err := tx.Exec(`DELETE FROM trending_label_suggestion`)

	if err != nil {
		return err
	}

	return nil
}

func removeTrendingLabelBotTasks(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing trending label bot tasks")

	_, err := tx.Exec(`DELETE FROM trending_label_bot_task`)

	if err != nil {
		return err
	}

	return nil
}

func removePendingImageDescriptions(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing pending image descriptions")

	_, err := tx.Exec(`DELETE FROM image_description WHERE state != 'unlocked'`)

	if err != nil {
		return err
	}

	return nil
}

func removeImageReports(tx *sql.Tx) error {
	log.Info("[Obfuscation] Removing image reports")

	_, err := tx.Exec(`DELETE FROM image_report`)

	if err != nil {
		return err
	}

	return nil
}

func obfuscateImageCollections(tx *sql.Tx) error {
	log.Info("[Obfuscation] Obfuscating image collections")

	_, err := tx.Exec(`UPDATE user_image_collection
						SET name = 'imagemonkey-collection-name-' || uuid_generate_v4(),
						description = 'imagemonkey-collection-description-' || uuid_generate_v4()`)
	if err != nil {
		return err
	}

	return nil
}

/*func changeMonkeyUserPassword(tx *sql.Tx) error {
	log.Info("[Obfuscation] Changing monkey password")

	_, err := tx.Exec(`ALTER ROLE monkey WITH PASSWORD 'imagemonkey'`)

	if err != nil {
		return err
	}
	return nil
}*/

func handleObfuscationError(tx *sql.Tx, err error) {
	log.Error("[Obfuscation] Couldn't obfuscate dataset: ", err.Error())
	err = tx.Rollback()
	if err != nil {
		log.Error("[Obfuscation] Couldn't rollback transaction: ", err.Error())
	}

	os.Exit(1)
}

func obfuscate(tx *sql.Tx) {
	if err := obfuscateUsernames(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeEmailAddresses(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeHashedPasswords(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeAccessTokens(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeApiTokens(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeDonationsInQuarantine(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeImageReports(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeTrendingLabelBotTasks(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeUnverifiedDonations(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeTrendingLabelSuggestions(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeLabelSuggestions(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removeBlogSubscriptions(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := removePendingImageDescriptions(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	if err := obfuscateImageCollections(tx); err != nil {
		handleObfuscationError(tx, err)
	}
	/*if err := changeMonkeyUserPassword(tx); err != nil {
		handleObfuscationError(tx, err)
	}*/
}
