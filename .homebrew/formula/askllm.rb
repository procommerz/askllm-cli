class Askllm < Formula
  desc "A simple CLI interface for asking questions to the OpenAI GPT API from command line."
  homepage "https://github.com/procommerz/askllm-cli"
  version "1.0.0"

  if OS.mac?
    if Hardware::CPU.intel?
      url "https://github.com/procommerz/askllm-cli/releases/download/v1.0.0-macos-amd64/askllm.zip"
      sha256 "8135ee739bc2498bcc650ce8c414696d76ef29a5c3107b0446ff3c501ac12cf2"
    elsif Hardware::CPU.arm?
      url "https://github.com/procommerz/askllm-cli/releases/download/v1.0.0-macos-arch64/askllm.zip"
      sha256 "a16637fd20b080d5f1234f8dc84f4dc69e48302fb61de2b87402839e1ad32b63"
    end
#   elsif OS.linux?
#     url "https://github.com/yourusername/yourapp/releases/download/v1.0.0/yourapp-1.0.0-linux.tar.gz"
#     sha256 "SHA256_OF_LINUX_TAR_GZ"
  end

  def install
    bin.install "askllm"
  end

  test do
    system "#{bin}/askllm", "--version"
  end
end
