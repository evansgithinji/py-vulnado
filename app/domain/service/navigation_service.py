from app.domain.entity.redirect_request import RedirectRequest
from app.domain.repository.redirect_policy_repository import RedirectPolicyRepository
from app.domain.usecase.url_resolver import UrlResolver


class NavigationService:
    def __init__(
        self,
        policy_repo: RedirectPolicyRepository,
        url_resolver: UrlResolver,
    ):
        self._policy_repo = policy_repo
        self._url_resolver = url_resolver

    def resolve_redirect(self, request: RedirectRequest) -> str:
        """
        Deep call graph for Open Redirect:
        Handler → NavigationService.resolve_redirect()
          → RedirectPolicyRepository.getPolicy() → RedirectPolicy
          → UrlResolver.resolveUrl() → finalUrl (VULNERABLE: allows external)
        """
        policy = self._policy_repo.get_policy(request.context)
        return self._url_resolver.resolve_url(request.target_url, policy)
