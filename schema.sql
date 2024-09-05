-- 用户表
create table users (
    id integer not null primary key,
    email text not null unique,
    password text not null,
    username text not null,
    avatar text not null default '',
    -- role 角色：0:普通用户，1:管理员
    role integer not null default 0,
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text
);

-- 文章属性表
create table attributes (
    id integer not null primary key autoincrement,
    attr_name text not null
);

-- 文章表
create table posts (
    id integer not null primary key,
    author_id integer not null,
    title text not null,
    content text not null,
    attr_id integer not null default 1,
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text,
    foreign key(author_id) references users(id),
    foreign key(attr_id) references features(id)
);

-- 资源(文件)表
create table assets (
    id integer not null primary key,
    user_id integer not null,
    data blob not null,
    -- 1:IMAGE|2:VIDEO|3:AUDIO|4:DOCUMENT
    filetype integer no null default 1, 
    -- 文件名，带后缀（如：xxx.png）
    filename text not null,
    -- 文件大小
    size integer not null default 0,
    description text not null default "",
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text,
    foreign key(user_id) references users(id)
);

-- 标签表
create table tags (
    id integer not null primary key autoincrement,
    tag_name text not null unique,
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime'))
);

-- 文章标签关联表
create table posts_tags (
    post_id integer not null,
    tag_id integer not null,
    foreign key(post_id) references posts(id),
    foreign key(tag_id) references tags(id),
    primary key (post_id, tag_id)
);

-- 插入数据
insert into attributes (attr_name)
values ("原创"),("随笔"),("置顶"),("热门"),("漫谈"),("日常"),("专栏"),("转摘"),("摘抄"),("翻译");

INSERT INTO tags (tag_name) 
VALUES ("Go"), ("SQLite"), ("MySQL"), ("PostgreSQL"), ("Vue"), ("React"), ("TypeScript"), 
("JavaScript"), ("HTML"), ("CSS"), ("SQL");

