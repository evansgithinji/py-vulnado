import subprocess
import requests


class NetworkUseCase:
    def ping_host(self, host: str) -> str:
        # VULNERABLE: Command Injection (CWE-78)
        output = subprocess.check_output(
            f"ping -c 1 {host}", shell=True, stderr=subprocess.STDOUT
        )
        return output.decode("utf-8")

    def fetch_url(self, url: str) -> dict:
        # VULNERABLE: SSRF (CWE-918)
        response = requests.get(url, timeout=10)
        return {
            "status_code": response.status_code,
            "headers": dict(response.headers),
            "body": response.text,
        }

    def test_webhook(self, url: str, payload: dict) -> dict:
        # VULNERABLE: SSRF (CWE-918)
        response = requests.post(url, json=payload, timeout=10)
        return {
            "status_code": response.status_code,
            "body": response.text,
        }

    def proxy_request(self, url: str) -> str:
        # VULNERABLE: SSRF (CWE-918)
        response = requests.get(url, timeout=10)
        return response.text

    def check_port(self, host: str, port: str) -> str:
        # VULNERABLE: Command Injection (CWE-78)
        output = subprocess.check_output(
            f"nc -zv -w 3 {host} {port} 2>&1 || true", shell=True
        )
        return output.decode("utf-8")
