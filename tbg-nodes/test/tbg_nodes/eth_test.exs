defmodule TbgNodes.ETHTest do
  use TbgNodes.DataCase
  @moduletag :ETH

  test "ETH.get_public_key(:hex, " do
    private_key = "66be853c4c2e43a749b403c4a15efc024ae0d0b448b8daf55350fefcb867786e"

    expected_public_key =
      "f0c0f59d310fc6994bf6b48b5d6f1d55b60963c91efb62d6b52cb47bd00a47a868c74d82a9bc42294a0333c63415fca089f4d27b8e0d156965c29af6470b38de"

    {:ok, result} = TbgNodes.ETH.get_public_key(:hex, private_key)
    assert result == expected_public_key
  end

  test "ETH.get_address(:hex, " do
    private_key = "66be853c4c2e43a749b403c4a15efc024ae0d0b448b8daf55350fefcb867786e"
    {:ok, result} = TbgNodes.ETH.get_address(:hex, private_key)

    assert result ==
             <<19, 6, 98, 81, 71, 207, 153, 144, 152, 178, 19, 234, 137, 76, 94, 92, 145, 55, 61,
               93>>
  end

  describe "test genesis" do
    tests = [
      %{
        label: "1",
        input: %{
          validator_count: 3
        },
        expected: %{}
      }
    ]

    for %{label: label, input: input, expected: expected} <- tests do
      @label label
      @input input
      @expected expected
      test "#{@label}:" do
        genesis = create_genesis_config()

        assert String.length(genesis.extraData) > 0
        assert genesis.config.chainId != 0
      end
    end
  end

  describe "test genesis chainId is random" do
    genesis1 = create_genesis_config()
    genesis2 = create_genesis_config()
    genesis3 = create_genesis_config()

    assert genesis1.config.chainId != genesis2.config.chainId
    assert genesis1.config.chainId != genesis3.config.chainId
    assert genesis2.config.chainId != genesis3.config.chainId
  end

  describe "test RLP" do
    tests = [
      %{
        label: "1",
        input: %{
          addresses: [
            "284583d5c009a34ddb8814f435d8b21c070c4196",
            "4646de9912997d27e088e0b274368a1e3cfb38b1",
            "735e0e3be059be0df69289372253735ac0958394",
            "742e1d3276cd1d02fd96e5c0c766fdc4eba8af22"
          ]
        },
        expected: %{
          extra_data:
            "0xf87ea00000000000000000000000000000000000000000000000000000000000000000f85494284583d5c009a34ddb8814f435d8b21c070c4196944646de9912997d27e088e0b274368a1e3cfb38b194735e0e3be059be0df69289372253735ac095839494742e1d3276cd1d02fd96e5c0c766fdc4eba8af22808400000000c0"
        }
      },
      %{
        label: "2",
        input: %{
          addresses: [
            "a61fde69f0936e113a80b68161de5bd361d8897b",
            "7c58dcd03717704bb86d260c30c9f38d83f3ac0a",
            "3264301d87cd1585f496a8b15cea4fbf492c0a6a",
            "8317eeb9c4fba32040a7236202161fb2504757e3"
          ]
        },
        expected: %{
          extra_data:
            "0xf87ea00000000000000000000000000000000000000000000000000000000000000000f85494a61fde69f0936e113a80b68161de5bd361d8897b947c58dcd03717704bb86d260c30c9f38d83f3ac0a943264301d87cd1585f496a8b15cea4fbf492c0a6a948317eeb9c4fba32040a7236202161fb2504757e3808400000000c0"
        }
      }
    ]

    for %{label: label, input: input, expected: expected} <- tests do
      @label label
      @input input
      @expected expected
      test "#{@label}:" do
        addresses =
          Enum.map(@input.addresses, fn address ->
            ETH.Utils.decode16(address)
          end)

        assert @expected.extra_data == TbgNodes.ETH.extra_data_from_addresses(addresses)
      end
    end
  end

  describe "test ETH" do
    tests = [
      %{
        label: "1",
        input: %{priv: "cd65c5dac8f2f894223dbb3350ce21ceca4ac17378cde72802e84f9a6b7b9dfb"},
        expected: %{
          pub:
            "7743a85c242db395ac6b52bd4664c16a9cd1d60c04e92bff15b1cbbfd2f9fc2fae6dc70457ebe4be7dfbdad9a36b8eae251db00744d5dc6abaeca14adec123b1",
          address: "735e0e3be059be0df69289372253735ac0958394"
        }
      },
      %{
        label: "2",
        input: %{priv: "835f0720923e39d2c2e0514045c4caae1e3563c6dec23c14887028ec05d6a747"},
        expected: %{
          pub:
            "87787c9f66c0bb7a045a18cadbdcc8256ad257c2c00b3d5b7734881431f4e74425d1729a4269ef53f2e4e88f5e582c36d951f1b0774db0e0ac90b00a432995ea",
          address: "742e1d3276cd1d02fd96e5c0c766fdc4eba8af22"
        }
      },
      %{
        label: "3",
        input: %{priv: "d59f4ac48a9b76dac6c076221f8d76564d271afa34d91975acfbe94847a0476d"},
        expected: %{
          pub:
            "424273b510f5bc00b97185449ea12cced01b3384e04a1a05257446d17b604fe8bbe771cb6f7825591f80c3e04b600d86691af6b6813f2bf66d61f73a6e136957",
          address: "4646de9912997d27e088e0b274368a1e3cfb38b1"
        }
      },
      %{
        label: "4",
        input: %{priv: "5ace269e616cbb4830f6a0bdaac2892a8735481bf64a4d38dc1e6cc0077fef9f"},
        expected: %{
          pub:
            "8ce52426e137a34b81716861af5b3272a55ddffc1e24b5b8966cb790f822f409d432dc9551d079481d9632c15e69b5e7dca3ae7db9a0977dbe7d4ba555fa904d",
          address: "284583d5c009a34ddb8814f435d8b21c070c4196"
        }
      }
    ]

    for %{label: label, input: input, expected: expected} <- tests do
      @label label
      @input input
      @expected expected
      test "#{@label}:" do
        "04" <> pub =
          @input.priv
          |> ETH.Utils.decode16()
          |> ETH.Utils.get_public_key()
          |> ETH.Utils.encode16()

        assert pub == @expected.pub

        address = ("04" <> pub) |> ETH.Utils.get_address()
        "0x" <> address = address |> String.downcase()

        assert address == @expected.address
      end
    end
  end
end
