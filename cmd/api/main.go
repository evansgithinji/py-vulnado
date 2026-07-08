package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"goapp/internal/adapter/handler"
	"goapp/internal/adapter/persistence"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
	"goapp/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Ensure directories exist
	ensureDirectories(cfg)

	// Initialize database
	db, err := persistence.InitDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize SQLite repositories
	userRepo := persistence.NewSQLiteUserRepository(db)
	productRepo := persistence.NewSQLiteProductRepository(db)
	orderRepo := persistence.NewSQLiteOrderRepository(db)
	reviewRepo := persistence.NewSQLiteReviewRepository(db)
	messageRepo := persistence.NewSQLiteMessageRepository(db)

	// Initialize use cases (kept for backwards compat where handlers still need them)
	authUseCase := usecase.NewAuthUseCase(userRepo)
	_ = usecase.NewProductUseCase(productRepo) // kept but no longer passed to handler
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	_ = usecase.NewReviewUseCase(reviewRepo) // kept but no longer passed to handler
	_ = usecase.NewMessageUseCase(messageRepo)
	_ = usecase.NewFileUseCase(cfg.UploadDir, cfg.FilesDir, cfg.ImagesDir)
	_ = usecase.NewNetworkUseCase()
	_ = usecase.NewTemplateUseCase()
	_ = usecase.NewAdminUseCase()
	xmlUseCase := usecase.NewXmlUseCase("/tmp")

	// === Deep Call Graph Wiring - Original 11 Vulnerability Types ===

	// SQL Injection - 6 layer call graph
	sqlPolicyRepo := repository.NewSqlQueryPolicyRepository()
	searchValidator := usecase.NewSearchQueryValidator()
	sqlQueryBuilder := usecase.NewSqlQueryBuilder()
	sqlQueryExecutor := usecase.NewSqlQueryExecutor(db)
	searchResultMapper := usecase.NewSearchResultMapper()
	catalogService := service.NewCatalogService(
		sqlPolicyRepo, searchValidator, sqlQueryBuilder, sqlQueryExecutor, searchResultMapper,
	)

	// Reflected XSS - reuses CatalogService + HtmlResponseBuilder
	htmlBuilder := usecase.NewHtmlResponseBuilder()

	// Command Injection - 5 layer call graph
	cmdPolicyRepo := repository.NewCommandPolicyRepository()
	cmdBuilder := usecase.NewCommandBuilder()
	shellExecutor := usecase.NewShellExecutor()
	systemCommandService := service.NewSystemCommandService(cmdPolicyRepo, cmdBuilder, shellExecutor)

	// Path Traversal - 5 layer call graph
	filePolicyRepo := repository.NewFileAccessPolicyRepository(cfg.FilesDir, cfg.UploadDir)
	pathResolver := usecase.NewPathResolver()
	fileReader := usecase.NewFileReader()
	fileService := service.NewFileService(filePolicyRepo, pathResolver, fileReader)

	// SSRF - 5 layer call graph
	urlPolicyRepo := repository.NewUrlPolicyRepository()
	requestBuilder := usecase.NewRequestBuilder()
	httpClientAdapter := usecase.NewHttpClientAdapter()
	externalRequestService := service.NewExternalRequestService(urlPolicyRepo, requestBuilder, httpClientAdapter)

	// Stored XSS - 5 layer call graph
	contentPolicyRepo := repository.NewContentPolicyRepository()
	contentProcessor := usecase.NewContentProcessor()
	contentService := service.NewContentService(contentPolicyRepo, contentProcessor, messageRepo, reviewRepo)

	// XXE - 5 layer call graph
	xmlConfigRepo := repository.NewXmlParserConfigRepository()
	xmlParserFactory := usecase.NewXmlParserFactory()
	xmlDocProcessor := usecase.NewXmlDocumentProcessor("/tmp")
	xmlProcessingService := service.NewXmlProcessingService(xmlConfigRepo, xmlParserFactory, xmlDocProcessor)

	// Open Redirect - 4 layer call graph
	redirectPolicyRepo := repository.NewRedirectPolicyRepository()
	urlResolver := usecase.NewUrlResolver()
	navigationService := service.NewNavigationService(redirectPolicyRepo, urlResolver)

	// SSTI - 5 layer call graph
	templateConfigRepo := repository.NewTemplateConfigRepository()
	templateCompiler := usecase.NewTemplateCompiler()
	templateEngine := usecase.NewTemplateEngine()
	notificationService := service.NewNotificationService(templateConfigRepo, templateCompiler, templateEngine)

	// Broken Auth / Info Disclosure - 5 layer call graph
	credentialRepo := repository.NewCredentialRepository(db)
	credentialValidator := usecase.NewCredentialValidator()
	sessionManager := usecase.NewSessionManager()
	authenticationService := service.NewAuthenticationService(credentialRepo, credentialValidator, sessionManager)

	// --- LDAP Injection: layered wiring ---
	ldapUserRepo := repository.NewInMemoryLdapUserRepository()
	ldapConnAdapter := persistence.NewLdapConnectionAdapter()
	ldapFilterBuilder := persistence.NewLdapFilterBuilder()
	directoryService := service.NewDirectoryService(ldapUserRepo, ldapConnAdapter, ldapFilterBuilder)

	// --- XPath Injection: layered wiring ---
	xmlDocRepo := repository.NewInMemoryXmlDocumentRepository("/tmp")
	xpathExprBuilder := usecase.NewXPathExpressionBuilder()
	xpathEvaluator := usecase.NewXPathEvaluator()
	xmlAuthService := service.NewXmlAuthService(xmlDocRepo, xpathExprBuilder, xpathEvaluator)

	// --- Header Injection: layered wiring ---
	headerPolicyRepo := repository.NewInMemoryHeaderPolicyRepository()
	localeRepo := repository.NewInMemoryLocaleRepository()
	headerProcessor := usecase.NewHeaderValueProcessor()
	headerWriter := usecase.NewResponseHeaderWriter()
	redirectBuilder := usecase.NewRedirectBuilder()
	cookieManager := usecase.NewCookieManager()
	customizationService := service.NewResponseCustomizationService(headerPolicyRepo, headerProcessor, headerWriter)
	localizationService := service.NewLocalizationService(localeRepo, redirectBuilder, cookieManager)

	// --- Log Injection: layered wiring ---
	auditPolicyRepo := repository.NewInMemoryAuditPolicyRepository()
	auditEnricher := usecase.NewAuditEventEnricher()
	logFormatter := usecase.NewLogFormatter()
	logWriter := usecase.NewLogWriter()
	logStorage := usecase.NewLogStorageAdapter(100)
	auditService := service.NewAuditService(auditPolicyRepo, auditEnricher, logFormatter, logWriter, logStorage)

	// --- NoSQL Injection: layered wiring ---
	docCollectionRepo := repository.NewInMemoryDocumentCollectionRepository()
	queryBuilderNoSql := usecase.NewQueryBuilder()
	docQueryExecutor := usecase.NewDocumentQueryExecutor()
	exprEvaluator := usecase.NewExpressionEvaluator()
	profileService := service.NewProfileService(docCollectionRepo, queryBuilderNoSql, docQueryExecutor, exprEvaluator)

	// --- Calculator / Code Injection: layered wiring ---
	ruleRepo := repository.NewInMemoryRuleRepository()
	exprPreprocessor := usecase.NewExpressionPreprocessor()
	formulaBuilder := usecase.NewFormulaBuilder()
	evaluatorSink := usecase.NewExpressionEvaluatorSink()
	resultFormatter := usecase.NewResultFormatter()
	pricingEngine := service.NewPricingEngine(ruleRepo, exprPreprocessor, formulaBuilder, evaluatorSink, resultFormatter)

	// Initialize handlers (deep call graph - original 11 types)
	authHandler := handler.NewAuthHandler(authenticationService, navigationService, htmlBuilder, authUseCase)
	productHandler := handler.NewProductHandler(catalogService, htmlBuilder)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	reviewHandler := handler.NewReviewHandler(contentService)
	messageHandler := handler.NewMessageHandler(contentService)
	fileHandler := handler.NewFileHandler(fileService, externalRequestService, systemCommandService, cfg.FilesDir, cfg.UploadDir)
	networkHandler := handler.NewNetworkHandler(systemCommandService, externalRequestService, htmlBuilder)
	templateHandler := handler.NewTemplateHandler(notificationService)
	adminHandler := handler.NewAdminHandler(systemCommandService, navigationService, cfg.BackupsDir)
	xmlHandler := handler.NewXmlHandler(xmlProcessingService, xmlUseCase)

	// Initialize handlers (7 SAST-correlated types - unchanged)
	logHandler := handler.NewLogHandler(auditService)
	headerHandler := handler.NewHeaderHandler(customizationService, localizationService)
	ldapHandler := handler.NewLdapHandler(directoryService)
	xpathHandler := handler.NewXPathHandler(xmlAuthService)
	nosqlHandler := handler.NewNoSqlHandler(profileService)
	calculatorHandler := handler.NewCalculatorHandler(pricingEngine)

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Register routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":       "Invicti Vulnerable Go App",
			"version":   "1.0.0",
			"framework": "Gin",
			"endpoints": gin.H{
				"health":            "/health",
				"login":             "POST /api/login",
				"users_search":      "GET /users/search?q=",
				"products":          "GET /api/products",
				"products_search":   "GET /api/products/search?q=",
				"orders":            "GET /api/orders?user_id=",
				"messages":          "GET /api/messages",
				"messages_board":    "GET /messages",
				"reviews":           "GET /api/products/{id}/reviews",
				"files":             "GET /api/files?filename=",
				"ping":              "POST /api/ping",
				"fetch":             "GET /api/fetch?url=",
				"proxy":             "GET /api/proxy?url=",
				"webhook":           "POST /api/webhook/test",
				"greeting":          "GET /api/greeting?message=",
				"template_render":   "POST /api/template/render",
				"xml_parse":         "POST /api/xml/parse",
				"xml_validate":      "POST /api/xml/validate",
				"backup":            "POST /api/backup",
				"debug":             "GET /api/debug",
				"redirect":          "GET /redirect?url=",
				"calculate":         "GET /api/calculate?expr=",
			},
		})
	})

	authHandler.RegisterRoutes(r)
	productHandler.RegisterRoutes(r)
	orderHandler.RegisterRoutes(r)
	reviewHandler.RegisterRoutes(r)
	messageHandler.RegisterRoutes(r)
	fileHandler.RegisterRoutes(r)
	networkHandler.RegisterRoutes(r)
	templateHandler.RegisterRoutes(r)
	adminHandler.RegisterRoutes(r)
	xmlHandler.RegisterRoutes(r)
	logHandler.RegisterRoutes(r)
	headerHandler.RegisterRoutes(r)
	ldapHandler.RegisterRoutes(r)
	xpathHandler.RegisterRoutes(r)
	nosqlHandler.RegisterRoutes(r)
	calculatorHandler.RegisterRoutes(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Starting Vulnerable Web Application on port %s", cfg.Port)
	log.Printf("Database: %s", cfg.DBPath)
	log.Printf("WARNING: This application is intentionally vulnerable!")

	if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func ensureDirectories(cfg *config.Config) {
	dirs := []string{
		cfg.UploadDir,
		cfg.FilesDir,
		cfg.ImagesDir,
		cfg.StaticDir,
		cfg.ExportsDir,
		cfg.BackupsDir,
		"/app/data",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Could not create directory %s: %v", dir, err)
		}
	}

	// Create sample files for path traversal demo
	createSampleFiles(cfg)
}

