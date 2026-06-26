from app.domain.repository.order_repository import OrderRepository
from app.domain.entity.order import Order


class OrderUseCase:
    def __init__(self, order_repo: OrderRepository):
        self.order_repo = order_repo

    def get_user_orders(self, user_id: str):
        return self.order_repo.find_by_user(user_id)

    def create_order(self, user_id: int, product_id: int, quantity: int, total: float):
        order = Order(
            user_id=user_id,
            product_id=product_id,
            quantity=quantity,
            total=total,
        )
        return self.order_repo.create(order)

    def get_all_orders(self):
        return self.order_repo.find_all()
