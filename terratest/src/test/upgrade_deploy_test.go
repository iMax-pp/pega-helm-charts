package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	k8sbatch "k8s.io/api/batch/v1"
)

const pegaHelmChartPath = "../../../charts/pega"

var options = &helm.Options{
	SetValues: map[string]string{
		"global.actions.execute": "upgrade-deploy",
		"global.provider":        "k8s",
	},
}

// VerifyUpgradeDeployActionShouldNotRenderDeployments - Tests all the skipped templates for action upgrade-deploy. These templates not supposed to be rendered for upgrade-deploy action.
func VerifyUpgradeActionSkippedTemplates(t *testing.T) {
	t.Parallel()
	output := helm.RenderTemplate(t, options, pegaHelmChartPath, []string{
		"templates/pega-action-validate.yaml",
		"charts/installer/templates/pega-install-environment-config.yaml",
	})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)
	// assert that above templates are not rendered
	require.Empty(t, deployment)
}

// ValidateUpgradeJobs - Tests Upgrade jobs yaml rendered with the values as provided in default values.yaml for action upgrade-deploy
func ValidateUpgradeJobs(t *testing.T) {
	var installerJobObj k8sbatch.Job
	var installerSlice = ReturnJobSlices(t, pegaHelmChartPath, options)
	println(len(installerSlice))
	var expectedJob pegaJob
	for index, installerInfo := range installerSlice {
		if index >= 1 && index <= 3 {
			if index == 1 {
				expectedJob = pegaJob{"pega-pre-upgrade", []string{}, "pega-upgrade-environment-config"}
			} else if index == 2 {
				expectedJob = pegaJob{"pega-db-upgrade", []string{"wait-for-pre-dbupgrade"}, "pega-upgrade-environment-config"}
			} else if index == 3 {
				expectedJob = pegaJob{"pega-post-upgrade", []string{"wait-for-pegaupgrade", "wait-for-rolling-updates"}, "pega-upgrade-environment-config"}
			}

			helm.UnmarshalK8SYaml(t, installerInfo, &installerJobObj)
			VerifyPegaJob(t, options, &installerJobObj, expectedJob)
		}

	}
}

// TestUpgradeDeployActions - Test all objects deployed for upgrade-deploy action with the values as provided in default values.yaml
func TestUpgradeDeployActions(t *testing.T) {
	VerifyUpgradeActionSkippedTemplates(t)
	ValidateUpgradeJobs(t)
	VerifyUpgradeEnvConfig(t, options, pegaHelmChartPath)
	VerifyInstallerConfigMaps(t, options, pegaHelmChartPath)
	VerifyInstallerRoleBinding(t, options, pegaHelmChartPath)
	VerifyInstallerRole(t, options, pegaHelmChartPath)
	VerifyPegaStandardTierDeployment(t, pegaHelmChartPath, options, []string{"wait-for-pegaupgrade"})
}
