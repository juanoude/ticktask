package sync

import (
	"ticktask/persistence"
	"ticktask/persistence/sync"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	SyncCmd.AddCommand(configCmd)
}

// configCmd sets up AWS credentials and S3 bucket information.
// Prompts for: region, bucket name, access key ID, and secret access key.
// Non-sensitive values (region, bucket) are stored in the local database.
// Sensitive values (credentials) are stored in the system keyring.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configure your s3 credentials and info",
	Long:  "here is the configuration for credentials and bucket info so we have a bucket to store your db",
	Run: func(cmd *cobra.Command, args []string) {
		region, wasCancelled := views.RunInput("Please insert your aws region: ", false)
		if wasCancelled {
			return
		}

		bucketName, wasCancelled := views.RunInput("Please insert your bucket name: ", false)
		if wasCancelled {
			return
		}

		accessKeyId, wasCancelled := views.RunInput("Please insert your aws access key id: ", true)
		if wasCancelled {
			return
		}

		secretAccessKey, wasCancelled := views.RunInput("Please insert your aws secretAccessKey: ", true)
		if wasCancelled {
			return
		}

		db := persistence.GetDB()
		db.StoreConfig(persistence.AWSRegionConfig, region)
		db.StoreConfig(persistence.AWSBucketNameConfig, bucketName)

		wallet := persistence.GetWallet()
		wallet.StoreKey(sync.AWSAccessKeyID, accessKeyId)
		wallet.StoreKey(sync.AWSSecretAccessKey, secretAccessKey)
	},
}
