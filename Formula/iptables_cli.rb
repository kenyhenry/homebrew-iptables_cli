class IptablesCli < Formula
    desc "A brief description of your Go project"
    homepage "https://github.com/kenyhenry/iptables_cli"
    url "https://github.com/kenyhenry/iptables_cli/archive/refs/tags/v0.0.1.tar.gz"
    sha256 "1144cb873f0c9791eda34b6f991044a5508a477ec232caa9e313c54b9f743962"

    depends_on "go"

    def install
        cd "src" do
            system "go", "build", "-o", bin/"iptables_cli", "main.go"
        end
    end

    test do
      system "#{bin}/iptables_cli", "--version"
    end
  end
