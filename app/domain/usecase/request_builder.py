from app.domain.entity.http_request_dto import HttpRequestDto, UrlPolicy, PreparedRequest


class RequestBuilder:
    def build_request(self, dto: HttpRequestDto, policy: UrlPolicy) -> PreparedRequest:
        # VULNERABLE: No URL validation - allows internal IPs, metadata endpoints
        return PreparedRequest(
            url=dto.url,
            method=dto.method,
            headers=dto.headers,
            body=dto.body,
            timeout=dto.timeout,
            policy=policy,
        )
