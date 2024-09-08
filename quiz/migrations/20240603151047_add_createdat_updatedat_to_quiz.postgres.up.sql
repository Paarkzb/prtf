ALTER TABLE IF EXISTS public.quiz ADD created_at timestamp;
ALTER TABLE IF EXISTS public.quiz ADD updated_at timestamp;

ALTER TABLE IF EXISTS public.quiz ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE IF EXISTS public.quiz ALTER COLUMN updated_at SET DEFAULT now();