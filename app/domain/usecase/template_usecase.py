import jinja2


class TemplateUseCase:
    def render_greeting(self, name: str) -> str:
        # VULNERABLE: SSTI (CWE-1336)
        # User input directly used as Jinja2 template string
        env = jinja2.Environment()
        template = env.from_string(f"Hello, {name}! Welcome to our store.")
        return template.render()

    def render_invoice(self, template_str: str, **kwargs) -> str:
        # VULNERABLE: SSTI (CWE-1336)
        env = jinja2.Environment()
        template = env.from_string(template_str)
        return template.render(**kwargs)

    def render_custom(self, template_str: str) -> str:
        # VULNERABLE: SSTI (CWE-1336)
        env = jinja2.Environment()
        template = env.from_string(template_str)
        return template.render()
