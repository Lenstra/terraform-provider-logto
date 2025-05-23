package resource_sign_in_experience

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SignInExperienceBuilder struct {
	ctx   context.Context
	model SignInExperienceModel
	diags diag.Diagnostics
}

func NewSignInExperienceBuilder(ctx context.Context) *SignInExperienceBuilder {
	return &SignInExperienceBuilder{ctx: ctx}
}

func (b *SignInExperienceBuilder) FromTfModel(tfModel *SignInExperienceModel) (*client.SignInExperienceModel, diag.Diagnostics) {
	b.model = *tfModel
	model := &client.SignInExperienceModel{}

	model.TenantId = tfModel.TenantId.ValueString()
	model.ID = tfModel.Id.ValueString()
	model.TermsOfUseUrl = tfModel.TermsOfUseUrl.ValueString()
	model.PrivacyPolicyUrl = tfModel.PrivacyPolicyUrl.ValueString()
	model.AgreeToTermsPolicy = tfModel.AgreeToTermsPolicy.ValueString()
	model.SignInMode = tfModel.SignInMode.ValueString()
	model.CustomCss = tfModel.CustomCss.ValueString()
	model.SingleSignOnEnabled = tfModel.SingleSignOnEnabled.ValueBool()
	model.SupportEmail = tfModel.SupportEmail.ValueString()
	model.SupportWebsiteUrl = tfModel.SupportWebsiteUrl.ValueString()
	model.UnknownSessionRedirectUrl = tfModel.UnknownSessionRedirectUrl.ValueString()

	model.Color = b.buildColor()
	model.Branding = b.buildBranding()
	model.LanguageInfo = b.buildLanguageInfo()
	model.SignIn = b.buildSignIn()
	model.SignUp = b.buildSignUp() // NEED TO ENABLE CONNECTORS TO USE IT
	model.SocialSignIn = b.buildSocialSignIn()

	slice, diagsList := convertListToSlice(b.ctx, b.model.SocialSignInConnectorTargets)
	b.addDiags(diagsList)

	model.SocialSignInConnectorTargets = slice

	model.CustomContent = b.buildCustomContent()
	model.PasswordPolicy = b.buildPasswordPolicy()
	model.Mfa = b.buildMfa()
	model.CaptchaPolicy = b.buildCaptchaPolicy()
	model.SentinelPolicy = b.buildSentinelPolicy()
	model.EmailBlocklistPolicy = b.buildEmailBlocklistPolicy()

	if b.diags.HasError() {
		return nil, b.diags
	}

	return model, b.diags
}

func (b *SignInExperienceBuilder) buildColor() *client.Color {
	if b.model.Color.IsNull() || b.model.Color.IsUnknown() {
		return nil
	}
	return &client.Color{
		PrimaryColor:      b.model.Color.PrimaryColor.ValueString(),
		IsDarkModeEnabled: b.model.Color.IsDarkModeEnabled.ValueBool(),
		DarkPrimaryColor:  b.model.Color.DarkPrimaryColor.ValueString(),
	}
}

func (b *SignInExperienceBuilder) buildBranding() *client.Branding {
	if b.model.Branding.IsNull() || b.model.Branding.IsUnknown() {
		return nil
	}
	return &client.Branding{
		LogoUrl:     b.model.Branding.LogoUrl.ValueString(),
		DarkLogoUrl: b.model.Branding.DarkLogoUrl.ValueString(),
		Favicon:     b.model.Branding.Favicon.ValueString(),
		DarkFavicon: b.model.Branding.DarkFavicon.ValueString(),
	}
}

func (b *SignInExperienceBuilder) buildLanguageInfo() *client.LanguageInfo {
	if b.model.LanguageInfo.IsNull() || b.model.LanguageInfo.IsUnknown() {
		return nil
	}
	return &client.LanguageInfo{
		AutoDetect:       b.model.LanguageInfo.AutoDetect.ValueBool(),
		FallbackLanguage: b.model.LanguageInfo.FallbackLanguage.ValueString(),
	}
}

