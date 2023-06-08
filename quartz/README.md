# quartz

[quartz](https://github.com/reugn/go-quartz) is scheduling library inspired by Java quartz scheduler.

In this demo we can see a basic usage.
You have to edit code to see messages from function displayJobs. 

## Missing features

Quartz use log package, so we cannot use another logging library - [https://github.com/reugn/go-quartz/issues/21](https://github.com/reugn/go-quartz/issues/21)

Missing retry on failure, when RunOnceTrigger is used.
