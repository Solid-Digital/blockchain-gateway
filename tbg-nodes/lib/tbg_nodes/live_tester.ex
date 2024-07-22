defmodule TbgNodes.LiveTester do
  @moduledoc false

  use GenServer
  use Phoenix.HTML
  alias TbgNodes.PermissionedEthereumNetworks
  alias TbgNodes.PermissionedEthereumNetworks.InfraAPI
  import TbgNodes.LiveTestHelpers
  require Logger

  @live_tester_ets_table :live_tester

  @spec start_link(any) :: :ignore | {:error, any} | {:ok, pid}
  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{})
  end

  @impl true
  @spec init(any()) :: {:ok, true}
  def init(_args) do
    # Schedule work to be performed at some point
    _ = init_ets()
    schedule_work(1000)

    {:ok, true}
  end

  @spec init_ets :: atom | :ets.tid()
  def init_ets do
    :ets.new(@live_tester_ets_table, [:named_table, :set, :public])
  end

  @spec get_loop_interval :: non_neg_integer()
  def get_loop_interval do
    Application.get_env(:tbg_nodes, TbgNodes.LiveTester)[:loop_interval]
  end

  @spec get_test_user_email :: String.t()
  def get_test_user_email do
    Application.get_env(:tbg_nodes, TbgNodes.LiveTester)[:test_user_email]
  end

  @spec get_test_user_password :: String.t()
  def get_test_user_password do
    Application.get_env(:tbg_nodes, TbgNodes.LiveTester)[:test_user_password]
  end

  @spec schedule_work(non_neg_integer()) :: reference()
  defp schedule_work(time) do
    Process.send_after(self(), :work, time)
  end

  @impl true
  def handle_info(:work, first_iteration?) do
    success = test_permissioned_network_creation()

    case maybe_create_notification(
           success,
           first_iteration?,
           "test_permissioned_network_creation"
         ) do
      nil -> nil
      notification -> send_notification(notification)
    end

    schedule_work(get_loop_interval())
    {:noreply, false}
  end

  def handle_info(msg, state) do
    Logger.info("#{__MODULE__} - unhandled event #{inspect(msg)}")
    {:noreply, state}
  end

  @spec maybe_create_notification(boolean(), boolean(), String.t()) :: nil | binary
  def maybe_create_notification(
        success,
        first_iteration?,
        test_name
      ) do
    if !success || first_iteration? do
      fields = [
        test_name: test_name,
        test_user_email: get_test_user_email(),
        version: Application.spec(:tbg_nodes, :vsn) |> to_string()
      ]

      color =
        case success do
          true -> "#27AE60"
          false -> "#c92100"
        end

      message =
        case success do
          true -> "Live Test succeeded"
          false -> "Live Test failed"
        end

      create_attachment(
        message,
        color,
        fields
      )
    end
  end

  @spec create_attachment(String.t(), String.t(), keyword()) :: String.t()
  def create_attachment(msg, color, fields) do
    fields =
      fields
      |> Enum.map(fn {key, value} -> %{title: titlecase(key), short: true, value: value} end)

    %{
      color: color,
      author_name: "TBG Alerts",
      title: "Alert",
      text: msg,
      as_user: true,
      fields: fields
    }
    |> List.wrap()
    |> Jason.encode!()
  end

  @spec titlecase(atom) :: String.t()
  defp titlecase(atom) do
    atom |> Atom.to_string() |> String.split("_") |> Enum.join(" ") |> String.capitalize()
  end

  @spec send_notification(String.t()) :: true | false
  def send_notification(attachments_json) do
    TbgNodesWeb.SlackNotifier.send_message(
      "",
      %{attachments: [attachments_json]}
    )
  end

  @valid_user_input %{
    network_name: "permissioned_network_1",
    number_besu_validators: 2,
    number_besu_normal_nodes: 3,
    number_besu_boot_nodes: 1,
    join_network: false,
    managed_by: "unchain",
    consensus: "clique"
  }

  @spec test_permissioned_network_creation((() -> InfraAPI)) :: boolean()
  def test_permissioned_network_creation(get_infra_api \\ &InfraAPI.get_infra_api!/0) do
    user =
      TbgNodes.LiveTestHelpers.user_fixture(%{
        email: get_test_user_email(),
        password: get_test_user_password()
      })

    valid_user_input =
      TbgNodes.LiveTestHelpers.permissioned_ethereum_network_user_input(@valid_user_input)

    {:ok, network} =
      TbgNodes.PermissionedEthereumNetworks.create_network(
        valid_user_input,
        user.id,
        get_infra_api
      )

    network = PermissionedEthereumNetworks.get_network_by_uuid!(network.uuid)

    try do
      interface =
        network.external_interfaces
        |> Enum.filter(fn interface ->
          interface.protocol == "http"
        end)
        |> List.first()

      creds = interface.basicauth_creds |> hd

      request =
        TbgNodesWeb.NetworksView.assemble_request(
          interface.protocol,
          interface.url,
          creds.username,
          creds.password
        )

      success =
        check_until(fn -> node_api_works?(request) end, max_iterations: 600, label: "live_tester")

      Logger.info("test_permissioned_network_creation success - #{success}")

      true =
        :ets.insert(
          @live_tester_ets_table,
          {:test_permissioned_ethereum_network, success}
        )

      success
    after
      :ok = delete_test_user_networks(user.id, get_infra_api)
    end
  end

  @spec node_api_works?(String.t()) :: boolean
  def node_api_works?(request) do
    {~s({\n  "jsonrpc" : "2.0",\n  "id" : 67,\n  "result" : "0x40"\n}), _}
    |> match?(System.cmd("/bin/sh", ["-c", request <> " -k 2> /dev/null"]))
  end

  @spec delete_test_user_networks(integer(), (() -> InfraAPI)) :: :ok
  def delete_test_user_networks(user_id, get_infra_api \\ &InfraAPI.get_infra_api!/0) do
    _ =
      PermissionedEthereumNetworks.list_networks_for_user(user_id)
      |> Enum.map(fn network ->
        Task.async(fn ->
          Logger.info("deleting network #{network.uuid}")
          {:ok, _} = PermissionedEthereumNetworks.delete_network(network.uuid, get_infra_api)
          Logger.info("deleted network #{network.uuid}")
        end)
      end)
      |> Task.await_many(10_000)

    :ok
  end
end
