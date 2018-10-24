package api

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func generateRandomBucketName() string {
	return fmt.Sprintf("scaleway-test-%d-%s", time.Now().Unix(), fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))))
}

func TestScalewaAPI_ManageBucket(t *testing.T) {
	client, err := buildClient("ams1")
	if err != nil {
		t.Fatalf("Unable to create a client for AMS1: %v", err.Error())
	}

	desiredBucketName := generateRandomBucketName()
	bucket, err := client.CreateBucket(&CreateBucketRequest{
		Name:         desiredBucketName,
		Organization: client.Organization,
	})
	if err != nil {
		t.Fatalf("Expected request to succeed, but didn't: %v", err.Error())
	}

	defer func() {
		client.DeleteBucket(desiredBucketName)
	}()

	if bucket.Name != desiredBucketName {
		t.Errorf("Expected bucket name of %s, but got %s", desiredBucketName, bucket.Name)
	}
	if bucket.NumberOfObjects != 0 || bucket.Size != 0 {
		t.Errorf("Expected to create a fresh bucket, but got %d items, %d size", bucket.NumberOfObjects, bucket.Size)
	}

	_, err = client.ListObjects(bucket.Name)
	if err != nil {
		t.Fatalf("Expected request to succeed, but didn't: %v", err.Error())
	}
}

func TestScalewaAPI_ManageObject(t *testing.T) {
	client, err := buildClient("ams1")
	if err != nil {
		t.Fatalf("Unable to create a client for AMS1: %v", err.Error())
	}

	bucket, err := client.CreateBucket(&CreateBucketRequest{
		Name:         generateRandomBucketName(),
		Organization: client.Organization,
	})
	if err != nil {
		t.Fatalf("test requires bucket, but creation failed: %v", err.Error())
	}

	defer func() {
		client.DeleteBucket(bucket.Name)
	}()

	desiredName := "just-a-test"
	t.Run("Put object", func(t *testing.T) {
		tmpfile := bytes.NewBuffer([]byte("temporary file's content"))

		object, err := client.PutObject(&PutObjectRequest{
			BucketName: bucket.Name,
			ObjectName: desiredName,
		}, tmpfile)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if object == nil {
			t.Fatalf("Expected object, but got %v", object)
		}
		if object.Name != desiredName {
			t.Fatalf("Expected %s as name, but got %s", desiredName, object.Name)
		}
	})
	t.Run("Get object", func(t *testing.T) {
		object, err := client.GetObject(bucket.Name, desiredName)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if object == nil {
			t.Fatalf("Expected object, but got %v", object)
		}
		if object.Name != desiredName {
			t.Fatalf("Expected %s as name, but got %s", desiredName, object.Name)
		}
	})
	t.Run("List objects", func(t *testing.T) {
		containers, err := client.ListObjects(bucket.Name)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(containers) != 1 {
			t.Fatalf("Expected 1 object, got %d", len(containers))
		}
	})
	t.Run("Delete objects", func(t *testing.T) {
		err := client.DeleteObject(bucket.Name, desiredName)
		if err != nil {
			t.Fatalf(err.Error())
		}

		containers, err := client.ListObjects(bucket.Name)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if len(containers) != 0 {
			t.Fatalf("Expected 0 object, got %d", len(containers))
		}
	})
}