func createSampleFiles(cfg *config.Config) {
	// Sample file in files directory
	sampleContent := "This is a sample file content.\nSensitive data: API_KEY=secret123\n"
	os.WriteFile(cfg.FilesDir+"/sample.txt", []byte(sampleContent), 0644)
	os.WriteFile(cfg.FilesDir+"/config.json", []byte(`{"secret": "db_password_123", "api_key": "sk-live-xxx"}`), 0644)

	// Sample static file
	os.WriteFile(cfg.StaticDir+"/readme.txt", []byte("Welcome to the static files directory!"), 0644)

	// Sample image placeholder
	os.WriteFile(cfg.ImagesDir+"/placeholder.txt", []byte("Image placeholder"), 0644)

	// Sample XML files
	sampleXML := `<?xml version="1.0"?>
<products>
	<product id="1">
		<name>Sample Product</name>
		<price>99.99</price>
		<category>Electronics</category>
	</product>
</products>`
	os.WriteFile(cfg.FilesDir+"/sample.xml", []byte(sampleXML), 0644)

	configXML := `<?xml version="1.0"?>
<config>
	<database>
		<host>localhost</host>
		<username>admin</username>
		<password>secret123</password>
	</database>
</config>`
	os.WriteFile(cfg.FilesDir+"/config.xml", []byte(configXML), 0644)
}
