class ObjectMapper:
    def deserialize(self, data: str, deserializer) -> object:
        # VULNERABLE: Delegates to unsafe deserializer (e.g. pickle)
        return deserializer.deserialize(data)

    def serialize(self, obj, deserializer) -> str:
        return deserializer.serialize(obj)
