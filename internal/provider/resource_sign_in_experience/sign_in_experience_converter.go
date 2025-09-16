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

	diags.Append(convertSignIn(ctx, signInExperience.SignIn, model)...)
	diags.Append(convertSignUp(ctx, signInExperience.SignUp, model)...)
	diags.Append(convertPasswordPolicy(ctx, signInExperience.PasswordPolicy, model)...)

	convertBranding(signInExperience.Branding, model)
	convertCaptchaPolicy(signInExperience.CaptchaPolicy, model)
	convertColor(signInExperience.Color, model)
	convertEmailBlocklistPolicy(signInExperience.EmailBlocklistPolicy, model)
	convertLanguageInfo(signInExperience.LanguageInfo, model)
	convertMfa(signInExperience.Mfa, model)
	convertSentinelPolicy(signInExperience.SentinelPolicy, model)
	convertSocialSignIn(signInExperience.SocialSignIn, model)

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
			state:   attr.ValueStateKnown,
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
			state:             attr.ValueStateKnown,
		}
	} else {
		tfModel.Color = NewColorValueNull()
	}
}

func convertEmailBlocklistPolicy(apiEmailBlocklistPolicy *client.EmailBlocklistPolicy, tfModel *SignInExperienceModel) {
	if apiEmailBlocklistPolicy != nil {

		tfModel.EmailBlocklistPolicy = EmailBlocklistPolicyValue{
			BlockDisposableAddresses: types.BoolValue(apiEmailBlocklistPolicy.BlockDisposableAddresses),
			BlockSubaddressing:       types.BoolValue(apiEmailBlocklistPolicy.BlockSubaddressing),
			CustomBlocklist:          stringSliceToList(apiEmailBlocklistPolicy.CustomBlocklist),
			state:                    attr.ValueStateKnown,
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
			state:            attr.ValueStateKnown,
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
			state:                         attr.ValueStateKnown,
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
			state:          attr.ValueStateKnown,
		}
	} else {
		tfModel.PasswordPolicy = NewPasswordPolicyValueNull()
	}

	return diags
}

func buildPasswordPolicyLength(apiLength *client.Length) LengthValue {
	if apiLength != nil {
		return LengthValue{
			Min:   types.NumberValue(big.NewFloat(float64(apiLength.Min))),
			Max:   types.NumberValue(big.NewFloat(float64(apiLength.Max))),
			state: attr.ValueStateKnown,
		}
	}

	return NewLengthValueNull()
}
func buildPasswordPolicyCharacterTypes(apiCharacterTypes *client.CharacterTypes) CharacterTypesValue {
	if apiCharacterTypes != nil {
		return CharacterTypesValue{
			Min:   types.NumberValue(big.NewFloat(float64(apiCharacterTypes.Min))),
			state: attr.ValueStateKnown,
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
			state:                 attr.ValueStateKnown,
		}
	}

	return NewRejectsValueNull()
}

func convertSentinelPolicy(apiSentinelPolicy *client.SentinelPolicy, tfModel *SignInExperienceModel) {
	if apiSentinelPolicy != nil {
		tfModel.SentinelPolicy = SentinelPolicyValue{
			MaxAttempts:     types.NumberValue(big.NewFloat(apiSentinelPolicy.MaxAttempts)),
			LockoutDuration: types.NumberValue(big.NewFloat(apiSentinelPolicy.LockoutDuration)),
			state:           attr.ValueStateKnown,
		}
	} else {
		tfModel.SentinelPolicy = NewSentinelPolicyValueNull()
	}
}

func convertSignIn(ctx context.Context, apiSignIn *client.SignIn, tfModel *SignInExperienceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	var methodsAttrValues []attr.Value
	if apiSignIn != nil && apiSignIn.Methods != nil {
		for _, m := range apiSignIn.Methods {
			methodValue := MethodsValue{
				Identifier:        types.StringValue(m.Identifier),
				IsPasswordPrimary: types.BoolValue(m.IsPasswordPrimary),
				Password:          types.BoolValue(m.Password),
				VerificationCode:  types.BoolValue(m.VerificationCode),
				state:             attr.ValueStateKnown,
			}
			methodsAttrValues = append(methodsAttrValues, methodValue)
		}
	}

	methodsList, diagsList := types.ListValueFrom(ctx, MethodsValue{}.Type(ctx), methodsAttrValues)
	diags.Append(diagsList...)
	if diags.HasError() {
		return diags
	}

	tfModel.SignIn = SignInValue{
		Methods: methodsList,
		state:   attr.ValueStateKnown,
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
				state:      attr.ValueStateKnown,
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
			state:                attr.ValueStateKnown,
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
			state:                   attr.ValueStateKnown,
		}
	} else {
		tfModel.SocialSignIn = NewSocialSignInValueNull()
	}
}
