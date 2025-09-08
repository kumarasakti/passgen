package backends

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// R2Backend implements StorageBackend for Cloudflare R2
type R2Backend struct {
	client      *s3.Client
	bucket      string
	storePrefix string
	accountID   string
	endpoint    string
}

// R2Config holds the configuration for R2 backend
type R2Config struct {
	AccountID       string `yaml:"account_id" json:"account_id"`
	AccessKeyID     string `yaml:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key"`
	BucketName      string `yaml:"bucket_name" json:"bucket_name"`
	Region          string `yaml:"region" json:"region"`
}

// Establishes Cloudflare R2 cloud storage with AWS S3-compatible interface
func NewR2Backend(cfg R2Config, storePrefix string) (*R2Backend, error) {
	// Cloudflare R2 endpoint format
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	// Create AWS config for R2
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: cfg.Region,
				}, nil
			},
		)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &R2Backend{
		client:      client,
		bucket:      cfg.BucketName,
		storePrefix: storePrefix,
		accountID:   cfg.AccountID,
		endpoint:    endpoint,
	}, nil
}

// SaveFile saves data to R2 with the given key
func (r *R2Backend) SaveFile(key string, data []byte) error {
	fullKey := r.getFullKey(key)

	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:               &r.bucket,
		Key:                  &fullKey,
		Body:                 bytes.NewReader(data),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		ContentType:          aws.String("application/octet-stream"),
	})

	if err != nil {
		return fmt.Errorf("failed to save file to R2: %w", err)
	}

	return nil
}

// LoadFile loads data from R2 with the given key
func (r *R2Backend) LoadFile(key string) ([]byte, error) {
	fullKey := r.getFullKey(key)

	resp, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    &fullKey,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load file from R2: %w", err)
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return buf.Bytes(), nil
}

// ListFiles lists all files with the given prefix
func (r *R2Backend) ListFiles(prefix string) ([]string, error) {
	fullPrefix := r.getFullKey(prefix)

	resp, err := r.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &r.bucket,
		Prefix: &fullPrefix,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files from R2: %w", err)
	}

	var files []string
	for _, obj := range resp.Contents {
		// Remove the store prefix to get relative key
		key := strings.TrimPrefix(*obj.Key, r.storePrefix+"/")
		if key != "" {
			files = append(files, key)
		}
	}

	return files, nil
}

// DeleteFile deletes a file from R2
func (r *R2Backend) DeleteFile(key string) error {
	fullKey := r.getFullKey(key)

	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key:    &fullKey,
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %w", err)
	}

	return nil
}

// Verifies object presence in R2 bucket using HEAD operation
func (r *R2Backend) FileExists(key string) (bool, error) {
	fullKey := r.getFullKey(key)

	_, err := r.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &r.bucket,
		Key:    &fullKey,
	})

	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "NotFound" || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// Establishes R2 store structure with metadata object and proper configuration
func (r *R2Backend) Initialize(storeName string) error {
	// Create a metadata file to mark the store as initialized
	metadataKey := ".passgen-store.json"
	metadata := fmt.Sprintf(`{
  "name": "%s",
  "backend": "r2",
  "created_at": "%s",
  "version": "1.0"
}`, storeName, time.Now().Format(time.RFC3339))

	return r.SaveFile(metadataKey, []byte(metadata))
}

// Verifies R2 store has been properly initialized with required metadata object
func (r *R2Backend) IsInitialized(storeName string) (bool, error) {
	return r.FileExists(".passgen-store.json")
}

// Sync is a no-op for R2 since it's always in sync
func (r *R2Backend) Sync() error {
	// R2 is always synchronized, no action needed
	return nil
}

// Identifies R2 cloud storage type for backend selection and configuration
func (r *R2Backend) GetBackendType() string {
	return "r2"
}

// Provides R2 configuration details for debugging and connection diagnostics
func (r *R2Backend) GetConnectionInfo() map[string]string {
	return map[string]string{
		"type":       "cloudflare-r2",
		"bucket":     r.bucket,
		"account_id": r.accountID,
		"endpoint":   r.endpoint,
		"prefix":     r.storePrefix,
	}
}

// Constructs prefixed object key for organized storage within R2 bucket
func (r *R2Backend) getFullKey(key string) string {
	if r.storePrefix == "" {
		return key
	}
	return r.storePrefix + "/" + key
}

// TestConnection tests the R2 connection
func (r *R2Backend) TestConnection() error {
	// Try to list objects to test connection
	_, err := r.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &r.bucket,
		MaxKeys: aws.Int32(1),
	})

	if err != nil {
		return fmt.Errorf("R2 connection test failed: %w", err)
	}

	return nil
}
