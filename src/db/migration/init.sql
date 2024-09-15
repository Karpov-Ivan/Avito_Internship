-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

-- Вставка данных в таблицу employee
INSERT INTO employee (id, username, first_name, last_name, created_at, updated_at) VALUES
    ('36f4d9d3-4c12-4b5f-bacc-ec4d54b66e82', 'jsmith', 'John', 'Smith', '2024-09-15 10:00:00', '2024-09-15 10:00:00'),
    ('4a2fdd28-9a30-4b51-bd4d-4a31414f2e7f', 'abrown', 'Alice', 'Brown', '2024-09-14 09:30:00', '2024-09-14 09:30:00'),
    ('58cf8676-f9ec-438c-bb09-fefb2e0e5c5e', 'mjohnson', 'Michael', 'Johnson', '2024-09-13 15:45:00', '2024-09-13 15:45:00'),
    ('c66e51fa-94a7-49f7-b31d-20e6e73352f6', 'klang', 'Karen', 'Lang', '2024-09-12 12:00:00', '2024-09-12 12:00:00'),
    ('75a2a0eb-1a8b-4b1c-bb68-ccab04547c4e', 'dwhite', 'David', 'White', '2024-09-11 14:20:00', '2024-09-11 14:20:00');

-- Вставка данных в таблицу organization
INSERT INTO organization (id, name, description, type, created_at, updated_at) VALUES
    ('e807f7a6-cc5f-4eeb-97de-1e48271bb028', 'Tech Innovations LLC', 'A leading company in tech solutions', 'LLC', '2024-09-15 10:00:00', '2024-09-15 10:00:00'),
    ('2b78a65e-4bb0-4916-b847-5fd1c4f299b8', 'Green Energy JSC', 'Focused on renewable energy sources', 'JSC', '2024-09-14 09:30:00', '2024-09-14 09:30:00'),
    ('6d93c0b9-1b3e-4620-b47d-cf2e7f62edc1', 'Health Solutions IE', 'Offers health-related services', 'IE', '2024-09-13 15:45:00', '2024-09-13 15:45:00'),
    ('2e0c066f-eede-4f3e-81c4-98a7bc4c5b56', 'Global Logistics LLC', 'Specializes in shipping and logistics', 'LLC', '2024-09-12 12:00:00', '2024-09-12 12:00:00'),
    ('b6b22d5f-18b3-4ee2-ae88-3db31e98bb7d', 'Innovative Education JSC', 'Provides educational resources', 'JSC', '2024-09-11 14:20:00', '2024-09-11 14:20:00');

-- Вставка данных в таблицу organization_responsible
INSERT INTO organization_responsible (id, organization_id, user_id) VALUES
    ('19f93e93-1cc6-4b46-a9d7-f9bfb2040c5d', 'e807f7a6-cc5f-4eeb-97de-1e48271bb028', '36f4d9d3-4c12-4b5f-bacc-ec4d54b66e82'),
    ('235b2458-4556-42cf-81d4-e8c91a84b4b7', '2b78a65e-4bb0-4916-b847-5fd1c4f299b8', '4a2fdd28-9a30-4b51-bd4d-4a31414f2e7f'),
    ('0beecfe4-cb35-41d7-9be1-4c53deaf3572', '6d93c0b9-1b3e-4620-b47d-cf2e7f62edc1', '58cf8676-f9ec-438c-bb09-fefb2e0e5c5e'),
    ('4c258682-3b2d-4eb3-9728-b3f3e35b0c1a', '2e0c066f-eede-4f3e-81c4-98a7bc4c5b56', '75a2a0eb-1a8b-4b1c-bb68-ccab04547c4e'),
    ('13ff3a93-6eae-41cb-9f1e-947182ab3e1f', 'b6b22d5f-18b3-4ee2-ae88-3db31e98bb7d', '75a2a0eb-1a8b-4b1c-bb68-ccab04547c4e');

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    service_type VARCHAR(50) NOT NULL,
    creator_username VARCHAR(50) NOT NULL

);

CREATE TABLE proposal (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);