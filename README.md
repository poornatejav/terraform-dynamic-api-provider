# terraform-dynamic-api-provider

This repository contains Terraform provider implementations that fetch data from any REST API and expose it as a Terraform data source. It demonstrates three different approaches to handling API responses:

Basic REST API Response Handling – Parses structured JSON and stores specific fields.
Dynamic JSON Handling – Allows dynamic JSON responses with key-value pairs.
Dynamic Struct Generation – Uses reflection to create Go structs dynamically based on the API response.



When to Use Which Implementation?
Approach	                      Use Case
Basic                           REST API Handling	When the API response structure is fixed and known.
Dynamic JSON Handling	          When the API returns varying JSON structures dynamically.
Dynamic Struct Handling	        When you need strongly typed access but don't know the structure beforehand.
