# user-service 测试报告

- **日期**：2026-05-07
- **服务**：user-service
- **执行人**：Kimi

## 测试环境

| 项目 | 值 |
|------|-----|
| Go 版本 | 1.24.1 |
| PostgreSQL | 18.3 |
| Redis | 8.6.3 |

## 测试结果

```
go test ./handler -v
=== RUN   TestWxLogin
--- PASS: TestWxLogin (0.74s)
=== RUN   TestGetProfile
--- PASS: TestGetProfile (0.10s)
=== RUN   TestRecharge
--- PASS: TestRecharge (0.07s)
=== RUN   TestDeductBalance
--- PASS: TestDeductBalance (0.06s)
=== RUN   TestRefundBalance
--- PASS: TestRefundBalance (0.06s)
=== RUN   TestListPets
--- PASS: TestListPets (0.06s)
=== RUN   TestAddPet
--- PASS: TestAddPet (0.06s)
PASS
```

## 通过的测试用例

- [x] TestWxLogin — Mock 微信登录
- [x] TestGetProfile — 获取用户信息（含会员等级）
- [x] TestRecharge — 充值成功 + 无效金额失败
- [x] TestDeductBalance — 余额扣款（内部接口）
- [x] TestRefundBalance — 余额退款（内部接口）
- [x] TestListPets — 宠物列表
- [x] TestAddPet — 添加宠物

## 失败的测试用例

- 无

## 已知问题

1. 微信支付参数暂未接入真实 API（Mock 阶段）
2. 充值赠送金额固定为 20%，后续需接入 activity-service 查询动态折扣

## 结论

- 服务是否可进入下一阶段：是
- 阻塞项：无
