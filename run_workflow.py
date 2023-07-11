import asyncio

from workflows import SayHello
from temporalio.client import Client

ID = "say-hello-workflow-id"
TASK_QUEUE = "say-hello-task-queue"

async def main():
    client = await Client.connect("localhost:7233")

    result = await client.execute_workflow(
        SayHello.run, "Temporal", id=ID, task_queue=TASK_QUEUE
    )

    print(f"Result: {result}")


if __name__ == "__main__":
    asyncio.run(main())
