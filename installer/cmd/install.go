package cmd

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os/user"
	"path"
)

var installCmdOptions struct {
	kube struct {
		namespace    string
		inCluster    bool
		context      string
		nodeSelector string
	}
	Argo struct {
		Host     string
		Username string
		Password string
	}
	Codefresh struct {
		Host        string
		Token       string
		Integration string
	}
}

func sendPrompt(msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	return result == "Y" || result == "y"
}

func ensureIntegration() error {
	err := holder.ApiHolder.CreateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, false)
	if err == nil {
		return nil
	}

	codefreshErr, ok := err.(*codefresh.CodefreshError)
	if !ok {
		return err
	}

	if codefreshErr.Status != 409 {
		return codefreshErr
	}

	needDelete := sendPrompt("You already have integration with this name or host, do you want to update it")
	if !needDelete {
		return fmt.Errorf("you should delete integration")
	}

	errEnsure := holder.ApiHolder.CreateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, true)

	if errEnsure != nil {
		return errEnsure
	}

	return nil
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install agent",
	Long:  `Install agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if installCmdOptions.Codefresh.Token == "" || installCmdOptions.Codefresh.Host == "" {
			config, err := cliconfig.GetCurrentConfig()
			if err != nil {
				return err
			}
			installCmdOptions.Codefresh.Token = config.Token
			installCmdOptions.Codefresh.Host = config.Url

		}

		holder.ApiHolder = codefresh.Api{
			Token:       installCmdOptions.Codefresh.Token,
			Host:        installCmdOptions.Codefresh.Host,
			Integration: installCmdOptions.Codefresh.Integration,
		}

		err := ensureIntegration()
		if err != nil {
			return err
		}

		fmt.Println("Integration updated")

		var kubeConfigPath string
		currentUser, _ := user.Current()
		if currentUser != nil {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}

		kubeOptions := installCmdOptions.kube

		if kubeOptions.context == "" {
			contexts, err := kube.GetAllContexts(kubeConfigPath)
			if err != nil {
				return err
			}

			prompt := promptui.Select{
				Label: "Select Kubernetes context",
				Items: contexts,
			}
			_, selectedContext, err := prompt.Run()
			kubeOptions.context = selectedContext
		}

		if kubeOptions.namespace == "" {
			prompt := promptui.Prompt{
				Label: "Kubernetes namespace to install",
			}

			kubeOptions.namespace, err = prompt.Run()

			if err != nil {
				return err
			}
		}

		cs, err := kube.ClientBuilder(kubeOptions.context, kubeOptions.namespace, kubeConfigPath, kubeOptions.inCluster).BuildClient()

		if err != nil {
			return err
		}

		installOptions := templates.InstallOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(installCmdOptions),
			Namespace:      kubeOptions.namespace,
			KubeClientSet:  cs,
		}

		err = templates.Install(&installOptions)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installCmdOptions.Argo.Host, "argo-host", "https://34.71.103.174", "")
	installCmd.Flags().StringVar(&installCmdOptions.Argo.Username, "argo-username", "admin", "")
	installCmd.Flags().StringVar(&installCmdOptions.Argo.Password, "argo-password", "newpassword", "")

	installCmd.Flags().StringVar(&installCmdOptions.Codefresh.Host, "codefresh-host", "", "")
	installCmd.Flags().StringVar(&installCmdOptions.Codefresh.Token, "codefresh-token", "", "")
	installCmd.Flags().StringVar(&installCmdOptions.Codefresh.Integration, "codefresh-integration", "test-integration", "")

	installCmd.Flags().StringVar(&installCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be installed [$KUBE_NAMESPACE]")
	installCmd.Flags().StringVar(&installCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be installed (default is current-context) [$KUBE_CONTEXT]")
	installCmd.Flags().BoolVar(&installCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if Argo agent is been installed from inside a cluster")

}
