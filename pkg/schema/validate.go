package schema

import (
	"context"

	"github.com/hashicorp/go-multierror"
)

func ValidateSchema(ctx context.Context, argData, dataData, resultData []byte) error {
	merr := &multierror.Error{}
	if err := ValidateArgSchema(ctx, argData); err != nil {
		merr = multierror.Append(merr, err)
	}
	if err := ValidateDataSchema(ctx, dataData); err != nil {
		merr = multierror.Append(merr, err)
	}
	if err := ValidateResultSchema(ctx, resultData); err != nil {
		merr = multierror.Append(merr, err)
	}
	return merr.ErrorOrNil()
}

func ValidateArgSchema(ctx context.Context, schemaData []byte) error {
	if len(schemaData) == 0 {
		return nil
	}
	keyErrs, err := ArgSchema.ValidateBytes(ctx, schemaData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}

func ValidateDataSchema(ctx context.Context, schemaData []byte) error {
	if len(schemaData) == 0 {
		return nil
	}
	keyErrs, err := DataSchema.ValidateBytes(ctx, schemaData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}

func ValidateResultSchema(ctx context.Context, schemaData []byte) error {
	if len(schemaData) == 0 {
		return nil
	}
	keyErrs, err := ResultSchema.ValidateBytes(ctx, schemaData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}
