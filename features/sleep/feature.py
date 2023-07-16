import asyncio

from temporalio import workflow

RESULT = "result"


@workflow.defn
class SleepWorkflow:
    @workflow.run
    async def run(self) -> str:
        await asyncio.sleep(2)
        return RESULT
