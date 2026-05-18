import os
import json
from flask import Flask, jsonify

from app.config import Config
from app.adapter.persistence.database import init_database
from app.adapter.persistence.sqlite_user_repository import SQLiteUserRepository
from app.adapter.persistence.sqlite_product_repository import SQLiteProductRepository
from app.adapter.persistence.sqlite_order_repository import SQLiteOrderRepository
from app.adapter.persistence.sqlite_message_repository import SQLiteMessageRepository
from app.adapter.persistence.sqlite_review_repository import SQLiteReviewRepository
from app.domain.usecase.auth_usecase import AuthUseCase
from app.domain.usecase.product_usecase import ProductUseCase
from app.domain.usecase.order_usecase import OrderUseCase

# Deep call graph imports - Original 11 vuln types
from app.domain.repository.sql_query_policy_repository import SqlQueryPolicyRepository
from app.domain.usecase.search_query_validator import SearchQueryValidator
from app.domain.usecase.sql_query_builder import SqlQueryBuilder
from app.domain.usecase.sql_query_executor import SqlQueryExecutor
from app.domain.usecase.search_result_mapper import SearchResultMapper
from app.domain.usecase.html_response_builder import HtmlResponseBuilder
from app.domain.service.catalog_service import CatalogService

from app.domain.repository.command_policy_repository import CommandPolicyRepository
from app.domain.usecase.command_builder import CommandBuilder
from app.domain.usecase.shell_executor import ShellExecutor
from app.domain.service.system_command_service import SystemCommandService

from app.domain.repository.file_access_policy_repository import FileAccessPolicyRepository
from app.domain.usecase.path_resolver import PathResolver
from app.domain.usecase.file_reader import FileReader
from app.domain.service.file_service import FileService

from app.domain.repository.url_policy_repository import UrlPolicyRepository
from app.domain.usecase.request_builder import RequestBuilder
from app.domain.usecase.http_client_adapter import HttpClientAdapter
from app.domain.service.external_request_service import ExternalRequestService

from app.domain.repository.content_policy_repository import ContentPolicyRepository
from app.domain.usecase.content_processor import ContentProcessor
from app.domain.service.content_service import ContentService

from app.domain.repository.xml_parser_config_repository import XmlParserConfigRepository
from app.domain.usecase.xml_parser_factory import XmlParserFactory
from app.domain.usecase.xml_document_processor import XmlDocumentProcessor
from app.domain.service.xml_processing_service import XmlProcessingService

from app.domain.repository.format_detector_repository import FormatDetectorRepository
from app.domain.usecase.deserializer_factory import DeserializerFactory
from app.domain.usecase.object_mapper import ObjectMapper
from app.domain.service.serialization_service import SerializationService

from app.domain.repository.redirect_policy_repository import RedirectPolicyRepository
from app.domain.usecase.url_resolver import UrlResolver
from app.domain.service.navigation_service import NavigationService

from app.domain.repository.template_config_repository import TemplateConfigRepository
from app.domain.usecase.template_compiler import TemplateCompiler
from app.domain.usecase.template_engine_adapter import TemplateEngineAdapter
from app.domain.service.notification_service import NotificationService

from app.domain.repository.credential_repository import CredentialRepository
from app.domain.usecase.credential_validator import CredentialValidator
from app.domain.usecase.session_manager import SessionManager
from app.domain.service.authentication_service import AuthenticationService

# Deep call graph imports - 7 SAST-correlated types (existing)
from app.domain.repository.ldap_user_repository import InMemoryLdapUserRepository
from app.domain.repository.xml_document_repository import InMemoryXmlDocumentRepository
from app.domain.repository.header_policy_repository import (
    InMemoryHeaderPolicyRepository,
    InMemoryLocaleRepository,
)
from app.domain.repository.audit_policy_repository import InMemoryAuditPolicyRepository
from app.domain.repository.document_collection_repository import InMemoryDocumentCollectionRepository
from app.domain.repository.rule_repository import InMemoryRuleRepository

from app.adapter.persistence.ldap_connection_adapter import LdapConnectionAdapter, LdapFilterBuilder
from app.domain.usecase.xpath_evaluator import XPathExpressionBuilder, XPathEvaluator
from app.domain.usecase.header_processor import (
    HeaderValueProcessor,
    ResponseHeaderWriter,
    CookieManager,
    RedirectBuilder,
)
from app.domain.usecase.log_formatter import (
    AuditEventEnricher,
    LogFormatter,
    LogWriter,
    LogStorageAdapter,
)
from app.domain.usecase.nosql_executor import QueryBuilder, DocumentQueryExecutor, ExpressionEvaluator
from app.domain.usecase.expression_evaluator import (
    ExpressionPreprocessor,
    ExpressionEvaluatorSink,
    FormulaBuilder,
    ResultFormatter,
)

