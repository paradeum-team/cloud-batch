package models

type Rnode struct {
	IP       string `yaml:"ip,omitempty" json:"ip,omitempty"`
	PublicIP string `yaml:"publicIP,omitempty" json:"publicIP,omitempty"`
	ID       int    `yaml:"id,omitempty" json:"id,omitempty"`
}

type BfsValues struct {
	ImagePullSecrets []map[string]string `yaml:"imagePullSecrets,omitempty" json:"imagePullSecrets,omitempty"`
	Cluster          struct {
		Name                string `yaml:"name,omitempty" json:"name,omitempty"`
		RnodeNonInteraction bool   `yaml:"rnodeNonInteraction,omitempty" json:"rnodeNonInteraction,omitempty"`
		PublicNetwork       bool   `yaml:"publicNetwork,omitempty" json:"publicNetwork,omitempty"`
		Test                struct {
			ActiveDeadlineSeconds int `yaml:"activeDeadlineSeconds,omitempty" json:"activeDeadlineSeconds,omitempty"`
		} `yaml:"test,omitempty" json:"test,omitempty"`
		Rnodes []*Rnode `yaml:"rnodes,omitempty" json:"rnodes,omitempty"`
	} `yaml:"cluster,omitempty"`
	Monitor struct {
		Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
		App       string `yaml:"app,omitempty" json:"app,omitempty"`
		Release   string `yaml:"release,omitempty" json:"release,omitempty"`
		Grafana   struct {
			Enabled                  bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			DefaultDashboardsEnabled bool `yaml:"defaultDashboardsEnabled,omitempty" json:"defaultDashboardsEnabled,omitempty"`
			Sidecar                  struct {
				Dashboards struct {
					Label string `yaml:"label,omitempty" json:"label,omitempty"`
				} `yaml:"dashboards,omitempty" json:"dashboards,omitempty"`
			} `yaml:"sidecar,omitempty" json:"sidecar,omitempty"`
		} `yaml:"grafana,omitempty" json:"grafana,omitempty"`
	} `yaml:"monitor,omitempty" json:"monitor,omitempty"`
	Tolerations struct {
	} `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	RnodeHttpd struct {
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Affinity struct {
		} `yaml:"affinity,omitempty" json:"affinity,omitempty"`
	} `yaml:"rnodeHttpd,omitempty" json:"rnodeHttpd,omitempty"`
	RnodeServer struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
	} `yaml:"rnodeServer,omitempty" json:"rnodeServer,omitempty"`
	RnodePkg struct {
		Arfs struct {
			Copies string `yaml:"copies,omitempty" json:"copies,omitempty"`
		} `yaml:"arfs,omitempty" json:"arfs,omitempty"`
		Afs struct {
			Copies string `yaml:"copies,omitempty" json:"copies,omitempty"`
		} `yaml:"afs,omitempty" json:"afs,omitempty"`
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
	} `yaml:"rnodePkg,omitempty" json:"rnodePkg,omitempty"`
	RnodeExporter struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Logger struct {
			Level string `yaml:"level,omitempty" json:"level,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		PodMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"podMonitor,omitempty" json:"podMonitor,omitempty"`
		PrometheusRule struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
		} `yaml:"prometheusRule,omitempty" json:"prometheusRule,omitempty"`
	} `yaml:"rnodeExporter,omitempty" json:"rnodeExporter,omitempty"`
	RnodeAPI struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		Logger struct {
			Level     string `yaml:"level,omitempty" json:"level,omitempty"`
			Maxageday string `yaml:"maxageday,omitempty" json:"maxageday,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Metrics struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
		ServiceMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"serviceMonitor,omitempty" json:"serviceMonitor,omitempty"`
	} `yaml:"rnodeApi,omitempty" json:"rnodeApi,omitempty"`
	TnodeAPI struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		Logger struct {
			Level     string `yaml:"level,omitempty" json:"level,omitempty"`
			Maxageday string `yaml:"maxageday,omitempty" json:"maxageday,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Metrics struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
		ServiceMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"serviceMonitor,omitempty" json:"serviceMonitor,omitempty"`
	} `yaml:"tnodeApi,omitempty" json:"tnodeApi,omitempty"`
	PnodeAPI struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty"`
		Logger struct {
			Level     string `yaml:"level,omitempty" json:"level,omitempty"`
			Maxageday string `yaml:"maxageday,omitempty" json:"maxageday,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Metrics struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
		ServiceMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"serviceMonitor,omitempty" json:"serviceMonitor,omitempty"`
		Test struct {
			Image struct {
				Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
				Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
			} `yaml:"image,omitempty" json:"image,omitempty"`
		} `yaml:"test,omitempty" json:"test,omitempty"`
	} `yaml:"pnodeApi,omitempty" json:"pnodeApi,omitempty"`
	DataDetection struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		Logger struct {
			Level string `yaml:"level,omitempty" json:"level,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Metrics struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
		ServiceMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"serviceMonitor,omitempty" json:"serviceMonitor,omitempty"`
	} `yaml:"dataDetection,omitempty" json:"dataDetection,omitempty"`
	DataDetectionProxy struct {
		Image struct {
			Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
			Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		} `yaml:"image,omitempty" json:"image,omitempty"`
		Listen struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"listen,omitempty" json:"listen,omitempty"`
		Logger struct {
			Level string `yaml:"level,omitempty" json:"level,omitempty"`
		} `yaml:"logger,omitempty" json:"logger,omitempty"`
		Metrics struct {
			Port string `yaml:"port,omitempty" json:"port,omitempty"`
		} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
		Storage        string      `yaml:"storage,omitempty" json:"storage,omitempty"`
		StorageClass   interface{} `yaml:"storageClass,omitempty" json:"storageClass,omitempty"`
		ServiceMonitor struct {
			Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
			Labels  struct {
				Release string `yaml:"release,omitempty" json:"release,omitempty"`
			} `yaml:"labels,omitempty" json:"labels,omitempty"`
		} `yaml:"serviceMonitor,omitempty" json:"serviceMonitor,omitempty"`
	} `yaml:"dataDetectionProxy,omitempty" json:"dataDetectionProxy,omitempty"`
}
