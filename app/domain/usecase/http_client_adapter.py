import requests
from app.domain.entity.http_request_dto import PreparedRequest


class HttpClientAdapter:
    def execute(self, prepared: PreparedRequest) -> dict:
        # VULNERABLE: Fetches any URL without restriction (SSRF)
        if prepared.method.upper() == "POST":
            response = requests.post(
                prepared.url,
                headers=prepared.headers,
                data=prepared.body,
                timeout=prepared.timeout,
            )
        else:
            response = requests.get(
                prepared.url,
                headers=prepared.headers,
                timeout=prepared.timeout,
            )
        return {
            "status_code": response.status_code,
            "headers": dict(response.headers),
            "body": response.text,
        }

    def execute_raw(self, prepared: PreparedRequest) -> str:
        # VULNERABLE: Returns raw response body (SSRF)
        response = requests.get(prepared.url, timeout=prepared.timeout)
        return response.text
