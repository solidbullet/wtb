--
-- PostgreSQL database dump
--

\restrict 5a0t90zNvHuBOAaOqr83TC45gWOf24SNKvyy5Aic11uzjGaBiTIqVfFYuQnDdE8

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
-- Name: order_items; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint,
    dish_id bigint,
    dish_name character varying(100),
    quantity bigint,
    unit_price bigint
);


ALTER TABLE public.order_items OWNER TO admin;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO admin;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: order_status_logs; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.order_status_logs (
    id bigint NOT NULL,
    order_id bigint,
    from_status character varying(20),
    to_status character varying(20),
    operator character varying(50) DEFAULT ''::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.order_status_logs OWNER TO admin;

--
-- Name: order_status_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.order_status_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_status_logs_id_seq OWNER TO admin;

--
-- Name: order_status_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.order_status_logs_id_seq OWNED BY public.order_status_logs.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    order_no character varying(32),
    seat_id character varying(50),
    user_id bigint,
    status character varying(20) DEFAULT 'pending'::character varying,
    total_amount bigint,
    discount_amount bigint DEFAULT 0,
    pay_amount bigint,
    remark character varying(500) DEFAULT ''::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.orders OWNER TO admin;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO admin;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: order_status_logs id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_status_logs ALTER COLUMN id SET DEFAULT nextval('public.order_status_logs_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.order_items (id, order_id, dish_id, dish_name, quantity, unit_price) FROM stdin;
1	1	1	红烧肉饭	2	3200
2	1	3	芝士蛋糕	1	2500
3	1	2	珍珠奶茶	1	1800
4	2	16	红烧肉饭	1	3200
5	3	16	红烧肉饭	1	2800
6	3	19	鸡腿饭	1	2500
7	4	16	红烧肉饭	1	2800
8	4	18	芝士蛋糕	1	2500
\.


--
-- Data for Name: order_status_logs; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.order_status_logs (id, order_id, from_status, to_status, operator, created_at) FROM stdin;
1	0				2026-05-07 17:44:18.963373+08
2	1		pending		2026-05-08 15:27:51.974461+08
3	2		pending		2026-05-11 15:23:49.177807+08
4	3		pending		2026-05-11 15:28:19.945486+08
5	4		pending		2026-05-11 15:29:24.775367+08
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.orders (id, order_no, seat_id, user_id, status, total_amount, discount_amount, pay_amount, remark, created_at, updated_at) FROM stdin;
1	WTB20260508152751	demo-seat	0	pending	10700	0	10700	测试订单	2026-05-08 15:27:51.968251+08	2026-05-08 15:27:51.968251+08
2	WTB20260511152349	demo-seat	0	pending	3200	0	3200		2026-05-11 15:23:49.169608+08	2026-05-11 15:23:49.169608+08
3	WTB20260511152819	demo-seat	11	pending	5300	0	5300	不要辣	2026-05-11 15:28:19.91047+08	2026-05-11 15:28:19.91047+08
4	WTB20260511152924	demo-seat	11	pending	5300	0	5300		2026-05-11 15:29:24.770901+08	2026-05-11 15:29:24.770901+08
\.


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.order_items_id_seq', 8, true);


--
-- Name: order_status_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.order_status_logs_id_seq', 5, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.orders_id_seq', 4, true);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: order_status_logs order_status_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_status_logs
    ADD CONSTRAINT order_status_logs_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: idx_orders_order_no; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_orders_order_no ON public.orders USING btree (order_no);


--
-- PostgreSQL database dump complete
--

\unrestrict 5a0t90zNvHuBOAaOqr83TC45gWOf24SNKvyy5Aic11uzjGaBiTIqVfFYuQnDdE8

