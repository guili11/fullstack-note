
sql 加速  事务  锁        
日志  sql引擎

sql 加速优化。     开发流程就是，日志发现问题，explain分析，去对sql或者索引优化。

1. 底层基础  原则
数据/索引怎么存在B+树中的，b+树的结构，特性     ， 数据是怎么存在硬盘里的？

查询一个数据的过程   sql执行顺序    数据们 怎么从硬盘中按照我们的需求，读出来的？

sql内优化器会 自己决定走不走索引， 自动决定where执行顺序  
优化器会根据「过滤性」选择索引     


2. 如何加速  索引    筛选，避免扫描  有序，避免排序
定一 优一
sql本身优化   索引优化

3. 观测，  如何发现问题
explain   日志
- 慢查询日志开启与配置
- 慢查询分析工具使用
- 如何从慢日志中定位核心问题SQL

long_query_time  1s
log_queries_not_using_indexes 测试开启

![alt text](image.png)
案例：
```
完全没问题！给你准备**2个【真实业务高频场景】的经典慢SQL案例**，严格按照你要的流程：
**日志发现慢SQL → EXPLAIN结果 → 分析问题 → 优化方案**
全是面试/工作必考，能彻底巩固你学的索引规则，直接抄录记录即可！

---

# 案例1：电商订单统计（高频：等值+范围+分组 + 函数导致索引失效）
## 基础信息
业务表：`order_info`（电商订单表）
字段：`user_id(用户ID)、order_time(下单时间)、order_money(订单金额)、status(订单状态)`
**原有错误索引**：`idx_user_time(user_id, order_time)`

---
### 步骤1：日志发现慢SQL
业务需求：统计**指定用户、2025年订单**的总金额
慢日志：执行耗时 2.8s，CPU 100%

SELECT user_id, SUM(order_money) total_money
FROM order_info
WHERE user_id = 1001 
  AND YEAR(order_time) = 2025
GROUP BY user_id;


### 步骤2：EXPLAIN 核心结果
| type | key  | rows  | Extra                              |
| ---- | ---- | ----- | ---------------------------------- |
| ALL  | NULL | 50000 | Using where; Using temporary; Using filesort |

### 步骤3：问题分析（你要记的核心）
1. **索引失效**：`YEAR(order_time)` 函数破坏索引，全表扫描
2. **性能双杀**：`Using temporary`(分组无序) + `Using filesort`(排序)
3. **索引设计不匹配**：筛选/分组顺序不符合「等值→范围→分组」规则

### 步骤4：优化方案（不改业务逻辑）
1. **SQL改写**：去掉函数，用时间范围
2. **索引重建**：遵循「等值→范围→分组」

-- 优化后SQL
SELECT user_id, SUM(order_money) total_money
FROM order_info
WHERE user_id = 1001 
  AND order_time BETWEEN '2025-01-01' AND '2025-12-31';

-- 最优索引
CREATE INDEX idx_user_time_money ON order_info(user_id, order_time, order_money);


### 步骤5：优化效果
type=ref，无临时表/无文件排序，耗时＜1ms

---

# 案例2：后台日志去重查询（高频：等值+范围+去重 + 临时表）
## 基础信息
业务表：`sys_log`（咱们一直用的系统日志表）
字段：`type(日志类型)、level(级别)、create_time(时间)、msg(内容)`
**原有索引**：`idx_type_level_time(type, level, create_time)`

---
### 步骤1：日志发现慢SQL
业务需求：查询**指定类型、级别>3**的**不重复日志时间**
慢日志：执行耗时 1.2s，触发临时表

SELECT DISTINCT create_time
FROM sys_log
WHERE type = 2 
  AND level > 3;


### 步骤2：EXPLAIN 核心结果
| type  | key                  | rows  | Extra                     |
| ----- | -------------------- | ----- | ------------------------- |
| range | idx_type_level_time  | 8000  | Using temporary; Using filesort |

### 步骤3：问题分析（你要记的核心）
1. **索引生效但有序性破坏**：`level>3` 范围查询，打乱 `create_time` 顺序
2. **必现临时表**：`DISTINCT` 去重需要有序数据，无序→被迫建临时表
3. **物理无法规避**：范围+右侧字段去重，是索引铁律

### 步骤4：优化方案（最小损耗优化）
不改SQL、不改业务，只做**覆盖索引优化**（降低开销）

-- 索引已存在，仅改写SQL为覆盖索引（无回表）
SELECT DISTINCT create_time
FROM sys_log
WHERE type = 2 
  AND level > 3;


### 步骤5：优化效果
保留 `Using temporary`，但新增 `Using index`（无回表），耗时降低80%

---

# 2个案例核心记忆点
1. **只要分组/去重 + 数据无序** → 必触发 `Using temporary`
2. **索引字段禁函数** → 函数一用，索引直接报废
3. **范围查询必破坏右侧有序性** → 无法彻底规避临时表/排序，只能降开销
```
