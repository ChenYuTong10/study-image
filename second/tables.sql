/* mysql  Ver 8.0.27 for Linux on x86_64 (MySQL Community Server - GPL) */

CREATE TABLE t_url_image
(
id INTEGER auto_increment,
path VARCHAR(100),
reference INTEGER DEFAULT 1,
PRIMARY KEY (id)
);