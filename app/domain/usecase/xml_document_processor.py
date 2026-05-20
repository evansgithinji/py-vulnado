from lxml import etree


class XmlDocumentProcessor:
    def parse(self, parser, content: str) -> str:
        # VULNERABLE: Parses XML with entity-resolving parser (XXE)
        root = etree.fromstring(content.encode("utf-8"), parser)
        return etree.tostring(root, pretty_print=True).decode("utf-8")

    def parse_to_dict(self, parser, content: str) -> dict:
        # VULNERABLE: Parses XML with entity-resolving parser (XXE)
        root = etree.fromstring(content.encode("utf-8"), parser)
        result = {}
        for element in root.iter():
            if element.text and element.text.strip():
                result[element.tag] = element.text.strip()
        return result

    def validate(self, parser, content: str) -> str:
        # VULNERABLE: Validates XML with entity-resolving parser (XXE)
        import subprocess
        import tempfile
        import os
        with tempfile.NamedTemporaryFile(mode="w", suffix=".xml", delete=False) as tmp:
            tmp.write(content)
            tmp_path = tmp.name
        try:
            output = subprocess.check_output(
                f"xmllint --noent {tmp_path}",
                shell=True,
                stderr=subprocess.STDOUT,
            )
            return output.decode("utf-8")
        except subprocess.CalledProcessError as e:
            return e.output.decode("utf-8")
        finally:
            os.unlink(tmp_path)
