package provider_logto

/*
TODO: Voir avec Rémi. J'ai mis ce fichier sous le format test_provider.go et non provider_test.go car Go traite différement les fichier
finissants par _test et lorsque je voulais utiliser les variables ProviderConfig et TestAccProtoV6ProviderFactories dans le package
resource_application_test, il n'arrivait pas a y accéder.

Dans Go, les fichiers avec le suffixe _test.go sont traités différemment par le compilateur - ils ne
sont compilés et exécutés que pendant les tests. Si vous avez des variables/constantes définies dans
un fichier provider_test.go, elles ne seront disponibles que pour les tests dans le même package.
Pour rendre ces variables accessibles à d'autres packages, vous devez les déclarer dans un fichier
sans le suffixe _test.go ou dans un fichier de test dans le même package que celui qui les utilise.

Le problème vient du fait qu'on garde la structure de fichier d'origine fournit par le générateur Terraform qui utilise le SDK V2 qui
est une version déprécié. Si le fichier provider.go provider_test.go et les resources était dans le même dossier il n'y aurait pas de
problème.
*/

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Logto client is properly configured.
	// It is also possible to use the LOGTO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	ProviderConfig = `
provider "logto" {
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"logto": providerserver.NewProtocol6WithError(New("test")()),
	}
)
