本质：llm没有记忆能力，为了让llm知道上下文，需要发送历史给llm
有点类似 jwt token哦，无状态，还不能llm里面存储

由此分出 短时记忆和长时记忆
## 1.短期记忆（Short-term）
- 保存当前任务或对话的上下文
- 作用：让 Agent 在多步操作中“记得之前发生了什么”
- 实现：Workflow state、Node 之间传递、每次调用 LLM 时把 history 作为 prompt
## 2.长期记忆（Long-term）
- 保存跨任务或跨会话的历史知识
- 作用：让 Agent 记住用户偏好、历史操作、经验数据
- 实现：数据库或向量检索，检索 relevant 历史再拼接到 prompt

embedding？