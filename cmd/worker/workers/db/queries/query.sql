-- name: DecrStock :execrows
UPDATE seckill_activities
SET stock = stock - 1
WHERE sku_id = ? AND stock > 0;

-- name: InsertOrder :exec
INSERT INTO seckill_orders (user_id, activity_id, sku_id, price)
VALUES (?, ?, ?, ?);