func (b *SignInExperienceBuilder) buildSignIn() *client.SignIn {
	if b.model.SignIn.IsNull() || b.model.SignIn.IsUnknown() {
		return nil
	}
	if b.model.SignIn.Methods.IsNull() || b.model.SignIn.Methods.IsUnknown() {
		return nil
	}

	var methodsValues []MethodsValue
	diagsList := b.model.SignIn.Methods.ElementsAs(b.ctx, &methodsValues, false)

	b.addDiags(diagsList)
	if b.diags.HasError() {
		return nil
	}

	var methods []client.Methods
	for _, methodVal := range methodsValues {
		methods = append(methods, client.Methods{
			Identifier:        methodVal.Identifier.ValueString(),
			IsPasswordPrimary: methodVal.IsPasswordPrimary.ValueBool(),
			Password:          methodVal.Password.ValueBool(),
			VerificationCode:  methodVal.VerificationCode.ValueBool(),
		})
	}

	return &client.SignIn{
		Methods: methods,
	}
}

func (b *SignInExperienceBuilder) buildSignUp() *client.SignUp {
	if b.model.SignUp.IsNull() || b.model.SignUp.IsUnknown() {
		return nil
	}

	list, diagsList := b.model.SignUp.Identifiers.ToListValue(b.ctx)
	b.addDiags(diagsList)
	if b.diags.HasError() {
		return nil
	}

	identifiers, diagsList := convertListToSlice(b.ctx, list)
	b.addDiags(diagsList)
	if b.diags.HasError() {
		return nil
	}

	attrs := b.model.SignUp.SecondaryIdentifiers.Attributes()

	var secondaryIdentifiers []client.SecondaryIdentifiers

	if !b.model.SignUp.SecondaryIdentifiers.IsNull() && !b.model.SignUp.SecondaryIdentifiers.IsUnknown() {
		identifier := ""
		verify := false

		if idVal, ok := attrs["identifier"].(basetypes.StringValue); ok && !idVal.IsNull() {
			identifier = idVal.ValueString()
		}

		if verifyVal, ok := attrs["verify"].(basetypes.BoolValue); ok && !verifyVal.IsNull() {
			v := verifyVal.ValueBool()
			verify = v
		}

		secondaryIdentifiers = append(secondaryIdentifiers, client.SecondaryIdentifiers{
			Identifier: identifier,
			Verify:     verify,
		})
	}

	return &client.SignUp{
		Identifiers:          identifiers,
		Password:             b.model.SignUp.Password.ValueBool(),
		Verify:               b.model.SignUp.Verify.ValueBool(),
		SecondaryIdentifiers: &secondaryIdentifiers,
	}
}

func (b *SignInExperienceBuilder) buildSocialSignIn() *client.SocialSignIn {
	if b.model.SocialSignIn.IsNull() || b.model.SocialSignIn.IsUnknown() {
		return nil
	}
	return &client.SocialSignIn{
		AutomaticAccountLinking: b.model.SocialSignIn.AutomaticAccountLinking.ValueBool(),
	}
}

func (b *SignInExperienceBuilder) buildCustomContent() map[string]string {
	customContent := make(map[string]string)
	if b.model.CustomContent.IsNull() || b.model.CustomContent.IsUnknown() || b.model.CustomContent.Elements() == nil {
		return customContent
	}
	for key, val := range b.model.CustomContent.Elements() {
		if v, ok := val.(types.String); ok {
			customContent[key] = v.ValueString()
		} else {
			b.diags.AddWarning("CustomContent Map Conversion", "Map element for key '"+key+"' is not a string type.")
		}
	}
	return customContent
}

func (b *SignInExperienceBuilder) buildPasswordPolicy() *client.PasswordPolicy {
	if b.model.PasswordPolicy.IsNull() || b.model.PasswordPolicy.IsUnknown() {
		return nil
	}

	p := &client.PasswordPolicy{}

	length, diagsLength := b.buildPasswordPolicyLength()
	b.addDiags(diagsLength)
	p.Length = length

	characterTypes, diagsCharTypes := b.buildPasswordPolicyCharacterTypes()
	b.addDiags(diagsCharTypes)
	p.CharacterTypes = characterTypes

	rejects, diagsRejects := b.buildPasswordPolicyRejects()
	b.addDiags(diagsRejects)
	p.Rejects = rejects

	if b.diags.HasError() {
		return nil
	}

	return p
}

