// Copyright © 2021 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cni

import "strings"

type Cilium struct {
	metadata MetaData
}

func (c Cilium) Manifests(template string) string {
	if template == "" {
		template = c.Template()
	}
	if c.metadata.CniRepo == "" || c.metadata.CniRepo == "k8s.gcr.io" {
		c.metadata.CniRepo = "quay.io/cilium"
	}
	//if c.metadata.Interface == "" {
	//	c.metadata.Interface = "interface=" + defaultInterface
	//}
	if c.metadata.CIDR == "" {
		c.metadata.CIDR = defaultCIDR
	}

	if c.metadata.K8sServiceHost == "" {
		c.metadata.K8sServiceHost = defaultK8sServiceHost
	}
	if c.metadata.K8sServicePort == "" {
		c.metadata.K8sServicePort = defaultK8sServicePort
	}
	c.metadata.K8sServiceHost = strings.Split(c.metadata.K8sServiceHost, ":")[0]

	return render(c.metadata, template)
}

func (c Cilium) Template() string {
	return CiliumManifests
}

// docs https://docs.cilium.io/en/stable/gettingstarted/minikube/
// quick install kubectl create -f https://raw.githubusercontent.com/cilium/cilium/1.9.3/install/kubernetes/quick-install.yaml
// experimental kubectl create -f https://raw.githubusercontent.com/cilium/cilium/1.9.3/install/kubernetes/experimental-install.yaml
// cilium kernel-check kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/1.9.3/examples/kubernetes/kernel-check/kernel-check.yaml
const CiliumManifests = `
---
# Source: cilium/templates/cilium-agent-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cilium
  namespace: kube-system
---
# Source: cilium/templates/cilium-operator-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cilium-operator
  namespace: kube-system
---
# Source: cilium/templates/cilium-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cilium-config
  namespace: kube-system
data:

  # Identity allocation mode selects how identities are shared between cilium
  # nodes by setting how they are stored. The options are "crd" or "kvstore".
  # - "crd" stores identities in kubernetes as CRDs (custom resource definition).
  #   These can be queried with:
  #     kubectl get ciliumid
  # - "kvstore" stores identities in a kvstore, etcd or consul, that is
  #   configured below. Cilium versions before 1.6 supported only the kvstore
  #   backend. Upgrades from these older cilium versions should continue using
  #   the kvstore by commenting out the identity-allocation-mode below, or
  #   setting it to "kvstore".
  identity-allocation-mode: crd
  cilium-endpoint-gc-interval: "5m0s"

  # If you want to run cilium in debug mode change this value to true
  debug: "false"
  # The agent can be put into the following three policy enforcement modes
  # default, always and never.
  # https://docs.cilium.io/en/latest/policy/intro/#policy-enforcement-modes
  enable-policy: "default"

  # Enable IPv4 addressing. If enabled, all endpoints are allocated an IPv4
  # address.
  enable-ipv4: "true"

  # Enable IPv6 addressing. If enabled, all endpoints are allocated an IPv6
  # address.
  enable-ipv6: "false"

  # -- Clean all eBPF datapath state from the initContainer of the cilium-agent
  # DaemonSet.
  #
  # WARNING: Use with care!
  clean-cilium-bpf-state: "false"

  # -- Clean all local Cilium state from the initContainer of the cilium-agent
  # DaemonSet. Implies cleanBpfState: true.
  #
  # WARNING: Use with care!
  clean-cilium-state: "false"

  # Users who wish to specify their own custom CNI configuration file must set
  # custom-cni-conf to "true", otherwise Cilium may overwrite the configuration.
  custom-cni-conf: "false"
  enable-bpf-clock-probe: "true"
  # If you want cilium monitor to aggregate tracing for packets, set this level
  # to "low", "medium", or "maximum". The higher the level, the less packets
  # that will be seen in monitor output.
  monitor-aggregation: medium

  # The monitor aggregation interval governs the typical time between monitor
  # notification events for each allowed connection.
  #
  # Only effective when monitor aggregation is set to "medium" or higher.
  monitor-aggregation-interval: 5s

  # The monitor aggregation flags determine which TCP flags which, upon the
  # first observation, cause monitor notifications to be generated.
  #
  # Only effective when monitor aggregation is set to "medium" or higher.
  monitor-aggregation-flags: all
  # Specifies the ratio (0.0-1.0) of total system memory to use for dynamic
  # sizing of the TCP CT, non-TCP CT, NAT and policy BPF maps.
  bpf-map-dynamic-size-ratio: "0.0025"
  # bpf-policy-map-max specifies the maximum number of entries in endpoint
  # policy map (per endpoint)
  bpf-policy-map-max: "16384"
  # bpf-lb-map-max specifies the maximum number of entries in bpf lb service,
  # backend and affinity maps.
  bpf-lb-map-max: "65536"
  # Pre-allocation of map entries allows per-packet latency to be reduced, at
  # the expense of up-front memory allocation for the entries in the maps. The
  # default value below will minimize memory usage in the default installation;
  # users who are sensitive to latency may consider setting this to "true".
  #
  # This option was introduced in Cilium 1.4. Cilium 1.3 and earlier ignore
  # this option and behave as though it is set to "true".
  #
  # If this value is modified, then during the next Cilium startup the restore
  # of existing endpoints and tracking of ongoing connections may be disrupted.
  # As a result, reply packets may be dropped and the load-balancing decisions
  # for established connections may change.
  #
  # If this option is set to "false" during an upgrade from 1.3 or earlier to
  # 1.4 or later, then it may cause one-time disruptions during the upgrade.
  preallocate-bpf-maps: "true"

  # Regular expression matching compatible Istio sidecar istio-proxy
  # container image names
  sidecar-istio-proxy-image: "cilium/istio_proxy"

  # Name of the cluster. Only relevant when building a mesh of clusters.
  cluster-name: default
  # Unique ID of the cluster. Must be unique across all conneted clusters and
  # in the range of 1 and 255. Only relevant when building a mesh of clusters.
  cluster-id: ""

  # Encapsulation mode for communication between nodes
  # Possible values:
  #   - disabled
  #   - vxlan (default)
  #   - geneve
  tunnel: vxlan
  # Enables L7 proxy for L7 policy enforcement and visibility
  enable-l7-proxy: "true"

  # wait-bpf-mount makes init container wait until bpf filesystem is mounted
  wait-bpf-mount: "false"

  masquerade: "true"
  enable-bpf-masquerade: "true"

  enable-xt-socket-fallback: "true"
  install-iptables-rules: "true"

  auto-direct-node-routes: "true"
  enable-bandwidth-manager: "false"
  enable-local-redirect-policy: "false"

  kube-proxy-replacement:  "strict"
  kube-proxy-replacement-healthz-bind-address: ""
  enable-health-check-nodeport: "true"
  node-port-bind-protection: "true"
  enable-auto-protect-node-port-range: "true"
  enable-session-affinity: "true"
  enable-endpoint-health-checking: "true"
  enable-health-checking: "true"
  enable-well-known-identities: "false"
  enable-remote-node-identity: "true"
  operator-api-serve-addr: "127.0.0.1:9234"
  # Enable Hubble gRPC service.
  enable-hubble: "true"
  # UNIX domain socket for Hubble server to listen to.
  hubble-socket-path:  "/var/run/cilium/hubble.sock"
  # An additional address for Hubble server to listen to (e.g. ":4244").
  hubble-listen-address: ":4244"
  hubble-disable-tls: "false"
  hubble-tls-cert-file: /var/lib/cilium/tls/hubble/server.crt
  hubble-tls-key-file: /var/lib/cilium/tls/hubble/server.key
  hubble-tls-client-ca-files: /var/lib/cilium/tls/hubble/client-ca.crt
  ipam: "cluster-pool"
  cluster-pool-ipv4-cidr: "{{ .CIDR }}"
  cluster-pool-ipv4-mask-size: "24"
  disable-cnp-status-updates: "true"
---
# Source: cilium/templates/cilium-agent-clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cilium
rules:
- apiGroups:
  - networking.k8s.io
  resources:
  - networkpolicies
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - namespaces
  - services
  - nodes
  - endpoints
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - pods
  - pods/finalizers
  verbs:
  - get
  - list
  - watch
  - update
  - delete
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
  - list
  - watch
  - update
- apiGroups:
  - ""
  resources:
  - nodes
  - nodes/status
  verbs:
  - patch
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  # Deprecated for removal in v1.10
  - create
  - list
  - watch
  - update

  # This is used when validating policies in preflight. This will need to stay
  # until we figure out how to avoid "get" inside the preflight, and then
  # should be removed ideally.
  - get
- apiGroups:
  - cilium.io
  resources:
  - ciliumnetworkpolicies
  - ciliumnetworkpolicies/status
  - ciliumnetworkpolicies/finalizers
  - ciliumclusterwidenetworkpolicies
  - ciliumclusterwidenetworkpolicies/status
  - ciliumclusterwidenetworkpolicies/finalizers
  - ciliumendpoints
  - ciliumendpoints/status
  - ciliumendpoints/finalizers
  - ciliumnodes
  - ciliumnodes/status
  - ciliumnodes/finalizers
  - ciliumidentities
  - ciliumidentities/finalizers
  - ciliumlocalredirectpolicies
  - ciliumlocalredirectpolicies/status
  - ciliumlocalredirectpolicies/finalizers
  verbs:
  - "*"
---
# Source: cilium/templates/cilium-operator-clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cilium-operator
rules:
- apiGroups:
  - ""
  resources:
  # to automatically delete [core|kube]dns pods so that are starting to being
  # managed by Cilium
  - pods
  verbs:
  - get
  - list
  - watch
  - delete
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  # to perform the translation of a CNP that contains ToGroup to its endpoints
  - services
  - endpoints
  # to check apiserver connectivity
  - namespaces
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - cilium.io
  resources:
  - ciliumnetworkpolicies
  - ciliumnetworkpolicies/status
  - ciliumnetworkpolicies/finalizers
  - ciliumclusterwidenetworkpolicies
  - ciliumclusterwidenetworkpolicies/status
  - ciliumclusterwidenetworkpolicies/finalizers
  - ciliumendpoints
  - ciliumendpoints/status
  - ciliumendpoints/finalizers
  - ciliumnodes
  - ciliumnodes/status
  - ciliumnodes/finalizers
  - ciliumidentities
  - ciliumidentities/status
  - ciliumidentities/finalizers
  - ciliumlocalredirectpolicies
  - ciliumlocalredirectpolicies/status
  - ciliumlocalredirectpolicies/finalizers
  verbs:
  - "*"
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - create
  - get
  - list
  - update
  - watch
# For cilium-operator running in HA mode.
#
# Cilium operator running in HA mode requires the use of ResourceLock for Leader Election
# between multiple running instances.
# The preferred way of doing this is to use LeasesResourceLock as edits to Leases are less
# common and fewer objects in the cluster watch "all Leases".
# The support for leases was introduced in coordination.k8s.io/v1 during Kubernetes 1.14 release.
# In Cilium we currently dont support HA mode for K8s version < 1.14. This condition make sure
# that we only authorize access to leases resources in supported K8s versions.
- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - create
  - get
  - update
---
# Source: cilium/templates/cilium-agent-clusterrolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cilium
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cilium
subjects:
- kind: ServiceAccount
  name: cilium
  namespace: kube-system
---
# Source: cilium/templates/cilium-operator-clusterrolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cilium-operator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cilium-operator
subjects:
- kind: ServiceAccount
  name: cilium-operator
  namespace: kube-system
---
# Source: cilium/templates/cilium-agent-daemonset.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: cilium
  name: cilium
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: cilium
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 2
    type: RollingUpdate
  template:
    metadata:
      annotations:
        # This annotation plus the CriticalAddonsOnly toleration makes
        # cilium to be a critical pod in the cluster, which ensures cilium
        # gets priority scheduling.
        # https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
        scheduler.alpha.kubernetes.io/critical-pod: ""
      labels:
        k8s-app: cilium
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: k8s-app
                operator: In
                values:
                - cilium
            topologyKey: kubernetes.io/hostname
      containers:
      - args:
        - --config-dir=/tmp/cilium/config-map
        command:
        - cilium-agent
        livenessProbe:
          httpGet:
            host: 127.0.0.1
            path: /healthz
            port: 9876
            scheme: HTTP
            httpHeaders:
            - name: "brief"
              value: "true"
          failureThreshold: 10
          # The initial delay for the liveness probe is intentionally large to
          # avoid an endless kill & restart cycle if in the event that the initial
          # bootstrapping takes longer than expected.
          initialDelaySeconds: 120
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            host: 127.0.0.1
            path: /healthz
            port: 9876
            scheme: HTTP
            httpHeaders:
            - name: "brief"
              value: "true"
          failureThreshold: 3
          initialDelaySeconds: 5
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 5
        env:
        - name: K8S_NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: CILIUM_K8S_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        - name: CILIUM_FLANNEL_MASTER_DEVICE
          valueFrom:
            configMapKeyRef:
              key: flannel-master-device
              name: cilium-config
              optional: true
        - name: CILIUM_FLANNEL_UNINSTALL_ON_EXIT
          valueFrom:
            configMapKeyRef:
              key: flannel-uninstall-on-exit
              name: cilium-config
              optional: true
        - name: CILIUM_CLUSTERMESH_CONFIG
          value: /var/lib/cilium/clustermesh/
        - name: CILIUM_CNI_CHAINING_MODE
          valueFrom:
            configMapKeyRef:
              key: cni-chaining-mode
              name: cilium-config
              optional: true
        - name: CILIUM_CUSTOM_CNI_CONF
          valueFrom:
            configMapKeyRef:
              key: custom-cni-conf
              name: cilium-config
              optional: true
{{- if .K8sServiceHost }}
        - name: KUBERNETES_SERVICE_HOST
          value: "{{ .K8sServiceHost  }}"
{{- end }}
{{- if .K8sServicePort }}
        - name: KUBERNETES_SERVICE_PORT
          value: "{{ .K8sServicePort  }}"
{{- end }}
        image: {{ .CniRepo }}/cilium:v1.9.3
        imagePullPolicy: IfNotPresent
        lifecycle:
          postStart:
            exec:
              command:
              - "/cni-install.sh"
              - "--enable-debug=false"
          preStop:
            exec:
              command:
              - /cni-uninstall.sh
        name: cilium-agent
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
            - SYS_MODULE
          privileged: true
        volumeMounts:
        - mountPath: /sys/fs/bpf
          name: bpf-maps
        - mountPath: /var/run/cilium
          name: cilium-run
        - mountPath: /host/opt/cni/bin
          name: cni-path
        - mountPath: /host/etc/cni/net.d
          name: etc-cni-netd
        - mountPath: /var/lib/cilium/clustermesh
          name: clustermesh-secrets
          readOnly: true
        - mountPath: /tmp/cilium/config-map
          name: cilium-config-path
          readOnly: true
          # Needed to be able to load kernel modules
        - mountPath: /lib/modules
          name: lib-modules
          readOnly: true
        - mountPath: /run/xtables.lock
          name: xtables-lock
        - mountPath: /var/lib/cilium/tls/hubble
          name: hubble-tls
          readOnly: true
      hostNetwork: true
      initContainers:
      - command:
        - /init-container.sh
        env:
        - name: CILIUM_ALL_STATE
          valueFrom:
            configMapKeyRef:
              key: clean-cilium-state
              name: cilium-config
              optional: true
        - name: CILIUM_BPF_STATE
          valueFrom:
            configMapKeyRef:
              key: clean-cilium-bpf-state
              name: cilium-config
              optional: true
        - name: CILIUM_WAIT_BPF_MOUNT
          valueFrom:
            configMapKeyRef:
              key: wait-bpf-mount
              name: cilium-config
              optional: true
{{- if .K8sServiceHost }}
        - name: KUBERNETES_SERVICE_HOST
          value: "{{ .K8sServiceHost }}"
{{- end }}
{{- if .K8sServicePort }}
        - name: KUBERNETES_SERVICE_PORT
          value: "{{ .K8sServicePort }}"
{{- end }}
        image: {{ .CniRepo }}/cilium:v1.9.3
        imagePullPolicy: IfNotPresent
        name: clean-cilium-state
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: true
        volumeMounts:
        - mountPath: /sys/fs/bpf
          name: bpf-maps
          mountPropagation: HostToContainer
        - mountPath: /var/run/cilium
          name: cilium-run
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
      restartPolicy: Always
      priorityClassName: system-node-critical
      serviceAccount: cilium
      serviceAccountName: cilium
      terminationGracePeriodSeconds: 1
      tolerations:
      - operator: Exists
      volumes:
        # To keep state between restarts / upgrades
      - hostPath:
          path: /var/run/cilium
          type: DirectoryOrCreate
        name: cilium-run
        # To keep state between restarts / upgrades for bpf maps
      - hostPath:
          path: /sys/fs/bpf
          type: DirectoryOrCreate
        name: bpf-maps
      # To install cilium cni plugin in the host
      - hostPath:
          path:  /opt/cni/bin
          type: DirectoryOrCreate
        name: cni-path
        # To install cilium cni configuration in the host
      - hostPath:
          path: /etc/cni/net.d
          type: DirectoryOrCreate
        name: etc-cni-netd
        # To be able to load kernel modules
      - hostPath:
          path: /lib/modules
        name: lib-modules
        # To access iptables concurrently with other processes (e.g. kube-proxy)
      - hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
        name: xtables-lock
        # To read the clustermesh configuration
      - name: clustermesh-secrets
        secret:
          defaultMode: 420
          optional: true
          secretName: cilium-clustermesh
        # To read the configuration from the config map
      - configMap:
          name: cilium-config
        name: cilium-config-path
      - name: hubble-tls
        projected:
          sources:
          - secret:
              name: hubble-server-certs
              items:
                - key: tls.crt
                  path: server.crt
                - key: tls.key
                  path: server.key
              optional: true
          - configMap:
              name: hubble-ca-cert
              items:
                - key: ca.crt
                  path: client-ca.crt
              optional: true
---
# Source: cilium/templates/cilium-operator-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    io.cilium/app: operator
    name: cilium-operator
  name: cilium-operator
  namespace: kube-system
spec:
  # We support HA mode only for Kubernetes version > 1.14
  # See docs on ServerCapabilities.LeasesResourceLock in file pkg/k8s/version/version.go
  # for more details.
  replicas: 1
  selector:
    matchLabels:
      io.cilium/app: operator
      name: cilium-operator
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      annotations:
      labels:
        io.cilium/app: operator
        name: cilium-operator
    spec:
      # In HA mode, cilium-operator pods must not be scheduled on the same
      # node as they will clash with each other.
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: io.cilium/app
                operator: In
                values:
                - operator
            topologyKey: kubernetes.io/hostname
      containers:
      - args:
        - --config-dir=/tmp/cilium/config-map
        - --debug=$(CILIUM_DEBUG)
        command:
        - cilium-operator-generic
        env:
        - name: K8S_NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: CILIUM_K8S_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        - name: CILIUM_DEBUG
          valueFrom:
            configMapKeyRef:
              key: debug
              name: cilium-config
              optional: true
{{- if .K8sServiceHost }}
        - name: KUBERNETES_SERVICE_HOST
          value: "{{ .K8sServiceHost }}"
{{- end }}
{{- if .K8sServicePort }}
        - name: KUBERNETES_SERVICE_PORT
          value: "{{ .K8sServicePort }}"
{{- end }}
        image: {{ .CniRepo }}/operator-generic:v1.9.3
        imagePullPolicy: IfNotPresent
        name: cilium-operator
        livenessProbe:
          httpGet:
            host: 127.0.0.1
            path: /healthz
            port: 9234
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 3
        volumeMounts:
        - mountPath: /tmp/cilium/config-map
          name: cilium-config-path
          readOnly: true
      hostNetwork: true
      restartPolicy: Always
      priorityClassName: system-cluster-critical
      serviceAccount: cilium-operator
      serviceAccountName: cilium-operator
      tolerations:
      - operator: Exists
      volumes:
        # To read the configuration from the config map
      - configMap:
          name: cilium-config
        name: cilium-config-path
`
