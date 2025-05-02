# Database Migrations

This directory contains database migration files for the PostgreSQL database used by Penfeel.

## Migration files

Each migration consists of an SQL file with a timestamp-based naming convention:
```
YYYYMMDDHHMMSS_migration_name.up.sql
```

For example:
```
20240428000001_create_users_table.up.sql
```

## Creating a New Migration

To create a new migration, use the `make` command:

```bash
make migration-new
```

This will prompt you for a migration name and create two files:
- `YYYYMMDDHHMMSS_your_migration_name.up.sql` - SQL to apply the migration
- `YYYYMMDDHHMMSS_your_migration_name.down.sql` - SQL to revert the migration

## Running Migrations

Migrations are automatically run during application startup. Both the `auth-service` and `document-service` will attempt to run migrations when they start.

If you want to manually run migrations, use:

```bash
make migration-up
```

## Rolling Back Migrations

To roll back migrations:

```bash
make migration-down
```

This will prompt you for the number of migrations to roll back.

## Checking Migration Status

To check the current migration status:

```bash
make migration-status
```

## Previewing Migrations

To see which migrations would be applied without actually running them:

```bash
make migration-plan
```

## Dumping Schema

To dump the current database schema to a file:

```bash
make dump-schema
```

This creates a `schema.sql` file containing the current database schema.

## Best Practices

1. **Keep migrations small and focused** - Each migration should do one thing and do it well.
2. **Migrations are forward-only** - While we have down migrations, in production we generally only migrate forward.
3. **Test migrations before applying to production** - Use the CI validation workflow.
4. **Never edit existing migrations** - Once committed, treat migrations as immutable.
5. **Include transactions when needed** - For data migrations, consider using transactions.
6. **Make migrations idempotent when possible** - Use `IF NOT EXISTS` and similar constructs.
7. **Add descriptive comments** - Document what the migration does and why.
8. **Coordinate schema changes** - When making breaking changes, coordinate with team members to prevent conflicts.

## Schema Validation in CI

Pull requests that modify migrations are automatically validated using a GitHub Actions workflow. This workflow:

1. Applies all migrations to a test database
2. Dumps the schema to verify it's valid
3. Compares the schema to the main branch to highlight changes

This helps prevent migration conflicts and ensures migrations can be applied successfully. 