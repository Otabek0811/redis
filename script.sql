


WITH market_product AS (
    SELECT
        JSON_AGG(
            JSON_BUILD_OBJECT (
                'id', p.id,
                'name', p.name,
                'price', p.price,
                'category_id', p.category_id,
                'created_at', p.created_at,
                'updated_at', p.updated_at
            )
        )  AS products,
        mpr.market_id AS market_id

    FROM product AS p
    JOIN market_product_relation AS mpr ON mpr.product_id = p.id
    WHERE mpr.market_id = '01e6e80c-a2e4-4323-bb34-bc33dcd7ae67'
    GROUP BY mpr.market_id
)
SELECT
    m.id,
    m.name,
    m.address,
    m.phone_number,
    m.created_at,
    m.updated_at,

    mp.products
    
FROM market AS m
JOIN market_product AS mp ON mp.market_id = m.id
WHERE m.id =  '01e6e80c-a2e4-4323-bb34-bc33dcd7ae67'






