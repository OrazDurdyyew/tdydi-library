-- Add font file and URL columns to site_settings table
ALTER TABLE site_settings 
ADD COLUMN font_file_path TEXT,
ADD COLUMN font_url TEXT,
ADD COLUMN font_uploaded_at TIMESTAMP WITH TIME ZONE;