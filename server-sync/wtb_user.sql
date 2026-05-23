--
-- PostgreSQL database dump
--

\restrict ycLc1Yf6tRBvb59Oe3inHjaKnkv9ulRuk6bt0NvLjdphS0mDUHtltIcfphuX72C

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
-- Name: balance_logs; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.balance_logs (
    id bigint NOT NULL,
    user_id bigint,
    type character varying(20),
    amount bigint,
    order_no character varying(64) DEFAULT ''::character varying,
    remark character varying(255) DEFAULT ''::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.balance_logs OWNER TO admin;

--
-- Name: balance_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.balance_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.balance_logs_id_seq OWNER TO admin;

--
-- Name: balance_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.balance_logs_id_seq OWNED BY public.balance_logs.id;


--
-- Name: consumption_records; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.consumption_records (
    id bigint NOT NULL,
    user_id bigint,
    order_id bigint,
    amount bigint,
    dish_count bigint DEFAULT 0,
    created_at timestamp with time zone
);


ALTER TABLE public.consumption_records OWNER TO admin;

--
-- Name: consumption_records_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.consumption_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.consumption_records_id_seq OWNER TO admin;

--
-- Name: consumption_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.consumption_records_id_seq OWNED BY public.consumption_records.id;


--
-- Name: pet_profiles; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.pet_profiles (
    id bigint NOT NULL,
    user_id bigint,
    name character varying(50),
    breed character varying(50) DEFAULT ''::character varying,
    weight numeric(5,2) DEFAULT 0,
    birthday date,
    created_at timestamp with time zone
);


ALTER TABLE public.pet_profiles OWNER TO admin;

--
-- Name: pet_profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.pet_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pet_profiles_id_seq OWNER TO admin;

--
-- Name: pet_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.pet_profiles_id_seq OWNED BY public.pet_profiles.id;


--
-- Name: recharge_records; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.recharge_records (
    id bigint NOT NULL,
    user_id bigint,
    amount bigint,
    gifted_amount bigint DEFAULT 0,
    channel character varying(20) DEFAULT 'wxpay'::character varying,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.recharge_records OWNER TO admin;

--
-- Name: recharge_records_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.recharge_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recharge_records_id_seq OWNER TO admin;

--
-- Name: recharge_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.recharge_records_id_seq OWNED BY public.recharge_records.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    openid character varying(64),
    unionid character varying(64),
    nickname character varying(100),
    avatar_url character varying(500),
    phone character varying(20),
    member_level smallint DEFAULT 0,
    balance bigint DEFAULT 0,
    total_consumption bigint DEFAULT 0,
    total_orders bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: balance_logs id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.balance_logs ALTER COLUMN id SET DEFAULT nextval('public.balance_logs_id_seq'::regclass);


--
-- Name: consumption_records id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.consumption_records ALTER COLUMN id SET DEFAULT nextval('public.consumption_records_id_seq'::regclass);


--
-- Name: pet_profiles id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.pet_profiles ALTER COLUMN id SET DEFAULT nextval('public.pet_profiles_id_seq'::regclass);


--
-- Name: recharge_records id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_records ALTER COLUMN id SET DEFAULT nextval('public.recharge_records_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: balance_logs; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.balance_logs (id, user_id, type, amount, order_no, remark, created_at) FROM stdin;
1	3	deduct	-5000	WTB20260101120000	订单扣款	2026-05-07 16:16:40.642844+08
2	4	refund	3000	WTB20260101120001	测试退款	2026-05-07 16:16:40.703578+08
\.


--
-- Data for Name: consumption_records; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.consumption_records (id, user_id, order_id, amount, dish_count, created_at) FROM stdin;
\.


--
-- Data for Name: pet_profiles; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.pet_profiles (id, user_id, name, breed, weight, birthday, created_at) FROM stdin;
1	0	旺财	金毛	28.50	2024-03-15	2026-05-07 16:16:40.825951+08
\.


--
-- Data for Name: recharge_records; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.recharge_records (id, user_id, amount, gifted_amount, channel, status, created_at) FROM stdin;
1	1	10000	2000	wxpay	pending	2026-05-07 16:16:40.580073+08
2	10	10000	2000	wxpay	pending	2026-05-11 10:07:23.237685+08
3	11	10000	2000	wxpay	pending	2026-05-11 10:42:04.491297+08
4	11	9900	1980	wxpay	pending	2026-05-11 10:42:04.515019+08
5	12	10000	2000	wxpay	pending	2026-05-11 10:44:47.320112+08
6	12	9900	1980	wxpay	pending	2026-05-11 10:44:47.3456+08
7	13	20000	4000	wxpay	success	2026-05-11 10:53:18.009442+08
8	14	20000	4000	wxpay	success	2026-05-11 10:53:52.279457+08
9	15	19900	0	wxpay	success	2026-05-11 11:12:24.836641+08
10	15	100000	20000	wxpay	success	2026-05-11 11:12:24.862099+08
11	16	19900	0	wxpay	success	2026-05-11 11:12:49.545976+08
12	17	19900	0	wxpay	success	2026-05-11 11:30:37.722532+08
13	11	19900	0	wxpay	success	2026-05-11 11:36:57.530626+08
14	11	100000	20000	wxpay	success	2026-05-11 15:23:27.270064+08
15	21	19900	0	wxpay	success	2026-05-11 19:35:06.753217+08
16	21	100000	20000	wxpay	success	2026-05-11 19:35:28.42233+08
17	22	19900	0	wxpay	success	2026-05-11 20:14:44.032339+08
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, openid, unionid, nickname, avatar_url, phone, member_level, balance, total_consumption, total_orders, created_at, updated_at) FROM stdin;
7	dev_openid_dev_test_001		开发用户			0	50000	0	0	2026-05-07 17:09:25.683535+08	2026-05-07 17:09:25.683535+08
8	dev_openid_dev_test_member		开发用户			0	0	0	0	2026-05-11 10:05:12.22507+08	2026-05-11 10:05:12.22507+08
9	dev_openid_dev_test_member2		开发用户			0	0	0	0	2026-05-11 10:05:23.005226+08	2026-05-11 10:05:23.005226+08
10	dev_openid_dev_test_v3		开发用户			2	0	0	0	2026-05-11 10:06:09.382278+08	2026-05-11 10:06:09.382278+08
12	dev_openid_dev_test_199_v2		开发用户			2	0	0	0	2026-05-11 10:44:47.282451+08	2026-05-11 10:44:47.282451+08
13	dev_openid_dev_test_full		开发用户			2	24000	0	0	2026-05-11 10:53:17.98151+08	2026-05-11 10:53:17.98151+08
14	dev_openid_dev_test_gateway		开发用户			2	24000	0	0	2026-05-11 10:53:52.242751+08	2026-05-11 10:53:52.242751+08
15	dev_openid_dev_member_test		开发用户			2	120000	0	0	2026-05-11 11:12:24.805784+08	2026-05-11 11:12:24.805784+08
16	dev_openid_dev_member_gateway		开发用户			1	0	0	0	2026-05-11 11:12:49.526619+08	2026-05-11 11:12:49.526619+08
17	dev_openid_dev_test_mock_12345		开发用户			1	0	0	0	2026-05-11 11:30:27.071632+08	2026-05-11 11:30:27.071632+08
11	dev_openid_dev_test_199		开发用户			2	120000	0	0	2026-05-11 10:42:04.462776+08	2026-05-11 10:42:04.462776+08
18	test_openid_12345		微信用户			0	0	0	0	2026-05-11 16:32:01.293378+08	2026-05-11 16:32:01.293378+08
19	mock_openid_0b3wSu1w3ban273tsR3w3eOv1y2wSu1B		微信用户			0	0	0	0	2026-05-11 16:38:48.091005+08	2026-05-11 16:38:48.091005+08
20	mock_openid_0c32Ho0w3YXg173cHf0w3Qr8jR22Ho0H		用户1403		13818571403	0	0	0	0	2026-05-11 16:38:58.206328+08	2026-05-11 16:38:58.206328+08
21	o5VbI5cPavKedSd3wydDpCRRChys		用户1403		13818571403	2	120000	0	0	2026-05-11 19:35:00.039238+08	2026-05-11 19:35:00.039238+08
22	o5VbI5eVeF9yS93PbbBtvCQbceT0		用户1403	wxfile://tmp_c00a3d27d5c30f5797aa7057afb95db9255b37ea91cc591084b856344ac1d066.jpeg	13818571403	1	0	0	0	2026-05-11 20:11:40.862547+08	2026-05-11 20:11:40.862547+08
\.


--
-- Name: balance_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.balance_logs_id_seq', 2, true);


--
-- Name: consumption_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.consumption_records_id_seq', 1, false);


--
-- Name: pet_profiles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.pet_profiles_id_seq', 1, true);


--
-- Name: recharge_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.recharge_records_id_seq', 17, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 22, true);


--
-- Name: balance_logs balance_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.balance_logs
    ADD CONSTRAINT balance_logs_pkey PRIMARY KEY (id);


--
-- Name: consumption_records consumption_records_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.consumption_records
    ADD CONSTRAINT consumption_records_pkey PRIMARY KEY (id);


--
-- Name: pet_profiles pet_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.pet_profiles
    ADD CONSTRAINT pet_profiles_pkey PRIMARY KEY (id);


--
-- Name: recharge_records recharge_records_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_records
    ADD CONSTRAINT recharge_records_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_users_open_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_users_open_id ON public.users USING btree (openid);


--
-- PostgreSQL database dump complete
--

\unrestrict ycLc1Yf6tRBvb59Oe3inHjaKnkv9ulRuk6bt0NvLjdphS0mDUHtltIcfphuX72C

