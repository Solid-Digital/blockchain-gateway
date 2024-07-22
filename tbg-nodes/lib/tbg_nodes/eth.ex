defmodule TbgNodes.ETH do
  @moduledoc """
  Extra data for ibft2.
  """

  def generate_private_key(:hex) do
    private_key_hex =
      ETH.Utils.get_private_key()
      |> ETH.Utils.encode16()

    {:ok, private_key_hex}
  end

  def get_public_key(:hex, private_key) do
    public_key =
      private_key
      |> ETH.Utils.get_public_key()
      |> ETH.Utils.encode16()
      |> format_pub_key()

    {:ok, public_key}
  end

  def get_address(:hex, private_key) do
    address =
      private_key
      |> ETH.Utils.get_address()
      |> format_address()
      |> ETH.Utils.decode16()

    {:ok, address}
  end

  def format_pub_key("04" <> pub_key) do
    # Remove the 0x04 prefix that public keys have because besu doesn't want it there
    pub_key
  end

  def format_address("0x" <> address) do
    address |> String.downcase()
  end

  def extra_data_from_addresses(addresses) do
    "0x" <>
      ([
         # vanity data
         <<0::256>>,
         # node addresses
         addresses,
         # votes
         0,
         # round
         <<0::32>>,
         # seals
         []
       ]
       |> ExRLP.encode(encoding: :hex))
  end

  # Truffle - a common tool for interacting with ethereum has certain
  # limitations on this value. There is a certain V value which equals
  # to something + 2 * chainID. And that V value has to be
  # less than 2 ^ 53 - 1 because that's the maximum value of the Number
  # type in Javascript
  # 2 ^ 48 - 1.
  @chain_id_max_value 281_474_976_710_655

  def genesis(addresses) do
    extra_data = extra_data_from_addresses(addresses)

    %{
      config: %{
        chainId: :rand.uniform(@chain_id_max_value),
        constantinoplefixblock: 0,
        contractSizeLimit: 2_147_483_647,
        ibft2: %{
          blockperiodseconds: 4,
          epochlength: 64_800,
          requesttimeoutseconds: 8
        }
      },
      nonce: "0x0",
      timestamp: "0x58ee40ba",
      gasLimit: "0x1fffffffffffff",
      difficulty: "0x1",
      mixHash: "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
      coinbase: "0x0000000000000000000000000000000000000000",
      extraData: extra_data
    }
  end
end
