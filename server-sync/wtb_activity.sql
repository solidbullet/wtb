--
-- PostgreSQL database dump
--

\restrict CiDaaecQg2b0FsjND0SgSEAFPqpnZ27jBJW4HXZkFgn35UJctwRrn5kkMgTYFDN

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
-- Name: activities; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.activities (
    id bigint NOT NULL,
    title character varying(100),
    description text,
    image character varying(500) DEFAULT ''::character varying,
    max_participants bigint DEFAULT '-1'::integer,
    current_participants bigint DEFAULT 0,
    event_time timestamp with time zone,
    location character varying(200) DEFAULT ''::character varying,
    status character varying(20) DEFAULT 'draft'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.activities OWNER TO admin;

--
-- Name: activities_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.activities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.activities_id_seq OWNER TO admin;

--
-- Name: activities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.activities_id_seq OWNED BY public.activities.id;


--
-- Name: activity_registrations; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.activity_registrations (
    id bigint NOT NULL,
    activity_id bigint,
    user_id bigint,
    name character varying(50) DEFAULT ''::character varying,
    phone character varying(20) DEFAULT ''::character varying,
    remark character varying(200) DEFAULT ''::character varying,
    status character varying(20) DEFAULT 'registered'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.activity_registrations OWNER TO admin;

--
-- Name: activity_registrations_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.activity_registrations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.activity_registrations_id_seq OWNER TO admin;

--
-- Name: activity_registrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.activity_registrations_id_seq OWNED BY public.activity_registrations.id;


--
-- Name: announcements; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.announcements (
    id bigint NOT NULL,
    title character varying(100),
    content text,
    type character varying(20) DEFAULT 'text'::character varying,
    image character varying(500) DEFAULT ''::character varying,
    link_type character varying(20) DEFAULT ''::character varying,
    link_target character varying(500) DEFAULT ''::character varying,
    sort_order bigint DEFAULT 0,
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    status smallint DEFAULT 1
);


ALTER TABLE public.announcements OWNER TO admin;

--
-- Name: announcements_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.announcements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.announcements_id_seq OWNER TO admin;

--
-- Name: announcements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.announcements_id_seq OWNED BY public.announcements.id;


--
-- Name: activities id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.activities ALTER COLUMN id SET DEFAULT nextval('public.activities_id_seq'::regclass);


--
-- Name: activity_registrations id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.activity_registrations ALTER COLUMN id SET DEFAULT nextval('public.activity_registrations_id_seq'::regclass);


--
-- Name: announcements id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.announcements ALTER COLUMN id SET DEFAULT nextval('public.announcements_id_seq'::regclass);


--
-- Data for Name: activities; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.activities (id, title, description, image, max_participants, current_participants, event_time, location, status, created_at) FROM stdin;
8	狗狗聚会日	狗狗们的盛大草坪聚会	http://localhost:8000/dog_party.png	30	18	2026-06-15 14:00:00+08	户外草坪	published	2026-05-08 10:28:23.29176+08
3	宠物摄影日			20	6	\N	户外花园	draft	2026-05-07 17:13:38.075085+08
2	周末宠物派对			20	6	\N	店内大厅	draft	2026-05-07 17:13:38.056802+08
4	测试活动			10	0	\N	大厅	draft	2026-05-07 17:44:18.909432+08
5	修复测试活动			15	0	\N	花园	draft	2026-05-07 17:49:49.448574+08
6	冒烟测试活动			10	0	\N	测试场地	draft	2026-05-07 17:52:18.163655+08
7	最终测试活动			-1	0	\N	测试	draft	2026-05-07 18:10:46.688238+08
9	六月生日派对	给六月份出生的狗狗庆生	http://localhost:8000/birthday_party.png	20	5	2026-06-20 16:00:00+08	室内A区	published	2026-05-08 10:28:23.29176+08
\.


--
-- Data for Name: activity_registrations; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.activity_registrations (id, activity_id, user_id, name, phone, remark, status, created_at) FROM stdin;
1	3	7				registered	2026-05-07 17:18:12.20987+08
2	2	7				registered	2026-05-07 17:22:40.950971+08
\.


--
-- Data for Name: announcements; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.announcements (id, title, content, type, image, link_type, link_target, sort_order, start_time, end_time, status) FROM stdin;
1	新店开业大优惠	全场8折，欢迎光临！	text				0	2026-05-06 17:16:46.488563+08	2026-06-06 17:16:46.488563+08	1
2	新增宠物寄养服务	现在可以预约宠物寄养啦！	text				0	2026-05-06 17:16:46.488563+08	2026-06-06 17:16:46.488563+08	1
3	测试公告	测试内容	text				0	0001-01-01 08:05:43+08:05:43	0001-01-01 08:05:43+08:05:43	1
4	冒烟测试公告	这是一条测试公告	text				0	0001-01-01 08:05:43+08:05:43	0001-01-01 08:05:43+08:05:43	1
\.


--
-- Name: activities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.activities_id_seq', 9, true);


--
-- Name: activity_registrations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.activity_registrations_id_seq', 2, true);


--
-- Name: announcements_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.announcements_id_seq', 4, true);


--
-- Name: activities activities_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);


--
-- Name: activity_registrations activity_registrations_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.activity_registrations
    ADD CONSTRAINT activity_registrations_pkey PRIMARY KEY (id);


--
-- Name: announcements announcements_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.announcements
    ADD CONSTRAINT announcements_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict CiDaaecQg2b0FsjND0SgSEAFPqpnZ27jBJW4HXZkFgn35UJctwRrn5kkMgTYFDN

