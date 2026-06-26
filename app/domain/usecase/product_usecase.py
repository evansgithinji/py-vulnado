from app.domain.repository.product_repository import ProductRepository


class ProductUseCase:
    def __init__(self, product_repo: ProductRepository):
        self.product_repo = product_repo

    def get_all_products(self):
        return self.product_repo.find_all()

    def get_product(self, product_id: int):
        return self.product_repo.find_by_id(product_id)

    def search_products(self, query: str):
        return self.product_repo.search(query)

    def search_by_category(self, category: str):
        return self.product_repo.search_by_category(category)
