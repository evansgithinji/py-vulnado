from abc import ABC, abstractmethod
from app.domain.entity.header_policy import HeaderPolicy, LocaleConfig


class HeaderPolicyRepository(ABC):
    @abstractmethod
    def get_policy(self, header_name: str) -> HeaderPolicy:
        ...

    @abstractmethod
    def get_all_policies(self) -> list:
        ...


class LocaleRepository(ABC):
    @abstractmethod
    def resolve_locale(self, lang: str) -> LocaleConfig:
        ...


class InMemoryHeaderPolicyRepository(HeaderPolicyRepository):
    def __init__(self):
        self._policies = {
            "X-Custom": HeaderPolicy(name="X-Custom", allowed_pattern=None, max_length=1024, allow_raw=True),
            "X-Request-Id": HeaderPolicy(name="X-Request-Id", allowed_pattern="^[a-zA-Z0-9-]+$", max_length=64),
            "X-Language": HeaderPolicy(name="X-Language", allowed_pattern=None, max_length=32, allow_raw=True),
        }
        self._default_policy = HeaderPolicy(name="default", allowed_pattern=None, max_length=2048, allow_raw=True)

    def get_policy(self, header_name: str) -> HeaderPolicy:
        return self._policies.get(header_name, self._default_policy)

    def get_all_policies(self) -> list:
        return list(self._policies.values())


class InMemoryLocaleRepository(LocaleRepository):
    def __init__(self):
        self._locales = {
            "en": LocaleConfig(code="en", name="English", direction="ltr"),
            "fr": LocaleConfig(code="fr", name="French", direction="ltr"),
            "ar": LocaleConfig(code="ar", name="Arabic", direction="rtl"),
            "ja": LocaleConfig(code="ja", name="Japanese", direction="ltr"),
        }

    def resolve_locale(self, lang: str) -> LocaleConfig:
        return self._locales.get(lang, LocaleConfig(code=lang, name=lang, direction="ltr"))
