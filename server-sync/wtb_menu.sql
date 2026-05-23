--
-- PostgreSQL database dump
--

\restrict kpzSGh0ilmLrxKobdsUEOuhbVoTqg1xVIC5thxpWD1AEljQpZ38GoAcW0BEC0df

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
-- Name: categories; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    name character varying(50),
    parent_id bigint DEFAULT 0,
    sort_order bigint DEFAULT 0,
    status smallint DEFAULT 1,
    created_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO admin;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO admin;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: dish_prices; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.dish_prices (
    id bigint NOT NULL,
    dish_id bigint,
    price_type character varying(20),
    price bigint,
    start_time timestamp with time zone,
    end_time timestamp with time zone
);


ALTER TABLE public.dish_prices OWNER TO admin;

--
-- Name: dish_prices_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.dish_prices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dish_prices_id_seq OWNER TO admin;

--
-- Name: dish_prices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.dish_prices_id_seq OWNED BY public.dish_prices.id;


--
-- Name: dish_stocks; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.dish_stocks (
    id bigint NOT NULL,
    dish_id bigint,
    daily_limit bigint DEFAULT '-1'::integer,
    sold_count bigint DEFAULT 0,
    date date
);


ALTER TABLE public.dish_stocks OWNER TO admin;

--
-- Name: dish_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.dish_stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dish_stocks_id_seq OWNER TO admin;

--
-- Name: dish_stocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.dish_stocks_id_seq OWNED BY public.dish_stocks.id;


--
-- Name: dishes; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.dishes (
    id bigint NOT NULL,
    category_id bigint,
    name character varying(100),
    subtitle character varying(200) DEFAULT ''::character varying,
    description text,
    images text,
    tags character varying(200) DEFAULT ''::character varying,
    status smallint DEFAULT 1,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.dishes OWNER TO admin;

--
-- Name: dishes_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.dishes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dishes_id_seq OWNER TO admin;

--
-- Name: dishes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.dishes_id_seq OWNED BY public.dishes.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: dish_prices id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dish_prices ALTER COLUMN id SET DEFAULT nextval('public.dish_prices_id_seq'::regclass);


--
-- Name: dish_stocks id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dish_stocks ALTER COLUMN id SET DEFAULT nextval('public.dish_stocks_id_seq'::regclass);


--
-- Name: dishes id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dishes ALTER COLUMN id SET DEFAULT nextval('public.dishes_id_seq'::regclass);


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.categories (id, name, parent_id, sort_order, status, created_at) FROM stdin;
24	主食	0	1	1	2026-05-08 11:49:05.539434+08
25	饮品	0	2	1	2026-05-08 11:49:05.539434+08
26	甜点	0	3	1	2026-05-08 11:49:05.539434+08
27	宠物专属	0	4	1	2026-05-08 11:49:05.539434+08
\.


--
-- Data for Name: dish_prices; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.dish_prices (id, dish_id, price_type, price, start_time, end_time) FROM stdin;
1	7	standard	4500	\N	\N
2	8	standard	5500	\N	\N
3	9	standard	9999	\N	\N
4	10	standard	100	\N	\N
5	11	standard	500	\N	\N
6	12	standard	1000	\N	\N
7	13	standard	19	\N	\N
8	14	standard	1800	\N	\N
9	15	standard	3200	\N	\N
20	16	standard	3200	\N	\N
21	17	standard	1800	\N	\N
22	18	standard	2800	\N	\N
23	19	standard	2900	\N	\N
24	20	standard	1500	\N	\N
25	21	standard	4500	\N	\N
26	22	standard	1200	\N	\N
27	16	member	2800	\N	\N
28	17	member	1500	\N	\N
29	18	member	2500	\N	\N
30	19	member	2500	\N	\N
31	20	member	1200	\N	\N
32	21	member	3800	\N	\N
33	22	member	1000	\N	\N
\.


--
-- Data for Name: dish_stocks; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.dish_stocks (id, dish_id, daily_limit, sold_count, date) FROM stdin;
1	7	80	0	2026-05-07
2	8	66	0	2026-05-07
3	9	99	0	2026-05-07
4	10	10	0	2026-05-07
5	11	10	0	2026-05-07
6	12	50	0	2026-05-07
\.


--
-- Data for Name: dishes; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.dishes (id, category_id, name, subtitle, description, images, tags, status, created_at, updated_at) FROM stdin;
16	24	红烧肉饭	招牌慢炖红烧肉	精选五花肉，慢火炖制2小时，入口即化	/images/hongshao.png	推荐,热门	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
19	24	鸡腿饭	脆皮大鸡腿	金黄酥脆大鸡腿，搭配时令蔬菜	/images/chicken_rice.png	推荐	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
17	25	珍珠奶茶	Q弹黑糖珍珠	现煮黑糖珍珠搭配香浓奶茶	/images/boba.png	推荐	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
20	25	鲜榨橙汁	100%鲜榨	现切鲜橙，不加水不加糖	/images/orange_juice.png	推荐	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
18	26	芝士蛋糕	每日现烤	新西兰进口芝士，绵密口感	/images/cheesecake.png	推荐,新品	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
21	27	宠物鲜粮套餐	狗狗专属	营养师配比，鲜肉蔬菜均衡搭配		推荐	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
22	27	狗狗饼干	磨牙小零食	天然食材烘焙，无添加剂		推荐	1	2026-05-08 11:49:28.863288+08	2026-05-08 11:49:28.863288+08
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.categories_id_seq', 27, true);


--
-- Name: dish_prices_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.dish_prices_id_seq', 33, true);


--
-- Name: dish_stocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.dish_stocks_id_seq', 6, true);


--
-- Name: dishes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.dishes_id_seq', 22, true);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: dish_prices dish_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dish_prices
    ADD CONSTRAINT dish_prices_pkey PRIMARY KEY (id);


--
-- Name: dish_stocks dish_stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dish_stocks
    ADD CONSTRAINT dish_stocks_pkey PRIMARY KEY (id);


--
-- Name: dishes dishes_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.dishes
    ADD CONSTRAINT dishes_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict kpzSGh0ilmLrxKobdsUEOuhbVoTqg1xVIC5thxpWD1AEljQpZ38GoAcW0BEC0df

