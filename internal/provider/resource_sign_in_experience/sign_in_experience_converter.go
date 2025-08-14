package resource_sign_in_experience

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func convertToTerraformModel(ctx context.Context, signInExperience *client.SignInExperienceModel, model *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	*model = SignInExperienceModel{
		Id:                        types.StringValue(signInExperience.ID),
		TenantId:                  types.StringValue(signInExperience.TenantId),
		TermsOfUseUrl:             types.StringValue(signInExperience.TermsOfUseUrl),
		PrivacyPolicyUrl:          types.StringValue(signInExperience.PrivacyPolicyUrl),
		AgreeToTermsPolicy:        types.StringValue(signInExperience.AgreeToTermsPolicy),
		SignInMode:                types.StringValue(signInExperience.SignInMode),
		CustomCss:                 types.StringValue(signInExperience.CustomCss),
		SingleSignOnEnabled:       types.BoolValue(signInExperience.SingleSignOnEnabled),
		SupportEmail:              types.StringValue(signInExperience.SupportEmail),
		SupportWebsiteUrl:         types.StringValue(signInExperience.SupportWebsiteUrl),
		UnknownSessionRedirectUrl: types.StringValue(signInExperience.UnknownSessionRedirectUrl),
	}

	diags.Append(convertCustomContent(signInExperience.CustomContent, model)...)
	diags.Append(convertSignInMethods(ctx, signInExperience.SignIn, model)...)

	convertBranding(ctx, signInExperience.Branding, model)
	convertCaptchaPolicy(ctx, signInExperience.CaptchaPolicy, model)
	convertColor(ctx, signInExperience.Color, model)
	convertEmailBlocklistPolicy(ctx, signInExperience.EmailBlocklistPolicy, model)
	convertLanguageInfo(ctx, signInExperience.LanguageInfo, model)
	convertMfa(ctx, signInExperience.Mfa, model)
	convertPasswordPolicy(ctx, signInExperience.PasswordPolicy, model)
	convertSentinelPolicy(ctx, signInExperience.SentinelPolicy, model)
	convertSignUp(ctx, signInExperience.SignUp, model)
	convertSocialSignIn(ctx, signInExperience.SocialSignIn, model)

	model.SocialSignInConnectorTargets = stringSliceToList(signInExperience.SocialSignInConnectorTargets)

	return diags
}

func convertBranding(ctx context.Context, apiBranding *client.Branding, tfModel *SignInExperienceModel) {
	if apiBranding != nil {
		tfModel.Branding = NewBrandingValueMust(
			BrandingValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"logo_url":      types.StringValue(apiBranding.LogoUrl),
				"dark_logo_url": types.StringValue(apiBranding.DarkLogoUrl),
				"favicon":       types.StringValue(apiBranding.Favicon),
				"dark_favicon":  types.StringValue(apiBranding.DarkFavicon),
			},
		)
	} else {
		tfModel.Branding = NewBrandingValueNull()
	}
}

func convertCaptchaPolicy(ctx context.Context, apiCaptchaPolicy *client.CaptchaPolicy, tfModel *SignInExperienceModel) {
	if apiCaptchaPolicy != nil {
		tfModel.CaptchaPolicy = NewCaptchaPolicyValueMust(
			CaptchaPolicyValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"enabled": types.BoolValue(apiCaptchaPolicy.Enabled),
			},
		)
	} else {
		tfModel.CaptchaPolicy = NewCaptchaPolicyValueNull()
	}
}

func convertColor(ctx context.Context, apiColor *client.Color, tfModel *SignInExperienceModel) {
	if apiColor != nil {
		tfModel.Color = NewColorValueMust(
			ColorValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"primary_color":        types.StringValue(apiColor.PrimaryColor),
				"dark_primary_color":   types.StringValue(apiColor.DarkPrimaryColor),
				"is_dark_mode_enabled": types.BoolValue(apiColor.IsDarkModeEnabled),
			},
		)
	} else {
		tfModel.Color = NewColorValueNull()
	}
}

