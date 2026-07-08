import subprocess
import tempfile
import os
from lxml import etree


class XmlUseCase:
    def parse_xml(self, xml_content: str) -> str:
        # VULNERABLE: XXE (CWE-611)
        # Entity resolution enabled, DTD loading enabled
        parser = etree.XMLParser(resolve_entities=True, load_dtd=True)
        root = etree.fromstring(xml_content.encode("utf-8"), parser)
        return etree.tostring(root, pretty_print=True).decode("utf-8")

    def validate_xml(self, xml_content: str) -> str:
        # VULNERABLE: XXE (CWE-611)
        # Uses xmllint with --noent which expands entities
        with tempfile.NamedTemporaryFile(
            mode="w", suffix=".xml", delete=False
        ) as tmp:
            tmp.write(xml_content)
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

    def transform_xml(self, xml_content: str) -> dict:
        # VULNERABLE: XXE (CWE-611)
        parser = etree.XMLParser(resolve_entities=True, load_dtd=True)
        root = etree.fromstring(xml_content.encode("utf-8"), parser)

        result = {}
        for element in root.iter():
            if element.text and element.text.strip():
                result[element.tag] = element.text.strip()

        return result
