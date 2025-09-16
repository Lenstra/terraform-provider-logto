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

// NewSignInExperienceBuilder.
func NewSignInExperienceBuilder(ctx context.Context) *SignInExperienceBuilder {
	return &SignInExperienceBuilder{ctx: ctx}
}

// DecodePlan convert a Terraform plan model into a client model.
func (b *SignInExperienceBuilder) DecodePlan(tfModel *SignInExperienceModel) (*client.SignInExperienceModel, diag.Diagnostics) {
	b.model = *tfModel
	model := &client.SignInExperienceModel{}

	model.ID = tfModel.Id.ValueString()
	model.TenantId = tfModel.TenantId.ValueString()
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

func (b *SignInExperienceBuilder) buildPasswordPolicy() *client.PasswordPolicy {
	if b.model.PasswordPolicy.IsNull() || b.model.PasswordPolicy.IsUnknown() {
		return nil
	}

	p := &client.PasswordPolicy{}

	length, diagsLength := b.buildPasswordPolicyLength()
	b.addDiags(diagsLength)
	p.Length = length

	characterTypes, diag := b.buildPasswordPolicyCharacterTypes()
	b.addDiags(diag)
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
	var diags diag.Diagnostics

	objectValue, diags := b.model.PasswordPolicy.Length.ToObjectValue(b.ctx)
	if diags.HasError() {
		return nil, diags
	}

	valuesMap := objectValue.Attributes()
	intValues := make(map[string]int64)

	for key, attrVal := range valuesMap {
		if attrVal.IsNull() {
			diags.AddError("Conversion Error", key+" attribute is null")
			continue
		}

		numberVal, ok := attrVal.(basetypes.NumberValue)
		if !ok {
			diags.AddError("Conversion Error", key+" attribute is not a NumberValue")
			continue
		}

		bf := numberVal.ValueBigFloat()
		if bf == nil {
			diags.AddError("Conversion Error", key+" attribute has no value")
			continue
		}

		i, accuracy := bf.Int64()
		if accuracy != big.Exact {
			diags.AddError("Conversion Error", key+" attribute must be an integer")
			continue
		}

		intValues[key] = i
	}

	minVal, minOk := intValues["min"]
	maxVal, maxOk := intValues["max"]
	if !minOk || !maxOk || minVal > maxVal {
		diags.AddError("Validation Error", "Password policy Min and Max values are not valid or coherent")
		return nil, diags
	}

	clientLength := &client.Length{
		Min: minVal,
		Max: maxVal,
	}

	return clientLength, diags
}

func (b *SignInExperienceBuilder) buildPasswordPolicyCharacterTypes() (*client.CharacterTypes, diag.Diagnostics) {
	var diags diag.Diagnostics

	objectValue, diags := b.model.PasswordPolicy.CharacterTypes.ToObjectValue(b.ctx)
	if diags.HasError() {
		return nil, diags
	}

	valuesMap := objectValue.Attributes()
	attrVal := valuesMap["min"]
	if attrVal.IsNull() {
		diags.AddError("Conversion Error", "Min attribute is null")
	}

	numberVal, ok := attrVal.(basetypes.NumberValue)
	if !ok {
		diags.AddError("Conversion Error", "PasswordPolicyCharacterTypes : Min attribute is not a NumberValue")
	}

	bf := numberVal.ValueBigFloat()
	if bf == nil {
		diags.AddError("Conversion Error", "PasswordPolicyCharacterTypes : Min attribute has no value")
	}

	i, accuracy := bf.Int64()
	if accuracy != big.Exact {
		diags.AddError("Conversion Error", "PasswordPolicyCharacterTypes : Min attribute must be an integer")
	}

	return &client.CharacterTypes{
		Min: i,
	}, diags
}

func (b *SignInExperienceBuilder) buildPasswordPolicyRejects() (*client.Rejects, diag.Diagnostics) {
	var diags diag.Diagnostics

	objectValue, diags := b.model.PasswordPolicy.Rejects.ToObjectValue(b.ctx)
	if diags.HasError() {
		return nil, diags
	}

	attrs := objectValue.Attributes()

	getBool := func(key string) bool {
		if val, ok := attrs[key]; ok && !val.IsNull() {
			if bVal, ok := val.(types.Bool); ok {
				return bVal.ValueBool()
			}
			diags.AddError("Conversion Error", key+" is not a bool")
		}
		return false
	}

	getStringList := func(key string) []string {
		result := []string{}
		if val, ok := attrs[key]; ok && !val.IsNull() {
			if listVal, ok := val.(types.List); ok {
				for _, e := range listVal.Elements() {
					if s, ok := e.(types.String); ok {
						result = append(result, s.ValueString())
					}
				}
			} else {
				diags.AddError("Conversion Error", key+" is not a list of strings")
			}
		}
		return result
	}

	clientRejects := &client.Rejects{
		Pwned:                 getBool("pwned"),
		RepetitionAndSequence: getBool("repetition_and_sequence"),
		UserInfo:              getBool("user_info"),
		Words:                 getStringList("words"),
	}

	return clientRejects, diags
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
