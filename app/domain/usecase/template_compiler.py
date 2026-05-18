from app.domain.entity.render_request import TemplateDefinition


class TemplateCompiler:
    def compile(self, definition: TemplateDefinition, user_input: str) -> str:
        # VULNERABLE: Injects user input directly into template string
        if definition.name == "greeting":
            return f"Hello, {user_input}! Welcome to our store."
        elif definition.name == "custom":
            return user_input  # VULNERABLE: entire template from user
        else:
            return definition.template_string