func (b *SignInExperienceBuilder) buildPasswordPolicyLength() (*client.Length, diag.Diagnostics) {
	var localDiags diag.Diagnostics

	if b.model.PasswordPolicy.Length.IsNull() || b.model.PasswordPolicy.Length.IsUnknown() {
		return nil, localDiags
	}
	lengthAttrs := b.model.PasswordPolicy.Length.Attributes()

	var clientMin, clientMax int64
	if minVal, ok := lengthAttrs["min"].(basetypes.NumberValue); ok && !minVal.IsNull() && !minVal.IsUnknown() {
		minBigFloat := minVal.ValueBigFloat()
		if minBigFloat == nil {
			localDiags.AddError("Invalid PasswordPolicy Length Min", "The 'min' value for password_policy.length is invalid (nil big.Float).")
			return nil, localDiags
		}
		var accuracy big.Accuracy
		clientMin, accuracy = minBigFloat.Int64()
		if accuracy == big.Below || accuracy == big.Above {
			localDiags.AddWarning("PasswordPolicy Length Min Precision", "The 'min' value for password_policy.length had precision loss during conversion to int64.")
		}
	} else {
		localDiags.AddError("Invalid PasswordPolicy Length Min", "The 'min' value for password_policy.length is missing or invalid.")
		return nil, localDiags
	}
	if maxVal, ok := lengthAttrs["max"].(basetypes.NumberValue); ok && !maxVal.IsNull() && !maxVal.IsUnknown() {
		maxBigFloat := maxVal.ValueBigFloat()
		if maxBigFloat == nil {
			localDiags.AddError("Invalid PasswordPolicy Length Max", "The 'max' value for password_policy.length is invalid (nil big.Float).")
			return nil, localDiags
		}
		var accuracy big.Accuracy
		clientMax, accuracy = maxBigFloat.Int64()
		if accuracy == big.Below || accuracy == big.Above {
			localDiags.AddWarning("PasswordPolicy Length Max Precision", "The 'max' value for password_policy.length had precision loss during conversion to int64.")
		}
	} else {
		localDiags.AddError("Invalid PasswordPolicy Length Max", "The 'max' value for password_policy.length is missing or invalid.")
		return nil, localDiags
	}
	return &client.Length{Min: clientMin, Max: clientMax}, localDiags
}

func (b *SignInExperienceBuilder) buildPasswordPolicyCharacterTypes() (*client.CharacterTypes, diag.Diagnostics) {
	var localDiags diag.Diagnostics

	if b.model.PasswordPolicy.CharacterTypes.IsNull() || b.model.PasswordPolicy.CharacterTypes.IsUnknown() {
		return nil, localDiags
	}
	charTypeAttrs := b.model.PasswordPolicy.CharacterTypes.Attributes()
	var clientMinChar int64

	if minVal, ok := charTypeAttrs["min"].(basetypes.NumberValue); ok && !minVal.IsNull() && !minVal.IsUnknown() {
		minBigFloat := minVal.ValueBigFloat()
		if minBigFloat == nil {
			localDiags.AddError("Invalid PasswordPolicy CharacterTypes Min", "The 'min' value for password_policy.character_types is invalid (nil big.Float).")
			return nil, localDiags
		}
		var accuracy big.Accuracy
		clientMinChar, accuracy = minBigFloat.Int64()
		if accuracy == big.Below || accuracy == big.Above {
			localDiags.AddWarning("PasswordPolicy CharacterTypes Min Precision", "The 'min' value for password_policy.character_types had precision loss during conversion to int64.")
		}
	} else {
		localDiags.AddError("Invalid PasswordPolicy CharacterTypes Min", "The 'min' value for password_policy.character_types is missing or invalid.")
		return nil, localDiags
	}
	return &client.CharacterTypes{Min: clientMinChar}, localDiags
}

func (b *SignInExperienceBuilder) buildPasswordPolicyRejects() (*client.Rejects, diag.Diagnostics) {
	var localDiags diag.Diagnostics

	if b.model.PasswordPolicy.Rejects.IsNull() || b.model.PasswordPolicy.Rejects.IsUnknown() {
		return nil, localDiags
	}
	rejectAttrs := b.model.PasswordPolicy.Rejects.Attributes()

	var clientPwned bool
	if pwnedVal, ok := rejectAttrs["pwned"].(basetypes.BoolValue); ok && !pwnedVal.IsNull() {
		clientPwned = pwnedVal.ValueBool()
	}

	var clientRepetitionAndSequence bool
	if repSeqVal, ok := rejectAttrs["repetition_and_sequence"].(basetypes.BoolValue); ok && !repSeqVal.IsNull() {
		clientRepetitionAndSequence = repSeqVal.ValueBool()
	}

	var clientUserInfo bool
	if userInfoVal, ok := rejectAttrs["user_info"].(basetypes.BoolValue); ok && !userInfoVal.IsNull() {
		clientUserInfo = userInfoVal.ValueBool()
	}

	var clientWords []string
	if wordsVal, ok := rejectAttrs["words"].(basetypes.ListValue); ok && !wordsVal.IsNull() && !wordsVal.IsUnknown() {
		wordsList, diagsList := convertListToSlice(b.ctx, wordsVal)

		localDiags.Append(diagsList...)
		if localDiags.HasError() {
			return nil, localDiags
		}

		clientWords = wordsList
	}

	return &client.Rejects{
		Pwned:                 clientPwned,
		RepetitionAndSequence: clientRepetitionAndSequence,
		UserInfo:              clientUserInfo,
		Words:                 clientWords,
	}, localDiags
}

