CREATE TABLE IF NOT EXISTS products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE, -- külső ID, API-hoz ideális

    -- Basic info
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    brand VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    gender VARCHAR(30) NOT NULL DEFAULT 'unisex',
    material VARCHAR(255),
    short_description TEXT,
    description TEXT,
    tags TEXT[],
    
    -- Pricing
    price BIGINT NOT NULL,
    discount_price BIGINT,
    discount_expires_at TIMESTAMPTZ,
    currency CHAR(3) NOT NULL DEFAULT 'HUF',

    -- Images
    main_image_url TEXT NOT NULL,
    gallery TEXT[],

    -- Flags
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT products_price_check
        CHECK (price > 0),

    CONSTRAINT products_discount_price_check
        CHECK (
            discount_price IS NULL
            OR (
                discount_price > 0
                AND discount_price < price
            )
        ),

    CONSTRAINT products_currency_format_check
        CHECK (currency ~ '^[A-Z]{3}$')
);

CREATE TABLE IF NOT EXISTS product_variants (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,

    product_id BIGINT NOT NULL,
    sku VARCHAR(100) NOT NULL UNIQUE,
    size VARCHAR(30) NOT NULL,
    color VARCHAR(100) NOT NULL,
    color_code VARCHAR(7),
    price BIGINT,
    barcode VARCHAR(100) UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT product_variants_product_fk
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT product_variants_price_check
        CHECK (price IS NULL OR price > 0),

    CONSTRAINT product_variants_color_code_check
        CHECK (
            color_code IS NULL
            OR color_code ~ '^#[0-9A-Fa-f]{6}$'
        ),

    CONSTRAINT product_variant_size_color_unique
        UNIQUE (product_id, size, color)
);

CREATE INDEX IF NOT EXISTS product_variants_product_id_idx
    ON product_variants(product_id);

CREATE INDEX IF NOT EXISTS products_active_created_at_idx
    ON products(created_at DESC)
    WHERE is_active = TRUE;