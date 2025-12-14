# Quick Setup Guide

## Issue Found

Your backend is failing to connect to Supabase due to IPv6 network issue.

## Fix Required

Add this to your `.env` file:

```bash
DB_SSLMODE=require
```

## Alternative: Use Supabase Connection Pooler

Instead of direct database connection, use Supabase's connection pooler which works better with Docker:

Change your `DB_HOST` in `.env` to use the pooler:

```bash
# Instead of:
# DB_HOST=db.aekaqrtibzgilkchnpia.supabase.co

# Use the pooler (IPv4):
DB_HOST=aws-0-ap-south-1.pooler.supabase.com
DB_PORT=6543
DB_SSLMODE=require
```

Get your exact pooler URL from:
Supabase Dashboard → Project Settings → Database → Connection Pooling

## Then restart:

```bash
docker-compose down
docker-compose up -d
```