from app.domain.service.directory_service import DirectoryService
from app.domain.service.xml_auth_service import XmlAuthService
from app.domain.service.response_customization_service import (
    ResponseCustomizationService,
    LocalizationService,
)
from app.domain.service.audit_service import AuditService
from app.domain.service.profile_service import ProfileService
from app.domain.service.pricing_engine import PricingEngine

from app.adapter.handler import (
    auth_handler,
    product_handler,
    order_handler,
    message_handler,
    review_handler,
    review_export,
    file_handler,
    network_handler,
    template_handler,
    xml_handler,
    admin_handler,
    deserialize_handler,
    calculator_handler,
    log_handler,
    header_handler,
    ldap_handler,
    xpath_handler,
    nosql_handler,
)


def ensure_directories(cfg: Config):
    dirs = [
        cfg.upload_dir,
        cfg.files_dir,
        cfg.images_dir,
        cfg.static_dir,
        cfg.exports_dir,
        cfg.backups_dir,
        os.path.dirname(cfg.db_path),
    ]
    for d in dirs:
        os.makedirs(d, exist_ok=True)

    create_sample_files(cfg)


def create_sample_files(cfg: Config):
    sample_txt = os.path.join(cfg.files_dir, "sample.txt")
    if not os.path.exists(sample_txt):
        with open(sample_txt, "w") as f:
            f.write("This is a sample file for testing file operations.\n")

    config_json = os.path.join(cfg.files_dir, "config.json")
    if not os.path.exists(config_json):
        with open(config_json, "w") as f:
            json.dump(
                {
                    "app_name": "Vulnerable Python App",
                    "version": "1.0.0",
                    "debug": True,
                    "database": {"host": "localhost", "port": 5432, "password": "db_secret_123"},
                    "api_keys": {"internal": "sk-internal-key-12345"},
                },
                f,
                indent=2,
            )

    readme_txt = os.path.join(cfg.static_dir, "readme.txt")
    if not os.path.exists(readme_txt):
        with open(readme_txt, "w") as f:
            f.write("Static files directory for the vulnerable Python application.\n")

    placeholder = os.path.join(cfg.images_dir, "placeholder.txt")
    if not os.path.exists(placeholder):
        with open(placeholder, "w") as f:
            f.write("Image placeholder file.\n")

    sample_xml = os.path.join(cfg.files_dir, "sample.xml")
    if not os.path.exists(sample_xml):
        with open(sample_xml, "w") as f:
            f.write(
                '<?xml version="1.0" encoding="UTF-8"?>\n'
                "<products>\n"
                "  <product><name>Laptop Pro</name><price>1299.99</price></product>\n"
                "  <product><name>Wireless Mouse</name><price>29.99</price></product>\n"
                "</products>\n"
            )

    config_xml = os.path.join(cfg.files_dir, "config.xml")
    if not os.path.exists(config_xml):
        with open(config_xml, "w") as f:
            f.write(
                '<?xml version="1.0" encoding="UTF-8"?>\n'
                "<config>\n"
                "  <database>\n"
                "    <host>db.internal</host>\n"
                "    <port>5432</port>\n"
                "    <username>admin</username>\n"
                "    <password>super_secret_db_pass</password>\n"
                "  </database>\n"
                "  <api-key>sk-prod-key-abc123</api-key>\n"
                "</config>\n"
            )


