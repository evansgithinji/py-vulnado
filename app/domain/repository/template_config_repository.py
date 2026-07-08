from app.domain.entity.render_request import TemplateDefinition


class TemplateConfigRepository:
    def __init__(self):
        self._templates = {
            "greeting": TemplateDefinition(
                name="greeting",
                template_string="Hello, {user_input}! Welcome to our store.",
                engine="jinja2",
            ),
            "invoice": TemplateDefinition(
                name="invoice",
                template_string="Invoice for {{ customer }}: ${{ amount }}",
                engine="jinja2",
            ),
            "custom": TemplateDefinition(
                name="custom",
                template_string="",
                engine="jinja2",
            ),
        }

    def load_template(self, name: str) -> TemplateDefinition:
        return self._templates.get(name, self._templates["greeting"])
