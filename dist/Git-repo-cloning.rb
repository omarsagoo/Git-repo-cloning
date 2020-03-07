# This file was generated by GoReleaser. DO NOT EDIT.
class GitRepoCloning < Formula
  desc "MakeClones to clone repos from a google sheet"
  homepage "https://github.com/omarsagoo/Git-repo-cloning"
  version "1.0.9"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/omarsagoo/Git-repo-cloning/releases/download/v1.0.9/Git-repo-cloning_1.0.9_macOS-64bit.tar.gz"
    sha256 "72d0cf261f306d8ae8db24bcf932e9cea29ec6be14fb74ae7eab0f17bd9fedf1"
  elsif OS.linux?
  end

  def install
    bin.install "git-repo-cloning"
  end
end
