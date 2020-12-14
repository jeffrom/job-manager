package schema

import (
	"context"
)

func ValidateSchema(ctx context.Context, b []byte) error {
	if len(b) == 0 {
		return nil
	}

	keyErrs, err := SelfSchema.ValidateBytes(ctx, b)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}
