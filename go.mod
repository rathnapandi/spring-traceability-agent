module github.com/rathnapandi/spring-traceability-agent

go 1.15

replace git.ecd.axway.org/apigov/apic_agents_sdk => /Users/rnatarajan/go/src/git.ecd.axway.org/apigov/apic_agents_sdk

require (
	4d63.com/embedfiles v1.0.0 // indirect
	git.ecd.axway.org/apigov/apic_agents_sdk v0.0.0-00010101000000-000000000000
	github.com/Shopify/sarama v1.26.4 // indirect
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/dop251/goja v0.0.0-20200526165454-f1752421c432 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/elastic/beats/v7 v7.7.1
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/klauspost/compress v1.10.9 // indirect
	github.com/miekg/dns v1.1.29 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/spf13/afero v1.3.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.5.0 // indirect
	k8s.io/client-go v0.17.0 // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.2.0+incompatible
	github.com/Shopify/sarama => github.com/elastic/sarama v0.0.0-20191122160421-355d120d0970
	github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20190717161051-705d9623b7c1+incompatible
	github.com/docker/go-plugins-helpers => github.com/elastic/go-plugins-helpers v0.0.0-20200207104224-bdf17607b79f
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
	github.com/fsnotify/fsnotify => github.com/adriansr/fsnotify v0.0.0-20180417234312-c9bbe1f46f1d
	github.com/google/gopacket => github.com/adriansr/gopacket v1.1.18-0.20200327165309-dd62abfa8a41
	github.com/insomniacslk/dhcp => github.com/elastic/dhcp v0.0.0-20200227161230-57ec251c7eb3 // indirect
	github.com/tonistiigi/fifo => github.com/containerd/fifo v0.0.0-20190816180239-bda0ff6ed73c
)
