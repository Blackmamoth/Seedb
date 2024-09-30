table "order_items" {
  schema = schema.public
  column "order_item_id" {
    null = false
    type = serial
  }
  column "order_id" {
    null = true
    type = integer
  }
  column "product_id" {
    null = true
    type = integer
  }
  column "quantity" {
    null = false
    type = integer
  }
  column "price_at_order" {
    null = false
    type = numeric(10,2)
  }
  primary_key {
    columns = [column.order_item_id]
  }
  foreign_key "order_items_order_id_fkey" {
    columns     = [column.order_id]
    ref_columns = [table.orders.column.order_id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "order_items_product_id_fkey" {
    columns     = [column.product_id]
    ref_columns = [table.products.column.product_id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "orders" {
  schema = schema.public
  column "order_id" {
    null = false
    type = serial
  }
  column "user_id" {
    null = true
    type = integer
  }
  column "order_date" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "total_amount" {
    null = true
    type = numeric(10,2)
  }
  primary_key {
    columns = [column.order_id]
  }
  foreign_key "orders_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.user_id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "products" {
  schema = schema.public
  column "product_id" {
    null = false
    type = serial
  }
  column "product_name" {
    null = false
    type = character_varying(100)
  }
  column "price" {
    null = false
    type = numeric(10,2)
  }
  primary_key {
    columns = [column.product_id]
  }
}
table "users" {
  schema = schema.public
  column "user_id" {
    null = false
    type = serial
  }
  column "username" {
    null = false
    type = character_varying(50)
  }
  column "email" {
    null = false
    type = character_varying(100)
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.user_id]
  }
  unique "users_email_key" {
    columns = [column.email]
  }
}
schema "public" {
  comment = "standard public schema"
}
