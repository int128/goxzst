class Goxzst < Formula
  desc "Go crossbuild, zip, shasum and templates"
  homepage "https://github.com/int128/goxzst"
  url "https://github.com/int128/goxzst/releases/download/v1.0.0/goxzst_darwin_amd64.zip"
  version "v1.0.0"
  sha256 "fe8df1a5a1980493ca9406ad3bb0e41297d979d90165a181fb39a5616a1c0789"
  def install
    bin.install "goxzst"
  end
  test do
    system "#{bin}/goxzst -h"
  end
end
