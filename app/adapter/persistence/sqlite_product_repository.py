import sqlite3
from app.domain.entity.product import Product
from app.domain.repository.product_repository import ProductRepository


class SQLiteProductRepository(ProductRepository):
    def __init__(self, conn: sqlite3.Connection):
        self.conn = conn

    def find_all(self) -> list[Product]:
        cursor = self.conn.execute("SELECT * FROM products")
        return [self._row_to_product(row) for row in cursor.fetchall()]

    def find_by_id(self, product_id: int) -> Product | None:
        cursor = self.conn.execute(
            "SELECT * FROM products WHERE id = ?", (product_id,)
        )
        row = cursor.fetchone()
        if row:
            return self._row_to_product(row)
        return None

    def search(self, query: str) -> list[Product]:
        # VULNERABLE: SQL Injection (CWE-89)
        sql = f"SELECT * FROM products WHERE name LIKE '%{query}%'"
        cursor = self.conn.execute(sql)
        return [self._row_to_product(row) for row in cursor.fetchall()]

    def search_by_category(self, category: str) -> list[Product]:
        # VULNERABLE: SQL Injection (CWE-89)
        sql = f"SELECT * FROM products WHERE category='{category}'"
        cursor = self.conn.execute(sql)
        return [self._row_to_product(row) for row in cursor.fetchall()]

    def _row_to_product(self, row) -> Product:
        return Product(
            id=row["id"],
            name=row["name"],
            price=row["price"],
            stock=row["stock"],
            category=row["category"],
        )
