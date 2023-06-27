import asyncio

from temporalio import activity, workflow
from temporalio.client import Client
from temporalio.worker import Worker

import activities
import workflows


async def main():
    client = await Client.connect("localhost:7233", namespace="default")
    worker = Worker(
        client, task_queue="hello-task-queue", workflows=[workflows.SayHello], activities=[activities.say_hello]
    )
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
