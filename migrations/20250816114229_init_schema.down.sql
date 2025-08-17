-- Drop schema (reverse order of creation)

ALTER TABLE IF EXISTS public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey;
ALTER TABLE IF EXISTS public.books DROP CONSTRAINT IF EXISTS books_user_id_fkey;
ALTER TABLE IF EXISTS public.book_user_profiles DROP CONSTRAINT IF EXISTS book_user_profiles_user_id_fkey;
ALTER TABLE IF EXISTS public.book_user_profiles DROP CONSTRAINT IF EXISTS book_user_profiles_book_id_fkey;

DROP TABLE IF EXISTS public.book_user_profiles CASCADE;
DROP TABLE IF EXISTS public.books CASCADE;
DROP TABLE IF EXISTS public.refresh_tokens CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;
DROP TABLE IF EXISTS public.alembic_version CASCADE;

DROP SEQUENCE IF EXISTS public.book_user_profiles_id_seq CASCADE;
DROP SEQUENCE IF EXISTS public.books_id_seq CASCADE;
DROP SEQUENCE IF EXISTS public.refresh_tokens_id_seq CASCADE;
DROP SEQUENCE IF EXISTS public.user_id_seq CASCADE;
