#!/usr/bin/env python
import asyncio
import logging
from typing import Callable, List, Type

from temporalio.client import Client
from temporalio.worker import Worker

import config
import lib

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def main(workflows: List[Type], activities: List[Callable]):
    client = await Client.connect(config.SERVER_ADDR, namespace=config.NAMESPACE)
    worker = Worker(
        client,
        task_queue=config.TASK_QUEUE,
        workflows=workflows,
        activities=activities,
    )
    await worker.run()


if __name__ == "__main__":
    workflows, activities = lib.get_all_workflows_and_activities()
    asyncio.run(main(workflows, activities))
