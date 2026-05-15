# Frontend Rules（强约束版）

## 一、总原则

- 所有流程必须经过 controller
- 所有数据必须经过 store
- 所有请求必须经过 service

## 二、分层结构

view / controller / store / service

## 三、职责定义

### 1. view

- 只负责 UI 渲染
- 只负责用户事件（onClick / onSubmit）
- 只调用 controller

禁止：

- 发请求
- 写业务逻辑
- 操作 store
- 直接调用 service

### 2. controller（核心层）

- 负责业务流程编排
- 负责调用 service
- 负责操作 store
- 负责管理局部状态（useState）
- 负责统一错误处理
- 所有副作用（请求）必须放在 controller
- 所有流程必须通过 controller

禁止：

- 写 JSX

### 3. store

- 负责存数据（全局状态）
- 只提供 get / set 能力
- 常用字段：list / loading / current / open

禁止：

- 请求
- 业务逻辑
- 副作用
- async

### 4. service

- 负责 API 请求
- 每个接口一个函数
- 只返回 data
- 错误直接 throw

禁止：

- UI 逻辑
- 状态管理
- 业务流程

## 四、数据流

view → controller → service → controller → store → view

## 五、错误流

- service：throw error
- controller：handle error
- view：不处理错误
- store：不接触错误

统一写法：

```ts
try {
  await xxx()
} catch (e) {
  handleError(e)
}
```

## 六、请求规范

- 所有请求必须通过 service
- 所有 service 必须通过 request.ts
- 组件中禁止出现 axios / fetch

## 七、控制流规范

- 页面加载：在 controller 中触发（useEffect / onMounted）
- 用户操作：view → controller
- 禁止 render 过程中请求
- 禁止多处重复请求同一接口

## 八、store 使用规范

- store 只做 get 数据、set 数据
- 禁止写 if/else 业务逻辑
- 禁止写 async

## 九、函数归属判断

- UI → view
- 流程 → controller
- 数据 → store
- 请求 → service

## 十、标准实现模式

### service

```ts
export async function xxx() {
  const res = await request.xxx()

  if (res.code !== 0) {
    throw new Error(res.message)
  }

  return res.data
}
```

### controller

```ts
async function action() {
  try {
    const data = await service()
    store.setData(data)
  } catch (e) {
    handleError(e)
  }
}
```

### view

```tsx
<button onClick={action} />
```

## 十一、禁止事项

- 禁止 view → service
- 禁止 store → service
- 禁止 store 写业务逻辑
- 禁止 controller 写 UI
- 禁止组件内写请求