func convertCustomContent(apiCustomContent map[string]string, tfModel *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if len(apiCustomContent) > 0 {
		elements := make(map[string]attr.Value)
		for k, v := range apiCustomContent {
			elements[k] = types.StringValue(v)
		}
		tfModel.CustomContent, diags = types.MapValue(types.StringType, elements)
	} else {
		tfModel.CustomContent = types.MapNull(types.StringType)
	}
	return diags
}

func convertEmailBlocklistPolicy(ctx context.Context, apiEmailBlocklistPolicy *client.EmailBlocklistPolicy, tfModel *SignInExperienceModel) {
	if apiEmailBlocklistPolicy != nil {
		tfModel.EmailBlocklistPolicy = NewEmailBlocklistPolicyValueMust(
			EmailBlocklistPolicyValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"block_disposable_addresses": types.BoolValue(apiEmailBlocklistPolicy.BlockDisposableAddresses),
				"block_subaddressing":        types.BoolValue(apiEmailBlocklistPolicy.BlockSubaddressing),
				"custom_blocklist":           stringSliceToList(apiEmailBlocklistPolicy.CustomBlocklist),
			},
		)
	} else {
		tfModel.EmailBlocklistPolicy = NewEmailBlocklistPolicyValueNull()
	}
}

func convertLanguageInfo(ctx context.Context, apiLanguageInfo *client.LanguageInfo, tfModel *SignInExperienceModel) {
	if apiLanguageInfo != nil {
		tfModel.LanguageInfo = NewLanguageInfoValueMust(
			LanguageInfoValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"auto_detect":       types.BoolValue(apiLanguageInfo.AutoDetect),
				"fallback_language": types.StringValue(apiLanguageInfo.FallbackLanguage),
			},
		)
	} else {
		tfModel.LanguageInfo = NewLanguageInfoValueNull()
	}
}

func convertMfa(ctx context.Context, apiMfa *client.Mfa, tfModel *SignInExperienceModel) {
	if apiMfa != nil {
		tfModel.Mfa = NewMfaValueMust(
			MfaValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"factors":                          stringSliceToList(apiMfa.Factors),
				"policy":                           types.StringValue(apiMfa.Policy),
				"organization_required_mfa_policy": types.StringValue(apiMfa.OrganizationRequiredMfaPolicy),
			},
		)
	} else {
		tfModel.Mfa = NewMfaValueNull()
	}
}

func convertPasswordPolicy(ctx context.Context, apiPasswordPolicy *client.PasswordPolicy, tfModel *SignInExperienceModel) {
	if apiPasswordPolicy != nil {
		lengthValue := buildPasswordPolicyLength(ctx, apiPasswordPolicy.Length)
		characterTypesValue := buildPasswordPolicyCharacterTypes(ctx, apiPasswordPolicy.CharacterTypes)
		rejectsValue := buildPasswordPolicyRejects(ctx, apiPasswordPolicy.Rejects)

		tfModel.PasswordPolicy = NewPasswordPolicyValueMust(
			PasswordPolicyValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"length":          lengthValue,
				"character_types": characterTypesValue,
				"rejects":         rejectsValue,
			},
		)
	} else {
		tfModel.PasswordPolicy = NewPasswordPolicyValueNull()
	}
}

func buildPasswordPolicyLength(ctx context.Context, apiLength *client.Length) types.Object {
	if apiLength != nil {
		return types.ObjectValueMust(
			LengthValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"min": types.NumberValue(big.NewFloat(float64(apiLength.Min))),
				"max": types.NumberValue(big.NewFloat(float64(apiLength.Max))),
			},
		)
	}

	return basetypes.NewObjectNull(LengthValue{}.AttributeTypes(ctx))
}
func buildPasswordPolicyCharacterTypes(ctx context.Context, apiCharacterTypes *client.CharacterTypes) types.Object {
	if apiCharacterTypes != nil {
		return types.ObjectValueMust(
			CharacterTypesValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"min": types.NumberValue(big.NewFloat(float64(apiCharacterTypes.Min))),
			},
		)
	}

	return basetypes.NewObjectNull(CharacterTypesValue{}.AttributeTypes(ctx))
}
func buildPasswordPolicyRejects(ctx context.Context, apiRejects *client.Rejects) types.Object {
	if apiRejects != nil {
		return types.ObjectValueMust(
			RejectsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"pwned":                   types.BoolValue(apiRejects.Pwned),
				"repetition_and_sequence": types.BoolValue(apiRejects.RepetitionAndSequence),
				"user_info":               types.BoolValue(apiRejects.UserInfo),
				"words":                   stringSliceToList(apiRejects.Words),
			},
		)
	}

	return basetypes.NewObjectNull(RejectsValue{}.AttributeTypes(ctx))
}

