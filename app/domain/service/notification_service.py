from app.domain.entity.render_request import RenderRequest
from app.domain.repository.template_config_repository import TemplateConfigRepository
from app.domain.usecase.template_compiler import TemplateCompiler
from app.domain.usecase.template_engine_adapter import TemplateEngineAdapter


class NotificationService:
    def __init__(
        self,
        template_repo: TemplateConfigRepository,
        compiler: TemplateCompiler,
        engine: TemplateEngineAdapter,
    ):
        self._template_repo = template_repo
        self._compiler = compiler
        self._engine = engine

    def render_notification(self, request: RenderRequest) -> str:
        """
        Deep call graph for SSTI:
        Handler → NotificationService.render_notification()
          → TemplateConfigRepository.loadTemplate() → TemplateDefinition
          → TemplateCompiler.compile() → compiled (VULNERABLE: user input in template)
          → TemplateEngineAdapter.render() → output (VULNERABLE: SSTI)
        """
        definition = self._template_repo.load_template(request.template_name)
        compiled = self._compiler.compile(definition, request.user_input)
        return self._engine.render(compiled, request.variables)

    def render_invoice(self, template_str: str, variables: dict) -> str:
        definition = self._template_repo.load_template("invoice")
        if template_str:
            definition.template_string = template_str
        compiled = self._compiler.compile(definition, template_str)
        return self._engine.render(compiled, variables)
