import asyncio
import logging

from temporalio.client import Client
from temporalio.worker import Worker

import activities
import workflows

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def main():
    logger.info("creating client...")
    client = await Client.connect("localhost:7233", namespace="default")
    logger.info("created client")
    worker = Worker(
        client, task_queue="hello-task-queue", workflows=[workflows.SayHello], activities=[activities.say_hello]
    )
    logger.info("created worker")
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
