package resource_sign_in_experience

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	diags.Append(convertSignUp(ctx, signInExperience.SignUp, model)...)

	convertBranding(signInExperience.Branding, model)
	convertCaptchaPolicy(signInExperience.CaptchaPolicy, model)
	convertColor(signInExperience.Color, model)
	convertEmailBlocklistPolicy(signInExperience.EmailBlocklistPolicy, model)
	convertLanguageInfo(signInExperience.LanguageInfo, model)
	convertMfa(signInExperience.Mfa, model)
	convertPasswordPolicy(ctx, signInExperience.PasswordPolicy, model)
	convertSentinelPolicy(signInExperience.SentinelPolicy, model)
	convertSocialSignIn(signInExperience.SocialSignIn, model)

	model.SocialSignInConnectorTargets = stringSliceToList(signInExperience.SocialSignInConnectorTargets)

	return diags
}

func convertBranding(apiBranding *client.Branding, tfModel *SignInExperienceModel) {
	if apiBranding != nil {
		tfModel.Branding = BrandingValue{
			LogoUrl:     types.StringValue(apiBranding.LogoUrl),
			DarkLogoUrl: types.StringValue(apiBranding.DarkLogoUrl),
			Favicon:     types.StringValue(apiBranding.Favicon),
			DarkFavicon: types.StringValue(apiBranding.DarkFavicon),
			state:       attr.ValueStateKnown,
		}
	} else {
		tfModel.Branding = NewBrandingValueNull()
	}
}

func convertCaptchaPolicy(apiCaptchaPolicy *client.CaptchaPolicy, tfModel *SignInExperienceModel) {
	if apiCaptchaPolicy != nil {
		tfModel.CaptchaPolicy = CaptchaPolicyValue{
			Enabled: types.BoolValue(apiCaptchaPolicy.Enabled),
		}
	} else {
		tfModel.CaptchaPolicy = NewCaptchaPolicyValueNull()
	}
}

