class HtmlResponseBuilder:
    def build_search_page(self, query: str, products: list) -> str:
        # VULNERABLE: Reflected XSS - no HTML escaping of query
        html = f"<h2>Products matching: {query}</h2>"
        html += "<table border='1'><tr><th>Name</th><th>Price</th><th>Category</th></tr>"
        for p in products:
            html += f"<tr><td>{p.name}</td><td>${p.price}</td><td>{p.category}</td></tr>"
        html += "</table>"
        return html

    def build_user_search_page(self, query: str, users: list) -> str:
        # VULNERABLE: Reflected XSS - no HTML escaping of query
        html = f"<h2>Search results for: {query}</h2><ul>"
        for u in users:
            html += f"<li>{u.username} - {u.email}</li>"
        html += "</ul>"
        return html

    def build_error_page(self, query: str, error: str) -> str:
        # VULNERABLE: Reflected XSS - no HTML escaping
        return f"<h2>Products matching: {query}</h2><p>Error: {error}</p>"

    def build_ping_page(self, host: str, output: str) -> str:
        # VULNERABLE: Reflected XSS - no HTML escaping
        return f"<h2>Ping results for: {host}</h2><pre>{output}</pre>"
