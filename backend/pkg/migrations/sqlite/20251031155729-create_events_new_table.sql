-- +migrate Up

-- 1. إنشاء جدول جديد بالشكل الصحيح
CREATE TABLE events_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id TEXT NOT NULL,
    title TEXT DEFAULT NULL,
    description TEXT NOT NULL,
    time DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- 2. نقل البيانات من الجدول القديم
INSERT INTO events_new (group_id, title, description, time, created_at)
SELECT group_id, title, description, time, created_at FROM events;

-- 3. حذف الجدول القديم
DROP TABLE events;

-- 4. إعادة تسمية الجدول الجديد
ALTER TABLE events_new RENAME TO events;
-- +migrate Down
DROP TABLE IF EXISTS events;