func convertColor(apiColor *client.Color, tfModel *SignInExperienceModel) {
	if apiColor != nil {
		tfModel.Color = ColorValue{
			PrimaryColor:      types.StringValue(apiColor.PrimaryColor),
			DarkPrimaryColor:  types.StringValue(apiColor.DarkPrimaryColor),
			IsDarkModeEnabled: types.BoolValue(apiColor.IsDarkModeEnabled),
		}
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

func convertEmailBlocklistPolicy(apiEmailBlocklistPolicy *client.EmailBlocklistPolicy, tfModel *SignInExperienceModel) {
	if apiEmailBlocklistPolicy != nil {

		tfModel.EmailBlocklistPolicy = EmailBlocklistPolicyValue{
			BlockDisposableAddresses: types.BoolValue(apiEmailBlocklistPolicy.BlockDisposableAddresses),
			BlockSubaddressing:       types.BoolValue(apiEmailBlocklistPolicy.BlockSubaddressing),
			CustomBlocklist:          stringSliceToList(apiEmailBlocklistPolicy.CustomBlocklist),
		}
	} else {
		tfModel.EmailBlocklistPolicy = NewEmailBlocklistPolicyValueNull()
	}
}

func convertLanguageInfo(apiLanguageInfo *client.LanguageInfo, tfModel *SignInExperienceModel) {
	if apiLanguageInfo != nil {
		tfModel.LanguageInfo = LanguageInfoValue{
			AutoDetect:       types.BoolValue(apiLanguageInfo.AutoDetect),
			FallbackLanguage: types.StringValue(apiLanguageInfo.FallbackLanguage),
		}
	} else {
		tfModel.LanguageInfo = NewLanguageInfoValueNull()
	}
}

func convertMfa(apiMfa *client.Mfa, tfModel *SignInExperienceModel) {
	if apiMfa != nil {
		tfModel.Mfa = MfaValue{
			Factors:                       stringSliceToList(apiMfa.Factors),
			Policy:                        types.StringValue(apiMfa.Policy),
			OrganizationRequiredMfaPolicy: types.StringValue(apiMfa.OrganizationRequiredMfaPolicy),
		}
	} else {
		tfModel.Mfa = NewMfaValueNull()
	}
}

func convertPasswordPolicy(ctx context.Context, apiPasswordPolicy *client.PasswordPolicy, tfModel *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiPasswordPolicy != nil {
		lengthValue, diags := buildPasswordPolicyLength(apiPasswordPolicy.Length).ToObjectValue(ctx)
		diags.Append(diags...)
		characterTypesValue, diags := buildPasswordPolicyCharacterTypes(apiPasswordPolicy.CharacterTypes).ToObjectValue(ctx)
		diags.Append(diags...)
		rejectsValue, diags := buildPasswordPolicyRejects(apiPasswordPolicy.Rejects).ToObjectValue(ctx)
		diags.Append(diags...)

		if diags.HasError() {
			return diags
		}

		tfModel.PasswordPolicy = PasswordPolicyValue{
			Length:         lengthValue,
			CharacterTypes: characterTypesValue,
			Rejects:        rejectsValue,
		}
	} else {
		tfModel.PasswordPolicy = NewPasswordPolicyValueNull()
	}

	return diags
}

func buildPasswordPolicyLength(apiLength *client.Length) LengthValue {
	if apiLength != nil {
		return LengthValue{
			Min: types.NumberValue(big.NewFloat(float64(apiLength.Min))),
			Max: types.NumberValue(big.NewFloat(float64(apiLength.Max))),
		}
	}

	return NewLengthValueNull()
}
func buildPasswordPolicyCharacterTypes(apiCharacterTypes *client.CharacterTypes) CharacterTypesValue {
	if apiCharacterTypes != nil {
		return CharacterTypesValue{
			Min: types.NumberValue(big.NewFloat(float64(apiCharacterTypes.Min))),
		}
	}

	return NewCharacterTypesValueNull()
}
func buildPasswordPolicyRejects(apiRejects *client.Rejects) RejectsValue {
	if apiRejects != nil {
		return RejectsValue{
			Pwned:                 types.BoolValue(apiRejects.Pwned),
			RepetitionAndSequence: types.BoolValue(apiRejects.RepetitionAndSequence),
			UserInfo:              types.BoolValue(apiRejects.UserInfo),
			Words:                 stringSliceToList(apiRejects.Words),
		}
	}

	return NewRejectsValueNull()
}

func convertSentinelPolicy(apiSentinelPolicy *client.SentinelPolicy, tfModel *SignInExperienceModel) {
	if apiSentinelPolicy != nil {
		tfModel.SentinelPolicy = SentinelPolicyValue{
			MaxAttempts:     types.NumberValue(big.NewFloat(apiSentinelPolicy.MaxAttempts)),
			LockoutDuration: types.NumberValue(big.NewFloat(apiSentinelPolicy.LockoutDuration)),
		}
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

func convertSignUp(ctx context.Context, apiSignUp *client.SignUp, tfModel *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiSignUp != nil {
		var secondaryIdentifiersAttr SecondaryIdentifiersValue

		if apiSignUp.SecondaryIdentifiers != nil && len(*apiSignUp.SecondaryIdentifiers) > 0 {
			secondaryIdentifier := (*apiSignUp.SecondaryIdentifiers)[0]
			secondaryIdentifiersAttr = SecondaryIdentifiersValue{
				Identifier: types.StringValue(secondaryIdentifier.Identifier),
				Verify:     types.BoolValue(secondaryIdentifier.Verify),
			}
		} else {
			secondaryIdentifiersAttr = NewSecondaryIdentifiersValueNull()
		}

		secondaryIdentifiersAttrObj, diags := secondaryIdentifiersAttr.ToObjectValue(ctx)
		if diags.HasError() {
			return diags
		}

		tfSignUp := SignUpValue{
			Identifiers:          stringSliceToList(apiSignUp.Identifiers),
			Password:             types.BoolValue(apiSignUp.Password),
			Verify:               types.BoolValue(apiSignUp.Verify),
			SecondaryIdentifiers: secondaryIdentifiersAttrObj,
		}

		tfModel.SignUp = tfSignUp
	} else {
		tfModel.SignUp = NewSignUpValueNull()
	}

	return diags
}
func convertSocialSignIn(apiSocialSignIn *client.SocialSignIn, tfModel *SignInExperienceModel) {
	if apiSocialSignIn != nil {
		tfModel.SocialSignIn = SocialSignInValue{
			AutomaticAccountLinking: types.BoolValue(apiSocialSignIn.AutomaticAccountLinking),
		}
	} else {
		tfModel.SocialSignIn = NewSocialSignInValueNull()
	}
}
