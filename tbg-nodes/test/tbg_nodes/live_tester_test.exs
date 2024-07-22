defmodule TbgNodes.LiveTesterTest do
  use TbgNodes.DataCase

  @spec get_infra_api_k8s :: TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s
  def get_infra_api_k8s do
    TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s
  end

  describe "live tests" do
    @tag timeout: :infinity
    @tag :k8s
    test "test_permissioned_network_creation" do
      TbgNodes.LiveTester.init_ets()

      Ecto.Adapters.SQL.Sandbox.unboxed_run(
        TbgNodes.Repo,
        fn ->
          assert TbgNodes.LiveTester.test_permissioned_network_creation(&get_infra_api_k8s/0) ==
                   true
        end
      )
    end

    test "maybe_create_notification" do
      assert TbgNodes.LiveTester.maybe_create_notification(
               false,
               true,
               "test_permissioned_network_creation"
             ) =~ "failed"

      assert TbgNodes.LiveTester.maybe_create_notification(
               false,
               false,
               "test_permissioned_network_creation"
             ) =~ "failed"

      assert TbgNodes.LiveTester.maybe_create_notification(
               true,
               true,
               "test_permissioned_network_creation"
             ) =~ "succeeded"

      assert TbgNodes.LiveTester.maybe_create_notification(
               true,
               false,
               "test_permissioned_network_creation"
             ) == nil
    end
  end
end
