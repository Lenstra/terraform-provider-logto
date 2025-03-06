package utils

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecretGetter interface {
	GetApplicationSecrets(ctx context.Context, applicationID string) (types.Map, diag.Diagnostics)
}

func GetSecrets(ctx context.Context, c *client.Client, applicationID string) (types.Map, diag.Diagnostics) {
	secrets, err := c.GetApplicationSecrets(ctx, applicationID)
	if err != nil {
		return types.MapNull(types.StringType), diag.Diagnostics{
			diag.NewErrorDiagnostic("Error getting secrets", err.Error()),
		}
	}

	secretsMap := make(map[string]attr.Value, len(secrets))
	for _, v := range secrets {
		secretsMap[v.Name] = types.StringValue(v.Value)
	}

	return types.MapValue(types.StringType, secretsMap)
}
