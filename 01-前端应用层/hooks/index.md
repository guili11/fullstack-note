# Hooks 核心概念

## 一、组件 vs Hook

除写 template + 按钮绑定交互，就是两个任务：**状态管理 + 副作用管理**

跟 Vue 不同，React 倾向于，尽可能将业务相关的所有状态和副作用都写在 Hook 里。

### 职责划分

```
组件 = 纯 UI 渲染函数
     （只负责：把状态变成页面，不负责管理状态）
     （有点函数在的也只是 button=>f，预留好位置，等着插入罢了）

Hook = 业务逻辑容器
     （只负责：管理状态 + 处理副作用，不负责渲染）

组件不关心逻辑怎么实现，只问 Hook 拿数据、调方法
```

---

## 二、自定义 Hook

### 定义

自定义 Hook 是一个以 `use` 开头的函数，内部可以调用任意内置 Hook，用来封装可复用的业务逻辑（状态 + 副作用）。

### 写法规范

1. 函数名必须 `use` 开头
2. 内部使用 React 内置 Hook
3. 返回需要的状态 / 方法（给组件使用）
4. 纯逻辑，不写 UI 代码

---

### 极简模板

```ts
// 1. 定义：use开头
export function useStudentController() {
  // 2. 内部：使用内置Hook（状态+副作用）
  const [studentList, setStudentList] = useState(null)

  useEffect(() => {
    // 副作用：请求接口
  }, [])

  // 3. 业务方法
  const deleteStudent = async (id) => {
    /* ... */
  }

  // 4. 返回：状态+方法
  return { studentList, deleteStudent }
}
```

---

### 组件中使用

```tsx
// 组件里直接调用，拿到状态和方法
function StudentPage() {
  const { studentList, deleteStudent } = useStudentController()

  return (
    <div>
      {studentList?.map((item) => (
        <div key={item.id}>{item.name}</div>
      ))}
    </div>
  )
}
```

---

## 三、实践原则

- **逻辑分 Hook**：每个 Hook 负责一个业务逻辑
- **Hook 之间通信**：使用参数或者回调函数，严禁在 Hook 之间为了通信而 useHook
- **聚合层 Hook**：唯一在 View 外可以 useHook 的是聚合层 Hook，本身没有任何业务逻辑，只为替代 View useXXX，然后整合多个子 Hook 给组件使用（还是子 Hook 的通信场所）

---

## 四、Hook 注意事项

### 状态隔离

- 一个组件内，同一个 Hook 只调用 1 次
- 每次调用自定义 Hook，都会创建一份全新的、独立的状态实例

> 跨组件共享状态：用 Context / 状态库，绝对不要靠「多次调用 Hook」共享状态，这种情况类似 store

---

### React 底层规则

**必须严格遵守：**

1. 只能在函数组件 / 自定义 Hook 中调用
   - 普通 JS 函数、类组件、回调函数里，不能直接用 Hook

2. 只能在顶层作用域调用
   - 不能放在 if / for / while / 嵌套函数里
   - React 靠调用顺序识别状态

3. 自定义 Hook 必须以 `use` 开头
   - 语法约定，ESLint 会强制校验

4. 禁止在条件 / 循环中调用 Hook
   - 保证每次渲染，Hook 的调用顺序完全一致

5. 依赖项数组必须完整声明
   - useEffect / useMemo / useCallback 用到的外部变量，必须放进依赖项
