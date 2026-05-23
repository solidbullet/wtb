--
-- PostgreSQL database dump
--

\restrict tCnb2HqfJLixt2uMzemZReGSaNTev2BrnwdV6Sd25jSwuD5f3rbrsgnYUGCGjxq

-- Dumped from database version 18.3 (Homebrew)
-- Dumped by pg_dump version 18.3 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: areas; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.areas (
    id bigint NOT NULL,
    name character varying(50),
    sort_order bigint DEFAULT 0,
    created_at timestamp with time zone
);


ALTER TABLE public.areas OWNER TO admin;

--
-- Name: areas_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.areas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.areas_id_seq OWNER TO admin;

--
-- Name: areas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.areas_id_seq OWNED BY public.areas.id;


--
-- Name: seat_status_logs; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.seat_status_logs (
    id bigint NOT NULL,
    seat_id bigint,
    old_status character varying(20),
    new_status character varying(20),
    order_id bigint,
    changed_at timestamp with time zone
);


ALTER TABLE public.seat_status_logs OWNER TO admin;

--
-- Name: seat_status_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seat_status_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seat_status_logs_id_seq OWNER TO admin;

--
-- Name: seat_status_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.seat_status_logs_id_seq OWNED BY public.seat_status_logs.id;


--
-- Name: seats; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.seats (
    id bigint NOT NULL,
    area_id bigint,
    name character varying(50),
    type character varying(20) DEFAULT 'normal'::character varying,
    capacity bigint DEFAULT 4,
    qrcode_url character varying(500) DEFAULT ''::character varying,
    status character varying(20) DEFAULT 'available'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.seats OWNER TO admin;

--
-- Name: seats_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seats_id_seq OWNER TO admin;

--
-- Name: seats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.seats_id_seq OWNED BY public.seats.id;


--
-- Name: areas id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.areas ALTER COLUMN id SET DEFAULT nextval('public.areas_id_seq'::regclass);


--
-- Name: seat_status_logs id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.seat_status_logs ALTER COLUMN id SET DEFAULT nextval('public.seat_status_logs_id_seq'::regclass);


--
-- Name: seats id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.seats ALTER COLUMN id SET DEFAULT nextval('public.seats_id_seq'::regclass);


--
-- Data for Name: areas; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.areas (id, name, sort_order, created_at) FROM stdin;
3	大厅区	0	2026-05-07 17:13:38.19069+08
\.


--
-- Data for Name: seat_status_logs; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.seat_status_logs (id, seat_id, old_status, new_status, order_id, changed_at) FROM stdin;
\.


--
-- Data for Name: seats; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.seats (id, area_id, name, type, capacity, qrcode_url, status, created_at, updated_at) FROM stdin;
3	1	A01	standard	4	http://localhost:8080/scan?seat_id=A01	available	\N	\N
4	1	A02	standard	4	http://localhost:8080/scan?seat_id=A02	available	\N	\N
5	1	A03	standard	4	http://localhost:8080/scan?seat_id=A03	available	\N	\N
\.


--
-- Name: areas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.areas_id_seq', 3, true);


--
-- Name: seat_status_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.seat_status_logs_id_seq', 1, false);


--
-- Name: seats_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.seats_id_seq', 5, true);


--
-- Name: areas areas_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.areas
    ADD CONSTRAINT areas_pkey PRIMARY KEY (id);


--
-- Name: seat_status_logs seat_status_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.seat_status_logs
    ADD CONSTRAINT seat_status_logs_pkey PRIMARY KEY (id);


--
-- Name: seats seats_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT seats_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict tCnb2HqfJLixt2uMzemZReGSaNTev2BrnwdV6Sd25jSwuD5f3rbrsgnYUGCGjxq