func (b *SignInExperienceBuilder) buildMfa() *client.Mfa {
	if b.model.Mfa.Policy.IsNull() || b.model.Mfa.Policy.IsUnknown() {
		return nil
	}

	policy := b.model.Mfa.Policy.ValueString()

	var factors []string
	factors, diagsList := convertListToSlice(b.ctx, b.model.Mfa.Factors)

	b.addDiags(diagsList)
	if b.diags.HasError() {
		return nil
	}

	var clientOrgRequiredMfaPolicy string
	if !b.model.Mfa.OrganizationRequiredMfaPolicy.IsNull() && !b.model.Mfa.OrganizationRequiredMfaPolicy.IsUnknown() {
		clientOrgRequiredMfaPolicy = b.model.Mfa.OrganizationRequiredMfaPolicy.ValueString()
	} else {
		clientOrgRequiredMfaPolicy = ""
	}

	return &client.Mfa{
		Factors:                       factors,
		Policy:                        policy,
		OrganizationRequiredMfaPolicy: clientOrgRequiredMfaPolicy,
	}
}

func (b *SignInExperienceBuilder) buildCaptchaPolicy() *client.CaptchaPolicy {
	if b.model.CaptchaPolicy.IsNull() || b.model.CaptchaPolicy.IsUnknown() {
		return &client.CaptchaPolicy{
			Enabled: true,
		}
	}

	enabledValue := b.model.CaptchaPolicy.Enabled

	apiEnabled := true
	if !enabledValue.IsUnknown() && !enabledValue.IsNull() {
		apiEnabled = enabledValue.ValueBool()
	}

	return &client.CaptchaPolicy{
		Enabled: apiEnabled,
	}
}

func (b *SignInExperienceBuilder) buildSentinelPolicy() *client.SentinelPolicy {
	if b.model.SentinelPolicy.IsNull() || b.model.SentinelPolicy.IsUnknown() {
		return nil
	}

	var maxAttempts float64
	if b.model.SentinelPolicy.MaxAttempts.IsNull() || b.model.SentinelPolicy.MaxAttempts.IsUnknown() {
		b.diags.AddError(
			"Missing SentinelPolicy MaxAttempts",
			"The 'max_attempts' value in sentinel_policy must be set and not unknown.",
		)
		return nil
	}
	maxAttempts, _ = b.model.SentinelPolicy.MaxAttempts.ValueBigFloat().Float64()

	var lockoutDuration float64
	if b.model.SentinelPolicy.LockoutDuration.IsNull() || b.model.SentinelPolicy.LockoutDuration.IsUnknown() {
		b.diags.AddError(
			"Missing SentinelPolicy LockoutDuration",
			"The 'lockout_duration' value in sentinel_policy must be set and not unknown.",
		)
		return nil
	}
	lockoutDuration, _ = b.model.SentinelPolicy.LockoutDuration.ValueBigFloat().Float64()

	return &client.SentinelPolicy{
		MaxAttempts:     maxAttempts,
		LockoutDuration: lockoutDuration,
	}
}

func (b *SignInExperienceBuilder) buildEmailBlocklistPolicy() *client.EmailBlocklistPolicy {
	if b.model.EmailBlocklistPolicy.IsNull() || b.model.EmailBlocklistPolicy.IsUnknown() {
		return nil
	}

	var customBlocklist []string
	customBlocklist, diagsList := convertListToSlice(b.ctx, b.model.EmailBlocklistPolicy.CustomBlocklist)

	b.addDiags(diagsList)
	if b.diags.HasError() {
		return nil
	}

	return &client.EmailBlocklistPolicy{
		BlockDisposableAddresses: b.model.EmailBlocklistPolicy.BlockDisposableAddresses.ValueBool(),
		BlockSubaddressing:       b.model.EmailBlocklistPolicy.BlockSubaddressing.ValueBool(),
		CustomBlocklist:          customBlocklist,
	}
}

func (b *SignInExperienceBuilder) addDiags(newDiags diag.Diagnostics) {
	b.diags.Append(newDiags...)
}
