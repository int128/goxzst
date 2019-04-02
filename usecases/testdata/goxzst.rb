class Goxzst < Formula
  desc "Go crossbuild, zip, shasum and templates"
  homepage "https://github.com/int128/goxzst"
  url "https://github.com/int128/goxzst/releases/download/{{ .version }}/goxzst_darwin_amd64.zip"
  version "{{ .version }}"
  sha256 "{{ .darwin_amd64_zip_sha256 }}"
  def install
    bin.install "goxzst"
  end
  test do
    system "#{bin}/goxzst -h"
  end
end