def create_app() -> Flask:
    cfg = Config()
    ensure_directories(cfg)

    conn = init_database(cfg.db_path)

    # Repositories (existing)
    user_repo = SQLiteUserRepository(conn)
    product_repo = SQLiteProductRepository(conn)
    order_repo = SQLiteOrderRepository(conn)
    message_repo = SQLiteMessageRepository(conn)
    review_repo = SQLiteReviewRepository(conn)

    # Use Cases (kept for backwards compat where needed)
    auth_uc = AuthUseCase(user_repo)
    product_uc = ProductUseCase(product_repo)
    order_uc = OrderUseCase(order_repo)

    # === Deep Call Graph Wiring - Original 11 Vulnerability Types ===

    # SQL Injection - 6 layer call graph
    sql_policy_repo = SqlQueryPolicyRepository()
    search_validator = SearchQueryValidator()
    sql_query_builder = SqlQueryBuilder()
    sql_executor = SqlQueryExecutor(conn)
    result_mapper = SearchResultMapper()
    catalog_service = CatalogService(
        sql_policy_repo, search_validator, sql_query_builder, sql_executor, result_mapper
    )

    # Reflected XSS - reuses CatalogService + HtmlResponseBuilder
    html_builder = HtmlResponseBuilder()

    # Command Injection - 5 layer call graph
    cmd_policy_repo = CommandPolicyRepository()
    cmd_builder = CommandBuilder()
    shell_executor = ShellExecutor()
    system_command_service = SystemCommandService(cmd_policy_repo, cmd_builder, shell_executor)

    # Path Traversal - 5 layer call graph
    file_policy_repo = FileAccessPolicyRepository(cfg.files_dir, cfg.upload_dir)
    path_resolver = PathResolver()
    file_reader = FileReader()
    file_service = FileService(file_policy_repo, path_resolver, file_reader)

    # SSRF - 5 layer call graph
    url_policy_repo = UrlPolicyRepository()
    request_builder = RequestBuilder()
    http_client = HttpClientAdapter()
    external_request_service = ExternalRequestService(url_policy_repo, request_builder, http_client)

    # Stored XSS - 5 layer call graph
    content_policy_repo = ContentPolicyRepository()
    content_processor = ContentProcessor()
    content_service = ContentService(content_policy_repo, content_processor, message_repo, review_repo)

    # XXE - 5 layer call graph
    xml_config_repo = XmlParserConfigRepository()
    xml_parser_factory = XmlParserFactory()
    xml_doc_processor = XmlDocumentProcessor()
    xml_processing_service = XmlProcessingService(xml_config_repo, xml_parser_factory, xml_doc_processor)

    # Insecure Deserialization - 5 layer call graph
    format_repo = FormatDetectorRepository()
    deserializer_factory = DeserializerFactory()
    object_mapper = ObjectMapper()
    serialization_service = SerializationService(format_repo, deserializer_factory, object_mapper)

    # Open Redirect - 4 layer call graph
    redirect_policy_repo = RedirectPolicyRepository()
    url_resolver = UrlResolver()
    navigation_service = NavigationService(redirect_policy_repo, url_resolver)

    # SSTI - 5 layer call graph
    template_config_repo = TemplateConfigRepository()
    template_compiler = TemplateCompiler()
    template_engine = TemplateEngineAdapter()
    notification_service = NotificationService(template_config_repo, template_compiler, template_engine)

    # Broken Auth / Info Disclosure - 5 layer call graph
    credential_repo = CredentialRepository(user_repo)
    credential_validator = CredentialValidator()
    session_manager = SessionManager()
    authentication_service = AuthenticationService(credential_repo, credential_validator, session_manager)

    # === Deep Call Graph Wiring - 7 SAST-Correlated Types (existing) ===

    # LDAP Injection - 6 layer call graph
    ldap_repo = InMemoryLdapUserRepository()
    ldap_conn_adapter = LdapConnectionAdapter()
    ldap_filter_builder = LdapFilterBuilder()
    directory_service = DirectoryService(ldap_repo, ldap_conn_adapter, ldap_filter_builder)

    # XPath Injection - 5 layer call graph
    xml_doc_repo = InMemoryXmlDocumentRepository()
    xpath_expr_builder = XPathExpressionBuilder()
    xpath_evaluator = XPathEvaluator()
    xml_auth_service = XmlAuthService(xml_doc_repo, xpath_expr_builder, xpath_evaluator)

    # Header Injection - 5 layer call graph
    header_policy_repo = InMemoryHeaderPolicyRepository()
    locale_repo = InMemoryLocaleRepository()
    header_value_processor = HeaderValueProcessor()
    header_writer = ResponseHeaderWriter()
    cookie_manager = CookieManager()
    redirect_builder_header = RedirectBuilder()
    customization_service = ResponseCustomizationService(
        header_policy_repo, header_value_processor, header_writer
    )
    localization_service = LocalizationService(locale_repo, cookie_manager, redirect_builder_header)

    # Log Injection - 6 layer call graph
    audit_policy_repo = InMemoryAuditPolicyRepository()
    audit_enricher = AuditEventEnricher()
    log_formatter = LogFormatter()
    log_writer = LogWriter()
    log_storage = LogStorageAdapter()
    audit_service = AuditService(
        audit_policy_repo, audit_enricher, log_formatter, log_writer, log_storage
    )

    # NoSQL Injection - 6 layer call graph
    doc_collection_repo = InMemoryDocumentCollectionRepository()
    nosql_query_builder = QueryBuilder()
    doc_executor = DocumentQueryExecutor()
    expression_evaluator = ExpressionEvaluator()
    profile_service = ProfileService(
        doc_collection_repo, nosql_query_builder, doc_executor, expression_evaluator
    )

    # Code Injection / Calculator - 5 layer call graph
    rule_repo = InMemoryRuleRepository()
    expr_preprocessor = ExpressionPreprocessor()
    expr_evaluator_sink = ExpressionEvaluatorSink()
    formula_builder = FormulaBuilder()
    result_formatter = ResultFormatter()
    pricing_engine = PricingEngine(
        rule_repo, expr_preprocessor, expr_evaluator_sink, formula_builder, result_formatter
    )

    # === Init Handlers (deep call graph - original 11 types) ===
    product_handler.init(catalog_service, html_builder)
    auth_handler.init(authentication_service, navigation_service, html_builder, auth_uc)
    network_handler.init(system_command_service, external_request_service, html_builder)
    admin_handler.init(system_command_service, navigation_service, cfg.backups_dir)
    file_handler.init(file_service, external_request_service, system_command_service, cfg.files_dir, cfg.upload_dir)
    message_handler.init(content_service)
    review_handler.init(content_service)
    review_export.init(content_service)
    xml_handler.init(xml_processing_service)
    deserialize_handler.init(serialization_service)
    template_handler.init(notification_service)

    # Init handlers (existing - unchanged)
    order_handler.init(order_uc)

    # Init handlers (deep call graph - 7 SAST-correlated types)
    calculator_handler.init(pricing_engine)
    log_handler.init(audit_service)
    header_handler.init(customization_service, localization_service)
    ldap_handler.init(directory_service)
    xpath_handler.init(xml_auth_service)
    nosql_handler.init(profile_service)

    # Flask app
    app = Flask(__name__)
    app.config["SECRET_KEY"] = cfg.secret_key

    # Register blueprints
    app.register_blueprint(auth_handler.bp)
    app.register_blueprint(product_handler.bp)
    app.register_blueprint(order_handler.bp)
    app.register_blueprint(message_handler.bp)
    app.register_blueprint(review_handler.bp)
    app.register_blueprint(review_export.bp)
    app.register_blueprint(file_handler.bp)
    app.register_blueprint(network_handler.bp)
    app.register_blueprint(template_handler.bp)
    app.register_blueprint(xml_handler.bp)
    app.register_blueprint(admin_handler.bp)
    app.register_blueprint(deserialize_handler.bp)
    app.register_blueprint(calculator_handler.bp)
    app.register_blueprint(log_handler.bp)
    app.register_blueprint(header_handler.bp)
    app.register_blueprint(ldap_handler.bp)
    app.register_blueprint(xpath_handler.bp)
    app.register_blueprint(nosql_handler.bp)

    @app.route("/health")
    def health():
        return jsonify({"status": "ok"})

    @app.route("/")
    def index():
        return jsonify({
            "app": "Invicti Vulnerable Python App",
            "version": "1.0.0",
            "framework": "Flask",
            "endpoints": {
                "health": "/health",
                "login": "POST /api/login",
                "users": "GET /api/users",
                "users_search": "GET /users/search?q=",
                "products": "GET /api/products",
                "products_search": "GET /api/products/search?q=",
                "orders": "GET /api/orders?user_id=",
                "messages": "GET /api/messages",
                "messages_board": "GET /messages",
                "reviews": "GET /api/products/{id}/reviews",
                "files": "GET /api/files?filename=",
                "files_list": "GET /api/files/list",
                "files_fetch": "GET /api/files/fetch?url=",
                "ping": "POST /api/ping",
                "fetch": "GET /api/fetch?url=",
                "proxy": "GET /api/proxy?url=",
                "webhook": "POST /api/webhook/test",
                "greeting": "GET /api/greeting?name=",
                "template_render": "POST /api/template/render",
                "xml_parse": "POST /api/xml/parse",
                "xml_validate": "POST /api/xml/validate",
                "backup": "POST /api/backup",
                "debug": "GET /api/debug",
                "redirect": "GET /redirect?url=",
                "calculate": "GET /api/calculate?expr=",
                "cart_export": "POST /api/cart/export",
                "cart_import": "POST /api/cart/import",
            },
        })

    return app


if __name__ == "__main__":
    cfg = Config()
    app = create_app()
    app.run(host="0.0.0.0", port=cfg.port, debug=False)
