defmodule Blyver.Utils.Objects do
  @moduledoc """
  Behaviour for S3-compatible object stores.
  """

  @doc """
  List objects in a bucket, takes an optional argument for directory
  """
  @callback list_objects(String.t(), String.t()) :: {:ok, term} | {:error, String.t()}

  @doc """
  Get an object based on the path
  """
  @callback get_object(String.t(), String.t()) :: {:ok, term} | {:error, String.t()}

  @doc """
  Put an object in a bucket, on a specified filepath and returns the MD5 hash of the file
  """
  @callback put_object(String.t(), String.t(), term) :: {:ok, term} | {:error, String.t()}
end
