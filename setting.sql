
CREATE DATABASE IF NOT EXISTS chatroom;

CREATE TABLE public.google_users IF NOT EXISTS(
	googleid varchar NOT NULL,
	first_name varchar NOT NULL,
	last_name varchar NOT NULL,
	email varchar NOT NULL,
	create_datetime timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT google_user_pk PRIMARY KEY (googleid)
);
CREATE INDEX google_user_email_idx ON public.google_users USING btree (email);

CREATE TABLE public.rooms IF NOT EXISTS(
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	admin_id varchar NOT NULL,
	"name" varchar NOT NULL,
	users_id _varchar NOT NULL,
	CONSTRAINT room_pk PRIMARY KEY (id),
	CONSTRAINT rooms_unique UNIQUE (name)
);
CREATE INDEX room_admin_idx ON public.rooms USING btree (admin_id);
CREATE INDEX room_users_idx ON public.rooms USING btree (users_id);

CREATE TABLE public.room_histories IF NOT EXISTS(
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	admin_id varchar NOT NULL,
	"name" varchar NOT NULL,
	users_id _varchar NOT NULL,
	create_datetime timestamptz DEFAULT now() NOT NULL
);
CREATE INDEX deleted_rooms_admin_id_idx ON public.room_histories USING btree (admin_id);
CREATE INDEX deleted_rooms_name_idx ON public.room_histories USING btree (name);

CREATE TABLE public.messages IF NOT EXISTS(
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"time" timestamptz DEFAULT now() NOT NULL,
	user_id varchar NOT NULL,
	"content" varchar NOT NULL,
	room_id uuid NOT NULL,
	deleted bool DEFAULT false NOT NULL,
	username varchar NULL,
	CONSTRAINT messages_pk PRIMARY KEY (id)
);
CREATE INDEX message_time_idx ON public.messages USING btree ("time");
CREATE INDEX message_user_idx ON public.messages USING btree (user_id);
