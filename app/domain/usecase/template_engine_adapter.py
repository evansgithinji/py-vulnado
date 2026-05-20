import jinja2


class TemplateEngineAdapter:
    def render(self, compiled_template: str, variables: dict = None) -> str:
        # VULNERABLE: Renders Jinja2 template with user-controlled content (SSTI)
        env = jinja2.Environment()
        template = env.from_string(compiled_template)
        return template.render(**(variables or {}))
