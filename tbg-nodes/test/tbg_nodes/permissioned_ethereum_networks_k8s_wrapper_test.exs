defmodule TbgNodes.PermissionedEthereumNetworksK8sWrapperTest do
  use TbgNodesWeb.LibCase

  alias TbgNodes.PermissionedEthereumNetworks
  import TbgNodes.TestHelpers

  @moduletag :PermissionedEthereumNetworksK8sWrapper

  setup [:create_k8s_conn, :create_test_ns, :delete_ns_on_exit]

  @tag :k8s
  test "create_configmap works with valid arguments", %{k8s_conn: k8s_conn, test_ns: test_ns} do
    test_params = %{
      name: "config",
      data: %{"key" => "value"},
      namespace: test_ns,
      annotations: %{"ann" => "okay"},
      labels: %{"test" => "true"}
    }

    assert {:ok, %{} = result_configmap} =
             PermissionedEthereumNetworks.K8sWrapper.create_configmap(test_params, k8s_conn)

    assert result_configmap["data"] == test_params.data
  end

  @tag :k8s
  test "create_secret works with valid arguments",
       %{k8s_conn: k8s_conn, test_ns: test_ns} do
    test_params = %{
      annotations: %{"ann" => "okay"},
      labels: %{"test" => "true"},
      name: "my-secret",
      namespace: test_ns,
      string_data: %{"key1" => "value1"}
    }

    assert {:ok, _r} =
             PermissionedEthereumNetworks.K8sWrapper.create_secret(test_params, k8s_conn)
  end

  @tag :k8s
  test "create_cluster_ip_service works with valid arguments",
       %{k8s_conn: k8s_conn, test_ns: test_ns} do
    ports = [
      %{
        "port" => 80,
        "targetPort" => 9376,
        "protocol" => "TCP"
      }
    ]

    test_params = %{
      annotations: %{"ann" => "okay"},
      labels: %{"label1" => "value1"},
      name: "my-service",
      namespace: test_ns,
      ports: ports,
      selector: %{"select-key-1" => "value"}
    }

    assert {:ok, _r} =
             PermissionedEthereumNetworks.K8sWrapper.create_cluster_ip_service(
               test_params,
               k8s_conn
             )
  end

  @tag :k8s
  test "create_ingress works with valid arguments",
       %{k8s_conn: k8s_conn, test_ns: test_ns} do
    service_name = "test-service"
    service_port = 80
    service_path = "/"
    service_host = "bar.foor.com"
    ingress_name = "test-ingress"
    labels = %{"test-key" => "test-value"}

    rules = [
      %{
        "host" => service_host,
        "http" => %{
          "paths" => [
            %{
              "path" => service_path,
              "backend" => %{
                "serviceName" => service_name,
                "servicePort" => service_port
              }
            }
          ]
        }
      }
    ]

    annotations = %{"key" => "value"}

    tls = [
      %{
        "secretName" => "tls-secret",
        "hosts" => ["some.host"]
      }
    ]

    params = %{
      "annotations" => annotations,
      "labels" => labels,
      "name" => ingress_name,
      "namespace" => test_ns,
      "rules" => rules,
      "tls" => tls
    }

    assert {:ok, ingress} =
             PermissionedEthereumNetworks.K8sWrapper.create_ingress(params, k8s_conn)

    assert ingress["metadata"]["name"] == ingress_name
    assert ingress["metadata"]["labels"] == labels
    assert ingress["metadata"]["namespace"] == test_ns
    assert ingress["metadata"]["annotations"] == annotations
    assert ingress["spec"]["rules"] == rules
    assert ingress["spec"]["tls"] == tls
  end

  @tag :k8s
  test "create_statefulset works",
       %{k8s_conn: k8s_conn, test_ns: test_ns} do
    pod_annotations = %{
      "prometheus.io/scrape" => "true",
      "prometheus.io/port" => "9545",
      "prometheus.io/path" => "/metrics"
    }

    command = [
      "/opt/besu/bin/besu"
    ]

    resources = %{
      "requests" => %{
        "cpu" => "100m",
        "memory" => "1024Mi"
      },
      "limits" => %{
        "cpu" => "500m",
        "memory" => "2048Mi"
      }
    }

    liveness_probe = %{
      "httpGet" => %{
        "path" => "/liveness",
        "port" => 8545
      },
      "initialDelaySeconds" => 60,
      "periodSeconds" => 30
    }

    labels = %{
      "app" => "test-label"
    }

    volumes = []
    env = []

    besu_data_path = "/opt/besu/data"

    args = [
      "--rpc-http-enabled",
      "--data-path=#{besu_data_path}",
      "--rpc-http-host=0.0.0.0",
      "--rpc-http-port=8545",
      "--rpc-http-api=ETH,NET,IBFT",
      "--graphql-http-enabled",
      "--graphql-http-host=0.0.0.0",
      "--graphql-http-port=8547",
      "--rpc-ws-enabled",
      "--rpc-ws-host=0.0.0.0",
      "--rpc-ws-port=8546",
      "--metrics-enabled=true",
      "--metrics-host=0.0.0.0",
      "--metrics-port=9545",
      "--host-whitelist=*"
    ]

    ports = [
      %{
        "containerPort" => 8545,
        "name" => "json-rpc",
        "protocol" => "TCP"
      },
      %{
        "containerPort" => 8546,
        "name" => "ws",
        "protocol" => "TCP"
      }
    ]

    name = "test-name"

    statefulset_storage_class =
      Application.get_env(
        :tbg_nodes,
        TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s
      )[:statefulset_storage_class]

    local_path_volume_claim_template = %{
      "metadata" => %{"name" => "local-storage"},
      "spec" => %{
        "accessModes" => ["ReadWriteOnce"],
        "storageClassName" => statefulset_storage_class,
        "resources" => %{"requests" => %{"storage" => "10Gi"}}
      }
    }

    local_path_volume_mount = %{
      "name" => local_path_volume_claim_template["metadata"]["name"],
      "mountPath" => besu_data_path
    }

    init_containers = [
      %{
        "name" => "init-container",
        "image" => "registry.unchain.io/unchainio/tools:v0.0.2",
        "securityContext" => %{
          "runAsUser" => 405,
          "runAsNonRoot" => true
        },
        "command" => [
          "sh",
          "-c",
          "echo init"
        ]
      }
    ]

    params = %{
      pod_annotations: pod_annotations,
      args: args,
      command: command,
      env: env,
      image: TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s.get_default_besu_image!(),
      labels: labels,
      liveness_probe: liveness_probe,
      init_containers: init_containers,
      name: name,
      namespace: test_ns,
      ports: ports,
      resources: resources,
      volume_mounts: [local_path_volume_mount],
      volume_claim_templates: [local_path_volume_claim_template],
      volumes: volumes
    }

    assert {:ok, _r} =
             PermissionedEthereumNetworks.K8sWrapper.create_statefulset(params, k8s_conn)

    # Allow pod to be created. Takes time because of the local storage pvc.
    :timer.sleep(10_000)

    check_pod = K8s.Client.get("v1", "Pod", namespace: test_ns, name: "#{name}-0")

    pod_crashed = fn val ->
      List.first(val)["restartCount"] > 0
    end

    opts = [find: ["status", "containerStatuses"], eval: pod_crashed, timeout: 10]

    # If the pod crashes, the command below will return {:ok, _}. If the container runs
    # for 10 seconds, the command below will return {:error., _}
    {:error, _} = K8s.Client.Runner.Wait.run(check_pod, k8s_conn, opts)
    {:ok, pod} = K8s.Client.run(check_pod, k8s_conn)

    assert pod["spec"]["containers"] |> Enum.any?(&match?(%{"name" => "tools"}, &1))
    assert pod["spec"]["containers"] |> Enum.any?(&match?(%{"name" => ^name}, &1))
    assert pod["spec"]["initContainers"] |> Enum.any?(&match?(%{"name" => "init-container"}, &1))
  end
end
