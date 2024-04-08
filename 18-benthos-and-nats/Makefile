# Create seed customers and put them in a sqlite database
customers: database.db
	benthos -w --chilled -c seed/create_customers.yaml -e .env

# Create seed products and put them in a sqlite database
products: database.db
	benthos -w --chilled -c seed/create_products.yaml -e .env

# Create seed ratings and put them in a kinesis stream
ratings:
	benthos -w --chilled -c seed/create_ratings.yaml -e .env

database.db:
	sqlite3 database.db < db/schema.sql

clean:
	rm database.db
