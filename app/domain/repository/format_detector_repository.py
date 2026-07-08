from app.domain.entity.deserialization_request import SerializationFormat


class FormatDetectorRepository:
    def __init__(self):
        self._formats = {
            "pickle": SerializationFormat(
                name="pickle",
                mime_type="application/x-pickle",
                is_safe=False,  # VULNERABLE
            ),
            "json": SerializationFormat(
                name="json",
                mime_type="application/json",
                is_safe=True,
            ),
        }

    def detect_format(self, data: str) -> SerializationFormat:
        # Simple heuristic - default to pickle (VULNERABLE)
        if data.strip().startswith("{") or data.strip().startswith("["):
            return self._formats["json"]
        return self._formats["pickle"]

    def get_format(self, name: str) -> SerializationFormat:
        return self._formats.get(name, self._formats["pickle"])
