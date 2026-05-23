--
-- PostgreSQL database dump
--

\restrict E2oVFBMl9dkRbXngoEAb2Jk8KhksKhz5WV8sn62dApA4EZ9N349hHlmSCeFLzcR

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
-- Name: exchange_goods; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.exchange_goods (
    id bigint NOT NULL,
    name character varying(100),
    image character varying(500) DEFAULT ''::character varying,
    points_price bigint,
    stock bigint DEFAULT 0,
    type character varying(20) DEFAULT 'physical'::character varying,
    status smallint DEFAULT 1
);


ALTER TABLE public.exchange_goods OWNER TO admin;

--
-- Name: exchange_goods_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.exchange_goods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.exchange_goods_id_seq OWNER TO admin;

--
-- Name: exchange_goods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.exchange_goods_id_seq OWNED BY public.exchange_goods.id;


--
-- Name: exchange_orders; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.exchange_orders (
    id bigint NOT NULL,
    user_id bigint,
    goods_id bigint,
    points_cost bigint,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.exchange_orders OWNER TO admin;

--
-- Name: exchange_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.exchange_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.exchange_orders_id_seq OWNER TO admin;

--
-- Name: exchange_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.exchange_orders_id_seq OWNED BY public.exchange_orders.id;


--
-- Name: points_logs; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.points_logs (
    id bigint NOT NULL,
    user_id bigint,
    type character varying(20),
    points bigint,
    source_id character varying(64) DEFAULT ''::character varying,
    remark character varying(200) DEFAULT ''::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.points_logs OWNER TO admin;

--
-- Name: points_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.points_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.points_logs_id_seq OWNER TO admin;

--
-- Name: points_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.points_logs_id_seq OWNED BY public.points_logs.id;


--
-- Name: points_rules; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.points_rules (
    id bigint NOT NULL,
    name character varying(50),
    type character varying(30),
    config_json text,
    status smallint DEFAULT 1
);


ALTER TABLE public.points_rules OWNER TO admin;

--
-- Name: points_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.points_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.points_rules_id_seq OWNER TO admin;

--
-- Name: points_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.points_rules_id_seq OWNED BY public.points_rules.id;


--
-- Name: user_points; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.user_points (
    id bigint NOT NULL,
    user_id bigint,
    total_points bigint DEFAULT 0,
    used_points bigint DEFAULT 0,
    frozen_points bigint DEFAULT 0,
    updated_at timestamp with time zone
);


ALTER TABLE public.user_points OWNER TO admin;

--
-- Name: user_points_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.user_points_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_points_id_seq OWNER TO admin;

--
-- Name: user_points_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.user_points_id_seq OWNED BY public.user_points.id;


--
-- Name: exchange_goods id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.exchange_goods ALTER COLUMN id SET DEFAULT nextval('public.exchange_goods_id_seq'::regclass);


--
-- Name: exchange_orders id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.exchange_orders ALTER COLUMN id SET DEFAULT nextval('public.exchange_orders_id_seq'::regclass);


--
-- Name: points_logs id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.points_logs ALTER COLUMN id SET DEFAULT nextval('public.points_logs_id_seq'::regclass);


--
-- Name: points_rules id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.points_rules ALTER COLUMN id SET DEFAULT nextval('public.points_rules_id_seq'::regclass);


--
-- Name: user_points id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_points ALTER COLUMN id SET DEFAULT nextval('public.user_points_id_seq'::regclass);


--
-- Data for Name: exchange_goods; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.exchange_goods (id, name, image, points_price, stock, type, status) FROM stdin;
3	宠物零食礼包	/images/gift.png	1000	50	physical	1
4	宠物牵引绳	/images/gift.png	2000	30	physical	1
2	宠物玩具球	/images/gift.png	500	98	physical	1
5	测试商品		300	50	physical	1
6	猫抓板		800	20	physical	1
7	测试商品		100	10	physical	1
8	冒烟测试商品		888	20	physical	1
9	最终测试商品		999	9	physical	1
\.


--
-- Data for Name: exchange_orders; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.exchange_orders (id, user_id, goods_id, points_cost, status, created_at) FROM stdin;
1	7	2	500	pending	2026-05-07 17:22:40.970689+08
2	7	2	500	pending	2026-05-07 17:22:55.368306+08
\.


--
-- Data for Name: points_logs; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.points_logs (id, user_id, type, points, source_id, remark, created_at) FROM stdin;
1	7	gain	5000	dev	开发测试积分	2026-05-07 17:21:00.162906+08
2	7	exchange	-500	1	积分兑换	2026-05-07 17:22:40.971715+08
3	7	exchange	-500	2	积分兑换	2026-05-07 17:22:55.368639+08
\.


--
-- Data for Name: points_rules; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.points_rules (id, name, type, config_json, status) FROM stdin;
\.


--
-- Data for Name: user_points; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.user_points (id, user_id, total_points, used_points, frozen_points, updated_at) FROM stdin;
1	7	5000	0	0	2026-05-07 17:21:00.160583+08
\.


--
-- Name: exchange_goods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.exchange_goods_id_seq', 9, true);


--
-- Name: exchange_orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.exchange_orders_id_seq', 2, true);


--
-- Name: points_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.points_logs_id_seq', 3, true);


--
-- Name: points_rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.points_rules_id_seq', 1, false);


--
-- Name: user_points_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.user_points_id_seq', 1, true);


--
-- Name: exchange_goods exchange_goods_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.exchange_goods
    ADD CONSTRAINT exchange_goods_pkey PRIMARY KEY (id);


--
-- Name: exchange_orders exchange_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.exchange_orders
    ADD CONSTRAINT exchange_orders_pkey PRIMARY KEY (id);


--
-- Name: points_logs points_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.points_logs
    ADD CONSTRAINT points_logs_pkey PRIMARY KEY (id);


--
-- Name: points_rules points_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.points_rules
    ADD CONSTRAINT points_rules_pkey PRIMARY KEY (id);


--
-- Name: user_points user_points_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_points
    ADD CONSTRAINT user_points_pkey PRIMARY KEY (id);


--
-- Name: idx_user_points_user_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_user_points_user_id ON public.user_points USING btree (user_id);


--
-- PostgreSQL database dump complete
--

\unrestrict E2oVFBMl9dkRbXngoEAb2Jk8KhksKhz5WV8sn62dApA4EZ9N349hHlmSCeFLzcR

