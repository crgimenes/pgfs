CREATE TABLE clients (
    id serial,
    name character varying(200),
    address text
);

COPY public.clients (id, name, address) FROM stdin;
1	Umbrella Corporation	545 S Birdneck RD STE 202B Virginia Beach, VA 23451
2	OCP Omni Consumer Products	Delta City (formerly Detroit) 
3	Weyland-Yutani Corporation	Weyland-Yutani Corporation HQ, Tokyo
\.
