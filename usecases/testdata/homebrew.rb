class Goxzst < Formula
  desc "Go crossbuild, zip, shasum and templates"
  homepage "https://github.com/int128/goxzst"
  url "https://github.com/int128/goxzst/releases/download/{{ env "VERSION" }}/goxzst_darwin_amd64.zip"
  version "{{ env "VERSION" }}"
  sha256 "{{ sha256 .darwin_amd64_archive }}"
  def install
    bin.install "goxzst"
  end
  test do
    system "#{bin}/goxzst -h"
  end
end
