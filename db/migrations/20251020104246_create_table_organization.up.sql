CREATE TABLE organization_members
(
    id             VARCHAR(100) PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    position       VARCHAR(100) NOT NULL,
    position_order INT          NOT NULL DEFAULT 99,
    order_index    INT                   DEFAULT 1,
    is_active      BOOLEAN               DEFAULT TRUE,
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);