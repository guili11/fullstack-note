# TypeScript 使用范例

## 核心问题

1. 拿到了 API 文档，TS 类型怎么写？
2. 有了类型，怎么在项目里使用？

---

## 一、项目标准结构

新建文件夹 `src/types/course.ts` → 专门放课程所有 TS 类型

**类型只写 3 种：**
- 请求参数
- 响应数据
- 表单数据

> 对于 props 的类型，写到对应的子组件中就可以

---

## 二、响应数据类型

axios 会自动 `res.json()` + 取出 data（去掉 header），一般剩下这样的结构：

```json
{
  "code": 0,
  "msg": "success",
  "data": ...
}
```

### 泛型响应体

```ts
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
```

### 课程列表响应

```ts
export interface CourseListResponse {
  list: Course[]
  total: number
  page: number
  pageSize: number
}
```

---

## 三、类型使用方式

1. 指明函数的参数和返回值类型（API）
2. 变量 = 函数 → 类型推断

---

## 四、全链路 3 步使用

### 第 1 步：接口请求中

约束入参 + 返回值

```ts
export const getCourseList = (params: CourseQueryParams): Promise<CourseListResponse> => {
  return request.get('/course/list', { params })
}
```

---

### 第 2 步：Controller 自定义 Hook 中

```ts
// ✅ 约束列表数据
const [courseList, setCourseList] = useState<CourseItem[]>([])

// 请求数据 → 自动提示、自动校验
const fetchCourseList = useCallback(async (params: CourseQueryParams) => {
  const res = await getCourseList(params)
  setCourseList(res.list)
  setTotal(res.total)
}, [])
```

---

### 第 3 步：View 页面中

变量 = 函数 → 类型推断

```ts
const { courseList, searchParams } = useCourseController()
```
