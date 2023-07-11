import asyncio
from datetime import timedelta
from temporalio import workflow

with workflow.unsafe.imports_passed_through():
    from activities import say_hello


@workflow.defn
class SayHello:
    def __init__(self) -> None:
        super().__init__()
        self.condition = False

    @workflow.signal
    def set_condition(self, condition: bool) -> None:
        self.condition = not self.condition

    @workflow.run
    async def run(self, name: str) -> str:
        results = []
        if self.condition:
            results.append(await workflow.execute_activity(
                say_hello, "1", start_to_close_timeout=timedelta(seconds=5)
            ))
        await asyncio.sleep(5) # cond <- True
        if self.condition:
            results.append(await workflow.execute_activity(
                say_hello, "2", start_to_close_timeout=timedelta(seconds=5)
            ))
        await asyncio.sleep(5) # cond <- False
        if self.condition:
            results.append(await workflow.execute_activity(
                say_hello, "3", start_to_close_timeout=timedelta(seconds=5)
            ))
        await asyncio.sleep(5) # cond <- True
        if self.condition:
            results.append(await workflow.execute_activity(
                say_hello, "4", start_to_close_timeout=timedelta(seconds=5)
            ))
        return '\n'.join(results)