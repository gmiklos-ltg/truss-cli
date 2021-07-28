package cmd

import (
	"os"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"

	"github.com/homeport/dyff/pkg/dyff"
	"github.com/gonvenience/ytbx"
)

var secretsViewCmd = &cobra.Command{
	Use:   "view [name] [kubeconfig]",
	Short: "Views a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "view")
		if err != nil {
			return err
		}

		_, err = secretCompare(sm, secret, true)
		return err
	},
}

// return true if same
func secretCompare(sm *truss.SecretsManager, secret truss.SecretConfig, localToRemote bool) (bool, error) {
	localContent, remoteContent, err := sm.View(secret)
	if err != nil {
		return false, err
	}

	localDocs, err := ytbx.LoadYAMLDocuments([]byte(localContent))
	if err != nil {
		return false, err
	}

	remoteDocs, err := ytbx.LoadYAMLDocuments([]byte(remoteContent))
	if err != nil {
		return false, err
	}

	localFile := ytbx.InputFile {
		Location: "local-changes",
		Documents: localDocs,
	}

	remoteFile := ytbx.InputFile {
		Location: "remote-changes",
		Documents: remoteDocs,
	}

	report, err := dyff.CompareInputFiles(remoteFile, localFile)

	reporter := dyff.HumanReport{
		Report:            report,
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		OmitHeader:        true,
	}

	println(len(report.Diffs))

	if err = reporter.WriteReport(os.Stdout); err != nil {
		return false, err
	}

	return remoteContent == localContent, nil
}
