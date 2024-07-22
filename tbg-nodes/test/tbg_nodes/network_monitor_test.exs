defmodule TbgNodes.NetworkMonitorTest do
  use TbgNodes.DataCase

  require Logger
  alias TbgNodes.NetworkMonitor.Status
  require TbgNodes.NetworkMonitor.Status

  alias TbgNodes.PermissionedEthereumNetworks

  @valid_user_input %{
    network_name: "permissioned_network_1",
    number_besu_validators: 2,
    number_besu_normal_nodes: 3,
    number_besu_boot_nodes: 1,
    join_network: false,
    managed_by: "unchain",
    consensus: "clique"
  }

  describe "send notification" do
    setup [
      :create_user,
      :create_public_ethereum_network_with_interfaces,
      :create_permissioned_ethereum_network,
      :create_permissioned_node
    ]

    tests = [
      %{
        label: "Status.down() -> Status.down() sends no notification",
        input: %{
          old_status: Status.down(),
          new_status: Status.down()
        },
        expected: %{
          is_nil: true,
          contains: ""
        }
      },
      %{
        label: "Status.fetching_status() -> Status.fetching_status() sends no notification",
        input: %{
          old_status: Status.fetching_status(),
          new_status: Status.fetching_status()
        },
        expected: %{
          is_nil: true,
          contains: ""
        }
      },
      %{
        label: "Status.live() -> Status.live() sends no notification",
        input: %{
          old_status: Status.live(),
          new_status: Status.live()
        },
        expected: %{
          is_nil: true,
          contains: ""
        }
      },
      %{
        label: "Status.ready() -> Status.ready() sends no notification",
        input: %{
          old_status: Status.ready(),
          new_status: Status.ready()
        },
        expected: %{
          is_nil: true,
          contains: ""
        }
      },
      %{
        label: "Status.fetching_status() -> Status.ready() sends no notification",
        input: %{
          old_status: Status.fetching_status(),
          new_status: Status.ready()
        },
        expected: %{
          is_nil: true,
          contains: ""
        }
      },
      %{
        label: "Status.fetching_status() -> Status.down() sends a notification",
        input: %{
          old_status: Status.fetching_status(),
          new_status: Status.down()
        },
        expected: %{
          is_nil: false,
          contains: "to `down`"
        }
      },
      %{
        label: "Status.fetching_status() -> Status.live() sends a notification",
        input: %{
          old_status: Status.fetching_status(),
          new_status: Status.live()
        },
        expected: %{
          is_nil: false,
          contains: "to `live`"
        }
      },
      %{
        label: "Status.live() -> Status.ready() sends a notification",
        input: %{
          old_status: Status.live(),
          new_status: Status.ready()
        },
        expected: %{
          is_nil: false,
          contains: "to `ready`"
        }
      },
      %{
        label: "Status.ready() -> Status.live() sends a notification",
        input: %{
          old_status: Status.ready(),
          new_status: Status.live()
        },
        expected: %{
          is_nil: false,
          contains: "to `live`"
        }
      }
    ]

    for %{label: label, input: input, expected: expected} <- tests do
      @label label
      @input input
      @expected expected
      test "#{@label}:", %{
        public_network: public_network,
        permissioned_network: permissioned_network,
        permissioned_node: permissioned_node
      } do
        [public_network, permissioned_network, permissioned_node]
        |> Enum.each(fn network_or_node ->
          notification =
            TbgNodes.NetworkMonitorInternal.maybe_create_notification(
              @input.old_status,
              @input.new_status,
              network_or_node,
              "liveness",
              "readiness"
            )

          if @expected.is_nil do
            assert notification == nil
          else
            assert notification =~ @expected.contains
          end
        end)
      end
    end
  end

  describe "NetworkMonitor.get_cached_status/1 shows the status of a network created by PermissionedEthereumNetworks.create_network/1 (integration test)" do
    setup [:create_user]

    def get_infra_api_k8s do
      TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s
    end

    @tag :k8s
    @tag timeout: 300_000
    test "valid user input", %{user: user} do
      valid_user_input = permissioned_ethereum_network_user_input(@valid_user_input)

      {:ok, network} =
        PermissionedEthereumNetworks.create_network(
          valid_user_input,
          user.id,
          &get_infra_api_k8s/0
        )

      delete_network_on_exit(network, &get_infra_api_k8s/0)

      query =
        from network in PermissionedEthereumNetworks.Network,
          where: network.uuid == ^network.uuid,
          select: network,
          preload: [:besu_nodes, external_interfaces: :basicauth_creds]

      network = Repo.one(query)

      with_network_monitor(%{})

      assert_until(fn -> network_status_ready?(network) end, label: "network_monitor")
    end
  end

  def network_status_ready?(network) do
    TbgNodes.NetworkMonitor.get_cached_status(network) == Status.ready()
  end

  describe "test query_status" do
    tests = [
      %{
        label: "Status.ready()",
        input: %{
          liveness_url: "liveness",
          readiness_url: "readiness",
          liveness_response: :up,
          readiness_response: :up
        },
        expected: %{
          status: Status.ready()
        }
      },
      %{
        label: "Status.live()",
        input: %{
          liveness_url: "liveness",
          readiness_url: "readiness",
          liveness_response: :up,
          readiness_response: :down
        },
        expected: %{
          status: Status.live()
        }
      },
      %{
        label: "Status.down()",
        input: %{
          liveness_url: "liveness",
          readiness_url: "readiness",
          liveness_response: :down,
          readiness_response: :down
        },
        expected: %{
          status: Status.down()
        }
      },
      %{
        label: "Status.not_available()",
        input: %{
          liveness_url: "liveness",
          readiness_url: "",
          readiness_response: :not_available,
          liveness_response: :not_available
        },
        expected: %{
          status: Status.not_available()
        }
      },
      %{
        label: "Status.down() with no liveness check",
        input: %{
          liveness_url: "",
          readiness_url: "readiness",
          liveness_response: :not_available,
          readiness_response: :down
        },
        expected: %{
          status: Status.down()
        }
      }
    ]

    for %{label: label, input: input, expected: expected} <- tests do
      @label label
      @input input
      @expected expected
      test "#{@label}:" do
        query_fn = fn url ->
          case url do
            "liveness" -> @input.liveness_response
            "readiness" -> @input.readiness_response
          end
        end

        assert TbgNodes.NetworkMonitorInternal.query_status(
                 @input.liveness_url,
                 @input.readiness_url,
                 query_fn
               ) ==
                 @expected.status
      end
    end
  end
end
