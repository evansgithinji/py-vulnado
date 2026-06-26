from app.domain.entity.xml_processing_request import XmlParserConfig


class XmlParserConfigRepository:
    def __init__(self):
        self._configs = {
            "default": XmlParserConfig(
                name="default",
                resolve_entities=True,  # VULNERABLE
                load_dtd=True,  # VULNERABLE
                no_network=False,  # VULNERABLE
            ),
            "validate": XmlParserConfig(
                name="validate",
                resolve_entities=True,
                load_dtd=True,
                no_network=False,
            ),
            "transform": XmlParserConfig(
                name="transform",
                resolve_entities=True,
                load_dtd=True,
                no_network=False,
            ),
        }

    def get_config(self, name: str) -> XmlParserConfig:
        return self._configs.get(name, self._configs["default"])
