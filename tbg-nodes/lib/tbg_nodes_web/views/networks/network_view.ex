defmodule TbgNodesWeb.NetworksView do
  use TbgNodesWeb, :view

  import TbgNodesWeb.ViewHelpers

  def format_protocol(protocol) do
    case protocol do
      "http" ->
        "HTTP"

      "websocket" ->
        "Websocket"
    end
  end

  @spec get_url(binary(), binary(), binary(), binary()) :: String.t()
  def get_url(tag, protocol, network_uuid, network_name) do
    url =
      TbgNodes.PublicEthereumNetworks.get_network_url(
        TbgNodes.PublicEthereumNetworks.get_network_url_config(),
        tag,
        protocol,
        network_uuid,
        network_name
      )

    Kernel.inspect(url)
    |> String.replace("\"", "")
  end

  def archive_data(network_configuration) do
    if String.contains?(network_configuration, "-archive") do
      "Enabled"
    else
      "Disabled"
    end
  end

  @spec assemble_request(String.t(), String.t(), String.t(), String.t()) ::
          String.t()
  def assemble_request(
        "http" = _protocol,
        url,
        username,
        password
      ) do
    url =
      URI.parse(url)
      |> Map.replace!(:userinfo, "#{username}:#{password}")
      |> URI.to_string()

    ~s(curl -X POST --data '{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[],"id":67}' #{
      url
    })
  end

  def assemble_request(
        "websocket" = _protocol,
        url,
        username,
        password
      ) do
    url = String.replace_prefix(url, "wss://", "")

    ~s(wscat -c wss://#{username}:#{password}@#{url})
  end

  def assemble_config_link(network, interface) do
    "/networks/public-ethereum/" <> to_string(network.uuid) <> "#" <> url_safe(interface.protocol)
  end

  def assemble_example_link(network, example) do
    "/networks/" <> to_string(network.uuid) <> "#" <> example
  end

  def url_safe(s) do
    s |> String.downcase() |> String.replace(" ", "-")
  end
end
