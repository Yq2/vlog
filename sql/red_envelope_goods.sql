
// 超卖现象的优化case语句
// 该方案会造成多余的数据查询
update red_envelope_goods r
set remain_quantity = case case_value
    when remain_quantity >= 扣减数量或金额 then
    remain_quantity-扣减数量或金额
    else remain_quantity end case
where id = ?


// update where 语句中添加条件判断
// 没有多余的查询语句
update red_envelope_goods r
set remain_quantity = remain_quantity-扣减数量或金额
where id = ? and remain_quantity >= 扣减数量或金额



