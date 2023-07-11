import asyncio

from temporalio.client import Client
from workflows import SayHello
from run_workflow import ID, TASK_QUEUE


async def main():
    client = await Client.connect("localhost:7233")
    handle = client.get_workflow_handle(ID)
    print(handle)
    await handle.signal(SayHello.set_condition, True)


if __name__ == "__main__":
    asyncio.run(main())
