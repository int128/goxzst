class Goxzst < Formula
  desc "Go crossbuild, zip, shasum and templates"
  homepage "https://github.com/int128/goxzst"
  url "https://github.com/int128/goxzst/releases/download/v1.0.0/goxzst_darwin_amd64.zip"
  version "v1.0.0"
  sha256 "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  def install
    bin.install "goxzst"
  end
  test do
    system "#{bin}/goxzst -h"
  end
end
