from datetime import timedelta

from temporalio import activity, workflow


@activity.defn
async def say_hello(name: str) -> str:
    return f"Hello, {name}!"


@workflow.defn
class SayHelloWorkflow:
    @workflow.run
    async def run(self) -> str:
        return await workflow.execute_activity(
            say_hello, "Temporal", start_to_close_timeout=timedelta(seconds=5)
        )
