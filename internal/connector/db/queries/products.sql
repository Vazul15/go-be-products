-- name: CreateProduct :one
INSERT INTO products (
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active
)
VALUES (
    sqlc.arg(public_id),
    sqlc.arg(name),
    sqlc.arg(slug),
    sqlc.arg(brand),
    sqlc.arg(category),
    sqlc.arg(gender),
    sqlc.arg(material),
    sqlc.arg(short_description),
    sqlc.arg(description),
    sqlc.arg(tags),
    sqlc.arg(price),
    sqlc.arg(discount_price),
    sqlc.arg(discount_expires_at),
    sqlc.arg(currency),
    sqlc.arg(main_image_url),
    sqlc.arg(gallery),
    sqlc.arg(type),
    sqlc.arg(is_featured),
    sqlc.arg(is_active)
)
RETURNING
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at;


-- name: GetProductByPublicID :one
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE public_id = sqlc.arg(public_id)
  AND is_active = TRUE;


-- name: CreateProductVariant :one
INSERT INTO product_variants (
    public_id,
    product_id,
    sku,
    size,
    color,
    color_code,
    price,
    barcode
)
VALUES (
    sqlc.arg(public_id),
    sqlc.arg(product_id),
    sqlc.arg(sku),
    sqlc.arg(size),
    sqlc.arg(color),
    sqlc.arg(color_code),
    sqlc.arg(price),
    sqlc.arg(barcode)
)
RETURNING
    id,
    public_id,
    product_id,
    sku,
    size,
    color,
    color_code,
    price,
    barcode,
    is_active,
    created_at,
    updated_at;


-- name: GetProductIDByPublicID :one
SELECT id
FROM products
WHERE public_id = sqlc.arg(public_id)
  AND is_active = TRUE;


-- name: GetProductVariantByPublicID :one
SELECT
    v.id,
    v.public_id,
    v.product_id,
    v.sku,
    v.size,
    v.color,
    v.color_code,
    v.price,
    v.barcode,
    v.is_active,
    v.created_at,
    v.updated_at
FROM product_variants AS v
JOIN products AS p
    ON p.id = v.product_id
WHERE v.public_id = sqlc.arg(public_id)
  AND v.is_active = TRUE
  AND p.is_active = TRUE;


-- name: UpdateProductVariant :one
UPDATE product_variants
SET
    sku = sqlc.arg(sku),
    size = sqlc.arg(size),
    color = sqlc.arg(color),
    color_code = sqlc.arg(color_code),
    price = sqlc.arg(price),
    barcode = sqlc.arg(barcode),
    updated_at = CURRENT_TIMESTAMP
WHERE product_variants.public_id = sqlc.arg(public_id)
  AND product_variants.is_active = TRUE
  AND EXISTS (
      SELECT 1
      FROM products
      WHERE products.id = product_variants.product_id
        AND products.is_active = TRUE
  )
RETURNING
    id,
    public_id,
    product_id,
    sku,
    size,
    color,
    color_code,
    price,
    barcode,
    is_active,
    created_at,
    updated_at;


-- name: SoftDeleteProductVariant :one
UPDATE product_variants
SET
    is_active = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE public_id = sqlc.arg(public_id)
  AND is_active = TRUE
RETURNING public_id;


-- name: RestoreProductVariant :one
UPDATE product_variants
SET
    is_active = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE public_id = sqlc.arg(public_id)
  AND is_active = FALSE
RETURNING
    id,
    public_id,
    product_id,
    sku,
    size,
    color,
    color_code,
    price,
    barcode,
    is_active,
    created_at,
    updated_at;


-- name: ListProductVariantsByProductPublicID :many
SELECT
    v.id,
    v.public_id,
    v.product_id,
    v.sku,
    v.size,
    v.color,
    v.color_code,
    v.price,
    v.barcode,
    v.is_active,
    v.created_at,
    v.updated_at
FROM product_variants AS v
JOIN products AS p
    ON p.id = v.product_id
WHERE p.public_id = sqlc.arg(product_public_id)
  AND p.is_active = TRUE
  AND v.is_active = TRUE
ORDER BY
    v.size,
    v.color;


-- name: CountActiveProducts :one
SELECT COUNT(*)
FROM products
WHERE is_active = TRUE;


