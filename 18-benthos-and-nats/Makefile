# Create seed customers and put them in a sqlite database
customers: database.db
	benthos -c seed/create_customers.yaml -e .env

# Create seed products and put them in a sqlite database
media:
	benthos -c seed/create_media.yaml -e .env

# Create seed ratings and put them in a kinesis stream
ratings:
	benthos -c seed/create_ratings.yaml -e .env

# Create seed views and put them in a kinesis stream
views:
	benthos -c seed/create_views.yaml -e .env

analytics:
	benthos -c analytics.yaml -e .env

database.db:
	sqlite3 database.db < db/schema.sql

clean:
	rm database.db
