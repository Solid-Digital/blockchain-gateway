-- +goose Up
--
-- PostgreSQL database dump
--

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (113, '2019-04-03 00:18:08.85447', NULL, '2019-04-03 00:18:08.85447', NULL, 'Jayden Davis', '$2a$10$nlUjRT1I2o35PY4VVYAqFeqJDxAadJAXbmmS7iCEq61YDFyZZQzF.', 'avasmith108@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (116, '2019-04-03 00:18:09.186298', NULL, '2019-04-03 00:18:09.186298', NULL, 'Alexander Smith', '$2a$10$LHswOUBFegv58RWz6l8L/eybpXCRuWMH9S9QPupjK1zThLOt9Li.K', 'michaelthompson160@example.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (148, '2019-04-03 00:18:36.974303', NULL, '2019-04-03 00:18:36.974303', NULL, 'Isabella Martinez', '$2a$10$faYAhtwSuZFSFMSwt3cILua2T0DkY.TUNWCSvFeCetUmycGAKnbhq', 'ethanjohnson415@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (149, '2019-04-03 00:18:37.039613', NULL, '2019-04-03 00:18:37.039613', NULL, 'Daniel Williams', '$2a$10$oxxxdTrwLYhC54fYsRUmberLxIteG66G2HHVTsVeCtVl7fxMc2oUK', 'miasmith630@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (150, '2019-04-03 00:18:37.170876', NULL, '2019-04-03 00:18:37.170876', NULL, 'Ethan Smith', '$2a$10$qhw3lnBsu5tTV/CtTHsPAueOqLFpKaAyDzONQVLz/RvNMtbwCh2si', 'avadavis421@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (122, '2019-04-03 00:18:34.265349', NULL, '2019-04-03 00:18:34.265349', NULL, 'Emily Robinson', '$2a$10$/W8ieSdYbhfsTC7/JkGVR./2bErzVYPYRnO.yV./aHIcVWWueyTSO', 'abigailbrown552@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (123, '2019-04-03 00:18:34.704658', NULL, '2019-04-03 00:18:34.704658', NULL, 'Sophia Thomas', '$2a$10$s8Ez.be/R/S.uyzOOtygVu6V4IdpBpZixRFXEP5FP8no0.3oBmmSy', 'sophiajackson360@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (124, '2019-04-03 00:18:34.922686', NULL, '2019-04-03 00:18:34.922686', NULL, 'Anthony Wilson', '$2a$10$6UzV3qOG5cwb0GeC22j0bu2qLlbqaV4yFCA9I17KbnZsC9w2Madx.', 'charlotteanderson665@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (125, '2019-04-03 00:18:34.992467', NULL, '2019-04-03 00:18:34.992467', NULL, 'James White', '$2a$10$.MKdxKm9yz2vUyN9oG8yi.HEC7WEiXP69AKutQsH4TG1.Ef5YwvEq', 'danieljones631@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (126, '2019-04-03 00:18:35.058444', NULL, '2019-04-03 00:18:35.058444', NULL, 'James Martinez', '$2a$10$iAIwU3ff1.KCIQUxW0GSZu6zwplT0Jei7Y2qF/reapCeMuxpBfgxy', 'benjaminmartin224@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (127, '2019-04-03 00:18:35.123427', NULL, '2019-04-03 00:18:35.123427', NULL, 'Joshua Jones', '$2a$10$fGegwbeLe5UzkLWxOIESv.7Yt/RZ0tsR0YsOTUu.l5h2MeBVJ2xWu', 'isabellawhite444@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (128, '2019-04-03 00:18:35.190405', NULL, '2019-04-03 00:18:35.190405', NULL, 'Natalie Anderson', '$2a$10$Iq/pO88LjOzmy/RQHOhSY.qeIvzFsxG5Gw2m3.eU2aQPCtESetxO6', 'zoeywhite362@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (129, '2019-04-03 00:18:35.321642', NULL, '2019-04-03 00:18:35.321642', NULL, 'Benjamin Harris', '$2a$10$0YAl1Uz6qaI25e0cZaiPY.R1GR/NCLLD00R10ioINTB/F4ffml08u', 'jaydenjackson378@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (130, '2019-04-03 00:18:35.411581', NULL, '2019-04-03 00:18:35.411581', NULL, 'Andrew Robinson', '$2a$10$qU5wmbibtTyPSjb/tDntFerGUssnHLkuvHEkqhP2YLJcxyJrfiEVC', 'madisonmoore383@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (131, '2019-04-03 00:18:35.479142', NULL, '2019-04-03 00:18:35.479142', NULL, 'Zoey Wilson', '$2a$10$Mg6e1BU2wfwuyp04yy83X.9WXi/giMG1bs6rFcetgmyj3.QXM2cSK', 'andrewsmith340@example.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (132, '2019-04-03 00:18:35.610477', NULL, '2019-04-03 00:18:35.610477', NULL, 'Olivia Davis', '$2a$10$dkViBBny/Kh8pjOcgI1mu.t/diTqXIb21ykkOf4H5lu6LrCtJWyUO', 'alexanderharris225@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (133, '2019-04-03 00:18:35.677454', NULL, '2019-04-03 00:18:35.677454', NULL, 'Lily Brown', '$2a$10$C9d7uofZ4HIYY1D7rmdP6uZ84gkX0MhugoJtLi7VCVzWurFI7cNrW', 'emilysmith204@example.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (134, '2019-04-03 00:18:35.74671', NULL, '2019-04-03 00:18:35.74671', NULL, 'Aiden Smith', '$2a$10$v2LpQ.kE30eIu3qfgOv9fugc5uIUqm/6g5iG2M8EVQXGH3uV9zyQS', 'natalieharris810@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (135, '2019-04-03 00:18:35.813698', NULL, '2019-04-03 00:18:35.813698', NULL, 'Natalie Martin', '$2a$10$aDbn5N7OvuVhaEva3dfvL.BVa5f9XycYuTkIDbOUly3kZw7V.AvWC', 'andrewjackson606@example.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (136, '2019-04-03 00:18:35.882908', NULL, '2019-04-03 00:18:35.882908', NULL, 'Mason Jackson', '$2a$10$tNmM9s1NNyfshyQ.9AmS9u5BCWrN7S2vuxqq0.srFyQDsOdS5Iuo2', 'zoeysmith277@example.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (137, '2019-04-03 00:18:35.950691', NULL, '2019-04-03 00:18:35.950691', NULL, 'Joseph White', '$2a$10$Dk24Y7RHN6ZgIV4fTI2fE.LH5saypOFTvIILbKqRRlrJRaJhMO3qm', 'zoeymartinez587@example.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (138, '2019-04-03 00:18:36.018074', NULL, '2019-04-03 00:18:36.018074', NULL, 'Jayden Wilson', '$2a$10$OoxkVPwk6uO9v8YYZD3AruEDCqSn4XbRtzRBEXVuWbC.sk2xbs3oW', 'ethanwilliams877@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (139, '2019-04-03 00:18:36.084761', NULL, '2019-04-03 00:18:36.084761', NULL, 'Avery Smith', '$2a$10$bYvRMTczRm9VLxDnddIRiedwAdoqIHsQB7ECWmHG63SwR25Jmd6VW', 'davidjohnson213@example.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (140, '2019-04-03 00:18:36.293993', NULL, '2019-04-03 00:18:36.293993', NULL, 'Liam Wilson', '$2a$10$HuEVDgN1lrfZS00TBBR3HOcNzCdN6zWXOPiKifOJSGBdXfzJ9J2fa', 'aubreyrobinson383@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (141, '2019-04-03 00:18:36.362115', NULL, '2019-04-03 00:18:36.362115', NULL, 'Noah Robinson', '$2a$10$ldgTfDFi7zZxd6W3CGyLHO4d0CJuhYRzLky8vIcqofuR4azb8hV.W', 'danielmartinez204@example.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (142, '2019-04-03 00:18:36.465119', NULL, '2019-04-03 00:18:36.465119', NULL, 'Matthew Garcia', '$2a$10$FQklWaZwyGIrGvvj5oO.weFT7O6DLTtMg85Xp5mRbaxEOBa/DbBGK', 'addisonmoore324@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (143, '2019-04-03 00:18:36.53443', NULL, '2019-04-03 00:18:36.53443', NULL, 'Benjamin Johnson', '$2a$10$aF.ANYtljpKl2Sfzt7C6EuMWqEGkpCXElaqtcpKmnicEwKmEJMbFq', 'emmasmith850@test.net', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (144, '2019-04-03 00:18:36.667489', NULL, '2019-04-03 00:18:36.667489', NULL, 'Benjamin Wilson', '$2a$10$iyFTm6iW82lBDSLTEmhmNOtDZq.749T8Mu/08y9tsAjqKXcUiBhSW', 'williamwhite382@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (145, '2019-04-03 00:18:36.721254', NULL, '2019-04-03 00:18:36.721254', NULL, '', '', 'liammiller431@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (146, '2019-04-03 00:18:36.812208', NULL, '2019-04-03 00:18:36.812208', NULL, 'Olivia Robinson', '$2a$10$sOTbMrG67pvVXI2s.whs0OCX9WlLX/.PANXVVBnFUoY/KbL.rtD1y', 'danieljackson431@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (147, '2019-04-03 00:18:36.881469', NULL, '2019-04-03 00:18:36.881469', NULL, 'Ethan Miller', '$2a$10$o0MHVPUIOo8rIPY9ri.cFuOZuqDluuYLm.jIxtml9BcZq9B6vCW8.', 'nataliethompson246@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (151, '2019-04-03 00:18:37.24599', NULL, '2019-04-03 00:18:37.24599', NULL, 'Addison Jackson', '$2a$10$Aou4nptuXD0k7iSDkzAnVe/okkf0olkbG/72XC5d7j.6yQIhGU5.u', 'elijahbrown827@example.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (152, '2019-04-03 00:18:37.379434', NULL, '2019-04-03 00:18:37.379434', NULL, 'David Martin', '$2a$10$CXpXNKAQSWL1NMZCyF/02uok076nf6uv4O8YJzM6AMwfe8y8EnDhm', 'isabellasmith732@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (153, '2019-04-03 00:18:37.469879', NULL, '2019-04-03 00:18:37.469879', NULL, 'Madison Taylor', '$2a$10$IBHkylHrMjwuLn4Ig3IA..YoiSbxcTGf.VwniHAWPYYihU5OIse.S', 'matthewanderson203@test.com', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (154, '2019-04-03 00:18:37.619534', NULL, '2019-04-03 00:18:37.619534', NULL, 'Aiden Jones', '$2a$10$ah5oIyLBm50oXzdS7rvX0.YZF2Q.V1B6.wUJRRDrTMZnnOf87z7ri', 'charlottemiller282@test.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (155, '2019-04-03 00:18:37.689209', NULL, '2019-04-03 00:18:37.689209', NULL, 'Olivia Jackson', '$2a$10$6nkgp/nIixjf8Y09pUdhheouQuMBpO0OufXPeUvQQrLzfBlXFCspe', 'lilyharris654@example.org', NULL);
INSERT INTO public.users (id, created_at, created_by_id, updated_at, updated_by_id, full_name, password_hash, email, default_organization_id) VALUES (156, '2019-04-03 00:18:37.762554', NULL, '2019-04-03 00:18:37.762554', NULL, 'Elizabeth Brown', '$2a$10$a4af/DUteFwT1C0nQ1F/m.TPAl5JaxNnsmLKa.0fbW4UwQhPGfQ76', 'addisonharris443@example.com', NULL);


--
-- Data for Name: organization; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (83, '2019-04-03 00:18:09.024691', 113, '2019-04-03 00:18:09.024691', 113, 'Pegasushickory', 'pegasushickory');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (84, '2019-04-03 00:18:09.188006', 116, '2019-04-03 00:18:09.189749', 116, 'some other name', 'eyehickory');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (90, '2019-04-03 00:18:34.332274', 122, '2019-04-03 00:18:34.332274', 122, 'Hornetluminous', 'hornetluminous');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (91, '2019-04-03 00:18:34.709087', 123, '2019-04-03 00:18:34.709087', 123, 'Gemtree', 'gemtree');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (92, '2019-04-03 00:18:34.926446', 124, '2019-04-03 00:18:34.926446', 124, 'Curtainrain', 'curtainrain');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (93, '2019-04-03 00:18:35.325441', 129, '2019-04-03 00:18:35.325441', 129, 'Samuraimulberry', 'samuraimulberry');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (37, '2019-04-03 00:20:14.572322', 1, '2019-04-03 00:20:14.572322', 1, 'demo3', 'demo3');
INSERT INTO public.organization (id, created_at, created_by_id, updated_at, updated_by_id, display_name, name) VALUES (144, '2019-04-03 18:26:53.60113', 1, '2019-04-03 18:26:53.60113', 1, 'another org', 'ares-v2-demo2');

--
-- Data for Name: action; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.action (id, developer_id, created_at, created_by_id, updated_at, updated_by_id, name, display_name, description, public) VALUES (1, 1, '2019-04-02 19:42:00.704598', 1, '2019-04-02 19:42:00.704598', 1, 'action-3', 'action-3', '', false);


--
-- Data for Name: action_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.action_version (id, action_id, created_at, created_by_id, updated_at, updated_by_id, version, public, example_config, input_schema, output_schema, description, readme, file_name, file_id) VALUES (1, 1, '2019-04-02 19:41:54.293767', 1, '2019-04-02 19:41:54.293767', 1, '0.0.1', false, '', '[]', '[]','', '', 'action-3.action.0.0.1.ares-v2-demo.so', 'ares-v2-demo/action/action-3/0.0.1/action-3.action.0.0.1.ares-v2-demo.so.tar.gz');


--
-- Data for Name: draft_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.draft_configuration (id, created_at, created_by_id, updated_at, updated_by_id, organization_id, revision) VALUES (1, '2019-04-02 19:42:17.469259', 1, '2019-04-03 15:42:00.333126', 1, 1, 13);
INSERT INTO public.draft_configuration (id, created_at, created_by_id, updated_at, updated_by_id, organization_id, revision) VALUES (2, '2019-04-03 18:28:00.715952', 1, '2019-04-03 18:28:00.715952', 1, 144, 0);


--
-- Data for Name: pipeline; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.pipeline (id, organization_id, created_at, created_by_id, updated_at, updated_by_id, display_name, name, status, description, draft_configuration_id) VALUES (1, 1, '2019-04-02 19:42:17.468081', 1, '2019-04-02 19:42:17.468081', 1, 'My pipeline also', 'my-pipeline', '', 'updated description', 1);
INSERT INTO public.pipeline (id, organization_id, created_at, created_by_id, updated_at, updated_by_id, display_name, name, status, description, draft_configuration_id) VALUES (3, 144, '2019-04-03 18:28:00.680697', 1, '2019-04-03 18:28:00.680697', 1, 'My pipeline', 'my-pipeline', '', 'my pipeline description', 2);


--
-- Data for Name: configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (2, '2019-04-03 00:28:13.991556', 1, '2019-04-03 00:28:13.991556', 1, 1, 1, 1, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (3, '2019-04-03 00:28:27.205826', 1, '2019-04-03 00:28:27.205826', 1, 1, 1, 2, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (4, '2019-04-03 00:28:28.474365', 1, '2019-04-03 00:28:28.474365', 1, 1, 1, 3, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (5, '2019-04-03 00:28:29.586466', 1, '2019-04-03 00:28:29.586466', 1, 1, 1, 4, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (6, '2019-04-03 00:31:10.584291', 1, '2019-04-03 00:31:10.584291', 1, 1, 1, 5, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (7, '2019-04-03 00:37:36.726324', 1, '2019-04-03 00:37:36.726324', 1, 1, 1, 6, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (8, '2019-04-03 00:37:38.064096', 1, '2019-04-03 00:37:38.064096', 1, 1, 1, 7, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (9, '2019-04-03 00:37:39.978464', 1, '2019-04-03 00:37:39.978464', 1, 1, 1, 8, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (1, '2019-04-03 14:30:22.814143', 1, '2019-04-03 14:30:22.814143', 1, 1, 1, 9, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (10, '2019-04-03 15:28:28.499409', 1, '2019-04-03 15:28:28.499409', 1, 1, 1, 10, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (11, '2019-04-03 15:41:43.293165', 1, '2019-04-03 15:41:43.293165', 1, 1, 1, 11, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (12, '2019-04-03 15:41:47.53425', 1, '2019-04-03 15:41:47.53425', 1, 1, 1, 12, '');
INSERT INTO public.configuration (id, created_at, created_by_id, updated_at, updated_by_id, pipeline_id, organization_id, revision, commit_message) VALUES (13, '2019-04-03 15:41:51.80818', 1, '2019-04-03 15:41:51.80818', 1, 1, 1, 13, '');


--
-- Data for Name: action_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (2, 2, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (3, 3, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (4, 4, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (5, 5, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (6, 6, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (7, 7, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (8, 8, 1, 'zzz', 0, 'zzz', '{}');
INSERT INTO public.action_configuration (id, configuration_id, version_id, name, index, config, message_config) VALUES (9, 9, 1, 'zzz', 0, 'zzz', '{}');


--
-- Data for Name: base_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (2, 2, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (3, 3, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (4, 4, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (5, 5, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (6, 6, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (7, 7, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (8, 8, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (9, 9, 1, 'asdf');
INSERT INTO public.base_configuration (id, configuration_id, version_id, config) VALUES (1, 1, 1, 'asdf');


--
-- Data for Name: base_draft_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.base_draft_configuration (id, draft_configuration_id, version_id, config) VALUES (14, 1, 1, 'asdf');
INSERT INTO public.base_draft_configuration (id, draft_configuration_id, version_id, config) VALUES (15, 2, 1, '');


--
-- Data for Name: trigger; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.trigger (id, developer_id, created_at, created_by_id, updated_at, updated_by_id, name, display_name, description, public) VALUES (1, 1, '2019-04-02 19:41:54.293767', 1, '2019-04-02 19:41:54.293767', 1, 'trigger-1', 'trigger-1', '', false);


--
-- Data for Name: trigger_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.trigger_version (id, trigger_id, created_at, created_by_id, updated_at, updated_by_id, version, description, public, example_config, input_schema, output_schema, readme, file_name, file_id) VALUES (1, 1, '2019-04-02 19:41:54.293767', 1, '2019-04-02 19:41:54.293767', 1, '0.0.10', '', false, '', '[]', '[]', '', 'trigger-1.trigger.0.0.1.ares-v2-demo.so', 'ares-v2-demo/trigger/trigger-1/0.0.1/trigger-1.trigger.0.0.1.ares-v2-demo.so.tar.gz');
INSERT INTO public.trigger_version (id, trigger_id, created_at, created_by_id, updated_at, updated_by_id, version, description, public, example_config, input_schema, output_schema, readme, file_name, file_id) VALUES (3, 1, '2019-04-02 19:41:54.293767', 1, '2019-04-02 19:41:54.293767', 1, '0.0.1', 'asdasd', false, '??', '[]', '[]', '??', 'trigger-1.trigger.0.0.1.ares-v2-demo.so', 'ares-v2-demo/trigger/trigger-1/0.0.1/trigger-1.trigger.0.0.1.ares-v2-demo.so.tar.gz');
INSERT INTO public.trigger_version (id, trigger_id, created_at, created_by_id, updated_at, updated_by_id, version, description, public, example_config, input_schema, output_schema, readme, file_name, file_id) VALUES (12, 1, '2019-04-02 19:41:54.293767', 1, '2019-04-02 19:41:54.293767', 1, '0.0.15', 'asdasd', false, '??', '[]', '[]', '??', 'trigger-1.trigger.0.0.15.ares-v2-demo.so', 'ares-v2-demo/trigger/trigger-1/0.0.15/trigger-1.trigger.0.0.15.ares-v2-demo.so.tar.gz');


--
-- Data for Name: trigger_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (2, 2, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (3, 3, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (4, 4, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (5, 5, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (6, 6, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (7, 7, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (8, 8, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (9, 9, 1, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (1, 1, 3, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (10, 10, 3, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (11, 11, 12, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (12, 12, 12, 'trigger-1', '', '{}');
INSERT INTO public.trigger_configuration (id, configuration_id, version_id, name, config, message_config) VALUES (13, 13, 12, 'trigger-1', '', '{}');


--
-- Data for Name: trigger_draft_configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.trigger_draft_configuration (id, draft_configuration_id, version_id, name, config, message_config) VALUES (12, 1, 12, 'trigger-1', '', '{}');
INSERT INTO public.trigger_draft_configuration (id, draft_configuration_id, version_id, name, config, message_config) VALUES (13, 2, NULL, 'trigger-1', '', '{}');


--
-- Data for Name: trigger_version_supported_bases; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user_organization; Type: TABLE DATA; Schema: public; Owner: postgres
--

-- INSERT INTO public.user_organization (user_id, organization_id) VALUES (1, 1);
INSERT INTO public.user_organization (user_id, organization_id) VALUES (1, 37);
INSERT INTO public.user_organization (user_id, organization_id) VALUES (1, 144);


--
-- Name: action_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.action_configuration_id_seq', 9, true);


--
-- Name: action_draft_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.action_draft_configuration_id_seq', 13, true);


--
-- Name: action_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.action_id_seq', 100, false);


--
-- Name: action_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.action_version_id_seq', 100, false);


--
-- Name: base_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.base_configuration_id_seq', 13, true);


--
-- Name: base_draft_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.base_draft_configuration_id_seq', 15, true);


--
-- Name: base_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.base_id_seq', 100, false);


--
-- Name: base_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.base_version_id_seq', 100, false);


--
-- Name: configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.configuration_id_seq', 13, true);


--
-- Name: deployment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.deployment_id_seq', 7, true);


--
-- Name: draft_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.draft_configuration_id_seq', 2, true);


--
-- Name: environment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.environment_id_seq', 4, true);


--
-- Name: organization_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.organization_id_seq', 155, true);


--
-- Name: pipeline_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pipeline_id_seq', 3, true);


--
-- Name: trigger_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.trigger_configuration_id_seq', 13, true);


--
-- Name: trigger_draft_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.trigger_draft_configuration_id_seq', 13, true);


--
-- Name: trigger_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.trigger_id_seq', 100, false);


--
-- Name: trigger_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.trigger_version_id_seq', 13, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 202, true);


--
-- PostgreSQL database dump complete
--



;