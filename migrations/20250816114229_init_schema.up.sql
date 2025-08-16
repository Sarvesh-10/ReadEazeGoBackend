-- Initial schema

CREATE TABLE public.alembic_version (
    version_num character varying(32) NOT NULL
);

CREATE TABLE public.book_user_profiles (
    id integer NOT NULL,
    user_id integer,
    book_id integer,
    book_name text NOT NULL,
    current_page integer DEFAULT 0,
    total_pages integer DEFAULT 0,
    read_percentage real GENERATED ALWAYS AS (
        CASE
            WHEN (total_pages > 0) THEN round((((current_page)::numeric / (total_pages)::numeric) * (100)::numeric), 2)
            ELSE (0)::numeric
        END
    ) STORED,
    mode text,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT book_user_profiles_mode_check CHECK ((mode = ANY (ARRAY['study'::text, 'casual'::text])))
);

CREATE SEQUENCE public.book_user_profiles_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.book_user_profiles_id_seq OWNED BY public.book_user_profiles.id;

CREATE TABLE public.books (
    id integer NOT NULL,
    user_id integer,
    name text NOT NULL,
    file_data bytea NOT NULL,
    uploaded_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    cover_image bytea
);

CREATE SEQUENCE public.books_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;

CREATE TABLE public.refresh_tokens (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token_hash text NOT NULL,
    revoked boolean DEFAULT false,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE SEQUENCE public.refresh_tokens_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    name text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE SEQUENCE public.user_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.user_id_seq OWNED BY public.users.id;

-- Defaults
ALTER TABLE ONLY public.book_user_profiles ALTER COLUMN id SET DEFAULT nextval('public.book_user_profiles_id_seq'::regclass);
ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);
ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);

-- Constraints
ALTER TABLE ONLY public.alembic_version ADD CONSTRAINT alembic_version_pkc PRIMARY KEY (version_num);
ALTER TABLE ONLY public.book_user_profiles ADD CONSTRAINT book_user_profiles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.books ADD CONSTRAINT books_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.refresh_tokens ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.refresh_tokens ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash);
ALTER TABLE ONLY public.users ADD CONSTRAINT user_email_key UNIQUE (email);
ALTER TABLE ONLY public.users ADD CONSTRAINT user_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users ADD CONSTRAINT user_username_key UNIQUE (name);

-- Foreign keys
ALTER TABLE ONLY public.book_user_profiles
    ADD CONSTRAINT book_user_profiles_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.book_user_profiles
    ADD CONSTRAINT book_user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
