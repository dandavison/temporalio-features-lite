import importlib
from pathlib import Path
from typing import Callable, List, Tuple, Type


def get_all_workflows_and_activities() -> Tuple[List[Type], List[Callable]]:
    workflows, activities = [], []
    for ww, aa in map(
        get_workflows_and_activities, Path("features").glob("**/feature.py")
    ):
        workflows.extend(ww)
        activities.extend(aa)
    return workflows, activities


def get_workflows_and_activities(feature: Path) -> Tuple[List[Type], List[Callable]]:
    module_name = str(feature).removesuffix(".py").replace("/", ".")
    module = importlib.import_module(module_name)
    workflows, activities = [], []
    for name in dir(module):
        obj = getattr(module, name)
        if hasattr(obj, "__temporal_workflow_definition"):
            workflows.append(obj)
        elif hasattr(obj, "__temporal_activity_definition"):
            activities.append(obj)

    return workflows, activities