-- name: CountProducts :one
SELECT COUNT(*)
FROM products
WHERE is_active = TRUE
  AND (
      sqlc.arg(search)::text = ''
      OR name ILIKE '%' || sqlc.arg(search)::text || '%'
      OR brand ILIKE '%' || sqlc.arg(search)::text || '%'
      OR category ILIKE '%' || sqlc.arg(search)::text || '%'
      OR description ILIKE '%' || sqlc.arg(search)::text || '%'
      OR array_to_string(tags, ' ') ILIKE '%' || sqlc.arg(search)::text || '%'
  )
  AND (
      sqlc.arg(category)::text = ''
      OR category = sqlc.arg(category)::text
  )
  AND (
      sqlc.arg(gender)::text = ''
      OR gender = sqlc.arg(gender)::text
  )
  AND (
      sqlc.arg(type)::text = ''
      OR type = sqlc.arg(type)::text
  );


-- name: ListProducts :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE is_active = TRUE
  AND (
      sqlc.arg(search)::text = ''
      OR name ILIKE '%' || sqlc.arg(search)::text || '%'
      OR brand ILIKE '%' || sqlc.arg(search)::text || '%'
      OR category ILIKE '%' || sqlc.arg(search)::text || '%'
      OR description ILIKE '%' || sqlc.arg(search)::text || '%'
      OR array_to_string(tags, ' ') ILIKE '%' || sqlc.arg(search)::text || '%'
  )
  AND (
      sqlc.arg(category)::text = ''
      OR category = sqlc.arg(category)::text
  )
  AND (
      sqlc.arg(gender)::text = ''
      OR gender = sqlc.arg(gender)::text
  )
  AND (
      sqlc.arg(type)::text = ''
      OR type = sqlc.arg(type)::text
  )
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);


-- name: ListActiveProducts :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);


-- name: ListActiveProductsByType :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE type = sqlc.arg(type)
  AND is_active = TRUE
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);


-- name: ListActiveProductsByCategory :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE category = sqlc.arg(category)
  AND is_active = TRUE
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);


-- name: ListActiveProductsByGender :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE gender = sqlc.arg(gender)
  AND is_active = TRUE
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);


-- name: SearchActiveProducts :many
SELECT
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at
FROM products
WHERE is_active = TRUE
  AND (
      name ILIKE '%' || sqlc.arg(search) || '%'
      OR brand ILIKE '%' || sqlc.arg(search) || '%'
      OR category ILIKE '%' || sqlc.arg(search) || '%'
      OR description ILIKE '%' || sqlc.arg(search) || '%'
      OR array_to_string(tags, ' ') ILIKE '%' || sqlc.arg(search) || '%'
  )
ORDER BY created_at DESC;


-- name: UpdateProduct :one
UPDATE products
SET
    name = sqlc.arg(name),
    slug = sqlc.arg(slug),
    brand = sqlc.arg(brand),
    category = sqlc.arg(category),
    gender = sqlc.arg(gender),
    material = sqlc.arg(material),
    short_description = sqlc.arg(short_description),
    description = sqlc.arg(description),
    tags = sqlc.arg(tags),
    price = sqlc.arg(price),
    discount_price = sqlc.arg(discount_price),
    discount_expires_at = sqlc.arg(discount_expires_at),
    currency = sqlc.arg(currency),
    main_image_url = sqlc.arg(main_image_url),
    gallery = sqlc.arg(gallery),
    type = sqlc.arg(type),
    is_featured = sqlc.arg(is_featured),
    is_active = sqlc.arg(is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE public_id = sqlc.arg(public_id)
RETURNING
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at;


-- name: SoftDeleteProduct :one
UPDATE products
SET
    is_active = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE public_id = sqlc.arg(public_id)
  AND is_active = TRUE
RETURNING public_id;


-- name: RestoreProduct :one
UPDATE products
SET
    is_active = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE public_id = sqlc.arg(public_id)
  AND is_active = FALSE
RETURNING
    id,
    public_id,
    name,
    slug,
    brand,
    category,
    gender,
    material,
    short_description,
    description,
    tags,
    price,
    discount_price,
    discount_expires_at,
    currency,
    main_image_url,
    gallery,
    type,
    is_featured,
    is_active,
    created_at,
    updated_at;