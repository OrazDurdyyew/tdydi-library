-- Remove font columns from site_settings table
ALTER TABLE site_settings 
DROP COLUMN font_file_path,
DROP COLUMN font_url,
DROP COLUMN font_uploaded_at;