defmodule TbgNodes.PermissionedEthereumNetworks.K8sWrapper do
  @moduledoc false

  @spec get_rancher_cluster_id! :: String.t()
  defp get_rancher_cluster_id! do
    Application.get_env(:tbg_nodes, :rancher)[:cluster_id] || ""
  end

  @spec get_rancher_project_id!() :: String.t()
  defp get_rancher_project_id! do
    Application.get_env(:tbg_nodes, :rancher)[:project_id] || ""
  end

  defp get_namespace_rancher_annotations do
    rancher_cluster_id = get_rancher_cluster_id!()
    rancher_project_id = get_rancher_project_id!()

    if rancher_cluster_id != "" && rancher_project_id != "" do
      %{
        "field.cattle.io/projectId" => "#{rancher_cluster_id}:#{rancher_project_id}"
      }
    else
      %{}
    end
  end

  defp get_namespace_annotations do
    %{}
    |> Map.merge(get_namespace_rancher_annotations())
  end

  defp get_namespace_rancher_labels do
    rancher_project_id = get_rancher_project_id!()

    if rancher_project_id != "" do
      %{
        "field.cattle.io/projectId" => "#{rancher_project_id}"
      }
    else
      %{}
    end
  end

  defp get_namespace_labels do
    %{}
    |> Map.merge(get_namespace_rancher_labels())
  end

  @spec create_namespace(%{name: String.t(), labels: map(), annotations: map()}, K8s.Conn.t()) ::
          K8s.Client.Runner.Base.result_t()
  def create_namespace(
        %{
          name: name,
          labels: labels,
          annotations: annotations
        } = _params,
        conn
      ) do
    annotations =
      annotations
      |> Map.merge(get_namespace_annotations())

    labels =
      labels
      |> Map.merge(get_namespace_labels())

    ns_k8s = %{
      "apiVersion" => "v1",
      "kind" => "Namespace",
      "metadata" => %{
        "name" => name,
        "labels" => labels,
        "annotations" => annotations
      }
    }

    create_ns = K8s.Client.create(ns_k8s)
    {:ok, _namespace} = K8s.Client.run(create_ns, conn)
  end

  @spec create_ingress(map(), K8s.Conn.t()) :: K8s.Client.Runner.Base.result_t()
  def create_ingress(
        %{
          "annotations" => %{} = annotations,
          "labels" => labels,
          "name" => name,
          "namespace" => namespace,
          "rules" => rules,
          "tls" => tls
        } = _params,
        conn
      ) do
    operation =
      K8s.Client.create(%{
        "apiVersion" => "networking.k8s.io/v1beta1",
        "kind" => "Ingress",
        "metadata" => %{
          "name" => name,
          "labels" => labels,
          "annotations" => annotations,
          "namespace" => namespace
        },
        "spec" => %{
          "tls" => tls,
          "rules" => rules
        }
      })

    {:ok, _ingress} = K8s.Client.run(operation, conn)
  end

  @spec create_cluster_ip_service(map(), K8s.Conn.t()) :: K8s.Client.Runner.Base.result_t()
  def create_cluster_ip_service(
        %{
          annotations: annotations,
          labels: labels,
          name: name,
          namespace: namespace,
          ports: ports,
          selector: selector
        } = _params,
        conn
      ) do
    operation =
      K8s.Client.create(%{
        "apiVersion" => "v1",
        "kind" => "Service",
        "metadata" => %{
          "name" => name,
          "labels" => labels,
          "namespace" => namespace,
          "annotations" => annotations
        },
        "spec" => %{
          "type" => "ClusterIP",
          "selector" => selector,
          "ports" => ports
        }
      })

    {:ok, _service} = K8s.Client.run(operation, conn)
  end

  @spec create_configmap(map(), K8s.Conn.t()) :: K8s.Client.Runner.Base.result_t()
  def create_configmap(
        %{
          annotations: %{} = annotations,
          data: %{} = data,
          labels: %{} = labels,
          name: name,
          namespace: namespace
        } = _params,
        conn
      ) do
    operation =
      K8s.Client.create(%{
        "apiVersion" => "v1",
        "kind" => "ConfigMap",
        "metadata" => %{
          "labels" => labels,
          "annotations" => annotations,
          "name" => name,
          "namespace" => namespace
        },
        "data" => data
      })

    {:ok, _configmap} = K8s.Client.run(operation, conn)
  end

  @spec create_secret(map(), K8s.Conn.t()) :: K8s.Client.Runner.Base.result_t()
  def create_secret(
        %{
          annotations: annotations,
          labels: labels,
          name: name,
          namespace: namespace,
          string_data: string_data
        } = _params,
        conn
      ) do
    operation =
      K8s.Client.create(%{
        "apiVersion" => "v1",
        "kind" => "Secret",
        "metadata" => %{
          "labels" => labels,
          "annotations" => annotations,
          "name" => name,
          "namespace" => namespace
        },
        "type" => "Opaque",
        "stringData" => string_data
      })

    {:ok, _secret} = K8s.Client.run(operation, conn)
  end

  @spec create_statefulset(map(), K8s.Conn.t()) :: K8s.Client.Runner.Base.result_t()
  def create_statefulset(
        %{
          pod_annotations: pod_annotations,
          args: args,
          command: command,
          env: env,
          image: image,
          labels: labels,
          liveness_probe: liveness_probe,
          name: name,
          init_containers: init_containers,
          namespace: namespace,
          ports: ports,
          resources: resources,
          volume_claim_templates: volume_claim_templates,
          volume_mounts: volume_mounts,
          volumes: volumes
        } = _params,
        conn
      ) do
    create_statefulset =
      K8s.Client.create(%{
        "apiVersion" => "apps/v1",
        "kind" => "StatefulSet",
        "metadata" => %{
          "name" => name,
          "labels" => labels,
          "namespace" => namespace
        },
        "spec" => %{
          "replicas" => 1,
          "selector" => %{
            "matchLabels" => labels
          },
          "volumeClaimTemplates" => volume_claim_templates,
          "template" => %{
            "metadata" => %{
              "labels" => labels,
              "annotations" => pod_annotations
            },
            "spec" => %{
              "securityContext" => %{
                "fsGroup" => 1000
              },
              "initContainers" => init_containers,
              "containers" => [
                %{
                  "name" => name,
                  "image" => image,
                  "imagePullPolicy" => "IfNotPresent",
                  "resources" => resources,
                  "volumeMounts" => volume_mounts,
                  "ports" => ports,
                  "env" => env,
                  "command" => command,
                  "args" => args,
                  "livenessProbe" => liveness_probe,
                  "securityContext" => %{
                    "runAsUser" => 405,
                    "runAsGroup" => 1000,
                    "runAsNonRoot" => true
                  }
                },
                %{
                  "name" => "tools",
                  "image" => "registry.unchain.io/unchainio/tools:v0.0.2",
                  "imagePullPolicy" => "IfNotPresent",
                  "resources" => %{
                    "requests" => %{
                      "cpu" => "100m",
                      "memory" => "100Mi"
                    },
                    "limits" => %{
                      "cpu" => "100m",
                      "memory" => "100Mi"
                    }
                  },
                  "env" => env
                }
              ],
              "volumes" => volumes
            }
          }
        }
      })

    {:ok, _statefulset} = K8s.Client.run(create_statefulset, conn)
  end
end
