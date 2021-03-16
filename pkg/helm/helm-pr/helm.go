package helmpr

import (
	"net/http"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func New() (*HelmClient, error) {
	return &HelmClient{
		settings: cli.New(),
	}, nil
}

func debug(format string, v ...interface{}) {}

func (hc *HelmClient) Install(args []string) (*release.Release, error) {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(hc.settings.RESTClientGetter(), hc.settings.Namespace(), os.Getenv("HELM_DRIVER"), debug)
	if err != nil {
		return nil, err
	}
	valueOpts := &values.Options{}

	install := action.NewInstall(actionConfig)
	name, appVersion, err := install.NameAndChart(args)
	if err != nil {
		return nil, err
	}
	install.ReleaseName = name
	install.Namespace = hc.settings.Namespace()

	rc, err := repo.NewChartRepository(&repo.Entry{
		Name:     os.Getenv("HELM_REPOSITORY_NAME"),
		URL:      os.Getenv("HELM_REPOSITORY_URL"),
		Username: os.Getenv(("HELM_REPOSITORY_USERNAME")),
		Password: os.Getenv("HELM_REPOSITORY_PASSWORD"),
	}, getter.All(hc.settings))

	if err != nil {
		return nil, err
	}
	indexPath, err := rc.DownloadIndexFile()
	if err != nil {
		return nil, err
	}
	index, err := repo.LoadIndexFile(indexPath)
	if err != nil {
		return nil, err
	}
	cv, err := index.Get(name, appVersion)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", cv.URLs[0], nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(
		os.Getenv("HELM_REPOSITORY_USERNAME"),
		os.Getenv("HELM_REPOSITORY_PASSWORD"),
	)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	chart, err := loader.LoadArchive(resp.Body)
	if err != nil {
		return nil, err
	}
	p := getter.All(hc.settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	return install.Run(chart, vals)
}

type HelmClient struct {
	settings *cli.EnvSettings
}
