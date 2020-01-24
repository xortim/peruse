package cmd

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xortim/peruse/k8sclient"
	"go.uber.org/zap"
)

func newServCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serv",
		Short: "Serves an HTML table",
		Long:  `Serves an HTML table`,
		RunE:  servRun,
	}
	return cmd
}

func servRun(cmd *cobra.Command, arts []string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/healthz", HealthHandler)
	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}

// HealthHandler serves /healthz and always returns 200
func HealthHandler(w http.ResponseWriter, req *http.Request) {
	zap.S().Debugf("%s - %s", req.RemoteAddr, req.RequestURI)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`OK`))
	return
}

// HomeHandler serves /
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	k8s, err := k8sclient.NewClient("", viper.GetString("kubeconfig"))
	if err != nil {
		zap.S().Errorf("Unable to authenticate: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`500 - unable to authenticate\n`))
		return
	}

	dips, err := k8sclient.GetDeploymentIngressPaths(k8s, viper.GetString("namespace"))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = dips
	t := dips.NewTable()
	t.SetHTMLCSSClass("table table-hover table-sm")
	w.Write([]byte(`
	<!doctype html>
	<html lang="en">
		<head>
			<!-- Required meta tags -->
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	
			<!-- Bootstrap CSS -->
			<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	
			<title>Peruse Deployments</title>
		</head>
		<body>
		<script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
	`))

	w.Write([]byte(t.RenderHTML()))

	w.Write([]byte(`
		</body>
	</html>
	`))
	return
}
