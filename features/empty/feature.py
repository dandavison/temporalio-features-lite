from temporalio import workflow

RESULT = "result"


@workflow.defn
class EmptyWorkflow:
    @workflow.run
    async def run(self) -> str:
        return RESULT
