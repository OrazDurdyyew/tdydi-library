-- Rollback initial schema

-- updated_at triggers
DROP TRIGGER IF EXISTS faq_items_updated_at ON faq_items;
DROP TRIGGER IF EXISTS faq_categories_updated_at ON faq_categories;
DROP TRIGGER IF EXISTS pages_updated_at ON pages;
DROP TRIGGER IF EXISTS posts_updated_at ON posts;
DROP TRIGGER IF EXISTS periodical_issues_updated_at ON periodical_issues;
DROP TRIGGER IF EXISTS periodicals_updated_at ON periodicals;
DROP TRIGGER IF EXISTS collections_updated_at ON collections;
DROP TRIGGER IF EXISTS books_updated_at ON books;
DROP TRIGGER IF EXISTS sections_updated_at ON sections;
DROP TRIGGER IF EXISTS site_settings_updated_at ON site_settings;
DROP TRIGGER IF EXISTS admins_updated_at ON admins;
DROP FUNCTION IF EXISTS set_updated_at();

-- view_logs
DROP TABLE IF EXISTS view_logs;
DROP TYPE IF EXISTS view_entity_type;

-- home
DROP TABLE IF EXISTS home_collections;
DROP TABLE IF EXISTS home_sections;
DROP TABLE IF EXISTS home_slides;

-- faq
DROP TABLE IF EXISTS faq_items;
DROP TABLE IF EXISTS faq_categories;

-- pages
DROP TABLE IF EXISTS pages;

-- posts
DROP TABLE IF EXISTS posts;
DROP TYPE IF EXISTS post_status;

-- periodicals
DROP TABLE IF EXISTS periodical_issues;
DROP TABLE IF EXISTS periodicals;
DROP TYPE IF EXISTS periodical_type;

-- collections
DROP TABLE IF EXISTS collection_books;
DROP TABLE IF EXISTS collections;

-- book associations
DROP TABLE IF EXISTS book_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS book_sections;
DROP TABLE IF EXISTS books;

-- permissions
DROP TABLE IF EXISTS admin_section_permissions;

-- sections (с триггером)
DROP TRIGGER IF EXISTS sections_max_depth_trigger ON sections;
DROP FUNCTION IF EXISTS sections_check_max_depth();
DROP TABLE IF EXISTS sections;

-- site settings
DROP TABLE IF EXISTS site_settings;

-- admins
DROP TABLE IF EXISTS admins;
DROP TYPE IF EXISTS admin_role;
