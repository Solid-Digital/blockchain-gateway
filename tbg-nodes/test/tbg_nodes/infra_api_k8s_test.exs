defmodule TbgNodes.PermissionedEthereumNetworksInfraAPIK8sTest do
  use TbgNodes.DataCase

  require Logger
  alias TbgNodes.ETH
  alias TbgNodes.PermissionedEthereumNetworks

  @moduletag :PermissionedEthereumNetworksInfraAPIK8s

  def create_curl_pod(k8s_conn, namespace) do
    op =
      K8s.Client.create(%{
        "apiVersion" => "v1",
        "kind" => "Pod",
        "metadata" => %{
          "name" => "curl-pod",
          "namespace" => namespace
        },
        "spec" => %{
          "containers" => [
            %{
              "name" => "curl-container",
              "image" => "curlimages/curl:7.70.0",
              "command" => ["sleep", "9999"],
              "securityContext" => %{
                "runAsUser" => 405,
                "runAsNonRoot" => true
              }
            }
          ]
        }
      })

    {:ok, _} = K8s.Client.run(op, k8s_conn)

    kube_config = Application.get_env(:k8s, :clusters).default.conn
    kubectl_args = ["exec", "--namespace", namespace, "curl-pod", "--"]

    curl_via_pod = fn curl_command ->
      all_args = kubectl_args ++ String.split(curl_command, " ")
      System.cmd("kubectl", all_args, env: [{"KUBECONFIG", kube_config}])
    end

    {:ok, curl_via_pod}
  end

  describe "InfraAPIK8s.deploy_network/1" do
    setup [:create_user, :create_k8s_conn]

    def stateful_set_exists?(node, namespace, k8s_conn) do
      {:ok, node_stateful_set_name} =
        PermissionedEthereumNetworks.InfraAPIK8s.get_statefulset_name(node)

      op =
        K8s.Client.get("apps/v1", "StatefulSet",
          namespace: namespace,
          name: node_stateful_set_name
        )

      case K8s.Client.run(op, k8s_conn) do
        {:ok, _} -> true
        _ -> false
      end
    end

    def date_from_string!(date_string) do
      {:ok, date, _} = DateTime.from_iso8601(date_string)

      date
    end

    def node_service_exists?(node, namespace, k8s_conn) do
      {:ok, node_service_name} =
        PermissionedEthereumNetworks.InfraAPIK8s.get_node_service_name(node)

      op =
        K8s.Client.get("v1", "Service",
          namespace: namespace,
          name: node_service_name
        )

      {:ok, _} = K8s.Client.run(op, k8s_conn)
    end

    def create_besu_node(name, type) do
      {:ok, private_key} = ETH.generate_private_key(:hex)
      {:ok, public_key} = ETH.get_public_key(:hex, private_key)

      node = %PermissionedEthereumNetworks.BesuNode{
        name: name,
        uuid: Ecto.UUID.generate(),
        public_key: public_key,
        private_key: private_key,
        node_type: type
      }

      {:ok, node}
    end

    @tag :k8s
    @tag timeout: 300_000
    test "InfraAPIK8s.deploy_network/1", %{k8s_conn: k8s_conn, user: user} do
      {:ok, besu_normal_node} = create_besu_node("normal-1", "normal")
      {:ok, besu_boot_node} = create_besu_node("boot-1", "boot")
      {:ok, besu_validator_node} = create_besu_node("validator-1", "validator")
      {:ok, besu_validator_address} = ETH.get_address(:hex, besu_validator_node.private_key)
      genesis = ETH.genesis([besu_validator_address])

      network = %PermissionedEthereumNetworks.Network{
        name: "some-network-name",
        uuid: "#{Ecto.UUID.generate()}",
        config: genesis,
        besu_nodes: [besu_normal_node, besu_boot_node, besu_validator_node],
        user: user
      }

      {:ok, namespace} = PermissionedEthereumNetworks.InfraAPIK8s.get_namespace_name(network.uuid)
      delete_ns_on_exit(%{test_ns: namespace, k8s_conn: k8s_conn})

      PermissionedEthereumNetworks.InfraAPIK8s.deploy_network(network)
      {:ok, curl_via_pod} = create_curl_pod(k8s_conn, namespace)

      op = K8s.Client.get("v1", "ConfigMap", namespace: namespace, name: "genesis-configmap")
      {:ok, _} = K8s.Client.run(op, k8s_conn)

      # Check that the namespace has the correct label.
      op = K8s.Client.get("v1", "namespace", name: namespace)
      {:ok, namespace_k8s} = K8s.Client.run(op, k8s_conn)

      deployment_target = PermissionedEthereumNetworks.InfraAPIK8s.get_deployment_target!()

      assert namespace_k8s["metadata"]["labels"]["tbg.unchain.io/deployment-target"] ==
               deployment_target

      for node <- network.besu_nodes do
        # Check node service
        assert node_service_exists?(node, namespace, k8s_conn)
        assert stateful_set_exists?(node, namespace, k8s_conn)

        assert_until(fn -> containers_started?(node, namespace, k8s_conn) end,
          label: "containers_started?"
        )

        {:ok, node_service_name} =
          PermissionedEthereumNetworks.InfraAPIK8s.get_node_service_name(node)

        assert_until(fn -> node_is_up?(curl_via_pod, node_service_name) end, label: "node_is_up?")

        assert_until(fn -> peers_are_valid?(curl_via_pod, node_service_name) end,
          label: "peers_are_valid?"
        )

        assert_until(fn -> block_number_is_valid?(curl_via_pod, node_service_name) end,
          label: "block_number_is_valid?"
        )
      end
    end
  end

  def containers_started?(node, namespace, k8s_conn) do
    {:ok, node_stateful_set_name} =
      PermissionedEthereumNetworks.InfraAPIK8s.get_statefulset_name(node)

    pod_name = "#{node_stateful_set_name}-0"

    check_pod = K8s.Client.get("v1", "Pod", namespace: namespace, name: pod_name)

    {:ok, pod} = K8s.Client.run(check_pod, k8s_conn)
    container_statuses = pod["status"]["containerStatuses"]

    container_statuses != nil &&
      container_statuses |> check_containers_started() &&
      container_statuses |> check_containers_ran_for_some_time(10) &&
      container_statuses |> check_containers_did_not_restart()
  end

  def check_containers_started(container_statuses) do
    container_statuses |> Enum.all?(fn %{"started" => started} -> started end)
  end

  def check_containers_ran_for_some_time(container_statuses, time) do
    container_statuses
    |> Enum.all?(fn %{"state" => %{"running" => %{"startedAt" => started_at}}} ->
      started_at
      |> date_from_string!()
      |> DateTime.add(time)
      |> DateTime.compare(DateTime.now!("Etc/UTC")) == :lt
    end)
  end

  def check_containers_did_not_restart(container_statuses) do
    container_statuses |> Enum.all?(fn %{"restartCount" => count} -> count == 0 end)
  end

  def node_is_up?(curl_via_pod, node_service_name) do
    # Check the besu node readiness.
    curl_command = "curl -s #{node_service_name}:8545/readiness"

    with {response, 0} <- curl_via_pod.(curl_command),
         {:ok, %{"status" => "UP"}} <- Jason.decode(response) do
      true
    else
      _ -> false
    end
  end

  def peers_are_valid?(curl_via_pod, node_service_name) do
    # Check that each node is connected to 2 peers.
    curl_command =
      ~s(curl -s -X POST --data {"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1} #{
        node_service_name
      }:8545)

    with {response, 0} <- curl_via_pod.(curl_command),
         {:ok, %{"result" => "0x2"}} <- Jason.decode(response) do
      true
    else
      _ -> false
    end
  end

  def block_number_is_valid?(curl_via_pod, node_service_name) do
    # Check that the blockchain height is not 0.
    curl_command =
      ~s(curl -s -X POST --data {"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1} #{
        node_service_name
      }:8545)

    with {response, 0} <- curl_via_pod.(curl_command),
         {:ok, %{"result" => block_number}} <- Jason.decode(response),
         true <- block_number != "0x0" do
      true
    else
      _ -> false
    end
  end

  describe "InfraAPIK8s.deploy_external_interface/1" do
    setup [:create_k8s_conn]

    def create_namespace(name, k8s_conn) do
      namespace_k8s = %{
        "apiVersion" => "v1",
        "kind" => "Namespace",
        "metadata" => %{"name" => name}
      }

      create_namespace_op = K8s.Client.create(namespace_k8s)
      {:ok, _namespace} = K8s.Client.run(create_namespace_op, k8s_conn)
    end

    for node_type <- ["validator", "normal", "boot"] do
      for protocol <- ["http", "websocket"] do
        @tag :k8s
        test "protocol=#{protocol} node_type=#{node_type}", %{
          k8s_conn: k8s_conn
        } do
          network_uuid = Ecto.UUID.generate()
          target_node_type = unquote(node_type)
          target = %{network_uuid: network_uuid, node_type: target_node_type}
          protocol = unquote(protocol)
          external_interface_uuid = Ecto.UUID.generate()

          external_interface = %PermissionedEthereumNetworks.ExternalInterface{
            protocol: protocol,
            target: target,
            uuid: external_interface_uuid
          }

          # Create namespace for external network to be deployed.
          {:ok, namespace_name} =
            PermissionedEthereumNetworks.InfraAPIK8s.get_namespace_name(network_uuid)

          {:ok, _} = create_namespace(namespace_name, k8s_conn)

          # Clean up once test is done and delete namespace used in test.
          delete_ns_on_exit(%{test_ns: namespace_name, k8s_conn: k8s_conn})

          # Deploy external interface.
          {:ok, {:url, result_url}} =
            PermissionedEthereumNetworks.InfraAPIK8s.deploy_external_interface(external_interface)

          list_services_op = K8s.Client.list("v1", "Service", namespace: namespace_name)

          {:ok, %{"items" => services}} = K8s.Client.run(list_services_op, k8s_conn)

          # Check that we only created 1 service.
          assert Enum.count(services) == 1
          [service] = services

          expected_resource_labels = %{
            "tbg.unchain.io/external-interface-uuid" => external_interface_uuid,
            "tbg.unchain.io/network-uuid" => network_uuid
          }

          assert %{
                   "metadata" => %{
                     "labels" => ^expected_resource_labels,
                     "name" => result_ingress_service_name
                   },
                   "spec" => %{
                     "selector" => %{
                       "tbg.unchain.io/network-uuid" => ^network_uuid,
                       "tbg.unchain.io/node-type" => ^target_node_type
                     }
                   }
                 } = service

          # Check ingress created.
          list_ingresses_op =
            K8s.Client.list("networking.k8s.io/v1beta1", "Ingress", namespace: namespace_name)

          {:ok, %{"items" => ingresses}} = K8s.Client.run(list_ingresses_op, k8s_conn)

          # Check that we only created 1 ingress.
          assert Enum.count(ingresses) == 1
          [ingress] = ingresses

          ingress_host = PermissionedEthereumNetworks.InfraAPIK8s.get_ingress_host!()

          assert %{
                   "metadata" => %{
                     "labels" => ^expected_resource_labels,
                     "annotations" => %{
                       "nginx.ingress.kubernetes.io/auth-url" => result_ingress_basicauth_url
                     }
                   },
                   "spec" => %{
                     "rules" => [
                       %{
                         "host" => ^ingress_host,
                         "http" => %{
                           "paths" => [
                             %{
                               "backend" => %{
                                 "serviceName" => ^result_ingress_service_name,
                                 "servicePort" => ^protocol
                               },
                               "path" => result_ingress_path
                             }
                           ]
                         }
                       }
                     ],
                     "tls" => [%{"hosts" => [^ingress_host]}]
                   }
                 } = ingress

          url_prefix =
            case protocol do
              "websocket" -> "wss"
              "http" -> "https"
            end

          # Check that the url returns by deploy_external_interface matches the ingress host + path
          url_regex = ~r/#{url_prefix}\:\/\/#{ingress_host}#{result_ingress_path}/
          assert String.match?("#{result_url}/liveness", url_regex) == true

          # Check that the ingress path includes external interface uuid
          ingress_path_regex = ~r/#{external_interface_uuid}/
          assert String.match?(result_ingress_path, ingress_path_regex) == true

          # Check that basicauth url includes basicauth host and external interface uuid
          ingress_basicauth_host =
            PermissionedEthereumNetworks.InfraAPIK8s.get_ingress_basicauth_host!()

          ingress_basicauth_url_regex = ~r/#{ingress_basicauth_host}.*#{external_interface_uuid}/
          assert String.match?(result_ingress_basicauth_url, ingress_basicauth_url_regex) == true
        end
      end
    end
  end

  describe "delete_network/1" do
    setup [:create_k8s_conn]

    @tag :k8s
    test "works with valid arguments", %{k8s_conn: k8s_conn} do
      network = %PermissionedEthereumNetworks.Network{
        uuid: Ecto.UUID.generate()
      }

      {:ok, namespace} = PermissionedEthereumNetworks.InfraAPIK8s.get_namespace_name(network.uuid)

      network_namespace = %{
        "apiVersion" => "v1",
        "kind" => "Namespace",
        "metadata" => %{
          "name" => namespace
        }
      }

      create_namespace_op = K8s.Client.create(network_namespace)
      {:ok, _ns} = K8s.Client.run(create_namespace_op, k8s_conn)

      # Delete network via InfraAPIK8s module.
      {:ok} = PermissionedEthereumNetworks.InfraAPIK8s.delete_network(network)

      get_namespace_op = K8s.Client.get(network_namespace)
      network_namespace = K8s.Client.run(get_namespace_op, k8s_conn)

      # Check that network namespace is being deleted/terminated.
      assert {:ok, %{"status" => %{"phase" => "Terminating"}}} = network_namespace
    end

    @tag :k8s
    test "returns error when fails" do
      network = %PermissionedEthereumNetworks.Network{
        uuid: Ecto.UUID.generate()
      }

      # Since no namespace exists for the network to delete, the call with fail.

      # Delete network via InfraAPIK8s module.
      assert {:error, _} = PermissionedEthereumNetworks.InfraAPIK8s.delete_network(network)
    end
  end

  describe "get_bootnodes_arg/1" do
    @describetag :get_bootnodes_arg

    test "returns correct string when 2 boot besu nodes" do
      besu_nodes = [
        %{
          name: "node1",
          uuid: "uuid1",
          public_key: "public-key",
          node_type: "boot"
        },
        %{
          name: "node2",
          uuid: "uuid2",
          public_key: "public-key",
          node_type: "boot"
        },
        %{
          name: "node3",
          uuid: "uuid3",
          public_key: "public-key",
          node_type: "validator"
        }
      ]

      expected =
        "enode://public-key@${NODE1_UUID1_SERVICE_HOST}:${NODE1_UUID1_SERVICE_PORT_DISCOVERY},enode://public-key@${NODE2_UUID2_SERVICE_HOST}:${NODE2_UUID2_SERVICE_PORT_DISCOVERY}"

      {:ok, result} = PermissionedEthereumNetworks.InfraAPIK8s.get_besu_bootnodes_arg(besu_nodes)
      assert result == expected
    end

    test "returns correct string when 1 boot node" do
      nodes = [
        %{
          name: "node1",
          uuid: "uuid1",
          public_key: "public-key",
          node_type: "boot"
        },
        %{
          name: "node2",
          uuid: "uuid2",
          public_key: "public-key",
          node_type: "normal"
        },
        %{
          name: "node3",
          uuid: "uuid3",
          public_key: "public-key",
          node_type: "validator"
        }
      ]

      expected =
        "enode://public-key@${NODE1_UUID1_SERVICE_HOST}:${NODE1_UUID1_SERVICE_PORT_DISCOVERY}"

      {:ok, result} = PermissionedEthereumNetworks.InfraAPIK8s.get_besu_bootnodes_arg(nodes)
      assert result == expected
    end
  end
end
