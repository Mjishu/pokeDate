CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum for notification statuses
DO $$ BEGIN
    CREATE TYPE notification_status AS ENUM ('accepted', 'denied', 'unseen');
    CREATE TYPE location_type AS ENUM ( 'city', 'state', 'country');
    CREATE TYPE sex_enum AS ENUM ( 'male', 'female', 'undefined');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Locations table
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100), 
    location_type location_type NOT NULL, 
    parent_id INT REFERENCES locations(id) ON DELETE SET NULL
);

-- Animals table
CREATE TABLE IF NOT EXISTS animals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL,
    species VARCHAR(100) NOT NULL, 
    date_of_birth DATE NOT NULL,
    sex sex_enum NOT NULL,
    price FLOAT DEFAULT 0.00,
    available BOOLEAN NOT NULL,
    breed VARCHAR(50) NOT NULL
);

-- Shots table
CREATE TABLE IF NOT EXISTS shots (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT 
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(40) NOT NULL,
    email VARCHAR(100),
    password TEXT NOT NULL,
    date_of_birth DATE,
    country_id INT REFERENCES locations(id) ON DELETE SET NULL,
    state_id INT REFERENCES locations(id) ON DELETE SET NULL,
    city_id INT REFERENCES locations(id) ON DELETE SET NULL,
    profile_picture_src TEXT,
    is_organization BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- Animal Images table
CREATE TABLE IF NOT EXISTS animal_images (
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
    image_src TEXT NOT NULL,
    priority INT NOT NULL, 
    image_public_id VARCHAR(120) NOT NULL,
    CONSTRAINT unique_priority_per_animal UNIQUE (animal_id, priority)
);

-- Animal Shots table
CREATE TABLE IF NOT EXISTS animal_shots (
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
    shots_id BIGINT REFERENCES shots(id) ON DELETE CASCADE NOT NULL,
    date_given DATE,
    next_due DATE,
    PRIMARY KEY (animal_id, shots_id)
);

-- Organization Animals table
CREATE TABLE IF NOT EXISTS organization_animals (
    organization_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (organization_id, animal_id)
);

-- User Animals table
CREATE TABLE IF NOT EXISTS user_animals (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL,
    liked BOOLEAN NOT NULL,
    PRIMARY KEY (user_id, animal_id)
);

-- Refresh Tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ DEFAULT NULL
);

-- Conversation table
CREATE TABLE IF NOT EXISTS conversation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_name VARCHAR(90),
    date_created TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Conversation Member table
CREATE TABLE IF NOT EXISTS conversation_member (
    member_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL ,
    conversation_id UUID REFERENCES conversation(id)ON DELETE CASCADE  NOT NULL ,
    joined_datetime TIMESTAMPTZ DEFAULT now() NOT NULL,
    left_datetime TIMESTAMPTZ 
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL ,
    message_text TEXT NOT NULL,
    sent_datetime TIMESTAMPTZ DEFAULT now() NOT NULL,
    conversation_id UUID REFERENCES conversation(id) ON DELETE CASCADE
);

-- Notification Types table
CREATE TABLE IF NOT EXISTS notification_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    notifier UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    entity_text TEXT,
    entity_type INT REFERENCES notification_types(id) ON DELETE SET NULL, 
    status notification_status,
    date_created TIMESTAMPTZ DEFAULT now() NOT NULL,
    date_seen TIMESTAMPTZ
);

-- Animal Groups table
CREATE TABLE IF NOT EXISTS animal_groups (
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE,
    notification_id UUID REFERENCES notifications(id) ON DELETE CASCADE,
    conversation_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    date_created TIMESTAMPTZ DEFAULT now() NOT NULL,
    date_updated TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- User Animals Seen table
CREATE TABLE IF NOT EXISTS users_animals_seen (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL ,
    animal_id UUID REFERENCES animals(id) ON DELETE CASCADE NOT NULL ,
    liked BOOLEAN NOT NULL DEFAULT FALSE,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    date_updated TIMESTAMPTZ
);
