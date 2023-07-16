#!/usr/bin/env python
import argparse
import asyncio
from typing import Type
from uuid import uuid4

from temporalio.client import Client

import config
from lib import get_workflows_and_activities


async def main(workflow_cls: Type):
    client = await Client.connect(config.SERVER_ADDR)

    result = await client.execute_workflow(
        workflow_cls.run,
        id=f"{workflow_cls.__temporal_workflow_definition.name}-{uuid4().hex[:4]}",
        task_queue=config.TASK_QUEUE,
    )

    print(f"Result: {result}")


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "feature", help="Path to a feature file, e.g., features/empty/feature.py"
    )
    return parser.parse_args()


if __name__ == "__main__":
    args = parse_args()
    [workflow_cls], _ = get_workflows_and_activities(args.feature)
    asyncio.run(main(workflow_cls))
