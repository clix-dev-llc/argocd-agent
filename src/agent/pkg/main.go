package main

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/store"
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

	token := argo.GetToken(argoUsername, argoPassword, argoHost)
	store.SetArgo(token, argoHost)

	store.SetCodefresh(codefreshHost, codefreshToken)

	argo.Watch()
}
