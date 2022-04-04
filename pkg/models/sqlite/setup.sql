-- User
CREATE TABLE IF NOT EXISTS "user" (
    "userid"  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "pass" TEXT NOT NULL
);

-- Post
CREATE TABLE IF NOT EXISTS "post" (
    "postid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "title" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "created" TEXT NOT NULL
    -- FOREIGN KEY ("authorid") REFERENCES "user" ("userid") ON DELETE CASCADE
);

-- Cookie
CREATE TABLE IF NOT EXISTS "cookie" (
    "userid" INTEGER NOT NULL, 
    "uuid" TEXT NOT NULL,
    "expirydate" TEXT NOT NULL,
    FOREIGN KEY ("userid") REFERENCES "user" ("userid") ON DELETE CASCADE
);
-- rate
CREATE TABLE IF NOT EXISTS "like" (
    "userid" INTEGER NOT NULL,
    "postid" INTEGER NOT NULL, 
    "symbol" TEXT NOT NULL,
    FOREIGN KEY ("userid") REFERENCES "user" ("userid") ON DELETE CASCADE,
    FOREIGN KEy ("postid") REFERENCES "post" ("postid") ON DELETE CASCADE
 );

-- comment
CREATE TABLE IF NOT EXISTS "comment" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "userid" INTEGER NOT NULL,
    "username" TEXT NOT NULL,
    "postid" INTEGER NOT NULL,
    "comment" TEXT NOT NULL,
    FOREIGN KEY ("userid") REFERENCES "user" ("userid") ON DELETE CASCADE,
    FOREIGN KEY ("postid") REFERENCES "post" ("postid") ON DELETE CASCADE
);

-- tag
CREATE TABLE IF NOT EXISTS "tag" (
    "postid" INTEGER NOT NULL,
    "tag" TEXT NOT NULL,
    FOREIGN KEY ("postid") REFERENCES "post" ("postid") ON DELETE CASCADE
);

-- rate comments
CREATE TABLE IF NOT EXISTS "rateComment" (
    "userid" INTEGER NOT NULL,
    "commentid" INTEGER NOT NULL, 
    "postid" INTEGER NOT NULL,
    "symbol" TEXT NOT NULL,
    FOREIGN KEY ("userid") REFERENCES "user" ("userid") ON DELETE CASCADE,
    FOREIGN KEY ("commentid") REFERENCES "comment" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("postid") REFERENCES "post" ("postid") ON DELETE CASCADE
);
