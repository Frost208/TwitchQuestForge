-- Twitch RPG Database Schema
-- Run this script to create the database structure

CREATE DATABASE IF NOT EXISTS twitch_rpg;
USE twitch_rpg;

-- Characters table - one per Twitch user
CREATE TABLE IF NOT EXISTS characters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    twitch_user_id VARCHAR(255) UNIQUE,
    level INT DEFAULT 1,
    experience INT DEFAULT 0,
    channel_points_spent INT DEFAULT 0,
    
    -- Base stats
    strength INT DEFAULT 10,
    agility INT DEFAULT 10,
    vitality INT DEFAULT 10,
    intelligence INT DEFAULT 10,
    
    -- Equipment slots (item IDs)
    boots_id INT DEFAULT NULL,
    pants_id INT DEFAULT NULL,
    armor_id INT DEFAULT NULL,
    helmet_id INT DEFAULT NULL,
    ring_id INT DEFAULT NULL,
    chain_id INT DEFAULT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (boots_id) REFERENCES items(id),
    FOREIGN KEY (pants_id) REFERENCES items(id),
    FOREIGN KEY (armor_id) REFERENCES items(id),
    FOREIGN KEY (helmet_id) REFERENCES items(id),
    FOREIGN KEY (ring_id) REFERENCES items(id),
    FOREIGN KEY (chain_id) REFERENCES items(id)
);

-- Items table - all equipment pieces
CREATE TABLE IF NOT EXISTS items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type ENUM('boots', 'pants', 'armor', 'helmet', 'ring', 'chain') NOT NULL,
    rarity ENUM('common', 'rare', 'epic', 'legendary') DEFAULT 'common',
    
    -- Stat bonuses
    strength_bonus INT DEFAULT 0,
    agility_bonus INT DEFAULT 0,
    vitality_bonus INT DEFAULT 0,
    intelligence_bonus INT DEFAULT 0,
    
    -- Special properties
    special_effect VARCHAR(500) DEFAULT NULL,
    value INT DEFAULT 100, -- Channel points value
    is_special BOOLEAN DEFAULT FALSE, -- For merchant items
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Character inventory - items owned by characters
CREATE TABLE IF NOT EXISTS character_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    item_id INT NOT NULL,
    quantity INT DEFAULT 1,
    acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id),
    UNIQUE KEY unique_character_item (character_id, item_id)
);

-- Combat logs for tracking fights
CREATE TABLE IF NOT EXISTS combat_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    attacker_id INT NOT NULL,
    defender_id INT NOT NULL,
    winner_id INT NOT NULL,
    
    -- Combat details
    attacker_power INT NOT NULL,
    defender_power INT NOT NULL,
    combat_log TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (attacker_id) REFERENCES characters(id),
    FOREIGN KEY (defender_id) REFERENCES characters(id),
    FOREIGN KEY (winner_id) REFERENCES characters(id)
);

-- Merchant events tracking
CREATE TABLE IF NOT EXISTS merchant_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_type ENUM('random_shop', 'special_trader') DEFAULT 'random_shop',
    available_items JSON, -- Array of item IDs available during this event
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- Quest system
CREATE TABLE IF NOT EXISTS quests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quest_type ENUM('daily', 'weekly', 'special') DEFAULT 'daily',
    requirements JSON, -- Flexible requirements system
    rewards JSON, -- Items, stats, or channel points
    channel_point_cost INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

-- Character quest progress
CREATE TABLE IF NOT EXISTS character_quests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    quest_id INT NOT NULL,
    progress JSON, -- Track quest progress
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (quest_id) REFERENCES quests(id),
    UNIQUE KEY unique_character_quest (character_id, quest_id)
);

-- Game events log for OBS integration
CREATE TABLE IF NOT EXISTS game_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_type ENUM('combat', 'merchant', 'level_up', 'item_acquired', 'quest_completed') NOT NULL,
    character_id INT,
    event_data JSON, -- Flexible event data for OBS
    obs_triggered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (character_id) REFERENCES characters(id)
);