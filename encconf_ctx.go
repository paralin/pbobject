package pbobject

import (
	"context"
)

var encryptionConfigCtxKey = &(struct{ encryptionConfigCtxKey string }{})

// WithEncryptionConf attaches a encryption config to the context.
func WithEncryptionConf(ctx context.Context, encryptionConfig *EncryptionConfig) context.Context {
	return context.WithValue(ctx, encryptionConfigCtxKey, encryptionConfig)
}

// GetEncryptionConf attempts to return the encryption config from the context.
func GetEncryptionConf(ctx context.Context) *EncryptionConfig {
	n := ctx.Value(encryptionConfigCtxKey)
	if n != nil {
		return n.(*EncryptionConfig)
	}
	return nil
}
