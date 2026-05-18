from dataclasses import dataclass


@dataclass
class XmlProcessingRequest:
    xml_content: str = ""
    operation: str = "parse"
    config_name: str = "default"


@dataclass
class XmlParserConfig:
    name: str = "default"
    resolve_entities: bool = True  # VULNERABLE: entities enabled
    load_dtd: bool = True  # VULNERABLE: DTD loading enabled
    no_network: bool = False  # VULNERABLE: allows network access
