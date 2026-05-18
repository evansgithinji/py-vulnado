import sqlite3
from app.domain.entity.order import Order
from app.domain.repository.order_repository import OrderRepository


class SQLiteOrderRepository(OrderRepository):
    def __init__(self, conn: sqlite3.Connection):
        self.conn = conn

    def find_by_user(self, user_id: str) -> list[Order]:
        # VULNERABLE: SQL Injection (CWE-89)
        sql = f"SELECT * FROM orders WHERE user_id={user_id}"
        cursor = self.conn.execute(sql)
        return [self._row_to_order(row) for row in cursor.fetchall()]

    def create(self, order: Order) -> Order:
        cursor = self.conn.execute(
            "INSERT INTO orders (user_id, product_id, quantity, total, status) VALUES (?, ?, ?, ?, ?)",
            (order.user_id, order.product_id, order.quantity, order.total, order.status),
        )
        self.conn.commit()
        order.id = cursor.lastrowid
        return order

    def find_all(self) -> list[Order]:
        cursor = self.conn.execute("SELECT * FROM orders")
        return [self._row_to_order(row) for row in cursor.fetchall()]

    def _row_to_order(self, row) -> Order:
        return Order(
            id=row["id"],
            user_id=row["user_id"],
            product_id=row["product_id"],
            quantity=row["quantity"],
            total=row["total"],
            status=row["status"],
        )
