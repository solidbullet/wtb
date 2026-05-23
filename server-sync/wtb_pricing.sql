--
-- PostgreSQL database dump
--

\restrict oFt6kMf0UAjctjo3SZ443UWGMSBID7EAwZhPEyVbddM5vQyNM9KDMY1sq24epe4

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
-- Name: combos; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.combos (
    id bigint NOT NULL,
    name character varying(100),
    price bigint,
    dish_list text,
    status smallint DEFAULT 1
);


ALTER TABLE public.combos OWNER TO admin;

--
-- Name: combos_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.combos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.combos_id_seq OWNER TO admin;

--
-- Name: combos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.combos_id_seq OWNED BY public.combos.id;


--
-- Name: price_rules; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.price_rules (
    id bigint NOT NULL,
    dish_id bigint,
    rule_type character varying(20),
    price bigint,
    start_time text,
    end_time text,
    status smallint DEFAULT 1
);


ALTER TABLE public.price_rules OWNER TO admin;

--
-- Name: price_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.price_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.price_rules_id_seq OWNER TO admin;

--
-- Name: price_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.price_rules_id_seq OWNED BY public.price_rules.id;


--
-- Name: promotions; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.promotions (
    id bigint NOT NULL,
    name character varying(100),
    type character varying(30),
    config_json text,
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    status smallint DEFAULT 1
);


ALTER TABLE public.promotions OWNER TO admin;

--
-- Name: promotions_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.promotions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.promotions_id_seq OWNER TO admin;

--
-- Name: promotions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.promotions_id_seq OWNED BY public.promotions.id;


--
-- Name: recharge_plans; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.recharge_plans (
    id bigint NOT NULL,
    name character varying(100),
    amount bigint,
    final_amount bigint,
    gift_amount bigint,
    sort_order bigint DEFAULT 0,
    status smallint DEFAULT 1,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.recharge_plans OWNER TO admin;

--
-- Name: recharge_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.recharge_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recharge_plans_id_seq OWNER TO admin;

--
-- Name: recharge_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.recharge_plans_id_seq OWNED BY public.recharge_plans.id;


--
-- Name: combos id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.combos ALTER COLUMN id SET DEFAULT nextval('public.combos_id_seq'::regclass);


--
-- Name: price_rules id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.price_rules ALTER COLUMN id SET DEFAULT nextval('public.price_rules_id_seq'::regclass);


--
-- Name: promotions id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.promotions ALTER COLUMN id SET DEFAULT nextval('public.promotions_id_seq'::regclass);


--
-- Name: recharge_plans id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_plans ALTER COLUMN id SET DEFAULT nextval('public.recharge_plans_id_seq'::regclass);


--
-- Data for Name: combos; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.combos (id, name, price, dish_list, status) FROM stdin;
\.


--
-- Data for Name: price_rules; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.price_rules (id, dish_id, rule_type, price, start_time, end_time, status) FROM stdin;
\.


--
-- Data for Name: promotions; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.promotions (id, name, type, config_json, start_time, end_time, status) FROM stdin;
\.


--
-- Data for Name: recharge_plans; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.recharge_plans (id, name, amount, final_amount, gift_amount, sort_order, status, created_at, updated_at) FROM stdin;
1	充100送10	10000	11000	1000	1	1	\N	\N
2	充200送40	20000	24000	4000	2	1	\N	\N
3	充500送150	50000	65000	15000	3	1	\N	\N
4	充1000送400	100000	140000	40000	4	1	\N	\N
\.


--
-- Name: combos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.combos_id_seq', 1, false);


--
-- Name: price_rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.price_rules_id_seq', 4, true);


--
-- Name: promotions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.promotions_id_seq', 1, false);


--
-- Name: recharge_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.recharge_plans_id_seq', 4, true);


--
-- Name: combos combos_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.combos
    ADD CONSTRAINT combos_pkey PRIMARY KEY (id);


--
-- Name: price_rules price_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.price_rules
    ADD CONSTRAINT price_rules_pkey PRIMARY KEY (id);


--
-- Name: promotions promotions_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.promotions
    ADD CONSTRAINT promotions_pkey PRIMARY KEY (id);


--
-- Name: recharge_plans recharge_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_plans
    ADD CONSTRAINT recharge_plans_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict oFt6kMf0UAjctjo3SZ443UWGMSBID7EAwZhPEyVbddM5vQyNM9KDMY1sq24epe4

