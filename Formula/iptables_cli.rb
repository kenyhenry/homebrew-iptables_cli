class IptablesCli < Formula
    desc "A brief description of your Go project"
    homepage "https://github.com/kenyhenry/iptables_cli"
    url "https://github.com/kenyhenry/iptables_cli/archive/refs/tags/v0.0.1.tar.gz"
    sha256 "52f6284505ab46e9f864882f2f991247388e53a9bde1bb32391cafb6be6ad5db"

    depends_on "go"

    def install
      system "go", "build", "-o", bin/"iptables_cli", "src/main.go"
    end

    test do
      system "#{bin}/iptables_cli", "--version"
    end
  end
