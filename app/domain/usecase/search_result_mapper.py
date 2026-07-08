from app.domain.entity.product import Product


class SearchResultMapper:
    def map_to_products(self, rows: list) -> list:
        products = []
        for row in rows:
            try:
                products.append(Product(
                    id=row.get("id", 0),
                    name=row.get("name", ""),
                    price=float(row.get("price", 0)),
                    category=row.get("category", ""),
                    stock=int(row.get("stock", 0)),
                ))
            except (ValueError, TypeError):
                continue
        return products
