# GraphQL Testing Guide

This document contains full example GraphQL queries and mutations for all services.

## Mutations

### Create Account
Mandatory fields: `email`, `password`. `name` is optional.

```graphql
mutation CreateAccount {
  createAccount(input: {
    name: "John Doe"
    email: "john@example.com"
    password: "password123"
  }) {
    id
    name
    email
  }
}
```

### Create Product
Create a new product in the catalog.

```graphql
mutation CreateProduct {
  createProduct(input: {
    name: "MacBook Pro"
    description: "M3 Max, 16-inch, 64GB RAM"
    price: 3499.99
  }) {
    id
    name
    description
    price
  }
}
```

### Create Order
Replace `ACCOUNT_ID` and `PRODUCT_ID` below.

```graphql
mutation CreateOrder {
  createOrder(input: {
    accountId: "ACCOUNT_ID"
    products: [
      {
        id: "PRODUCT_ID"
        quantity: 1
      }
    ]
  }) {
    id
    createdAt
    totalPrice
    products {
      id
      name
      quantity
      price
    }
  }
}
```

## Queries

### List Accounts
```graphql
query ListAccounts {
  accounts(pagination: { skip: 0, take: 10 }) {
    id
    name
    email
  }
}
```

### Get Account by ID
```graphql
query GetAccountByID {
  accounts(id: "ACCOUNT_ID") {
    id
    name
    email
  }
}
```

### List Products
```graphql
query ListProducts {
  products(pagination: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

### Search Products
```graphql
query SearchProducts {
  products(query: "MacBook") {
    id
    name
    description
    price
  }
}
```

### Get Order by ID
```graphql
query GetOrderByID {
  order(id: "ORDER_ID") {
    id
    createdAt
    totalPrice
    accountId
    products {
      id
      name
      quantity
      price
    }
  }
}
```

### List Orders for Account
```graphql
query GetOrdersForAccount {
  ordersForAccount(accountId: "ACCOUNT_ID") {
    id
    createdAt
    totalPrice
    products {
      id
      name
      quantity
      price
    }
  }
}
```

