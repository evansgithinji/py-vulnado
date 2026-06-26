from app.domain.entity.deserialization_request import DeserializationRequest
from app.domain.repository.format_detector_repository import FormatDetectorRepository
from app.domain.usecase.deserializer_factory import DeserializerFactory
from app.domain.usecase.object_mapper import ObjectMapper


class SerializationService:
    def __init__(
        self,
        format_repo: FormatDetectorRepository,
        deserializer_factory: DeserializerFactory,
        object_mapper: ObjectMapper,
    ):
        self._format_repo = format_repo
        self._deserializer_factory = deserializer_factory
        self._object_mapper = object_mapper

    def deserialize(self, request: DeserializationRequest):
        """
        Deep call graph for Insecure Deserialization:
        Handler → SerializationService.deserialize()
          → FormatDetectorRepository.getFormat() → SerializationFormat
          → DeserializerFactory.getDeserializer() → deserializer (VULNERABLE)
          → ObjectMapper.deserialize() → object (VULNERABLE: pickle/eval)
        """
        fmt = self._format_repo.get_format(request.format)
        deserializer = self._deserializer_factory.get_deserializer(fmt)
        return self._object_mapper.deserialize(request.data, deserializer)

    def serialize(self, items, format_name: str = "pickle") -> str:
        fmt = self._format_repo.get_format(format_name)
        deserializer = self._deserializer_factory.get_deserializer(fmt)
        return self._object_mapper.serialize(items, deserializer)
