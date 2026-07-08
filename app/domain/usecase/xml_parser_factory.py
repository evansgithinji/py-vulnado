from lxml import etree
from app.domain.entity.xml_processing_request import XmlParserConfig


class XmlParserFactory:
    def create_parser(self, config: XmlParserConfig):
        # VULNERABLE: Creates XML parser with entities enabled
        return etree.XMLParser(
            resolve_entities=config.resolve_entities,
            load_dtd=config.load_dtd,
            no_network=config.no_network,
        )
