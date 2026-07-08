from abc import ABC, abstractmethod
from app.domain.entity.calculation_rule import CalculationRule


class RuleRepository(ABC):
    @abstractmethod
    def get_rules(self, mode: str) -> list:
        ...


class InMemoryRuleRepository(RuleRepository):
    def __init__(self):
        self._rules = {
            "standard": [
                CalculationRule(name="identity", pattern=".*", transform=""),
            ],
            "discount": [
                CalculationRule(name="percent_off", pattern="percent", transform="price * (1 - {value}/100)"),
                CalculationRule(name="flat_off", pattern="flat", transform="price - {value}"),
                CalculationRule(name="custom", pattern="custom", transform="{formula}"),
            ],
            "scientific": [
                CalculationRule(name="allow_math", pattern="math\\.", transform=""),
                CalculationRule(name="allow_trig", pattern="sin|cos|tan", transform=""),
            ],
        }

    def get_rules(self, mode: str) -> list:
        return self._rules.get(mode, self._rules["standard"])
