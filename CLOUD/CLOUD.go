package CLOUD

import (
	"io"
	"github.com/Esseh/notorious-dev/CORE"
	"google.golang.org/appengine/log"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// Adds a file to the cloud.
func AddFile(ctx context.Context, filename, contentType string, freader io.Reader) error {
	client, clientErr := storage.NewClient(ctx)
	log.Infof(ctx, "storage.newclient error: ", clientErr)
	if clientErr != nil {
		return clientErr
	}
	defer client.Close()
	csWriter := client.Bucket(CORE.GCSBucket).Object(filename).NewWriter(ctx)
	// Cloud Storage Writer - Permissions
	csWriter.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}
	csWriter.CacheControl = "max-age=300"
	csWriter.ContentType = contentType
	if _, copyErr := io.Copy(csWriter, freader); copyErr != nil {
		csWriter.Close()
		log.Infof(ctx, "io.copy error: ", copyErr)
		return copyErr
	}
	return csWriter.Close()
}

// Removes a file from the cloud.
func RemoveFile(ctx context.Context, filename string) error {
	client, clientErr := storage.NewClient(ctx)
	if clientErr != nil { return clientErr }
	defer client.Close()
	return client.Bucket(CORE.GCSBucket).Object(filename).Delete(ctx)
}
