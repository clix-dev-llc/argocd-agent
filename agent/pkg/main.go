package main

import (
	"errors"
	argo2 "github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	store2 "github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"os"
)

func main() {

	argoHost, argoHostExistence := os.LookupEnv("ARGO_HOST")
	if !argoHostExistence {
		panic(errors.New("ARGO_HOST variable doesnt exist"))
	}

	argoUsername, argoUsernameExistence := os.LookupEnv("ARGO_USERNAME")
	if !argoUsernameExistence {
		panic(errors.New("ARGO_USERNAME variable doesnt exist"))
	}

	argoPassword, argoPasswordExistence := os.LookupEnv("ARGO_PASSWORD")
	if !argoPasswordExistence {
		panic(errors.New("ARGO_PASSWORD variable doesnt exist"))
	}

	codefreshToken, codefreshTokenExistence := os.LookupEnv("CODEFRESH_TOKEN")
	if !codefreshTokenExistence {
		panic(errors.New("CODEFRESH_TOKEN variable doesnt exist"))
	}

	codefreshHost, codefreshHostExistance := os.LookupEnv("CODEFRESH_HOST")
	if !codefreshHostExistance {
		codefreshHost = "https://g.codefresh.io"
	}

	codefreshIntegrationName, codefreshIntegrationNameExistence := os.LookupEnv("CODEFRESH_INTEGRATION")
	if !codefreshIntegrationNameExistence {
		panic(errors.New("CODEFRESH_INTEGRATION variable doesnt exist"))
	}

	token := argo2.GetToken(argoUsername, argoPassword, argoHost)

	store2.SetArgo(token, argoHost)

	store2.SetCodefresh(codefreshHost, codefreshToken, codefreshIntegrationName)

	err := codefresh.EnsureIntegration(codefreshIntegrationName, argoHost, argoUsername, argoPassword)

	if err != nil {
		panic(err)
	}

	argo2.Watch()
}
