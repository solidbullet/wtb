--
-- PostgreSQL database dump
--

\restrict fGBhlwE0r3bGsrrNXqVtFlh2BjWFT6qctzrGJbdjc0WUXWf4Y8cl4ml04DWejIj

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
-- Name: payment_orders; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.payment_orders (
    id bigint NOT NULL,
    order_no character varying(32),
    out_trade_no character varying(32),
    user_id bigint,
    amount bigint,
    channel character varying(20),
    status character varying(20) DEFAULT 'pending'::character varying,
    wx_prepay_id character varying(64) DEFAULT ''::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.payment_orders OWNER TO admin;

--
-- Name: payment_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.payment_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_orders_id_seq OWNER TO admin;

--
-- Name: payment_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.payment_orders_id_seq OWNED BY public.payment_orders.id;


--
-- Name: payment_records; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.payment_records (
    id bigint NOT NULL,
    payment_order_id bigint,
    channel character varying(20),
    amount bigint,
    transaction_id character varying(64) DEFAULT ''::character varying,
    paid_at timestamp with time zone
);


ALTER TABLE public.payment_records OWNER TO admin;

--
-- Name: payment_records_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.payment_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_records_id_seq OWNER TO admin;

--
-- Name: payment_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.payment_records_id_seq OWNED BY public.payment_records.id;


--
-- Name: recharge_orders; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.recharge_orders (
    id bigint NOT NULL,
    user_id bigint,
    amount bigint,
    gifted_amount bigint DEFAULT 0,
    discount_rate numeric(3,2) DEFAULT 1,
    final_amount bigint,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.recharge_orders OWNER TO admin;

--
-- Name: recharge_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.recharge_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recharge_orders_id_seq OWNER TO admin;

--
-- Name: recharge_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.recharge_orders_id_seq OWNED BY public.recharge_orders.id;


--
-- Name: refund_records; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.refund_records (
    id bigint NOT NULL,
    payment_order_id bigint,
    refund_no character varying(32),
    amount bigint,
    reason character varying(200) DEFAULT ''::character varying,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone
);


ALTER TABLE public.refund_records OWNER TO admin;

--
-- Name: refund_records_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.refund_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refund_records_id_seq OWNER TO admin;

--
-- Name: refund_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.refund_records_id_seq OWNED BY public.refund_records.id;


--
-- Name: payment_orders id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.payment_orders ALTER COLUMN id SET DEFAULT nextval('public.payment_orders_id_seq'::regclass);


--
-- Name: payment_records id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.payment_records ALTER COLUMN id SET DEFAULT nextval('public.payment_records_id_seq'::regclass);


--
-- Name: recharge_orders id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_orders ALTER COLUMN id SET DEFAULT nextval('public.recharge_orders_id_seq'::regclass);


--
-- Name: refund_records id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refund_records ALTER COLUMN id SET DEFAULT nextval('public.refund_records_id_seq'::regclass);


--
-- Data for Name: payment_orders; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.payment_orders (id, order_no, out_trade_no, user_id, amount, channel, status, wx_prepay_id, created_at) FROM stdin;
\.


--
-- Data for Name: payment_records; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.payment_records (id, payment_order_id, channel, amount, transaction_id, paid_at) FROM stdin;
\.


--
-- Data for Name: recharge_orders; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.recharge_orders (id, user_id, amount, gifted_amount, discount_rate, final_amount, status, created_at) FROM stdin;
\.


--
-- Data for Name: refund_records; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.refund_records (id, payment_order_id, refund_no, amount, reason, status, created_at) FROM stdin;
\.


--
-- Name: payment_orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.payment_orders_id_seq', 1, true);


--
-- Name: payment_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.payment_records_id_seq', 1, false);


--
-- Name: recharge_orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.recharge_orders_id_seq', 1, false);


--
-- Name: refund_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.refund_records_id_seq', 1, false);


--
-- Name: payment_orders payment_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.payment_orders
    ADD CONSTRAINT payment_orders_pkey PRIMARY KEY (id);


--
-- Name: payment_records payment_records_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.payment_records
    ADD CONSTRAINT payment_records_pkey PRIMARY KEY (id);


--
-- Name: recharge_orders recharge_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.recharge_orders
    ADD CONSTRAINT recharge_orders_pkey PRIMARY KEY (id);


--
-- Name: refund_records refund_records_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.refund_records
    ADD CONSTRAINT refund_records_pkey PRIMARY KEY (id);


--
-- Name: idx_payment_orders_out_trade_no; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_payment_orders_out_trade_no ON public.payment_orders USING btree (out_trade_no);


--
-- Name: idx_refund_records_refund_no; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_refund_records_refund_no ON public.refund_records USING btree (refund_no);


--
-- PostgreSQL database dump complete
--

\unrestrict fGBhlwE0r3bGsrrNXqVtFlh2BjWFT6qctzrGJbdjc0WUXWf4Y8cl4ml04DWejIj

