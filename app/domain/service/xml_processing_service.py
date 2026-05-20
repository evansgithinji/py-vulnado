from app.domain.entity.xml_processing_request import XmlProcessingRequest
from app.domain.repository.xml_parser_config_repository import XmlParserConfigRepository
from app.domain.usecase.xml_parser_factory import XmlParserFactory
from app.domain.usecase.xml_document_processor import XmlDocumentProcessor


class XmlProcessingService:
    def __init__(
        self,
        config_repo: XmlParserConfigRepository,
        parser_factory: XmlParserFactory,
        document_processor: XmlDocumentProcessor,
    ):
        self._config_repo = config_repo
        self._parser_factory = parser_factory
        self._document_processor = document_processor

    def process_xml(self, request: XmlProcessingRequest) -> str:
        """
        Deep call graph for XXE:
        Handler → XmlProcessingService.process_xml()
          → XmlParserConfigRepository.getConfig() → XmlParserConfig
          → XmlParserFactory.createParser() → parser (VULNERABLE: entities enabled)
          → XmlDocumentProcessor.parse() → result (VULNERABLE: XXE)
        """
        config = self._config_repo.get_config(request.config_name)
        parser = self._parser_factory.create_parser(config)
        return self._document_processor.parse(parser, request.xml_content)

    def validate_xml(self, request: XmlProcessingRequest) -> str:
        config = self._config_repo.get_config("validate")
        parser = self._parser_factory.create_parser(config)
        return self._document_processor.validate(parser, request.xml_content)

    def transform_xml(self, request: XmlProcessingRequest) -> dict:
        config = self._config_repo.get_config("transform")
        parser = self._parser_factory.create_parser(config)
        return self._document_processor.parse_to_dict(parser, request.xml_content)
