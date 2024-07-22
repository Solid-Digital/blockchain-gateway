defmodule TbgNodes.Repo.Migrations.ChangeNetworkToPublicEthereumNetwork do
  use Ecto.Migration

  def change do
    rename table(:networks), to: table(:public_ethereum_networks)

    rename table(:network_interface_creds), to: table(:public_ethereum_networks_basicauth_creds)

    rename table(:network_interfaces),
      to: table(:public_ethereum_networks_network_external_interfaces)

    alter table(:public_ethereum_networks) do
      remove(:protocol)
    end

    rename table(:public_ethereum_networks_basicauth_creds), :network_interface_id,
      to: :public_ethereum_networks_network_external_interface_id

    alter table(:public_ethereum_networks_basicauth_creds) do
      modify :public_ethereum_networks_network_external_interface_id,
             references(:public_ethereum_networks_network_external_interfaces)
    end
  end
end