func convertSentinelPolicy(ctx context.Context, apiSentinelPolicy *client.SentinelPolicy, tfModel *SignInExperienceModel) {
	if apiSentinelPolicy != nil {
		tfModel.SentinelPolicy = NewSentinelPolicyValueMust(
			SentinelPolicyValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"max_attempts":     types.NumberValue(big.NewFloat(apiSentinelPolicy.MaxAttempts)),
				"lockout_duration": types.NumberValue(big.NewFloat(apiSentinelPolicy.LockoutDuration)),
			},
		)
	} else {
		tfModel.SentinelPolicy = NewSentinelPolicyValueNull()
	}
}

func convertSignInMethods(ctx context.Context, apiSignIn *client.SignIn, tfModel *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiSignIn != nil && apiSignIn.Methods != nil && len(apiSignIn.Methods) > 0 {
		var methodsAttrValues []attr.Value
		for _, m := range apiSignIn.Methods {
			methodValue := NewMethodsValueMust(
				MethodsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"identifier":          types.StringValue(m.Identifier),
					"is_password_primary": types.BoolValue(m.IsPasswordPrimary),
					"password":            types.BoolValue(m.Password),
					"verification_code":   types.BoolValue(m.VerificationCode),
				},
			)
			methodsAttrValues = append(methodsAttrValues, methodValue)
		}

		methodsList, diagsList := types.ListValueFrom(ctx, MethodsValue{}.Type(ctx), methodsAttrValues)
		diags.Append(diagsList...)
		if diags.HasError() {
			return diags
		}

		tfModel.SignIn = NewSignInValueMust(
			SignInValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"methods": methodsList,
			},
		)
	} else {
		tfModel.SignIn = NewSignInValueNull()
	}

	return diags
}

func convertSignUp(ctx context.Context, apiSignUp *client.SignUp, tfModel *SignInExperienceModel) {
	if apiSignUp != nil {
		var secondaryIdentifiersAttr attr.Value

		if apiSignUp.SecondaryIdentifiers != nil && len(*apiSignUp.SecondaryIdentifiers) > 0 {
			secondaryIdentifier := (*apiSignUp.SecondaryIdentifiers)[0]
			secondaryIdentifiersAttr = types.ObjectValueMust(
				SecondaryIdentifiersValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"identifier": types.StringValue(secondaryIdentifier.Identifier),
					"verify":     types.BoolValue(secondaryIdentifier.Verify),
				},
			)
		} else {
			secondaryIdentifiersAttr = types.ObjectNull(SecondaryIdentifiersValue{}.AttributeTypes(ctx))
		}

		tfSignUp := NewSignUpValueMust(
			SignUpValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"identifiers":           stringSliceToList(apiSignUp.Identifiers),
				"password":              types.BoolValue(apiSignUp.Password),
				"verify":                types.BoolValue(apiSignUp.Verify),
				"secondary_identifiers": secondaryIdentifiersAttr,
			},
		)

		tfModel.SignUp = tfSignUp
	} else {
		tfModel.SignUp = NewSignUpValueNull()
	}
}
func convertSocialSignIn(ctx context.Context, apiSocialSignIn *client.SocialSignIn, tfModel *SignInExperienceModel) {
	if apiSocialSignIn != nil {
		tfModel.SocialSignIn = NewSocialSignInValueMust(
			SocialSignInValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"automatic_account_linking": types.BoolValue(apiSocialSignIn.AutomaticAccountLinking),
			},
		)
	} else {
		tfModel.SocialSignIn = NewSocialSignInValueNull()
	}
}
