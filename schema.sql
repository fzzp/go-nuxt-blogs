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


-- ==============================================后续添加=============================================


-- ========添加列表=========

-- 给用户增加一个个人介绍字段(bio)
alter table users add column bio text not null default "这个家伙很懒，什么都没有写～";

-- 给文章增加一个阅读数(views)
alter table posts add column views integer not null default 0;


-- 创建博客设置表
create table settings (
    id integer not null primary key autoincrement,
    total_posts integer not null default 0, -- 总文章数
    total_views integer not null default 0, -- 文章总阅读数
    total_comments integer not null default 0 -- 评论总数
);

insert into 
    settings (id, total_posts, total_views, total_comments)
values(1, 0, 0, 0);


-- ========添加触发器=========

-- 触发器参考：https://www.sqlitetutorial.net/sqlite-trigger/

-- 创建一个触发器增加 settings.total_posts
create trigger if not exists calc_post_count after insert on posts 
begin 
    update settings set total_posts = total_posts + 1 where id = 1;
end;

-- 创建一个触发器增加 settings.total_views
create trigger if not exists calc_post_view 
    after update on posts
    when old.views <> new.views 
begin
    update settings set total_views = total_views + 1 where id = 1;
end;

-- 查看触发器
select * from sqlite_master where type = 'trigger';