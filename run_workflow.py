import asyncio

from workflows import SayHello
from temporalio.client import Client


async def main():
    client = await Client.connect("localhost:7233")

    result = await client.execute_workflow(
        SayHello.run, "Temporal", id="hello-workflow-id-1", task_queue="hello-task-queue"
    )

    print(f"Result: {result}")


if __name__ == "__main__":
    asyncio.run(main())
