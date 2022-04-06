/* mysql  Ver 8.0.27 for Linux on x86_64 (MySQL Community Server - GPL) */

CREATE TABLE base64_store_image
(
    id INTEGER auto_increment,
    name VARCHAR(40), -- 上传图片时的名称
    code TEXT,  -- 图片base64编码
    PRIMARY KEY (id)
);