import os


class Config:
    def __init__(self):
        self.port = int(os.environ.get("PORT", "5001"))
        self.db_path = os.environ.get("DB_PATH", "/app/data/app.db")
        self.upload_dir = os.environ.get("UPLOAD_DIR", "/app/uploads")
        self.files_dir = os.environ.get("FILES_DIR", "/app/files")
        self.images_dir = os.environ.get("IMAGES_DIR", "/app/images")
        self.static_dir = os.environ.get("STATIC_DIR", "/app/static")
        self.exports_dir = os.environ.get("EXPORTS_DIR", "/app/exports")
        self.backups_dir = os.environ.get("BACKUPS_DIR", "/app/backups")
        self.secret_key = os.environ.get("SECRET_KEY", "vulnerable-secret-key-12345")
        self.api_key = os.environ.get("API_KEY", "ak_vuln_python_test_key_98765")
        self.database_password = os.environ.get("DATABASE_PASSWORD", "db_p@ssw0rd_python!")
        self.admin_password = os.environ.get("ADMIN_PASSWORD", "admin_super_secret")
