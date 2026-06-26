from app.domain.entity.http_request_dto import HttpRequestDto
from app.domain.repository.url_policy_repository import UrlPolicyRepository
from app.domain.usecase.request_builder import RequestBuilder
from app.domain.usecase.http_client_adapter import HttpClientAdapter


class ExternalRequestService:
    def __init__(
        self,
        policy_repo: UrlPolicyRepository,
        request_builder: RequestBuilder,
        http_client: HttpClientAdapter,
    ):
        self._policy_repo = policy_repo
        self._request_builder = request_builder
        self._http_client = http_client

    def execute_request(self, dto: HttpRequestDto) -> dict:
        """
        Deep call graph for SSRF:
        Handler → ExternalRequestService.execute_request()
          → UrlPolicyRepository.getPolicy() → UrlPolicy
          → RequestBuilder.buildRequest() → PreparedRequest (VULNERABLE: no URL validation)
          → HttpClientAdapter.execute() → response (VULNERABLE: fetches any URL)
        """
        policy = self._policy_repo.get_policy("fetch")
        prepared = self._request_builder.build_request(dto, policy)
        return self._http_client.execute(prepared)

    def proxy_request(self, url: str) -> str:
        dto = HttpRequestDto(url=url, method="GET")
        policy = self._policy_repo.get_policy("proxy")
        prepared = self._request_builder.build_request(dto, policy)
        return self._http_client.execute_raw(prepared)

    def test_webhook(self, url: str, payload: dict) -> dict:
        import json
        dto = HttpRequestDto(url=url, method="POST", body=json.dumps(payload))
        policy = self._policy_repo.get_policy("webhook")
        prepared = self._request_builder.build_request(dto, policy)
        return self._http_client.execute(prepared)

    def fetch_remote_file(self, url: str) -> str:
        dto = HttpRequestDto(url=url, method="GET")
        policy = self._policy_repo.get_policy("remote_file")
        prepared = self._request_builder.build_request(dto, policy)
        return self._http_client.execute_raw(prepared)
