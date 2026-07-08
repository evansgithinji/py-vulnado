from dataclasses import dataclass, field


@dataclass
class RenderRequest:
    user_input: str = ""
    template_name: str = "greeting"
    variables: dict = field(default_factory=dict)


@dataclass
class TemplateDefinition:
    name: str = ""
    template_string: str = ""
    engine: str = "jinja2"
