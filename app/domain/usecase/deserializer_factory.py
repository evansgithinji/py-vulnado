import pickle
import json
import base64
from app.domain.entity.deserialization_request import SerializationFormat


class DeserializerFactory:
    def get_deserializer(self, fmt: SerializationFormat):
        # VULNERABLE: Returns unsafe deserializer for pickle format
        if fmt.name == "json":
            return JsonDeserializer()
        return PickleDeserializer()  # VULNERABLE: unsafe by default


class PickleDeserializer:
    def deserialize(self, data: str):
        # VULNERABLE: pickle.loads on untrusted input
        raw = base64.b64decode(data)
        return pickle.loads(raw)

    def serialize(self, obj) -> str:
        data = pickle.dumps(obj)
        return base64.b64encode(data).decode("utf-8")


class JsonDeserializer:
    def deserialize(self, data: str):
        return json.loads(data)

    def serialize(self, obj) -> str:
        return json.dumps(obj)
