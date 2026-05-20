import pickle
import base64
import json


class DeserializeUseCase:
    def export_cart(self, items: list) -> str:
        # Serialize cart items using pickle and base64 encode
        data = pickle.dumps(items)
        return base64.b64encode(data).decode("utf-8")

    def import_cart(self, encoded_data: str) -> list:
        # VULNERABLE: Insecure Deserialization (CWE-502)
        # pickle.loads on untrusted input allows arbitrary code execution
        raw = base64.b64decode(encoded_data)
        return pickle.loads(raw)

    def export_cart_json(self, items: list) -> str:
        return json.dumps(items)

    def import_cart_json(self, json_data: str) -> list:
        return json.loads(json_data)
