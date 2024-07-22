defmodule TbgNodes.NetworkMonitorInternal do
  @moduledoc false

  alias TbgNodes.NetworkMonitor.Status
  require Status

  @network_monitor_ets_table :network_monitor

  @type network_or_node ::
          %TbgNodes.PublicEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.BesuNode{}

  @spec init :: atom() | :ets.tid()
  def init do
    :ets.new(@network_monitor_ets_table, [:named_table, :set, :public])
  end

  @spec cache_all_statuses :: :ok
  def cache_all_statuses do
    _ =
      (TbgNodes.PublicEthereumNetworks.list_networks() ++
         TbgNodes.PermissionedEthereumNetworks.list_networks() ++
         TbgNodes.PermissionedEthereumNetworks.list_nodes())
      |> Enum.map(&get_urls_tuple(&1))
      |> Enum.uniq_by(fn {_, liveness_url, readiness_url} -> {liveness_url, readiness_url} end)
      |> Enum.map(fn {network_or_node, liveness_url, readiness_url} ->
        Task.async(fn -> cache_url_task(network_or_node, liveness_url, readiness_url) end)
      end)
      |> Task.await_many(100_000)

    :ok
  end

  @spec get_urls_tuple(network_or_node()) :: {network_or_node(), String.t(), String.t()}

  defp get_urls_tuple(%TbgNodes.PublicEthereumNetworks.Network{} = network) do
    {
      network,
      TbgNodes.PublicEthereumNetworks.get_liveness_url(network),
      TbgNodes.PublicEthereumNetworks.get_readiness_url(network)
    }
  end

  defp get_urls_tuple(%TbgNodes.PermissionedEthereumNetworks.Network{} = network) do
    {
      network,
      TbgNodes.PermissionedEthereumNetworks.get_liveness_url(network),
      TbgNodes.PermissionedEthereumNetworks.get_readiness_url(network)
    }
  end

  defp get_urls_tuple(%TbgNodes.PermissionedEthereumNetworks.BesuNode{} = node) do
    {
      node,
      TbgNodes.PermissionedEthereumNetworks.get_liveness_url(node),
      TbgNodes.PermissionedEthereumNetworks.get_readiness_url(node)
    }
  end

  @spec cache_url_task(network_or_node(), String.t(), String.t()) :: boolean()
  def cache_url_task(network_or_node, liveness_url, readiness_url) do
    old_status = get_cached_status_from_urls(liveness_url, readiness_url)
    new_status = query_status(liveness_url, readiness_url)

    result = cache_status(liveness_url, readiness_url, new_status)

    case maybe_create_notification(
           old_status,
           new_status,
           network_or_node,
           liveness_url,
           readiness_url
         ) do
      nil -> nil
      notification -> send_notification(notification)
    end

    result
  end

  @spec maybe_create_notification(
          Status.t(),
          Status.t(),
          network_or_node(),
          String.t(),
          String.t()
        ) :: nil | String.t()
  def maybe_create_notification(
        old_status,
        new_status,
        network_or_node,
        liveness_url,
        readiness_url
      ) do
    # Only send a notification if the new status is different than the previous one.
    # Don't send a notification if the old status is fetching_status and the new one is ready,
    # because this usually happens at network startup
    if new_status != old_status && new_status != Status.not_available() &&
         !(old_status == Status.fetching_status() && new_status == Status.ready()) do
      fields =
        case network_or_node do
          network = %TbgNodes.PublicEthereumNetworks.Network{} ->
            [type: "Public Network", network_name: network.network_configuration]

          network = %TbgNodes.PermissionedEthereumNetworks.Network{} ->
            [
              type: "Permissioned Network",
              network_name: network.name,
              network_uuid: network.uuid,
              creator: network.user.email
            ]

          node = %TbgNodes.PermissionedEthereumNetworks.BesuNode{} ->
            [
              type: "Permissioned Node",
              node_name: node.name,
              node_uuid: node.uuid,
              network_name: node.network.name,
              network_uuid: node.network.uuid,
              creator: node.network.user.email
            ]
        end ++
          [
            liveness_url: liveness_url,
            readiness_url: readiness_url,
            version: Application.spec(:tbg_nodes, :vsn) |> to_string()
          ]

      create_attachment(
        "Status went from `#{old_status}` to `#{new_status}`!",
        new_status,
        fields
      )
    end
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

  @spec create_attachment(String.t(), Status.t(), keyword()) :: String.t()
  def create_attachment(msg, status, fields) do
    attachment =
      case status do
        Status.not_available() -> %{color: "#828282"}
        Status.fetching_status() -> %{color: "#828282"}
        Status.live() -> %{color: "#ffcc00"}
        Status.ready() -> %{color: "#27AE60"}
        Status.down() -> %{color: "#c92100"}
      end

    fields =
      fields
      |> Enum.map(fn {key, value} -> %{title: titlecase(key), short: true, value: value} end)

    attachment =
      Map.merge(
        attachment,
        %{
          author_name: "TBG Alerts",
          title: "Alert",
          text: msg,
          as_user: true,
          fields: fields
        }
      )
      |> List.wrap()
      |> Jason.encode!()

    attachment
  end

  @spec cache_status(String.t(), String.t(), Status.t()) :: boolean()
  def cache_status(liveness_url, readiness_url, status) do
    case {liveness_url, readiness_url} do
      {"", ""} ->
        false

      _ ->
        true =
          :ets.insert(
            @network_monitor_ets_table,
            {{liveness_url, readiness_url}, status}
          )
    end
  end

  @spec query_status(
          liveness_url :: String.t(),
          readiness_url :: String.t(),
          get_status_fn :: (url :: String.t() -> :up | :down)
        ) :: status :: Status.t()
  def query_status(liveness_url, readiness_url, get_status_fn \\ &get_status/1) do
    with {:readiness, :not_empty} <- {:readiness, check_empty(readiness_url)},
         {:readiness, :down} <- {:readiness, get_status_fn.(readiness_url)},
         {:liveness, :not_empty} <- {:liveness, check_empty(liveness_url)},
         {:liveness, :down} <- {:liveness, get_status_fn.(liveness_url)} do
      Status.down()
    else
      {:readiness, :empty} -> Status.not_available()
      {:readiness, :up} -> Status.ready()
      {:liveness, :empty} -> Status.down()
      {:liveness, :up} -> Status.live()
    end
  end

  @spec get_status(String.t()) :: :up | :down
  def get_status(url) do
    case HTTPoison.get(url) do
      {:ok, %HTTPoison.Response{status_code: 200}} -> :up
      _ -> :down
    end
  end

  @spec check_empty(String.t()) :: :empty | :not_empty
  defp check_empty(""), do: :empty
  defp check_empty(_), do: :not_empty

  @spec get_cached_status(network_or_node()) :: Status.t()
  def get_cached_status(%TbgNodes.PublicEthereumNetworks.Network{} = network) do
    get_cached_status_from_urls(
      TbgNodes.PublicEthereumNetworks.get_liveness_url(network),
      TbgNodes.PublicEthereumNetworks.get_readiness_url(network)
    )
  end

  def get_cached_status(%TbgNodes.PermissionedEthereumNetworks.Network{} = network) do
    get_cached_status_from_urls(
      TbgNodes.PermissionedEthereumNetworks.get_liveness_url(network),
      TbgNodes.PermissionedEthereumNetworks.get_readiness_url(network)
    )
  end

  def get_cached_status(%TbgNodes.PermissionedEthereumNetworks.BesuNode{} = node) do
    get_cached_status_from_urls(
      TbgNodes.PermissionedEthereumNetworks.get_liveness_url(node),
      TbgNodes.PermissionedEthereumNetworks.get_readiness_url(node)
    )
  end

  @spec get_cached_status_from_urls(liveness_url :: String.t(), readiness_url :: String.t()) ::
          status :: Status.t()
  def get_cached_status_from_urls(liveness_url, readiness_url) do
    case :ets.lookup(
           @network_monitor_ets_table,
           {liveness_url, readiness_url}
         ) do
      [] -> Status.fetching_status()
      [{_, status}] -> status
    end
  end
end
