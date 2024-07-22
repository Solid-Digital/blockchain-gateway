defmodule TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s do
  @moduledoc false
  alias TbgNodes.PermissionedEthereumNetworks

  @behaviour PermissionedEthereumNetworks.InfraAPI

  def get_default_besu_image! do
    Application.get_env(:tbg_nodes, TbgNodes.Networks)[:besu_image]
  end

  @default_env []

  @default_container_ports [
    %{
      "containerPort" => 8545,
      "name" => "http",
      "protocol" => "TCP"
    },
    %{
      "containerPort" => 8546,
      "name" => "websocket",
      "protocol" => "TCP"
    },
    %{
      "containerPort" => 8547,
      "name" => "graphql",
      "protocol" => "TCP"
    },
    %{
      "containerPort" => 30_303,
      "name" => "rlpx",
      "protocol" => "TCP"
    },
    %{
      "containerPort" => 30_303,
      "name" => "discovery",
      "protocol" => "UDP"
    }
  ]

  @default_resources %{
    "requests" => %{
      "cpu" => "100m",
      "memory" => "1024Mi"
    },
    "limits" => %{
      "cpu" => "500m",
      "memory" => "2048Mi"
    }
  }

  @genesis_file %{
    configmap_name: "genesis-configmap",
    mount_path: "/config",
    volume_name: "genesis-volume",
    filename: "genesis.json"
  }

  # --node-private-key-file=/secrets/key \
  @private_key_file %{
    filename: "private.key",
    mount_path: "/private_key",
    secret_name_base: "node-private-key",
    volume_name: "private-key",
    volume_name_base: "node-private"
  }

  @besu_data_path "/opt/besu/data"

  @default_args [
    "exec /opt/besu/bin/besu",
    "--genesis-file=#{@genesis_file.mount_path}/#{@genesis_file.filename}",
    "--data-path=#{@besu_data_path}",
    "--graphql-http-enabled",
    "--graphql-http-host=0.0.0.0",
    "--graphql-http-port=8547",
    "--host-whitelist=*",
    "--metrics-enabled=true",
    "--metrics-host=0.0.0.0",
    "--metrics-port=9545",
    "--min-gas-price=0",
    "--node-private-key-file=#{@private_key_file.mount_path}/#{@private_key_file.filename}",
    "--revert-reason-enabled=true",
    "--rpc-http-api=ETH,NET,IBFT",
    "--rpc-http-cors-origins=*",
    "--rpc-http-enabled",
    "--rpc-http-host=0.0.0.0",
    "--rpc-http-port=8545",
    "--rpc-ws-enabled",
    "--rpc-ws-host=0.0.0.0",
    "--rpc-ws-port=8546"
  ]

  @websocket_service_port %{
    "port" => 8546,
    "targetPort" => 8546,
    "protocol" => "TCP",
    "name" => "websocket"
  }

  @http_service_port %{
    "port" => 8545,
    "targetPort" => 8545,
    "protocol" => "TCP",
    "name" => "http"
  }

  @default_service_ports [
                           %{
                             "port" => 8547,
                             "targetPort" => 8547,
                             "protocol" => "TCP",
                             "name" => "graphql"
                           },
                           %{
                             "port" => 30_303,
                             "targetPort" => 30_303,
                             "protocol" => "TCP",
                             "name" => "rlpx"
                           },
                           %{
                             "port" => 30_303,
                             "targetPort" => 30_303,
                             "protocol" => "UDP",
                             "name" => "discovery"
                           }
                         ] ++
                           [
                             @websocket_service_port,
                             @http_service_port
                           ]

  @default_liveness_probe %{
    "httpGet" => %{
      "path" => "/liveness",
      "port" => 8545
    },
    "initialDelaySeconds" => 60,
    "periodSeconds" => 30
  }

  @prometheus_pod_scrape_annotations %{
    "prometheus.io/path" => "/metrics",
    "prometheus.io/port" => "9545",
    "prometheus.io/scrape" => "true"
  }

  @spec get_statefulset_storage_class! :: String.t()
  def get_statefulset_storage_class! do
    Application.get_env(:tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s)[
      :statefulset_storage_class
    ] ||
      raise "No config value found for :statefulset_storage_class."
  end

  @spec get_cloud_storage_volume_claim_template(String.t()) :: {:ok, map()}
  def get_cloud_storage_volume_claim_template(statefulset_storage_class) do
    {:ok,
     %{
       "metadata" => %{"name" => "cloud-storage"},
       "spec" => %{
         "accessModes" => ["ReadWriteOnce"],
         "storageClassName" => statefulset_storage_class,
         "resources" => %{"requests" => %{"storage" => "10Gi"}}
       }
     }}
  end

  @spec get_deployment_target!() :: String.t()
  def get_deployment_target! do
    Application.get_env(:tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s)[
      :deployment_target
    ] ||
      raise "No config value found for :deployment_target."
  end

  @spec get_ingress_host!() :: String.t()
  def get_ingress_host! do
    Application.get_env(:tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s)[
      :ingress_host
    ] ||
      raise "No config value found for :ingress_host."
  end

  @spec get_ingress_basicauth_host!() :: String.t()
  def get_ingress_basicauth_host! do
    Application.get_env(:tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s)[
      :ingress_basicauth_host
    ] ||
      raise "No config value found for :ingress_basicauth_host."
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def delete_network(%PermissionedEthereumNetworks.Network{uuid: network_uuid} = _network) do
    {:ok, conn} = K8s.Conn.lookup(:default)
    {:ok, namespace} = get_namespace_name(network_uuid)

    network_namespace = %{
      "apiVersion" => "v1",
      "kind" => "Namespace",
      "metadata" => %{
        "name" => namespace
      }
    }

    delete_namespace = K8s.Client.delete(network_namespace)

    case K8s.Client.run(delete_namespace, conn) do
      {:ok, _} -> {:ok}
      {:error, message} -> {:error, message}
    end
  end

  defp get_conn! do
    {:ok, conn} = K8s.Conn.lookup(:default)
    conn
  end

  defp get_network_labels(network) do
    %{
      "tbg.unchain.io/deployment-target" => get_deployment_target!(),
      "tbg.unchain.io/network-type" => "permissioned",
      "tbg.unchain.io/network-name" => network.name,
      "tbg.unchain.io/network-uuid" => network.uuid,
      "tbg.unchain.io/author-id" => "u-#{network.user.id}"
    }
  end

  defp get_node_labels(node) do
    %{
      "tbg.unchain.io/node-uuid" => node.uuid,
      "tbg.unchain.io/node-type" => node.node_type
    }
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  @spec deploy_network(%PermissionedEthereumNetworks.Network{}) :: {:ok} | {:error, String.t()}
  def deploy_network(
        %PermissionedEthereumNetworks.Network{
          uuid: network_uuid,
          config: %{} = genesis_data,
          besu_nodes: besu_nodes
        } = network
      ) do
    conn = get_conn!()
    {:ok, namespace_name} = get_namespace_name(network_uuid)

    network_labels = get_network_labels(network)

    {:ok, _} =
      PermissionedEthereumNetworks.K8sWrapper.create_namespace(
        %{name: namespace_name, labels: network_labels, annotations: %{}},
        conn
      )

    # Give time for namespace to be created.
    :timer.sleep(3000)

    {:ok, _} =
      PermissionedEthereumNetworks.K8sWrapper.create_configmap(
        %{
          annotations: %{},
          data: %{@genesis_file.filename => Jason.encode!(genesis_data, pretty: true)},
          labels: network_labels,
          name: @genesis_file.configmap_name,
          namespace: namespace_name
        },
        conn
      )

    {:ok, besu_bootnodes_arg} = get_besu_bootnodes_arg(besu_nodes)
    non_boot_node_args = @default_args ++ ["--bootnodes=#{besu_bootnodes_arg}"]
    {:ok, bootnode_liveness_url} = get_besu_bootnode_liveness_url_for_init_container(besu_nodes)

    non_boot_node_init_containers = [
      %{
        "name" => "init-bootnode",
        "image" => "registry.unchain.io/unchainio/tools:v0.0.2",
        "securityContext" => %{
          "runAsUser" => 405,
          "runAsNonRoot" => true
        },
        "command" => [
          "sh",
          "-c",
          "curl -X GET --connect-timeout 30 --max-time 10 --retry 6 --retry-delay 0 --retry-max-time 300 #{
            bootnode_liveness_url
          }"
        ]
      }
    ]

    # Besu nodes
    for node <- besu_nodes do
      node_labels =
        %{}
        |> Map.merge(network_labels)
        |> Map.merge(get_node_labels(node))

      # Node Services
      {:ok, node_service_name} = get_node_service_name(node)

      {:ok, _} =
        PermissionedEthereumNetworks.K8sWrapper.create_cluster_ip_service(
          %{
            annotations: %{},
            labels: node_labels,
            name: node_service_name,
            namespace: namespace_name,
            ports: @default_service_ports,
            selector: node_labels
          },
          conn
        )

      # Node Private Key Secret
      private_key_secret_name = "#{@private_key_file.secret_name_base}-#{node.uuid}"

      {:ok, _} =
        PermissionedEthereumNetworks.K8sWrapper.create_secret(
          %{
            annotations: %{},
            labels: node_labels,
            name: private_key_secret_name,
            namespace: namespace_name,
            string_data: %{@private_key_file.filename => node.private_key}
          },
          conn
        )

      args =
        if node.node_type == "boot" do
          @default_args
        else
          non_boot_node_args
        end

      init_containers =
        if node.node_type == "boot" do
          nil
        else
          non_boot_node_init_containers
        end

      # Add NAT arguments
      {:ok, node_service_port_env} = get_service_port_env(node_service_name, :discovery)
      {:ok, service_host_env} = get_service_host_env(node_service_name)

      args =
        args
        |> Kernel.++(["--nat-method=NONE"])
        |> Kernel.++(["--p2p-host=${#{service_host_env}}"])
        |> Kernel.++([
          "--p2p-port=${#{node_service_port_env}}"
        ])

      args_string = Enum.join(args, " ")

      {:ok, genesis_volume} =
        get_genesis_volume(
          @genesis_file.volume_name,
          @genesis_file.configmap_name,
          @genesis_file.filename
        )

      {:ok, genesis_volume_mount} =
        get_genesis_volume_mount(@genesis_file.volume_name, @genesis_file.mount_path)

      {:ok, statefulset_name} = get_statefulset_name(node)

      private_key_volume_name = "#{@private_key_file.volume_name_base}-#{node.uuid}"

      {:ok, private_key_volume} =
        get_private_key_volume(
          @private_key_file.filename,
          private_key_secret_name,
          private_key_volume_name
        )

      {:ok, private_key_volume_mount} =
        get_private_key_volume_mount(private_key_volume_name, @private_key_file.mount_path)

      statefulset_storage_class = get_statefulset_storage_class!()

      {:ok, cloud_storage_volume_claim_template} =
        get_cloud_storage_volume_claim_template(statefulset_storage_class)

      cloud_storage_volume_mount = %{
        "name" => cloud_storage_volume_claim_template["metadata"]["name"],
        "mountPath" => @besu_data_path
      }

      {:ok, _} =
        PermissionedEthereumNetworks.K8sWrapper.create_statefulset(
          %{
            name: statefulset_name,
            namespace: namespace_name,
            pod_annotations: @prometheus_pod_scrape_annotations,
            labels: node_labels,
            liveness_probe: @default_liveness_probe,
            init_containers: init_containers,
            image: get_default_besu_image!(),
            resources: @default_resources,
            volume_claim_templates: [cloud_storage_volume_claim_template],
            volumes: [
              genesis_volume,
              private_key_volume
            ],
            volume_mounts: [
              genesis_volume_mount,
              private_key_volume_mount,
              cloud_storage_volume_mount
            ],
            env: @default_env,
            command: ["/bin/sh", "-c"],
            args: [args_string],
            ports: @default_container_ports
          },
          conn
        )
    end

    {:ok}
  end

  @spec get_private_key_volume(String.t(), String.t(), String.t()) :: {:ok, %{}}
  defp get_private_key_volume(filename, secret_name, volume_name) do
    volume = %{
      "name" => volume_name,
      "secret" => %{
        "secretName" => secret_name,
        "items" => [
          %{
            "key" => filename,
            "path" => filename
          }
        ]
      }
    }

    {:ok, volume}
  end

  @spec get_private_key_volume_mount(String.t(), String.t()) :: {:ok, %{}}
  defp get_private_key_volume_mount(volume_name, mount_path) do
    volume_mount = %{
      "name" => volume_name,
      "mountPath" => mount_path,
      "readOnly" => true
    }

    {:ok, volume_mount}
  end

  @spec get_genesis_volume_mount(String.t(), String.t()) :: {:ok, %{}}
  defp get_genesis_volume_mount(volume_name, mount_path) do
    volume_mount = %{
      "name" => volume_name,
      "mountPath" => mount_path,
      "readOnly" => true
    }

    {:ok, volume_mount}
  end

  @spec get_genesis_volume(String.t(), String.t(), String.t()) :: {:ok, %{}}
  defp get_genesis_volume(volume_name, configmap_name, filename) do
    volume = %{
      "name" => volume_name,
      "configMap" => %{
        "name" => configmap_name,
        "items" => [
          %{
            "key" => filename,
            "path" => filename
          }
        ]
      }
    }

    {:ok, volume}
  end

  @spec get_namespace_name(String.t()) :: {:ok, String.t()}
  def get_namespace_name(network_uuid), do: {:ok, "network-#{network_uuid}"}

  @spec get_node_service_name(%PermissionedEthereumNetworks.BesuNode{}) :: {:ok, String.t()}
  def get_node_service_name(node), do: {:ok, "#{node.name}-#{node.uuid}"}

  @spec get_external_interface_service_name(Ecto.UUID.t()) ::
          {:ok, String.t()}
  def get_external_interface_service_name(external_interface_uuid),
    do: {:ok, "external-interface-#{external_interface_uuid}"}

  @spec get_service_host_env(String.t()) :: {:ok, String.t()}
  def get_service_host_env(service_name) do
    # Example output: NODE_NODE1_SERVICE_SERVICE_HOST=10.43.39.236
    service_host_env =
      service_name
      |> String.upcase()
      |> String.replace("-", "_")
      |> (fn s -> "#{s}_SERVICE_HOST" end).()

    {:ok, service_host_env}
  end

  @spec get_service_port_env(String.t(), :discovery | :http) :: {:ok, String.t()}
  defp get_service_port_env(service_name, port_name) do
    # Example output: NODE_NODE1_SERVICE_SERVICE_PORT_DISCOVERY=30303
    service_port_env =
      service_name
      |> String.replace("-", "_")
      |> (fn s -> "#{s}_SERVICE_PORT_#{Atom.to_string(port_name)}" end).()
      |> String.upcase()

    {:ok, service_port_env}
  end

  @spec get_statefulset_name(%PermissionedEthereumNetworks.BesuNode{}) :: {:ok, String.t()}
  def get_statefulset_name(besu_node), do: {:ok, besu_node.name}

  def get_besu_bootnodes_arg(besu_nodes) do
    bootnodes_arg =
      besu_nodes
      |> Enum.filter(fn node -> node.node_type == "boot" end)
      |> Enum.map(fn node ->
        {:ok, node_service_name} = get_node_service_name(node)
        {:ok, node_service_host_env} = get_service_host_env(node_service_name)
        {:ok, node_service_port_env} = get_service_port_env(node_service_name, :discovery)
        "enode://#{node.public_key}@${#{node_service_host_env}}:${#{node_service_port_env}}"
      end)
      |> Enum.join(",")

    {:ok, bootnodes_arg}
  end

  @spec get_besu_bootnode_liveness_url_for_init_container([
          %TbgNodes.PermissionedEthereumNetworks.BesuNode{}
        ]) :: {:ok, String.t()}
  def get_besu_bootnode_liveness_url_for_init_container(besu_nodes) do
    boot_node =
      besu_nodes
      |> Enum.find(fn node -> node.node_type == "boot" end)

    {:ok, node_service_name} = get_node_service_name(boot_node)
    {:ok, node_service_host_env} = get_service_host_env(node_service_name)
    {:ok, node_service_port_env} = get_service_port_env(node_service_name, :http)

    {:ok, "${#{node_service_host_env}}:${#{node_service_port_env}}/liveness"}
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def deploy_external_interface(
        %PermissionedEthereumNetworks.ExternalInterface{
          protocol: protocol,
          uuid: uuid,
          target: %{network_uuid: network_uuid} = target
        } = _external_interface
      ) do
    conn = get_conn!()
    {:ok, namespace} = get_namespace_name(network_uuid)

    # External Interface Service

    {:ok, service_name} = get_external_interface_service_name(uuid)

    ports =
      case protocol do
        "http" -> [@http_service_port]
        "websocket" -> [@websocket_service_port]
      end

    {:ok, selector} = get_selector(target)

    labels = %{
      "tbg.unchain.io/network-uuid" => network_uuid,
      "tbg.unchain.io/external-interface-uuid" => uuid
    }

    {:ok, _} =
      PermissionedEthereumNetworks.K8sWrapper.create_cluster_ip_service(
        %{
          annotations: %{},
          labels: labels,
          name: service_name,
          namespace: namespace,
          ports: ports,
          selector: selector
        },
        conn
      )

    # External Interface Ingress

    ingress_name = "network-#{protocol}-#{uuid}"
    ingress_host = get_ingress_host!()

    path =
      case protocol do
        "http" -> "v0/#{uuid}"
        "websocket" -> "v0/ws/#{uuid}"
      end

    ingress_rules = [
      %{
        "host" => ingress_host,
        "http" => %{
          "paths" => [
            %{
              "path" => "/#{path}/?(.*)",
              "backend" => %{
                "serviceName" => service_name,
                "servicePort" => protocol
              }
            }
          ]
        }
      }
    ]

    ingress_annotations = %{
      "kubernetes.io/ingress.class" => "nginx",
      "nginx.ingress.kubernetes.io/rewrite-target" => "/$1",
      "nginx.ingress.kubernetes.io/proxy-body-size" => "200m",
      "nginx.ingress.kubernetes.io/auth-url" => "#{get_ingress_basicauth_host!()}/#{uuid}"
    }

    {:ok, _} =
      PermissionedEthereumNetworks.K8sWrapper.create_ingress(
        %{
          "annotations" => ingress_annotations,
          "labels" => labels,
          "name" => ingress_name,
          "namespace" => namespace,
          "rules" => ingress_rules,
          "tls" => [
            %{
              "secretName" => ingress_host,
              "hosts" => [
                ingress_host
              ]
            }
          ]
        },
        conn
      )

    url =
      case protocol do
        "http" -> "https://#{ingress_host}/#{path}"
        "websocket" -> "wss://#{ingress_host}/#{path}"
      end

    {:ok, {:url, url}}
  end

  @spec get_selector(%{network_uuid: Ecto.UUID.t(), node_type: String.t()}) ::
          {:ok, %{String.t() => Ecto.UUID.t(), String.t() => String.t()}}
  defp get_selector(%{network_uuid: uuid, node_type: type} = _target) do
    {:ok, %{"tbg.unchain.io/network-uuid" => uuid, "tbg.unchain.io/node-type" => type}}
  end

  @spec get_selector(%{network_uuid: Ecto.UUID.t()}) :: {:ok, %{String.t() => Ecto.UUID.t()}}
  defp get_selector(%{network_uuid: uuid} = _target) do
    {:ok, %{"tbg.unchain.io/network-uuid" => uuid}}
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @impl PermissionedEthereumNetworks.InfraAPI
  def get_liveness_url(%PermissionedEthereumNetworks.BesuNode{} = node) do
    {:ok, service_name} = get_node_service_name(node)
    {:ok, namespace_name} = get_namespace_name(node.network.uuid)

    res =
      EEx.eval_string(
        TbgNodes.PublicEthereumNetworks.get_network_url_config()["permissioned"]["liveness"],
        service_name: service_name,
        namespace_name: namespace_name
      )

    if is_binary(res) do
      {:ok, res}
    else
      {:error, "the resulting url is not a string"}
    end
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @impl PermissionedEthereumNetworks.InfraAPI
  def get_liveness_url(%PermissionedEthereumNetworks.Network{} = network) do
    external_interface =
      network.external_interfaces
      |> Enum.filter(fn interface ->
        interface.protocol == "http"
      end)
      |> List.first()

    {:ok, service_name} = get_external_interface_service_name(external_interface.uuid)
    {:ok, namespace_name} = get_namespace_name(network.uuid)

    res =
      EEx.eval_string(
        TbgNodes.PublicEthereumNetworks.get_network_url_config()["permissioned"]["liveness"],
        service_name: service_name,
        namespace_name: namespace_name
      )

    if is_binary(res) do
      {:ok, res}
    else
      {:error, "the resulting url is not a string"}
    end
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @impl PermissionedEthereumNetworks.InfraAPI
  def get_readiness_url(%PermissionedEthereumNetworks.BesuNode{} = node) do
    {:ok, service_name} = get_node_service_name(node)
    {:ok, namespace_name} = get_namespace_name(node.network.uuid)

    res =
      EEx.eval_string(
        TbgNodes.PublicEthereumNetworks.get_network_url_config()["permissioned"]["readiness"],
        service_name: service_name,
        namespace_name: namespace_name
      )

    if is_binary(res) do
      {:ok, res}
    else
      {:error, "the resulting url is not a string"}
    end
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @impl PermissionedEthereumNetworks.InfraAPI
  def get_readiness_url(%PermissionedEthereumNetworks.Network{} = network) do
    external_interface =
      network.external_interfaces
      |> Enum.filter(fn interface ->
        interface.protocol == "http"
      end)
      |> List.first()

    {:ok, service_name} = get_external_interface_service_name(external_interface.uuid)
    {:ok, namespace_name} = get_namespace_name(network.uuid)

    res =
      EEx.eval_string(
        TbgNodes.PublicEthereumNetworks.get_network_url_config()["permissioned"]["readiness"],
        service_name: service_name,
        namespace_name: namespace_name
      )

    if is_binary(res) do
      {:ok, res}
    else
      {:error, "the resulting url is not a string"}
    end
  end
end
