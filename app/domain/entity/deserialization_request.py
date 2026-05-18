from dataclasses import dataclass


@dataclass
class DeserializationRequest:
    data: str = ""
    format: str = "pickle"


@dataclass
class SerializationFormat:
    name: str = ""
    mime_type: str = ""
    is_safe: bool = False  # VULNERABLE: unsafe formats allowed
